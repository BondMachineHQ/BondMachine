package bmbuilder

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func generatorsExec(b *BMBuilder) error {

	if b.debug {
		fmt.Println(purple("Pass") + ": " + getPassFunctionName()[passGeneratorsExec])
	}

	globalMetaData := new(bmline.BasmElement)
	globalMetaData.SetValue("globalmetadata")

	for key, value := range b.global.LoopMeta() {
		globalMetaData.BasmMeta = globalMetaData.SetMeta(key, value)
	}

	// Loop over the blocks
	for bName, block := range b.blocks {

		if b.debug {
			fmt.Println(yellow("\tGenerating block: " + bName))
		}

		b.currentBlock = bName

		// Allocate the BMs array for the block
		block.blockBMs = make([]*bondmachine.Bondmachine, 0)

		metaData := new(bmline.BasmElement)
		metaData.SetValue("metadata")

		for key, value := range globalMetaData.LoopMeta() {
			metaData.BasmMeta = metaData.SetMeta(key, value)
		}

		for key, value := range block.blockBody.LoopMeta() {
			metaData.BasmMeta = metaData.SetMeta(key, value)
		}

		// Loop over the lines
		body := block.blockBody
		for _, line := range body.Lines {
			operand := line.Operation.GetValue()

			if g, ok := b.generators[operand]; !ok {
				return fmt.Errorf("Generator for operation %s not found", operand)
			} else {
				if b.debug {
					fmt.Println(yellow("\t\tGenerating bm from line: " + fmt.Sprint(line) + " with" + fmt.Sprint(metaData)))
				}
				if bm, err := g(b, metaData, line); err != nil {
					return err
				} else {
					block.blockBMs = append(block.blockBMs, bm)
				}
			}
		}
	}

	return nil
}
