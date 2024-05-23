package basm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (bi *BasmInstance) ParseAssemblyDefault(filePath string) error {
	return bi.ParseAssemblyFile(filePath, basmParser)
}

func (bi *BasmInstance) ParseAssemblyStringDefault(text string) error {
	return bi.ParseAssemblyString(text, basmParser)
}

// ParseAssemblyFile opens the actual assembly file and call a parse function on every line with the underline BasmInstance loaded
func (bi *BasmInstance) ParseAssemblyFile(filePath string, parseFunction func(*BasmInstance, string, uint32) error) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	if bi.debug {
		fmt.Println(purple("Phase 0") + ": " + red("Reading Assembly file "+filePath))
	}

	lineNo := uint32(0)
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if err := parseFunction(bi, scanner.Text(), lineNo); err != nil {
			return err
		}
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// ParseAssemblyString opens the actual assembly file and call a parse function on every line with the underline BasmInstance loaded
func (bi *BasmInstance) ParseAssemblyString(text string, parseFunction func(*BasmInstance, string, uint32) error) error {

	if bi.debug {
		fmt.Println(purple("Synth") + ": " + red("Reading Assembly stream"))
	}

	lineNo := uint32(0)
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		if err := parseFunction(bi, scanner.Text(), lineNo); err != nil {
			return err
		}
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
