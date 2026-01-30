package basm

import (
	"errors"
	"strings"
)

// The bondbits rules are used to manipulate fragments and stuff inside a BasmInstance
// This function is the main entry point for applying bondbits rules.
// Some of the rules may need to run multiple steps before being fully applied.
// The rules are colon-separated strings with a specific format.
// Possible rules:
//   - merge:source1:source2:target
//     Merges two fragments (source1 and source2) into a new fragment (target).
//     The two source fragments has to exist and the target fragment must not exist.
//   - mapsinglecp:source:target
//     Maps a single fragment (source) to a cp (target).
//     The source fragment has to exist and the target cp must not exist.

func (bi *BasmInstance) ApplyBondBitsRule(rule string) error {
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
		return bi.mergeFragments(source1, source2, target)
	case "mapsinglecp":
		if partsNun != 3 {
			return errors.New("invalid mapsinglecp rule, expected 3 parts, got " + string(partsNun))
		}
		source := parts[1]
		target := parts[2]
		return bi.mapSingleFragmentToCP(source, target)
	default:
		return errors.New("unknown bondbits rule type: " + parts[0])
	}
}
