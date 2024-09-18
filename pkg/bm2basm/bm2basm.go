package bm2basm

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

type Bm2Basm struct {
	debug bool
}

func (b *Bm2Basm) Convert(bm *bondmachine.Bondmachine) (string, error) {
	result := ""

	globalIoMode := ""

	for _, cpId := range bm.Processors {

		// Start by dissasembling the CP code
		machId := bm.Processors[cpId]
		mach := bm.Domains[machId]

		code, _ := mach.Disassembler()

		// Get the number of inputs and outputs
		cpInputs := int(mach.N)
		cpOutputs := int(mach.M)

		// mapping I/O resources
		inputMap := make(map[string]string)
		outputMap := make(map[string]string)

		// Create the input/output mapping by filling the maps with empty strings
		for i := 0; i < cpInputs; i++ {
			inputMap["i"+strconv.Itoa(i)] = ""
		}
		for i := 0; i < cpOutputs; i++ {
			outputMap["o"+strconv.Itoa(i)] = ""
		}

		// Get the first N free registers, N being the number of inputs + outputs
		// A proper implementation would need to check if the registers are actually free by analyzing the code
		// and checking if the registers are used

		freeRegs := FirstNFreeRegisters(code, cpInputs+cpOutputs)

		// Map the registers to the inputs and outputs
		resIn := ""
		for i := 0; i < cpInputs; i++ {
			inputMap["i"+strconv.Itoa(i)] = freeRegs[i]
			resIn += freeRegs[i]
			if i < cpInputs-1 {
				resIn += ":"
			}
		}
		resOut := ""
		for i := 0; i < cpOutputs; i++ {
			outputMap["o"+strconv.Itoa(i)] = freeRegs[cpInputs+i]
			resOut += freeRegs[cpInputs+i]
			if i < cpOutputs-1 {
				resOut += ":"
			}
		}

		syncOrAsync := ""
		if IsAllSync(code) {
			syncOrAsync = "iomode:sync"
			if globalIoMode == "" {
				globalIoMode = "iomode:sync"
			} else if globalIoMode == "iomode:async" {
				globalIoMode = "failed"
			}
		} else if IsAllAsync(code) {
			syncOrAsync = "iomode:async"
			if globalIoMode == "" {
				globalIoMode = "iomode:async"
			} else if globalIoMode == "iomode:sync" {
				globalIoMode = "failed"
			}
		}

		result += "%fragment cp" + strconv.Itoa(cpId) + "fragment " + syncOrAsync + "\n"
		result += "%meta literal resin " + resIn + "\n"
		result += "%meta literal resout " + resOut + "\n"
		spCode := strings.Split(code, "\n")

		// Remove empty lines at the end
		for {
			if spCode[len(spCode)-1] == "" {
				spCode = spCode[:len(spCode)-1]
			} else {
				break
			}
		}
		// Remove j 0 at the end
		if spCode[len(spCode)-1] == "j 0" {
			spCode = spCode[:len(spCode)-1]
		}

		for _, line := range spCode {
			if line != "" {
				tokens := strings.Split(line, " ")
				switch tokens[0] {
				case "i2r":
					result += "\tcpy " + tokens[1] + "," + inputMap[tokens[2]] + "\n"
				case "r2o":
					result += "\tcpy " + outputMap[tokens[2]] + "," + tokens[1] + "\n"
				case "i2rw":
					result += "\tcpy " + tokens[1] + "," + inputMap[tokens[2]] + "\n"
				case "r2owa":
					result += "\tcpy " + outputMap[tokens[2]] + "," + tokens[1] + "\n"
				default:
					result += "\t" + tokens[0] + " " + strings.Join(tokens[1:], ",") + "\n"
				}
			}
		}
		result += "%endfragment\n"
		result += "\n"
	}

	result += "\n"

	// Meta code
	result += "%meta bmdef global registersize:" + strconv.Itoa(int(bm.Rsize)) + "\n"
	if globalIoMode != "failed" {
		result += "%meta bmdef global " + globalIoMode + "\n"
	}

	for _, cpId := range bm.Processors {
		result += "%meta fidef cp" + strconv.Itoa(cpId) + "fi fragment:cp" + strconv.Itoa(cpId) + "fragment\n"
	}

	result += "\n"

	for i, bond := range bm.Links {
		result += "%meta filinkdef bond" + strconv.Itoa(i) + " type:fl\n"
		endPoint1 := bm.Internal_inputs[i]
		endPoint2 := bm.Internal_outputs[bond]
		e1 := ""
		e2 := ""
		switch endPoint1.Map_to {
		case bondmachine.BMOUTPUT:
			e1 = "fi:ext, type:output, index:" + strconv.Itoa(endPoint1.Res_id)
		case bondmachine.CPINPUT:
			e1 = "fi:cp" + strconv.Itoa(endPoint1.Res_id) + "fi, type:input, index:" + strconv.Itoa(endPoint1.Ext_id)
		}

		switch endPoint2.Map_to {
		case bondmachine.BMINPUT:
			e2 = "fi:ext, type:input, index:" + strconv.Itoa(endPoint2.Res_id)
		case bondmachine.CPOUTPUT:
			e2 = "fi:cp" + strconv.Itoa(endPoint2.Res_id) + "fi, type:output, index:" + strconv.Itoa(endPoint2.Ext_id)
		}
		result += "%meta filinkatt bond" + strconv.Itoa(i) + " " + e1 + "\n"
		result += "%meta filinkatt bond" + strconv.Itoa(i) + " " + e2 + "\n"
		result += "\n"

	}

	for _, cpId := range bm.Processors {
		result += "%meta cpdef cp" + strconv.Itoa(cpId) + " fragcollapse:cp" + strconv.Itoa(cpId) + "fi\n"
	}

	return result, nil
}

func FirstNFreeRegisters(code string, n int) []string {
	usedRegs := make(map[string]struct{})
	for _, line := range strings.Split(code, "\n") {
		for _, token := range strings.Split(line, " ") {
			re := regexp.MustCompile(`^r[0-9]+$`)
			if re.MatchString(token) {
				usedRegs[token] = struct{}{}
			}
		}
	}

	result := make([]string, n)
	for i, j := 0, 0; i < n; i++ {
		for {
			rNumS := "r" + strconv.Itoa(j)
			if _, ok := usedRegs[rNumS]; !ok {
				result[i] = rNumS
				j++
				break
			}
			j++
		}
	}

	return result
}

func IsAllSync(code string) bool {
	for _, line := range strings.Split(code, "\n") {
		if strings.Contains(line, "i2r") {
			return false
		}
		if strings.Contains(line, "r2o") {
			return false
		}
	}
	return true
}

func IsAllAsync(code string) bool {
	for _, line := range strings.Split(code, "\n") {
		if strings.Contains(line, "i2rw") {
			return false
		}
		if strings.Contains(line, "r2owa") {
			return false
		}
	}
	return true
}
