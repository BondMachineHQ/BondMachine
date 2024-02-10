package bmbuilder

import (
	"fmt"
	"strings"
)

// metaExtractor is the first pass of the BMBuilder, it extracts the metadata from lines and places it in the BMBuilderBlock
func metaExtractor(b *BMBuilder) error {

	if b.debug {
		fmt.Println(purple("Pass") + ": " + getPassFunctionName()[passMetaExtractor])
	}

	// Loop over the blocks
	for bName, block := range b.blocks {

		if b.debug {
			fmt.Println(yellow("\tAnalyzing block: " + bName))
		}

		// Loop over the lines
		body := block.blockBody
		for i, line := range body.Lines {
			operand := line.Operation.GetValue()

			switch operand {
			case "qbits":
				qbits := make([]string, 0)
				qbitCheck := make(map[string]struct{})
				for _, element := range line.Elements {
					if _, ok := qbitCheck[element.GetValue()]; ok {
						return fmt.Errorf("Qbit %s already defined", element.GetValue())
					}
					qbits = append(qbits, element.GetValue())
					qbitCheck[element.GetValue()] = struct{}{}
				}
				meta := strings.Join(qbits, ":")
				block.blockBody.BasmMeta = block.blockBody.BasmMeta.SetMeta("qbits", meta)

				// Remove the line from the block
				block.blockBody.Lines = append(block.blockBody.Lines[:i], block.blockBody.Lines[i+1:]...)
			}
		}
	}

	return nil
}
