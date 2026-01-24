package basm

import "os"

func (bi *BasmInstance) ExportBasmFiles(outputDir string) error {
	// Check directory existence
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return err
	}

	// Export fragments one on each file named fragment_<name>.basm
	for fragName, fragment := range bi.fragments {
		filename := outputDir + string(os.PathSeparator) + "fragment_" + fragName + ".basm"
		if err := bi.ExportFragmentBasmFile(fragment, filename); err != nil {
			return err
		}
	}

	// TODO: finish with the others objects
	return nil
}
