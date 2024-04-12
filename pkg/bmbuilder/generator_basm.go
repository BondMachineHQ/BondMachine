package bmbuilder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/basm"
	"github.com/BondMachineHQ/BondMachine/pkg/bmconfig"
	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func BasmGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {

	if b.debug {
		fmt.Println(green("\t\t\tBasmGenerator - Start"))
		defer fmt.Println(green("\t\t\tBasmGenerator - End"))
	}

	bi := new(basm.BasmInstance)
	if b.debug {
		bi.SetDebug()
	}
	if b.verbose {
		bi.SetVerbose()
	}

	bi.BMinfo = new(bminfo.BMinfo)

	// Load the BMinfo file if it exists as metadata
	bmInfoFile := l.GetMeta("bminfofile")

	if bmInfoFile != "" {
		bmInfoFile := l.GetMeta("bminfofile")
		if bmInfoJSON, err := os.ReadFile(bmInfoFile); err == nil {
			if err := json.Unmarshal(bmInfoJSON, bi.BMinfo); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	bi.BasmInstanceInit(nil)

	if l.GetMeta("disabledynamicalmatching") != "" {
		bi.Activate(bmconfig.DisableDynamicalMatching)
	}
	if l.GetMeta("chooserminwordsize") != "" {
		bi.Activate(bmconfig.ChooserMinWordSize)
	}
	if l.GetMeta("chooserforcesamename") != "" {
		bi.Activate(bmconfig.ChooserForceSameName)
	}

	// TODO: The the passes control from metadata

	if l.GetMeta("optimizations") != "" {
		optA := strings.Split(l.GetMeta("optimizations"), ":")
		for _, o := range optA {
			if err := bi.ActivateOptimization(o); err != nil {
				return nil, err
			}
		}
	}

	startAssembling := false

	if l.GetMeta("basmfiles") == "" {
		return nil, fmt.Errorf("No BASM files to assemble")
	}

	basmFiles := strings.Split(l.GetMeta("basmfiles"), ":")

	for _, asmFile := range basmFiles {
		startAssembling = true

		// Get the file extension
		ext := filepath.Ext(asmFile)

		switch ext {

		case ".basm":
			err := bi.ParseAssemblyDefault(asmFile)
			if err != nil {
				return nil, errors.New("Error while parsing file" + err.Error())
			}
		case ".ll":
			err := bi.ParseAssemblyLLVM(asmFile)
			if err != nil {
				return nil, errors.New("Error while parsing file" + err.Error())
			}
		default:
			// Default is .basm
			err := bi.ParseAssemblyDefault(asmFile)
			if err != nil {
				return nil, errors.New("Error while parsing file" + err.Error())
			}
		}
	}

	if !startAssembling {
		return nil, errors.New("no BASM files to assemble")
	}

	if err := bi.RunAssembler(); err != nil {
		bi.Alert(err)
		return nil, err
	}

	if err := bi.Assembler2BondMachine(); err != nil {
		return nil, errors.New("Error while converting assembler to BondMachine" + err.Error())
	}

	bMach := bi.GetBondMachine()

	// TODO: bminfofile write ?

	// TODO: dumprequirements ?

	// TODO: bcof ?

	return bMach, nil
}
