package bmbuilder

import (
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

const (
	blockSequential = uint8(0) + iota
)

type Generator func(*BMBuilder, *bmline.BasmElement, *bmline.BasmLine) (*bondmachine.Bondmachine, error)
type Connector func(*BMBuilder, *bmline.BasmElement, *bmline.BasmLine) (*BMConnections, error)

// The instance
type BMBuilder struct {
	verbose       bool
	debug         bool
	isWithinBlock string
	isSymbolled   string
	currentBlock  string
	passes        uint64
	lineMeta      string
	blocks        map[string]*BMBuilderBlock
	global        *bmline.BasmElement
	generators    map[string]Generator
	connectors    map[string]Connector
	result        *bondmachine.Bondmachine
}

// Sections and Macros

type BMBuilderBlock struct {
	blockName string
	blockType uint8
	blockBody *bmline.BasmBody
	blockBMs  []*bondmachine.Bondmachine
	blockConn []*BMConnections // It is a len(blockBMs) + 1 length array
}

func (bld *BMBuilder) BMBuilderInit() {

	bld.verbose = false
	bld.debug = false
	bld.isWithinBlock = ""
	bld.blocks = make(map[string]*BMBuilderBlock)

	if bld.debug {
		fmt.Println(purple("Init") + ":")
	}
	bld.global = new(bmline.BasmElement)
	bld.global.SetValue("global")

	bld.passes = passMetaExtractor |
		passGeneratorsExec |
		0

	generators := make(map[string]Generator)
	generators["basm"] = BasmGenerator
	generators["h"] = HadamardGenerator
	generators["cx"] = CnotGenerator
	generators["zero"] = ZeroGenerator
	generators["maxpool"] = MaxPoolGenerator

	connectors := make(map[string]Connector)
	connectors["connector"] = ConnectorDefault

	bld.generators = generators
	bld.connectors = connectors

}

// SetVerbose sets the verbose flag on a BMBuilder
func (bld *BMBuilder) SetVerbose() {
	bld.verbose = true
}

// SetDebug sets the debug flag on a BMBuilder
func (bld *BMBuilder) SetDebug() {
	bld.debug = true
}

func (m *BMBuilderBlock) writeText() string {
	result := ""
	for _, line := range m.blockBody.Lines {
		result += line.Operation.GetValue() + " "
		for _, element := range line.Elements {
			result += element.GetValue()
			if element != line.Elements[len(line.Elements)-1] {
				result += ", "
			}
		}
		result += "\n"
	}
	return result
}

func (m *BMBuilderBlock) String() string {
	result := "\t\t" + red(m.blockName)
	result += m.blockBody.String()
	return result
}

func (bi *BMBuilder) String() string {
	result := purple("Instance Dump:") + "\n"
	if bi.global != nil {
		result += purple("\tBM meta") + ":\n"
		result += fmt.Sprintf("\t\t") + bi.global.String() + "\n"
	}
	if len(bi.blocks) > 0 {
		result += purple("\tBlocks") + ":\n"
		for _, block := range bi.blocks {
			result += block.String() + "\n"
		}
	}
	return result
}

func (bi *BMBuilder) RunBuilder() error {

	if bi.debug {
		fmt.Println(purple("RunBuilder") + ":")
	}

	step := uint64(1)
	names := getPassFunctionName()

	if bi.debug {
		fmt.Println(purple("Pre phase 1 ") + bi.String())
	}

	passes := getPassFunction()
	for i := 0; i < 64; i++ {
		if bi.ActivePass(step) {
			if bi.debug {
				fmt.Println(purple("Phase "+strconv.Itoa(i+1)) + ": " + red(names[step], " started"))
			}
			currentPass := passes[step]
			if err := currentPass(bi); err != nil {
				return err
			} else {
				if bi.debug {
					fmt.Println(purple("Phase "+strconv.Itoa(i+1)) + ": " + red(names[step], " completed"))
					fmt.Println(purple("Post phase "+strconv.Itoa(i+1)) + " " + bi.String())
				}
			}
		} else {
			if bi.debug {
				fmt.Println(purple("Phase "+strconv.Itoa(i+1)) + ": " + red(names[step], " is disabled"))
			}
		}
		if step == LAST_PASS {
			break
		}
		step = step << 1
	}

	return nil
}

func (bi *BMBuilder) ExportBasmBody() (*bmline.BasmBody, error) {
	// Get the main block
	mainBlock := bi.global.GetMeta("main")
	if mainBlock == "" {
		return nil, fmt.Errorf("No main block found")
	}

	// Get the main block
	return bi.blocks[mainBlock].blockBody, nil
}

func (bi *BMBuilder) GetBondMachine() *bondmachine.Bondmachine {
	return bi.result
}
