package procbuilder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

type DynFloPoCoFixedPoint struct{}

//------------------------------------------------------------------------------
// Public identity
//------------------------------------------------------------------------------
func (d DynFloPoCoFixedPoint) GetName() string { return "dyn_flopoco" }

//------------------------------------------------------------------------------
// 1. Accept names of the form  multflpfps<s>f<f>   or   addflpfps<s>f<f>
//------------------------------------------------------------------------------
func (d DynFloPoCoFixedPoint) MatchName(name string) bool {
	re := regexp.MustCompile(`^(mult|add)flpfps[0-9]+f[0-9]+$`)
	return re.MatchString(name)
}

//------------------------------------------------------------------------------
// 2. Build the operator, parse the report, return an Opcode
//------------------------------------------------------------------------------
func (d DynFloPoCoFixedPoint) CreateInstruction(name string) (Opcode, error) {

	//----------------------------------------------------------------------
	// Decode <op>, <s>, <f> from the instruction name
	//----------------------------------------------------------------------
	re := regexp.MustCompile(`^(mult|add)flpfps(?P<s>[0-9]+)f(?P<f>[0-9]+)$`)
	m := re.FindStringSubmatch(name)
	if m == nil {
		return nil, fmt.Errorf("unsupported operator name %q", name)
	}
	op, sStr, fStr := m[1], m[2], m[3]

	s, _ := strconv.Atoi(sStr) // integer bits (without sign)
	f, _ := strconv.Atoi(fStr) // fractional bits
	if s < 1 || f < 0 || s+f+1 > 64 {
		return nil, fmt.Errorf("bad fixed-point format in %q", name)
	}
	w := s + f + 1 // total word length (sign + int + frac)
	regSize := w

	//----------------------------------------------------------------------
	// Helper to compose the FloPoCo command
	//----------------------------------------------------------------------
	buildCmd := func(op string, w int) []string {
		switch op {
		case "mult":
			return []string{
				"flopoco",
				"IntMultiplier",
				fmt.Sprintf("wX=%d", w), fmt.Sprintf("wY=%d", w),
				"signedIO=1",
				"frequency=300", "pipeline=yes",
			}
		case "add":
			return []string{
				"flopoco",
				"IntAdder",
				fmt.Sprintf("wIn=%d", w),
				"frequency=300", "pipeline=yes",
			}
		default:
			return nil
		}
	}
	runCommand := buildCmd(op, w)

	//----------------------------------------------------------------------
	// Run FloPoCo in a temporary directory
	//----------------------------------------------------------------------
	dir, err := os.MkdirTemp("", "flopoco")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command(runCommand[0], runCommand[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("FloPoCo: %v\n%s", err, out)
	}
	resultReport := string(out)

	//----------------------------------------------------------------------
	// Parse entity names and pipeline depth from the report
	//----------------------------------------------------------------------
	var (
		entities []string
		pipeline int
	)
	for _, line := range strings.Split(resultReport, "\n") {
		if m := regexp.MustCompile(`Entity\s+(\w+)`).FindStringSubmatch(line); m != nil {
			entities = append(entities, m[1])
		}
		if m := regexp.MustCompile(`Pipeline depth = (\d+)`).FindStringSubmatch(line); m != nil {
			pipeline, _ = strconv.Atoi(m[1])
		}
	}
	if len(entities) == 0 {
		return nil, errors.New("could not find entity names in FloPoCo report")
	}
	topEntity := entities[len(entities)-1]

	//----------------------------------------------------------------------
	// Read the generated VHDL
	//----------------------------------------------------------------------
	vhdlBytes, err := os.ReadFile(filepath.Join(dir, "flopoco.vhdl"))
	if err != nil {
		return nil, err
	}
	vHDL := string(vhdlBytes)

	//----------------------------------------------------------------------
	// Return the wrapper that the rest of the tool-chain expects
	//----------------------------------------------------------------------
	return FloPoCo{
		floPoCoName: name,
		regSize:     regSize,
		vHDL:        vHDL,
		topEntity:   topEntity,
		entities:    entities,
		pipeline:    pipeline,
	}, nil
}

//------------------------------------------------------------------------------
// 3.  assembler helpers (left unchanged)
//------------------------------------------------------------------------------
func (d DynFloPoCoFixedPoint) HLAssemblerGeneratorMatch(bmc *bmconfig.BmConfig) []string {
	return []string{}
}

func (d DynFloPoCoFixedPoint) HLAssemblerGeneratorList(bmc *bmconfig.BmConfig, bl *bmline.BasmLine) []string {
	return []string{}
}
