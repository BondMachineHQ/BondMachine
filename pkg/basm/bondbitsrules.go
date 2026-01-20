package basm

import (
	"errors"
	"strings"
)

func (bi *BasmInstance) ApplyBondBitsRule(rule string, saveDirectory string) error {
	// Trim spaces for left and right
	rule = strings.Trim(rule, " ")

	// Split by :
	parts := strings.Split(rule, ":")
	partsNun := len(parts)
	if partsNun == 0 {
		return errors.New("invalid bondbits rule: " + rule)
	}

	// Determine the rule type based on the first part
	switch parts[0] {
	case "merge":
		if partsNun != 4 {
			return errors.New("invalid merge rule, expected 4 parts, got " + string(partsNun))
		}
		source1 := parts[1]
		source2 := parts[2]
		target := parts[3]
		return bi.mergeFragments(source1, source2, target, saveDirectory)
	default:
		return errors.New("unknown bondbits rule type: " + parts[0])
	}
}
