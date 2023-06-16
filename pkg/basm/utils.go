package basm

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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
		result := make([]string, 0)
		if len(obj) > 2 {
			if strings.HasPrefix(obj, "0x") {
				// Hex
				hexstring := obj[2:]
				if len(obj)%2 == 1 {
					hexstring = "0" + hexstring
				}
				for i := 0; i < len(hexstring); i = i + 2 {
					result = append(result, "0x"+hexstring[i:i+2])
				}
				return result, nil
			} else if obj[:2] == "0b" {
				// TODO
			} else if obj[:2] == "0d" {
				// TODO
			} else if obj[:2] == "0f" {
				// TODO
			}
		}
		return []string{}, errors.New("Unknown number format " + obj)
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
	currElem := ""
	currNumeric := true

	for i := 0; i < len(trimmed); i++ {
		ch := trimmed[i]
		switch ch {
		case '"':
			withinString = !withinString
			if withinString {
				currNumeric = false
			}
		case ',':
			if withinString {
				currElem += string(ch)
			} else {
				if currElem != "" {
					if converted, err := object2Bytes(currElem, currNumeric); err != nil {
						return []string{}, err
					} else {
						for _, el := range converted {
							result = append(result, el)
						}
						currElem = ""
						currNumeric = true
					}
				}
			}
		default:
			if withinString {
				currElem += string(ch)
			} else {
				currElem += strings.TrimSpace(string(ch))
			}
		}
	}

	if withinString {
		return []string{}, errors.New("A string has not been closed")
	}

	if converted, err := object2Bytes(currElem, currNumeric); err != nil {
		return []string{}, err
	} else {
		for _, el := range converted {
			result = append(result, el)
		}
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
