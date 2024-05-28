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

			if g, ok := b.generators[operand]; ok {
				if b.debug {
					fmt.Println(yellow("\t\tGenerating bm from line: " + fmt.Sprint(line) + " with" + fmt.Sprint(metaData)))
				}
				if bm, err := g(b, metaData, line); err != nil {
					return err
				} else {
					block.blockBMs = append(block.blockBMs, bm)
					if len(block.blockConn) < len(block.blockBMs) {
						block.blockConn = append(block.blockConn, nil)
					}
					continue
				}
			}

			if c, ok := b.connectors[operand]; ok {
				if len(block.blockConn) > len(block.blockBMs) {
					return fmt.Errorf("Connector cannot be called one after another, it should be called after a BM")
				}
				if b.debug {
					fmt.Println(yellow("\t\tConnecting bm from line: " + fmt.Sprint(line) + " with" + fmt.Sprint(metaData)))
				}
				if conn, err := c(b, metaData, line); err != nil {
					return err
				} else {
					block.blockConn = append(block.blockConn, conn)
					continue
				}
			}
			return fmt.Errorf("Unknown operation: %s", operand)
		}

		if len(block.blockConn) == len(block.blockBMs) {
			block.blockConn = append(block.blockConn, nil)
		}

		if b.debug {
			fmt.Println(yellow("\tBlock: " + bName + " has been generated"))
			fmt.Println(yellow("\t\tBMs: " + fmt.Sprint(block.blockBMs)))
			fmt.Println(yellow("\t\tConns: " + fmt.Sprint(block.blockConn)))

		}

	}

	return nil
}
