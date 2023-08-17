package basm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func getPairs(couples string) []string {
	// TODO Horrific and temporary code, a proper parser is desireable
	space := regexp.MustCompile(`\s+`)

	stripped := strings.TrimSpace(couples)
	stripped = space.ReplaceAllString(stripped, "")

	result := strings.Split(stripped, ",")

	return result
}

func (bi *BasmInstance) literalProcessor(meta string) error {

	stripped := strings.TrimSpace(strings.TrimSpace(meta)[7:]) // Remove the %meta: part and the literal

	re := regexp.MustCompile(`^(?P<meta>\w+)\s+(?P<value>.+)$`)
	if re.MatchString(stripped) {
		name := re.ReplaceAllString(stripped, "${meta}")
		value := re.ReplaceAllString(stripped, "${value}")

		if bi.isWithinFragment != "" {
			if bi.debug {
				fmt.Println("Adding %meta literal: ||" + name + "|| = ||" + value + "|| to " + bi.isWithinFragment)
			}
			bi.fragments[bi.isWithinFragment].fragmentBody.BasmMeta = bi.fragments[bi.isWithinFragment].fragmentBody.SetMeta(name, value)
		}

	} else {
		return errors.New("Invalid %meta: " + meta)
	}

	return nil
}

func (bi *BasmInstance) metaProcessor(meta string) error {

	// TODO Horrific and temporary code, a proper parser is desireable
	space := regexp.MustCompile(`\s+`)

	stripped := strings.TrimSpace(meta)
	stripped = space.ReplaceAllString(stripped, " ")

	splitted := strings.Split(stripped, " ")

	if splitted[0] == "literal" {
		return bi.literalProcessor(meta)
	}

	if len(splitted) < 3 {
		return errors.New("Invalid %meta entry")
	}

	metaCommand := splitted[0]
	metaObjId := splitted[1]
	metaValues := strings.Join(splitted[2:], "")

	metaValues = space.ReplaceAllString(metaValues, "")

	pairs := strings.Split(metaValues, ",")

	if bi.debug {
		fmt.Printf("Meta command: %s, object: %s, values: %s\n", metaCommand, metaObjId, metaValues)
	}

	if len(pairs) == 0 {
		return errors.New("Wrong format for %meta entry: " + metaCommand + " " + metaValues)
	}

	resultDict := make(map[string]string)

	// Subsequent pairs are key=value and : is the separator between them. Multiple : are for multivalued keys or value containing :
	for _, pair := range pairs {
		keyVal := strings.Split(pair, ":")
		if len(keyVal) < 2 {
			return errors.New("Wrong format for %meta entry: " + metaCommand + " " + metaValues)
		} else {
			resultDict[keyVal[0]] = strings.Join(keyVal[1:], ":")
		}
	}
	exCheck := false
	switch metaCommand {
	case "bmdef":
		switch metaObjId {
		case "global":
			for key, value := range resultDict {
				if err := bi.filteredMetaAdd(bi.global, key, value, "global"); err != nil {
					return err
				}
			}
		}
	case "cpdef":
		for _, cp := range bi.cps {
			if cp.GetValue() == metaObjId {
				exCheck = true
				for key, value := range resultDict {
					if err := bi.filteredMetaAdd(cp, key, value, metaCommand); err != nil {
						return err
					}
				}
				break
			}
		}
	case "sodef":
		for _, so := range bi.sos {
			if so.GetValue() == metaObjId {
				exCheck = true
				for key, value := range resultDict {
					if err := bi.filteredMetaAdd(so, key, value, metaCommand); err != nil {
						return err
					}
				}
				break
			}
		}
	case "iodef":
		for _, io := range bi.ios {
			if io.GetValue() == metaObjId {
				exCheck = true
				for key, value := range resultDict {
					if err := bi.filteredMetaAdd(io, key, value, metaCommand); err != nil {
						return err
					}
				}
				break
			}
		}
	case "fidef":
		for _, fi := range bi.fis {
			if fi.GetValue() == metaObjId {
				exCheck = true
				for key, value := range resultDict {
					if err := bi.filteredMetaAdd(fi, key, value, metaCommand); err != nil {
						return err
					}
				}
				break
			}
		}
	case "filinkdef":
		for _, fil := range bi.fiLinks {
			if fil.GetValue() == metaObjId {
				exCheck = true
				for key, value := range resultDict {
					if err := bi.filteredMetaAdd(fil, key, value, metaCommand); err != nil {
						return err
					}
				}
				break
			}
		}
	case "soatt":
	case "ioatt":
	case "filinkatt":
	default:
		return errors.New("Unknown %meta entry: " + metaCommand)
	}

	if !exCheck {
		newElem := new(bmline.BasmElement)
		newElem.SetValue(metaObjId)
		for key, value := range resultDict {
			if err := bi.filteredMetaAdd(newElem, key, value, metaCommand); err != nil {
				return err
			}
		}
		switch metaCommand {
		case "cpdef":
			bi.cps = append(bi.cps, newElem)
		case "sodef":
			bi.sos = append(bi.sos, newElem)
		case "iodef":
			bi.ios = append(bi.ios, newElem)
		case "fidef":
			bi.fis = append(bi.fis, newElem)
		case "filinkdef":
			bi.fiLinks = append(bi.fiLinks, newElem)
		case "soatt":
			// TODO Include checks for consistent attach
			bi.soAttach = append(bi.soAttach, newElem)
		case "ioatt":
			// TODO Include checks for consistent attach
			bi.ioAttach = append(bi.ioAttach, newElem)
		case "filinkatt":
			// TODO Include checks for consistent attach
			bi.fiLinkAttach = append(bi.fiLinkAttach, newElem)
		}
	}

	return nil
}
func (bi *BasmInstance) filteredMetaAdd(el *bmline.BasmElement, key string, value string, metaType string) error {
	switch metaType {
	case "global":
		switch key {
		case "registersize":
		case "iomode":
		case "defaultexecmode":
		default:
			return errors.New("Unknown global %meta: " + key)
		}
	case "cpdef":
		switch key {
		case "romcode":
		case "romdata":
		case "ramcode":
		case "ramdata":
		case "romsize":
		case "ramsize":
		case "execmode":
		case "fragcollapse":
		default:
			// If there in an unknown key, it is a user defined key that eventually will be used in a template.
			// Setting the CP as templated and the resolver will create a new instance of the code
			el.BasmMeta = el.SetMeta("templated", "true")
			if bi.debug {
				fmt.Println("Setting templated to true for " + el.GetValue())
			}
		}
	case "sodef":
		switch key {
		case "type":
		case "constraint":
		default:
			return errors.New("Unknown sodef %meta: " + key)
		}
	case "iodef":
		switch key {
		case "type":
		default:
			return errors.New("Unknown iodef %meta: " + key)
		}
	case "fidef":
		switch key {
		case "type":
		case "fragment":
		default:
			// If there in an unknown key, it is a user defined key that eventually will be used in a template.
			// Setting the CP as templated and the resolver will create a new instance of the code
			el.BasmMeta = el.SetMeta("templated", "true")
			if bi.debug {
				fmt.Println("Setting templated to true for " + el.GetValue())
			}
		}
	case "filinkdef":
		switch key {
		case "type":
		default:
			return errors.New("Unknown filinkdef %meta: " + key)
		}
	case "ioatt":
		switch key {
		case "cp":
		case "type":
		case "index":
		default:
			return errors.New("Unknown ioatt %meta: " + key)
		}
	case "soatt":
		switch key {
		case "cp":
		case "index":
		default:
			return errors.New("Unknown soatt %meta: " + key)
		}
	case "filinkatt":
		switch key {
		case "fi":
		case "index":
		case "type":
		default:
			return errors.New("Unknown filinkatt %meta: " + key)
		}
	}
	el.BasmMeta = el.SetMeta(key, value)
	return nil
}
