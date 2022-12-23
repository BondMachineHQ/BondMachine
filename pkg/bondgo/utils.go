package bondgo

import (
	"regexp"
	"strconv"
)

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
		//case "r":
		//	types = C_REGISTER
		case "i":
			types = C_INPUT
		case "o":
			types = C_OUTPUT
			//case "ch":
			//	types = C_CHANNEL
		}
	}

	return result, types
}
