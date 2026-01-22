package basm

import (
	"fmt"
)

const (
	OB_CP = uint8(0) + iota
	OB_ROMTEXT
	OB_RAMTEXT
	OB_FRAGMENT
)

var callInstructions = []string{
	"call",
	"call(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)",
	"callo(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)",
	"calla(?P<stacksize>[0-9]+)(?P<stackname>[a-zA-Z_]+)",
}

type dependencies struct {
	depList []depObj
}

type depObj struct {
	obType  uint8
	name    string
	depList []depObj
}

func sameObj(a depObj, b depObj) bool {
	if a.obType == b.obType && a.name == b.name {
		return true
	}
	return false
}

func (d *dependencies) inDepList(obj depObj) bool {
	for _, dep := range d.depList {
		if sameObj(dep, obj) {
			return true
		}
	}
	return false
}

// section entry points detection, the pass detects the symbol used as entry point of the section and sign it as metadata.
func dependencyResolver(bi *BasmInstance) error {
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\tProcessing section: " + sectName))
			}
			obj := depObj{
				name:    sectName,
				depList: make([]depObj, 0),
			}
			if section.sectionType == sectRomText {
				obj.obType = OB_ROMTEXT
			} else {
				obj.obType = OB_RAMTEXT
			}
			if !bi.deps.inDepList(obj) {
				bi.deps.depList = append(bi.deps.depList, obj)
			}
		}
	}
	for fragName, _ := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\tProcessing fragment: " + fragName))
		}
		obj := depObj{
			name:    fragName,
			depList: make([]depObj, 0),
			obType:  OB_FRAGMENT,
		}
		if !bi.deps.inDepList(obj) {
			bi.deps.depList = append(bi.deps.depList, obj)
		}
	}
	return nil
}

// TODO : implement call dependency resolution
