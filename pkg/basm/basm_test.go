package basm

import (
	"fmt"
	"testing"
)

func TestLineReader(t *testing.T) {
	bi := new(BasmInstance)
	bi.BasmInstanceInit(nil)
	bi.SetDebug()
	if bi.debug || bi.verbose {
		bi.PrintInit()
	}

	err := bi.ParseAssemblyDefault("test1.basm")
	if err != nil {
		bi.Alert("Error while parsing file", err)
		return
	}

	err = bi.ParseAssemblyDefault("test2.basm")
	if err != nil {
		bi.Alert("Error while parsing file", err)
		return
	}
	fmt.Print(bi.RunAssembler())

}
