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
	if bi.debug {
		fmt.Println("\t" + green("Building Sequential block"))
	}

	BMNum := len(block.blockBMs)
	if BMNum == 0 {
		return errors.New("no BMs in the block")
	}

	if BMNum == 1 {
		bi.result = block.blockBMs[0]
		return nil
	}

	bm1 := block.blockBMs[0]

	for i := 1; i < BMNum; i++ {
		bm2 := block.blockBMs[i]

		if bm1.Outputs != bm2.Inputs {
			return errors.New("BM outputs and inputs do not match")
		}

		con := new(BMConnections)
		con.Links = make([]BMLink, 0)

		// Set the bm1 inputs as final inputs
		for j := 0; j < bm1.Inputs; j++ {
			e1 := BMEndPoint{BType: FinalBMInput, BNum: j}
			e2 := BMEndPoint{BType: FirstBMInput, BNum: j}
			con.Links = append(con.Links, BMLink{E1: e1, E2: e2})
		}

		// Set the bm1 to bm2 connections
		for j := 0; j < bm1.Outputs; j++ {
			e1 := BMEndPoint{BType: FirstBMOutput, BNum: j}
			e2 := BMEndPoint{BType: SecondBMInput, BNum: j}
			con.Links = append(con.Links, BMLink{E1: e1, E2: e2})
		}

		// Set the bm2 outputs as final outputs
		for j := 0; j < bm2.Outputs; j++ {
			e1 := BMEndPoint{BType: SecondBMOutput, BNum: j}
			e2 := BMEndPoint{BType: FinalBMOutput, BNum: j}
			con.Links = append(con.Links, BMLink{E1: e1, E2: e2})
		}

		// Merge the two bondmachines
		if bm, err := bi.BMMerge(bm1, bm2, con); err == nil {
			bm1 = bm
		} else {
			return err
		}

	}

	bi.result = bm1

	return nil
}
