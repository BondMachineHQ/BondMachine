package basm

import (
	"fmt"
	"os"
)

func (bi *BasmInstance) ExportFragmentBasmFile(f *BasmFragment, fileName string) error {
	// Check for file existence, raise error if found
	if _, err := os.Stat(fileName); err == nil {
		return fmt.Errorf("file %s already exists", fileName)
	}

	// Export fragment to file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the header
	if _, err := file.WriteString("%fragment " + f.fragmentName + " "); err != nil {
		return err
	}

	body := f.fragmentBody

	for metaKey, metaValue := range body.LoopMeta() {
		if _, err := file.WriteString(metaKey + ":" + metaValue + " "); err != nil {
			return err
		}
	}

	if _, err := file.WriteString("\n"); err != nil {
		return err
	}

	// Write each line of the fragment
	for _, line := range body.Lines {

		// Line metadata
		lineMeta := line.LoopMeta()
		if symbol, ok := lineMeta["symbol"]; ok {
			if _, err := file.WriteString(symbol + ":"); err != nil {
				return err
			}
			delete(lineMeta, "symbol")
			for metaKey, metaValue := range lineMeta {
				if _, err := file.WriteString(" " + metaKey + ":" + metaValue); err != nil {
					return err
				}
			}
			if _, err := file.WriteString("\n"); err != nil {
				return err
			}
		} else {
			if len(lineMeta) > 0 {
				if _, err := file.WriteString(":"); err != nil {
					return err
				}
				for metaKey, metaValue := range lineMeta {
					if _, err := file.WriteString(" " + metaKey + ":" + metaValue); err != nil {
						return err
					}
				}
				if _, err := file.WriteString("\n"); err != nil {
					return err
				}
			}
		}

		// Line operation
		operation := line.Operation
		if _, err := file.WriteString("\t" + operation.GetValue() + "\t"); err != nil {
			return err
		}

		// Line arguments
		args := line.Elements
		for i, arg := range args {
			if i > 0 {
				if _, err := file.WriteString(", "); err != nil {
					return err
				}
			}
			if _, err := file.WriteString(arg.GetValue()); err != nil {
				return err
			}
		}
		if _, err := file.WriteString("\n"); err != nil {
			return err
		}
	}

	// Write the footer
	if _, err := file.WriteString("%endfragment\n"); err != nil {
		return err
	}

	return nil
}
