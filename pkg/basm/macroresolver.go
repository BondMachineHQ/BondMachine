package basm

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func macroResolver(bi *BasmInstance) error {

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}
			body := section.sectionBody
			if err := bi.bodyMacros(body); err != nil {
				return err
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	// Loop over the fragments
	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ") + fragName)
		}
		body := fragment.fragmentBody
		if err := bi.bodyMacros(body); err != nil {
			return err
		}
	}

	return nil

}

func (bi *BasmInstance) bodyMacros(body *bmline.BasmBody) error {
	for i := 0; i < len(body.Lines); i++ {
		line := body.Lines[i]
		op := line.Operation.GetValue()
		if macro, ok := bi.macros[op]; ok {
			if bi.debug {
				fmt.Println(yellow("\t\t\tMacro: ") + op)
			}
			expArgs := macro.macroArgs
			args := len(line.Elements)

			if args != expArgs {
				return fmt.Errorf("macro %s expects %d arguments, got %d", op, expArgs, args)
			}

			macroLines := bi.expandMacro(macro, line)
			// Insert lines starting at i
			if bi.debug {
				fmt.Println(yellow("\t\t\t\tExpanding macro: ") + op)
				for _, macroLine := range macroLines {
					fmt.Println(yellow("\t\t\t\t\t") + macroLine.String())
				}
			}
			body.Lines = append(body.Lines[:i], append(macroLines, body.Lines[i+1:]...)...)
			i += len(macroLines)
		}
	}
	return nil
}

func (bi *BasmInstance) expandMacro(macro *BasmMacro, line *bmline.BasmLine) []*bmline.BasmLine {
	return macro.macroBody.Lines
}
