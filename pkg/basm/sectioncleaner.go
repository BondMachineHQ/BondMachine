package basm

import (
	"fmt"
	"strings"
)

func sectionCleaner(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tPruning unused sections:"))
	}

	usefulSections := make(map[string]struct{})

	for _, cp := range bi.cps {
		romCodeName := cp.GetMeta("romcode")
		if romCodeName != "" {
			usefulSections[romCodeName] = struct{}{}
		}
		romDataName := cp.GetMeta("romdata")
		if romDataName != "" {
			usefulSections[romDataName] = struct{}{}
		}
		ramCodeName := cp.GetMeta("ramcode")
		if ramCodeName != "" {
			usefulSections[ramCodeName] = struct{}{}
		}
		ramDataName := cp.GetMeta("ramdata")
		if ramDataName != "" {
			usefulSections[ramDataName] = struct{}{}
		}
		altRomCodeName := cp.GetMeta("romalternatives")
		if altRomCodeName != "" {
			splitAlts := strings.Split(altRomCodeName, ":")
			for _, alt := range splitAlts {
				usefulSections[alt] = struct{}{}
			}
		}
		altRamCodeName := cp.GetMeta("ramalternatives")
		if altRamCodeName != "" {
			splitAlts := strings.Split(altRamCodeName, ":")
			for _, alt := range splitAlts {
				usefulSections[alt] = struct{}{}
			}
		}
		altRomDataName := cp.GetMeta("romdataalternatives")
		if altRomDataName != "" {
			splitAlts := strings.Split(altRomDataName, ":")
			for _, alt := range splitAlts {
				usefulSections[alt] = struct{}{}
			}
		}
		altRamDataName := cp.GetMeta("ramdataalternatives")
		if altRamDataName != "" {
			splitAlts := strings.Split(altRamDataName, ":")
			for _, alt := range splitAlts {
				usefulSections[alt] = struct{}{}
			}
		}
	}

	// Remove the unused sections
	for sectionName, _ := range bi.sections {
		if _, ok := usefulSections[sectionName]; !ok {
			if bi.debug {
				fmt.Println(yellow("\t\tRemoving section: ") + sectionName)
			}
			delete(bi.sections, sectionName)
		}
	}
	return nil
}
