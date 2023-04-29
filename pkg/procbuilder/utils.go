package procbuilder

import (
	"regexp"
	"sort"
	"strconv"
	"unsafe"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

const (
	O_REGISTER = uint8(0) + iota
	O_INPUT
	O_OUTPUT
	O_CHANNEL
)

// TODO Maybe two letters registers are not enough, maybe something like r5 or r546 is more preferreble
//
//	func Get_register_name(i int) string {
//		start_0 := 97
//		start_1 := 97
//
//		div := i / 26
//		mod := i % 26
//
//		start_0 = start_0 + mod
//		start_1 = start_1 + div
//
//		return string(start_1) + string(start_0)
//	}

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

func Get_register_name(i int) string {
	return "r" + strconv.Itoa(i)
}

func Get_channel_name(i int) string {
	return "ch" + strconv.Itoa(i)
}

func Get_input_name(i int) string {
	result := "i" + strconv.Itoa(i)
	return result
}

func Get_output_name(i int) string {
	result := "o" + strconv.Itoa(i)
	return result
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

func get_id(intr string) int {
	result := 0
	for i := 0; i < len(intr); i++ {
		if intr[len(intr)-1-i] == '1' {
			result = result + 1<<uint8(i)
		}
	}
	return result
}

func get_binary(i int) string {
	result := strconv.FormatInt(int64(i), 2)
	return result
}

func Process_number(input string) (string, error) {

	if bmNumber, err := bmnumbers.ImportString(input); err == nil {
		return bmNumber.ExportBinary(false)
	} else {
		return "", err
	}

	// if len(input) > 2 {
	// 	if input[:2] == "0x" {
	// 		// Hex
	// 		hexstring := input[2:]
	// 		cut := false
	// 		if len(input)%2 == 1 {
	// 			hexstring = "0" + hexstring
	// 			cut = true
	// 		}
	// 		if decoded, err := hex.DecodeString(hexstring); err == nil {
	// 			result := ""
	// 			for j := len(decoded) - 1; j >= 0; j-- {
	// 				ch := decoded[j]
	// 				for i := 0; i < 8; i++ {
	// 					if ch%2 == 0 {
	// 						result = "0" + result
	// 					} else {
	// 						result = "1" + result
	// 					}
	// 					ch = ch >> 1
	// 				}
	// 			}
	// 			if cut {
	// 				return result[4:], nil
	// 			} else {
	// 				return result, nil
	// 			}
	// 		} else {
	// 			return "", Prerror{"Unknown hex number " + input}
	// 		}
	// 	} else if input[:2] == "0b" {
	// 		// Binary
	// 		for _, ch := range input[2:] {
	// 			if ch != '0' && ch != '1' {
	// 				return "", Prerror{"Unknown binary number " + input}
	// 			}
	// 		}
	// 		return input[2:], nil
	// 	} else if input[:2] == "0d" {
	// 		// Decimal (also the default)
	// 		input = input[2:]
	// 	} else if input[:2] == "0f" {
	// 		// Float32
	// 		if s, err := strconv.ParseFloat(input[2:], 32); err == nil {
	// 			result := ""
	// 			buf := new(bytes.Buffer)
	// 			converted := float32(s)
	// 			err := binary.Write(buf, binary.LittleEndian, converted)
	// 			if err != nil {
	// 				return "", Prerror{"Unknown float32 number " + input}
	// 			}
	// 			for _, ch := range buf.Bytes() {
	// 				for i := 0; i < 8; i++ {
	// 					if ch%2 == 0 {
	// 						result = "0" + result
	// 					} else {
	// 						result = "1" + result
	// 					}
	// 					ch = ch >> 1
	// 				}
	// 			}
	// 			return result, nil
	// 		} else {
	// 			return "", Prerror{"Unknown float32 number " + input}
	// 		}
	// 	}
	// }

	// if result_int, err := strconv.Atoi(input); err == nil {
	// 	return get_binary(result_int), nil
	// }

	// return "", Prerror{"Unknown number format " + input}
}

func Process_input(iregname string, input_num int) (string, error) {
	for i := 0; i < input_num; i++ {
		if Get_input_name(i) == iregname {
			return get_binary(i), nil
		}
	}
	return "", Prerror{"Unknown input register name"}
}

func Process_output(iregname string, input_num int) (string, error) {
	for i := 0; i < input_num; i++ {
		if Get_output_name(i) == iregname {
			return get_binary(i), nil
		}
	}
	return "", Prerror{"Unknown input register name"}
}

func Process_shared(soshort string, soname string, num int) (string, error) {
	for i := 0; i < num; i++ {
		if soshort+strconv.Itoa(i) == soname {
			return get_binary(i), nil
		}
	}
	return "", Prerror{"Unknown shared object name"}
}

func Sequence_to_0(start string) ([]string, uint8) {

	var result []string
	var types uint8 = 255

	re := regexp.MustCompile("(?P<obj>(o|i|r|ch))(?P<value>[0-9]+)")
	if re.MatchString(start) {
		obj := re.ReplaceAllString(start, "${obj}")

		value_string := re.ReplaceAllString(start, "${value}")
		tempvalue, _ := strconv.Atoi(value_string)
		result = make([]string, tempvalue+1)
		for i := 0; i < tempvalue+1; i++ {
			result[i] = obj + strconv.Itoa(i)
		}
		switch obj {
		case "r":
			types = O_REGISTER
		case "i":
			types = O_INPUT
		case "o":
			types = O_OUTPUT
		case "ch":
			types = O_CHANNEL
		}
	}

	return result, types
}

func (arch *Arch) OnlyOne(curOp string, ops []string) bool {
	sorted := make([]string, len(ops))
	copy(sorted, ops)
	sort.Strings(sorted)
	for _, op := range sorted {
		if op == curOp {
			return true
		} else {
			for _, op2 := range arch.Conproc.Op {
				if op2.Op_get_name() == op {
					return false
				}
			}
		}
	}
	return false
}

func Int8bits(f int8) uint8 {
	return *(*uint8)(unsafe.Pointer(&f))
}

func Int16bits(f int16) uint16 {
	return *(*uint16)(unsafe.Pointer(&f))
}

func Int32bits(f int32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}

func Int64bits(f int64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func Int8FromBits(f uint8) int8 {
	return *(*int8)(unsafe.Pointer(&f))
}

func Int16FromBits(f uint16) int16 {
	return *(*int16)(unsafe.Pointer(&f))
}

func Int32FromBits(f uint32) int32 {
	return *(*int32)(unsafe.Pointer(&f))
}

func Int64FromBits(f uint64) int64 {
	return *(*int64)(unsafe.Pointer(&f))
}
