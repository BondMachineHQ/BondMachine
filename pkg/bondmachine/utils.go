package bondmachine

import (
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

func (bmach *Bondmachine) Get_so_name(id int) (string, bool) {

	if len(bmach.Shared_objects) > id {

		seq := make(map[string]int)

		for i, so := range bmach.Shared_objects {
			sname := so.Shortname()
			if _, ok := seq[sname]; ok {
				seq[sname]++
			} else {
				seq[sname] = 0
			}
			if i == id {
				return sname + strconv.Itoa(seq[sname]), true
			}
		}
	}
	return "", false
}

func (bmach *Bondmachine) GetInternalInputName(id int) (string, bool) {

	if len(bmach.Internal_inputs) > id {
		bond := bmach.Internal_inputs[id]
		switch bond.Map_to {
		case BMINPUT:
			return "i" + strconv.Itoa(bond.Res_id), true
		case BMOUTPUT:
			return "o" + strconv.Itoa(bond.Res_id), true
		case CPINPUT:
			return "p" + strconv.Itoa(bond.Res_id) + "i" + strconv.Itoa(bond.Ext_id), true
		case CPOUTPUT:
			return "p" + strconv.Itoa(bond.Res_id) + "o" + strconv.Itoa(bond.Ext_id), true
		}
	}
	return "", false
}

func (bmach *Bondmachine) GetInternalOutputName(id int) (string, bool) {

	if len(bmach.Internal_outputs) > id {
		bond := bmach.Internal_outputs[id]
		switch bond.Map_to {
		case BMINPUT:
			return "i" + strconv.Itoa(bond.Res_id), true
		case BMOUTPUT:
			return "o" + strconv.Itoa(bond.Res_id), true
		case CPINPUT:
			return "p" + strconv.Itoa(bond.Res_id) + "i" + strconv.Itoa(bond.Ext_id), true
		case CPOUTPUT:
			return "p" + strconv.Itoa(bond.Res_id) + "o" + strconv.Itoa(bond.Ext_id), true
		}
	}
	return "", false
}

func Get_input_name(i int) string {
	result := "i" + strconv.Itoa(i)
	return result
}

func Get_output_name(i int) string {
	result := "o" + strconv.Itoa(i)
	return result
}

func Needed_bits(num int) int {
	if num > 0 {
		for bits := 1; true; bits++ {
			if 1<<uint8(bits) >= num {
				return bits
			}
		}
	}
	return 0
}

func zeros_prefix(num int, value string) string {
	result := value
	for i := 0; i < num-len(value); i++ {
		result = "0" + result
	}
	return result
}

func zeros_suffix(num int, value string) string {
	result := value
	for i := 0; i < num-len(value); i++ {
		result = result + "0"
	}
	return result
}

func get_binary(i int) string {
	result := strconv.FormatInt(int64(i), 2)
	return result
}

func ImportNumber(c *Config, input string) (uint64, error) {

	if bmNumber, err := bmnumbers.ImportString(input); err == nil {
		return bmNumber.ExportUint64()
	} else {
		return 0, err
	}
}
