package bmnumbers

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type LinearQuantizer struct {
	linearQuantizerName string
}

func (d LinearQuantizer) GetName() string {
	return d.linearQuantizerName
}

func (d LinearQuantizer) getInfo() string {
	return ""
}

func (d LinearQuantizer) getSize() int {
	return 1 // TODO
}

func (d LinearQuantizer) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0lq<(?P<e>[0-9]+)>(?P<number>.+)$"] = floPoCoImport

	return result
}

func (d LinearQuantizer) Convert(n *BMNumber) error {
	convertFrom := n.nType.GetName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.GetName())
	}
	return nil
}

func linearQuantizerImport(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")
	es := re.ReplaceAllString(input, "${e}")
	e, _ := strconv.Atoi(es)

	EventuallyCreateType("flpe"+es+"f"+fs, nil)

	binNum := line
	var nextByte string
	newNumber := BMNumber{}
	newNumber.number = make([]byte, (len(binNum)-1)/8+1)
	newNumber.bits = len(binNum)
	newNumber.nType = LinearQuantizer{linearQuantizerName: "flpe" + es + "f" + fs}

	for i := 0; ; i++ {

		if len(binNum) > 8 {
			nextByte = binNum[len(binNum)-8:]
			binNum = binNum[:len(binNum)-8]
		} else if len(binNum) > 0 {
			nextByte = binNum
			binNum = ""
		} else {
			break
		}

		if val, err := strconv.ParseUint(nextByte, 2, 8); err == nil {
			decoded := byte(val)
			newNumber.number[i] = decoded
		} else {
			return nil, errors.New("invalid binary number" + input)
		}
	}
	return &newNumber, nil

	return nil, nil
}

func (d LinearQuantizer) ExportString(n *BMNumber) (string, error) {
	e := d.e
	f := d.f
	es := strconv.Itoa(e)
	fs := strconv.Itoa(f)
	number, _ := n.ExportBinary(false)

	// Add leading zeros
	for len(number) < e+f+3 {
		number = "0" + number
	}

	runCommand := []string{"bin2fp", es, fs, number}
	// fmt.Println(runCommand)
	// Create a temporary directory for the LinearQuantizer files
	dir, err := os.MkdirTemp("", "bin2fp")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command(runCommand[0], runCommand[1:]...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", err
	} else {
		resultReport := string(out)

		// Parse every line of the report
		for _, line := range strings.Split(resultReport, "\n") {
			result := "0flp<" + es + "." + fs + ">" + line
			return result, nil
		}
	}
	return "", errors.New("invalid binary number" + number)
}

func (d LinearQuantizer) ShowInstructions() map[string]string {
	return d.instructions
}

func (d LinearQuantizer) ShowPrefix() string {
	return "0flp<" + strconv.Itoa(d.e) + "." + strconv.Itoa(d.f) + ">"
}
