package basm

import (
	"fmt"
)

func clusterChecker(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}
	return nil
}
