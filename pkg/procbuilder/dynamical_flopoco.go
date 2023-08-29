package procbuilder

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

type DynFloPoCo struct{}

func (d DynFloPoCo) GetName() string {
	return "dyn_flopoco"
}

func (d DynFloPoCo) MatchName(name string) bool {
	re := regexp.MustCompile("multflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("addflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		return true
	}
	re = regexp.MustCompile("divflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	return re.MatchString(name)
}

func (d DynFloPoCo) CreateInstruction(name string) (Opcode, error) {
	var regSize int
	var runCommand []string
	var resultReport string
	var vHDL string
	entities := make([]string, 0)
	var topEntity string
	var pipeline int

	// Create a temporary directory for the FloPoCo files
	dir, err := os.MkdirTemp("", "flopoco")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	re := regexp.MustCompile("multflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		es := re.ReplaceAllString(name, "${e}")
		fs := re.ReplaceAllString(name, "${f}")
		e, _ := strconv.Atoi(es)
		f, _ := strconv.Atoi(fs)
		runCommand = []string{"flopoco", "FPMult", "we=" + es, "wf=" + fs}
		regSize = e + f + 3
	}
	re = regexp.MustCompile("addflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		es := re.ReplaceAllString(name, "${e}")
		fs := re.ReplaceAllString(name, "${f}")
		e, _ := strconv.Atoi(es)
		f, _ := strconv.Atoi(fs)
		runCommand = []string{"flopoco", "FPAdd", "we=" + es, "wf=" + fs}
		regSize = e + f + 3
	}
	re = regexp.MustCompile("divflpe(?P<e>[0-9]+)f(?P<f>[0-9]+)")
	if re.MatchString(name) {
		es := re.ReplaceAllString(name, "${e}")
		fs := re.ReplaceAllString(name, "${f}")
		e, _ := strconv.Atoi(es)
		f, _ := strconv.Atoi(fs)
		runCommand = []string{"flopoco", "FPDiv", "we=" + es, "wf=" + fs}
		regSize = e + f + 3
	}

	cmd := exec.Command(runCommand[0], runCommand[1:]...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, err
	} else {
		resultReport = string(out)
		// Parse every line of the report
		for _, line := range strings.Split(resultReport, "\n") {
			re = regexp.MustCompile(`.*Entity (?P<entity>\w+).*`)
			if re.MatchString(line) {
				newEnt := re.ReplaceAllString(line, "${entity}")
				entities = append(entities, newEnt)
			}
			re = regexp.MustCompile(`.*Pipeline depth = (?P<depth>[0-9]+).*`)
			if re.MatchString(line) {
				depth := re.ReplaceAllString(line, "${depth}")
				pipeline, _ = strconv.Atoi(depth)
			}
		}

		topEntity = entities[len(entities)-1]

		// Read the VHDL file
		f, err := os.ReadFile(dir + "/flopoco.vhdl")
		if err != nil {
			return nil, err
		}
		vHDL = string(f)

	}

	return FloPoCo{floPoCoName: name, regSize: regSize, vHDL: vHDL, topEntity: topEntity, entities: entities, pipeline: pipeline}, nil
}

func (d DynFloPoCo) HLAssemblerGeneratorMatch(c *DynConfig) []string {
	result := make([]string, 0)
	return result
}

func (d DynFloPoCo) HLAssemblerGeneratorList(c *DynConfig, line *bmline.BasmLine) []string {
	result := make([]string, 0)
	return result
}
