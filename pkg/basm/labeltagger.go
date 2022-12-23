package basm

import (
	"errors"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func labelTagger(bi *BasmInstance) error {

	// Filter the matchers to select only the label based one
	filteredMatchers := make([]*bmline.BasmLine, 0)

	if bi.debug {
		fmt.Println(green("\tFiltering matchers:"))
	}

	for _, matcher := range bi.matchers {
		if bmline.FilterMatcher(matcher, "label") {
			filteredMatchers = append(filteredMatchers, matcher)
			if bi.debug {
				fmt.Println(red("\t\tActive matcher:") + matcher.String())
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tInactive matcher:") + matcher.String())
			}
		}
	}

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == setcRomText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			// Map all the labels
			labels := make(map[string]struct{})

			body := section.sectionBody
			for _, line := range body.Lines {
				if label := line.GetMeta("label"); label != "" {
					if _, exists := labels[label]; exists {
						return errors.New("label is specified multiple time: " + label)
					} else {
						labels[label] = struct{}{}
					}
				}
			}

			for _, line := range body.Lines {

				for _, arg := range line.Elements {
					if _, ok := labels[arg.GetValue()]; ok {
						arg.BasmMeta = arg.SetMeta("type", "label")
					}
				}
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	// Loop over the fragments
	for fragName, fragment := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ") + fragName)
		}

		// Map all the labels
		labels := make(map[string]struct{})

		body := fragment.fragmentBody
		for _, line := range body.Lines {
			if label := line.GetMeta("label"); label != "" {
				if _, exists := labels[label]; exists {
					return errors.New("label is specified multiple time: " + label)
				} else {
					labels[label] = struct{}{}
				}
			}
		}

		for _, line := range body.Lines {

			for _, arg := range line.Elements {
				if _, ok := labels[arg.GetValue()]; ok {
					arg.BasmMeta = arg.SetMeta("type", "label")
				}
			}
		}
	}

	return nil
}
