package bmnumbers

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type FloPoCo struct {
	floPoCoName string
	e           int
	f           int
}

func (d FloPoCo) getName() string {
	return d.floPoCoName
}

func (d FloPoCo) getInfo() string {
	return ""
}

func (d FloPoCo) importMatchers() map[string]ImportFunc {
	result := make(map[string]ImportFunc)

	result["^0flp<(?P<e>[0-9]+),(?P<f>[0-9]+)>(?P<number>.+)$"] = floPoCoImport

	return result
}

func (d FloPoCo) Convert(n *BMNumber) error {
	convertFrom := n.nType.getName()

	switch convertFrom {
	default:
		return errors.New("cannot convert from " + convertFrom + " to " + d.getName())
	}
	return nil
}

func floPoCoImport(re *regexp.Regexp, input string) (*BMNumber, error) {
	number := re.ReplaceAllString(input, "${number}")
	es := re.ReplaceAllString(input, "${e}")
	fs := re.ReplaceAllString(input, "${f}")
	e, _ := strconv.Atoi(es)
	f, _ := strconv.Atoi(fs)

	EventuallyCreateType("flpe"+es+"f"+fs, nil)

	runCommand := []string{"fp2bin", es, fs, number}

	// Create a temporary directory for the FloPoCo files
	dir, err := os.MkdirTemp("", "fp2bin")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command(runCommand[0], runCommand[1:]...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, err
	} else {
		resultReport := string(out)

		// Parse every line of the report
		for _, line := range strings.Split(resultReport, "\n") {
			re = regexp.MustCompile(`[0-1]+`)
			if re.MatchString(line) {
				binNum := line
				var nextByte string
				newNumber := BMNumber{}
				newNumber.number = make([]byte, (len(binNum)-1)/8+1)
				newNumber.bits = len(binNum)
				newNumber.nType = FloPoCo{floPoCoName: "flpe" + es + "f" + fs, e: e, f: f}

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
			}
		}
	}

	return nil, nil
}
