package bmbuilder

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (bld *BMBuilder) ParseBuilderDefault(filePath string) error {
	return bld.ParseBuilderFile(filePath, basmParser)
}

// ParseBuilderFile opens the actual assembly file and call a parse function on every line with the underline BMBuilder loaded
func (bld *BMBuilder) ParseBuilderFile(filePath string, parseFunction func(*BMBuilder, string, uint32) error) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	if bld.debug {
		fmt.Println(purple("Phase 0") + ": " + red("Reading Builder file "+filePath))
	}

	lineNo := uint32(0)
	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		if err := parseFunction(bld, scanner.Text(), lineNo); err != nil {
			return err
		}
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// ParseBuilderString opens the actual assembly file and call a parse function on every line with the underline BMBuilder loaded
func (bld *BMBuilder) ParseBuilderString(text string, parseFunction func(*BMBuilder, string, uint32) error) error {

	if bld.debug {
		fmt.Println(purple("Synth") + ": " + red("Reading Builder stream"))
	}

	lineNo := uint32(0)
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		if err := parseFunction(bld, scanner.Text(), lineNo); err != nil {
			return err
		}
		lineNo++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
