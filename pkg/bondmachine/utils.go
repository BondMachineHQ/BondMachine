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

	// if len(input) > 2 {
	// 	if input[:2] == "0x" {
	// 		// TODO Hex
	// 	} else if input[:2] == "0b" {
	// 		// TODO Binary
	// 	} else if input[:2] == "0d" {
	// 		// Decimal (also the default)
	// 		if s, err := strconv.Atoi(input[2:]); err == nil {
	// 			return uint64(s), nil
	// 		} else {
	// 			return 0, errors.New("invalid number" + input)
	// 		}
	// 	} else if input[:2] == "0f" {
	// 		// Float32
	// 		if s, err := strconv.ParseFloat(input[2:], 32); err == nil {
	// 			return uint64(math.Float32bits(float32(s))), nil
	// 		} else {
	// 			return 0, errors.New("unknown float32 number " + input)
	// 		}
	// 	} else {
	// 		// Decimal (also the default)
	// 		if s, err := strconv.Atoi(input); err == nil {
	// 			return uint64(s), nil
	// 		} else {
	// 			return 0, errors.New("invalid number" + input)
	// 		}

	// 	}
	// } else {
	// 	// Decimal (also the default)
	// 	if s, err := strconv.Atoi(input); err == nil {
	// 		return uint64(s), nil
	// 	} else {
	// 		return 0, errors.New("invalid number" + input)
	// 	}
	// }

	// return 0, errors.New("unknown number format " + input)
}
