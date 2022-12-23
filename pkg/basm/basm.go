package basm

import (
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

const (
	sectRomData = uint8(0) + iota
	setcRomText
)

// The instance

type BasmInstance struct {
	*bminfo.BMinfo
	verbose          bool
	debug            bool
	isWithinMacro    string
	isWithinSection  string
	isWithinFragment string
	isWithinChunk    string
	isLabelled       string
	lineMeta         string
	macros           map[string]*BasmMacro
	sections         map[string]*BasmSection
	fragments        map[string]*BasmFragment
	chunks           map[string]*BasmChunk
	passes           uint64
	matchers         []*bmline.BasmLine
	matchersOps      []procbuilder.Opcode
	bm               *bondmachine.Bondmachine
	result           *bondmachine.Bondmachine
	rg               *bmreqs.ReqRoot
	global           *bmline.BasmElement
	cps              []*bmline.BasmElement
	sos              []*bmline.BasmElement
	ios              []*bmline.BasmElement
	ioAttach         []*bmline.BasmElement
	soAttach         []*bmline.BasmElement
	fis              []*bmline.BasmElement // Fragment instances
	fiLinks          []*bmline.BasmElement // Fragment links
	fiLinkAttach     []*bmline.BasmElement // Fragment link attachments
}

// Sections and Macros

type BasmSection struct {
	sectionName string
	sectionType uint8
	sectionBody *bmline.BasmBody
}

type BasmFragment struct {
	fragmentName string
	fragmentBody *bmline.BasmBody
}

type BasmChunk struct {
	chunkName string
	chunkBody *bmline.BasmBody
}

type BasmMacro struct {
	macroName string
	macroArgs int
	macroBody *bmline.BasmBody
}

// Initialization of the instance

func (bi *BasmInstance) BasmInstanceInit(bm *bondmachine.Bondmachine) {

	// bi.verbose = false
	// bi.debug = false
	bi.isWithinMacro = ""
	bi.isWithinSection = ""
	bi.isWithinFragment = ""
	bi.isLabelled = ""
	bi.macros = make(map[string]*BasmMacro)
	bi.sections = make(map[string]*BasmSection)
	bi.fragments = make(map[string]*BasmFragment)
	bi.chunks = make(map[string]*BasmChunk)
	bi.passes = uint64(4095)
	bi.matchers = make([]*bmline.BasmLine, 0)
	bi.matchersOps = make([]procbuilder.Opcode, 0)

	bi.rg = bmreqs.NewReqRoot()

	bi.cps = make([]*bmline.BasmElement, 0)
	bi.sos = make([]*bmline.BasmElement, 0)
	bi.ios = make([]*bmline.BasmElement, 0)
	bi.ioAttach = make([]*bmline.BasmElement, 0)
	bi.soAttach = make([]*bmline.BasmElement, 0)
	bi.fis = make([]*bmline.BasmElement, 0)
	bi.fiLinks = make([]*bmline.BasmElement, 0)
	bi.fiLinkAttach = make([]*bmline.BasmElement, 0)
	bi.global = new(bmline.BasmElement)
	bi.global.SetValue("global")

	if bi.debug {
		fmt.Println(purple("Init") + ":")
	}
	for _, op := range procbuilder.Allopcodes {
		if bi.debug {
			fmt.Println(purple("\tExamining opcode: ") + op.Op_get_name())
		}
		for _, line := range op.HLAssemblerMatch(nil) {
			if mt, err := bmline.Text2BasmLine(line); err == nil {
				bi.matchers = append(bi.matchers, mt)
				bi.matchersOps = append(bi.matchersOps, op)
			} else {
				bi.Warning(err)
			}
		}
	}
}

func (bi *BasmInstance) PrintInit() {
	fmt.Println(purple("Init: Reading Matchers"))
	for i, line := range bi.matchers {
		op := bi.matchersOps[i]
		fmt.Println(" ", line, yellow("--> "+op.Op_get_name()+" operand"))
	}
}

// SetVerbose sets the verbose flag on a BasmInstance
func (bi *BasmInstance) SetVerbose() {
	bi.verbose = true
}

// SetDebug sets the debug flag on a BasmInstance
func (bi *BasmInstance) SetDebug() {
	bi.debug = true
}

// Run the assembler
func (bi *BasmInstance) RunAssembler() error {
	step := uint64(1)
	names := getPassFunctionName()

	bi.rg.Requirement(bmreqs.ReqRequest{Node: "/", T: bmreqs.ObjectSet, Name: "code", Value: "romtexts", Op: bmreqs.OpAdd})
	bi.rg.Requirement(bmreqs.ReqRequest{Node: "/", T: bmreqs.ObjectSet, Name: "code", Value: "ramtexts", Op: bmreqs.OpAdd})

	if bi.debug {
		fmt.Println(purple("Pre phase 1 ") + bi.String())
	}

	passes := getPassFunction()
	for i := 0; i < 64; i++ {
		if activePass(bi.passes, step) {
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

	if bi.debug {
		fmt.Println(bi.rg.Requirement(bmreqs.ReqRequest{Node: "/", Op: bmreqs.OpDump}))
	}

	return nil
}

func (m *BasmMacro) String() string {
	result := "\t\t" + red(m.macroName)
	result += m.macroBody.String()
	return result
}

func (m *BasmSection) String() string {
	result := "\t\t" + red(m.sectionName)
	switch m.sectionType {
	case sectRomData:
		result += yellow("[.romdata]")
	case setcRomText:
		result += yellow("[.romtext]")
	}
	result += m.sectionBody.String()
	return result
}

func (m *BasmSection) writeText() string {
	result := ""
	for _, line := range m.sectionBody.Lines {
		if line.GetMeta("label") != "" {
			result += line.GetMeta("label") + ":\n"
		}
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

func (m *BasmFragment) writeText() string {
	result := ""
	for _, line := range m.fragmentBody.Lines {
		if line.GetMeta("label") != "" {
			result += line.GetMeta("label") + ":\n"
		}
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

func (m *BasmFragment) String() string {
	result := "\t\t" + red(m.fragmentName)
	result += m.fragmentBody.String()
	return result
}

func (m *BasmChunk) String() string {
	result := "\t\t" + red(m.chunkName)
	result += m.chunkBody.String()
	return result
}

func (bi *BasmInstance) String() string {
	result := purple("Instance Dump:") + "\n"
	if bi.global != nil {
		result += purple("\tBM meta") + ":\n"
		result += fmt.Sprintf("\t\t") + bi.global.String() + "\n"
	}
	if len(bi.cps) > 0 {
		result += purple("\tCPs meta") + ":\n"
		for i, elem := range bi.cps {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.sos) > 0 {
		result += purple("\tSOs meta") + ":\n"
		for i, elem := range bi.sos {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.ios) > 0 {
		result += purple("\tIOs meta") + ":\n"
		for i, elem := range bi.ios {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.soAttach) > 0 {
		result += purple("\tSO Attach") + ":\n"
		for i, elem := range bi.soAttach {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.ioAttach) > 0 {
		result += purple("\tIO Attach") + ":\n"
		for i, elem := range bi.ioAttach {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.macros) > 0 {
		result += purple("\tMacros") + ":\n"
		for _, macro := range bi.macros {
			result += macro.String() + "\n"
		}
	}
	if len(bi.fragments) > 0 {
		result += purple("\tFragments") + ":\n"
		for _, fragment := range bi.fragments {
			result += fragment.String() + "\n"
		}
	}
	if len(bi.fis) > 0 {
		result += purple("\tFragment instances") + ":\n"
		for i, elem := range bi.fis {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.fiLinks) > 0 {
		result += purple("\tFragment instance links") + ":\n"
		for i, elem := range bi.fiLinks {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.fiLinkAttach) > 0 {
		result += purple("\tFragment instance link attach") + ":\n"
		for i, elem := range bi.fiLinkAttach {
			result += fmt.Sprintf("\t\t%d:", i) + elem.String() + "\n"
		}
	}
	if len(bi.chunks) > 0 {
		result += purple("\tChunks") + ":\n"
		for _, chunk := range bi.chunks {
			result += chunk.String() + "\n"
		}
	}
	if len(bi.sections) > 0 {
		result += purple("\tSections") + ":\n"
		for _, section := range bi.sections {
			result += section.String() + "\n"
		}
	}
	if len(bi.matchers) > 0 {
		result += purple("\tMatchers") + ":\n"
		for _, matcher := range bi.matchers {
			result += "\t\t" + matcher.String() + "\n"
		}
	}
	return result
}
