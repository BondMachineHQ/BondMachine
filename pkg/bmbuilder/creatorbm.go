package bmbuilder

import (
	"errors"
	"fmt"
	"strconv"
)

func (bi *BMBuilder) BuildBondMachine() error {
	if bi.debug {
		fmt.Println("\t" + green("BondMachine metadata"))
	}

	registerSize := bi.global.GetMeta("registersize")

	if registerSize == "" {
		return errors.New("register size not specified")
	}

	var rSize uint8
	if size, err := strconv.Atoi(registerSize); err == nil {
		if 0 < size && size < 256 {
			rSize = uint8(size)
		} else {
			return errors.New("wrong value for register size")
		}
	} else {
		return errors.New("register size not valid")
	}
	if bi.debug {
		fmt.Println("\t\t"+green("register size:"), rSize)
	}

	// Get the main block
	mainBlock := bi.global.GetMeta("main")

	if mainBlock == "" {
		return errors.New("main block not specified")
	}

	if block, ok := bi.blocks[mainBlock]; ok {
		if bi.debug {
			fmt.Println("\t\t"+green("main block:"), mainBlock)
		}

		switch block.blockType {
		case blockSequential:
			return bi.BuildSequentialBlock(block)
		default:
			return errors.New("unknown block type for main block")
		}

	} else {
		return errors.New("main block not found")
	}

	return nil
}

func (bi *BMBuilder) BuildSequentialBlock(block *BMBuilderBlock) error {
	return nil
}
