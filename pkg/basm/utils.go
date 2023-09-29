package basm

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

func purple(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[35m" + ins + "\033[0m"
}
func green(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[32m" + ins + "\033[0m"
}
func yellow(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[33m" + ins + "\033[0m"
}
func blue(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[34m" + ins + "\033[0m"
}
func red(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[31m" + ins + "\033[0m"
}
func cyan(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[36m" + ins + "\033[0m"
}
func gray(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[37m" + ins + "\033[0m"
}
func white(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[97m" + ins + "\033[0m"
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

// Debug helper functrion
func (bi *BasmInstance) Debug(logline ...interface{}) {
	if bi.debug {
		log.Println(purple("[Debug]")+" -", logline)
	}
}

// Info helper function
func (bi *BasmInstance) Info(logline ...interface{}) {
	if bi.verbose || bi.debug {
		log.Println(green("[Info]")+" -", logline)
	}
}

// Warning helper function
func (bi *BasmInstance) Warning(logline ...interface{}) {
	log.Println(yellow("[Warn]")+" -", logline)
}

// Alert helper function
func (bi *BasmInstance) Alert(logline ...interface{}) {
	log.Println(red("[Alert]")+" -", logline)
}

func object2Bytes(obj string, numeric bool) ([]string, error) {
	if numeric {
		if bmNumber, err := bmnumbers.ImportString(obj); err != nil {
			return []string{}, err
		} else {
			b := bmNumber.GetBytes()
			result := make([]string, len(b))
			for i, ch := range b {
				result[i] = fmt.Sprintf("0x%x", ch)
			}
			return result, nil
		}
	} else {
		result := make([]string, len(obj))
		for i, ch := range obj {
			result[i] = fmt.Sprintf("0x%x", ch)
		}
		return result, nil
	}
}

func dbDataConverter(line string) ([]string, error) {
	trimmed := strings.TrimSpace(line)

	result := make([]string, 0)

	withinString := false
	curElem := ""
	curNumeric := true

	for i := 0; i < len(trimmed); i++ {
		ch := trimmed[i]
		switch ch {
		case '"':
			withinString = !withinString
			if withinString {
				curNumeric = false
			}
		case ',':
			if withinString {
				curElem += string(ch)
			} else {
				if curElem != "" {
					if converted, err := object2Bytes(curElem, curNumeric); err != nil {
						return []string{}, err
					} else {
						result = append(result, converted...)
						curElem = ""
						curNumeric = true
					}
				}
			}
		default:
			if withinString {
				curElem += string(ch)
			} else {
				curElem += strings.TrimSpace(string(ch))
			}
		}
	}

	if withinString {
		return []string{}, errors.New("a string has not been closed")
	}

	if converted, err := object2Bytes(curElem, curNumeric); err != nil {
		return []string{}, err
	} else {
		result = append(result, converted...)
	}

	return result, nil
}

func (bi *BasmInstance) GetMeta(req string) (string, error) {
	// TODO: Temporary
	splitRequest := strings.Split(req, ".")

	if len(splitRequest) > 0 {
		switch splitRequest[0] {
		case "sodef":
			for _, so := range bi.sos {
				fmt.Print(so.GetMeta("constraint"))
			}

		default:
			return "", errors.New("unknown or unimplemented meta request " + req)
		}
	}

	return "", nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func nextRes(res string) string {
	ty := res[0]
	num, _ := strconv.Atoi(res[1:])
	return string(ty) + strconv.Itoa(num+1)
}

func compareStrings(str1, str2 string) bool {
	// Extract the numbers from the strings
	num1, err1 := extractNumber(str1)
	num2, err2 := extractNumber(str2)

	// If both strings contain numbers, compare them as numbers
	if err1 == nil && err2 == nil {
		return num1 < num2
	}

	// Otherwise, compare them as strings
	return str1 < str2
}

func extractNumber(str string) (int, error) {
	// Find the index of the first digit
	tokens := strings.FieldsFunc(str, func(r rune) bool {
		return !('0' <= r && r <= '9')
	})

	// If no numeric tokens were found, return an error
	if len(tokens) == 0 {
		return 0, fmt.Errorf("No numeric tokens found in %s", str)
	}

	// Otherwise, return the last token as an integer
	return strconv.Atoi(tokens[len(tokens)-1])
}
