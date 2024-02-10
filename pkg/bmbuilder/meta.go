package bmbuilder

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

func (bi *BMBuilder) literalProcessor(meta string) error {

	stripped := strings.TrimSpace(strings.TrimSpace(meta)[7:]) // Remove the %meta: part and the literal

	re := regexp.MustCompile(`^(?P<meta>\w+)\s+(?P<value>.+)$`)
	if re.MatchString(stripped) {
		// name := re.ReplaceAllString(stripped, "${meta}")
		// value := re.ReplaceAllString(stripped, "${value}")

		// if bi.isWithinFragment != "" {
		// 	if bi.debug {
		// 		fmt.Println("Adding %meta literal: ||" + name + "|| = ||" + value + "|| to " + bi.isWithinFragment)
		// 	}
		// 	bi.fragments[bi.isWithinFragment].fragmentBody.BasmMeta = bi.fragments[bi.isWithinFragment].fragmentBody.SetMeta(name, value)
		// }

	} else {
		return errors.New("Invalid %meta: " + meta)
	}

	return nil
}

func (bi *BMBuilder) metaProcessor(meta string) error {

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
	// exCheck := false
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
	default:
		return errors.New("Unknown %meta entry: " + metaCommand)
	}

	// if !exCheck {
	// 	newElem := new(bmline.BasmElement)
	// 	newElem.SetValue(metaObjId)
	// 	for key, value := range resultDict {
	// 		if err := bi.filteredMetaAdd(newElem, key, value, metaCommand); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	// switch metaCommand {
	// 	// }
	// }

	return nil
}
func (bi *BMBuilder) filteredMetaAdd(el *bmline.BasmElement, key string, value string, metaType string) error {
	switch metaType {
	case "global":
		switch key {
		case "main":
		default:
			return errors.New("Unknown global %meta: " + key)
		}

		el.BasmMeta = el.SetMeta(key, value)
		return nil
	}
	return errors.New("Unknown %meta: " + metaType)
}
