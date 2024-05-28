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
		var connDelay *bmline.BasmLine
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
					if len(block.blockConn) < len(block.blockBMs)-1 {
						if connDelay == nil {
							// Fill with the default connector if it is not defined
							block.blockConn = append(block.blockConn, nil)
						} else {
							connOperand := connDelay.Operation.GetValue()
							if c, ok := b.connectors[connOperand]; ok {
								if b.debug {
									fmt.Println(yellow("\t\tGenerating connector from line: " + fmt.Sprint(connDelay) + " with" + fmt.Sprint(metaData)))
								}
								if conn, err := c(b, metaData, connDelay); err != nil {
									return err
								} else {
									block.blockConn = append(block.blockConn, conn)
								}
							} else {
								return fmt.Errorf("unknown connector: %s", connOperand)
							}
							connDelay = nil
						}
					}
					continue
				}
			}

			if _, ok := b.connectors[operand]; ok {
				if connDelay != nil {
					return fmt.Errorf("Connector cannot be called one after another, it should be called after a BM")
				}
				connDelay = line
				continue
			}

			return fmt.Errorf("unknown operation: %s", operand)
		}

		if connDelay != nil {
			return fmt.Errorf("Connector cannot be called as last line in the block")
		}

		if b.debug {
			fmt.Println(yellow("\tBlock: " + bName + " has been generated"))
			fmt.Println(yellow("\t\tBMs: " + fmt.Sprint(block.blockBMs)))
			fmt.Println(yellow("\t\tConns: " + fmt.Sprint(block.blockConn)))

		}

	}

	return nil
}
