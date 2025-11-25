package bondmachine

import "fmt"

func (sr *SimReport) String() string {
	if sr == nil {
		return "Empty SimReport"
	}

	if len(sr.Reportables) == 0 && len(sr.Showables) == 0 && len(sr.AbsGet) == 0 && len(sr.PerGet) == 0 {
		return "No reportable data"
	}

	result := "SimReport\n"
	result += "=========\n\n"

	// List reportable elements (GET)
	if len(sr.Reportables) > 0 {
		result += fmt.Sprintf("Reportable elements (GET): %d\n", len(sr.Reportables))
		for i, elem := range sr.Reportables {
			name := ""
			if i < len(sr.ReportablesNames) {
				name = sr.ReportablesNames[i]
			}
			elemType := "unsigned"
			if i < len(sr.ReportablesTypes) {
				elemType = sr.ReportablesTypes[i]
			}
			result += fmt.Sprintf("  [%d] %s (%s) @ %p\n", i, name, elemType, elem)
		}
		result += "\n"
	}

	// List showable elements (SHOW)
	if len(sr.Showables) > 0 {
		result += fmt.Sprintf("Showable elements (SHOW): %d\n", len(sr.Showables))
		for i, elem := range sr.Showables {
			name := ""
			if i < len(sr.ShowablesNames) {
				name = sr.ShowablesNames[i]
			}
			elemType := "unsigned"
			if i < len(sr.ShowablesTypes) {
				elemType = sr.ShowablesTypes[i]
			}
			result += fmt.Sprintf("  [%d] %s (%s) @ %p\n", i, name, elemType, elem)
		}
		result += "\n"
	}

	// Report absolute GET values
	if len(sr.AbsGet) > 0 {
		// Sort ticks
		absTicks := make([]uint64, 0, len(sr.AbsGet))
		for tick := range sr.AbsGet {
			absTicks = append(absTicks, tick)
		}
		for i := 0; i < len(absTicks); i++ {
			for j := i + 1; j < len(absTicks); j++ {
				if absTicks[i] > absTicks[j] {
					absTicks[i], absTicks[j] = absTicks[j], absTicks[i]
				}
			}
		}

		result += fmt.Sprintf("Absolute GET captures at %d tick(s):\n", len(sr.AbsGet))
		for _, tick := range absTicks {
			tickData := sr.AbsGet[tick]
			result += fmt.Sprintf("\nTick %d:\n", tick)

			// Get sorted indices
			indices := make([]int, 0, len(tickData))
			for idx := range tickData {
				indices = append(indices, idx)
			}
			for i := 0; i < len(indices); i++ {
				for j := i + 1; j < len(indices); j++ {
					if indices[i] > indices[j] {
						indices[i], indices[j] = indices[j], indices[i]
					}
				}
			}

			for _, idx := range indices {
				value := tickData[idx]
				name := ""
				if idx < len(sr.ReportablesNames) {
					name = sr.ReportablesNames[idx]
				}
				elemType := "unsigned"
				if idx < len(sr.ReportablesTypes) {
					elemType = sr.ReportablesTypes[idx]
				}

				result += fmt.Sprintf("  %s [%d]: ", name, idx)
				result += sr.formatValue(value, elemType) + "\n"
			}
		}
		result += "\n"
	}

	// Report periodic GET values
	if len(sr.PerGet) > 0 {
		// Sort periods
		perPeriods := make([]uint64, 0, len(sr.PerGet))
		for period := range sr.PerGet {
			perPeriods = append(perPeriods, period)
		}
		for i := 0; i < len(perPeriods); i++ {
			for j := i + 1; j < len(perPeriods); j++ {
				if perPeriods[i] > perPeriods[j] {
					perPeriods[i], perPeriods[j] = perPeriods[j], perPeriods[i]
				}
			}
		}

		result += fmt.Sprintf("Periodic GET captures (%d period(s)):\n", len(sr.PerGet))
		for _, period := range perPeriods {
			periodData := sr.PerGet[period]
			result += fmt.Sprintf("\nPeriod %d:\n", period)

			// Get sorted indices
			indices := make([]int, 0, len(periodData))
			for idx := range periodData {
				indices = append(indices, idx)
			}
			for i := 0; i < len(indices); i++ {
				for j := i + 1; j < len(indices); j++ {
					if indices[i] > indices[j] {
						indices[i], indices[j] = indices[j], indices[i]
					}
				}
			}

			for _, idx := range indices {
				value := periodData[idx]
				name := ""
				if idx < len(sr.ReportablesNames) {
					name = sr.ReportablesNames[idx]
				}
				elemType := "unsigned"
				if idx < len(sr.ReportablesTypes) {
					elemType = sr.ReportablesTypes[idx]
				}

				result += fmt.Sprintf("  %s [%d]: ", name, idx)
				result += sr.formatValue(value, elemType) + "\n"
			}
		}
		result += "\n"
	}

	// Report absolute SHOW values
	if len(sr.AbsShow) > 0 {
		// Sort ticks
		absTicks := make([]uint64, 0, len(sr.AbsShow))
		for tick := range sr.AbsShow {
			absTicks = append(absTicks, tick)
		}
		for i := 0; i < len(absTicks); i++ {
			for j := i + 1; j < len(absTicks); j++ {
				if absTicks[i] > absTicks[j] {
					absTicks[i], absTicks[j] = absTicks[j], absTicks[i]
				}
			}
		}

		result += fmt.Sprintf("Absolute SHOW triggers at %d tick(s):\n", len(sr.AbsShow))
		for _, tick := range absTicks {
			tickData := sr.AbsShow[tick]
			result += fmt.Sprintf("\nTick %d: ", tick)

			count := 0
			for idx := range tickData {
				if count > 0 {
					result += ", "
				}
				name := ""
				if idx < len(sr.ShowablesNames) {
					name = sr.ShowablesNames[idx]
				}
				if name != "" {
					result += fmt.Sprintf("%s[%d]", name, idx)
				} else {
					result += fmt.Sprintf("[%d]", idx)
				}
				count++
			}
			result += "\n"
		}
		result += "\n"
	}

	// Report periodic SHOW values
	if len(sr.PerShow) > 0 {
		// Sort periods
		perPeriods := make([]uint64, 0, len(sr.PerShow))
		for period := range sr.PerShow {
			perPeriods = append(perPeriods, period)
		}
		for i := 0; i < len(perPeriods); i++ {
			for j := i + 1; j < len(perPeriods); j++ {
				if perPeriods[i] > perPeriods[j] {
					perPeriods[i], perPeriods[j] = perPeriods[j], perPeriods[i]
				}
			}
		}

		result += fmt.Sprintf("Periodic SHOW triggers (%d period(s)):\n", len(sr.PerShow))
		for _, period := range perPeriods {
			periodData := sr.PerShow[period]
			result += fmt.Sprintf("\nPeriod %d: ", period)

			count := 0
			for idx := range periodData {
				if count > 0 {
					result += ", "
				}
				name := ""
				if idx < len(sr.ShowablesNames) {
					name = sr.ShowablesNames[idx]
				}
				if name != "" {
					result += fmt.Sprintf("%s[%d]", name, idx)
				} else {
					result += fmt.Sprintf("[%d]", idx)
				}
				count++
			}
			result += "\n"
		}
		result += "\n"
	}

	return result
}

func (sr *SimReport) formatValue(value interface{}, elemType string) string {
	switch v := value.(type) {
	case uint8:
		switch elemType {
		case "hex":
			return fmt.Sprintf("0x%02X", v)
		case "binary":
			return fmt.Sprintf("0b%08b", v)
		case "signed":
			return fmt.Sprintf("%d", int8(v))
		default: // unsigned
			return fmt.Sprintf("%d", v)
		}
	case uint16:
		switch elemType {
		case "hex":
			return fmt.Sprintf("0x%04X", v)
		case "binary":
			return fmt.Sprintf("0b%016b", v)
		case "signed":
			return fmt.Sprintf("%d", int16(v))
		default: // unsigned
			return fmt.Sprintf("%d", v)
		}
	case uint32:
		switch elemType {
		case "hex":
			return fmt.Sprintf("0x%08X", v)
		case "binary":
			return fmt.Sprintf("0b%032b", v)
		case "signed":
			return fmt.Sprintf("%d", int32(v))
		default: // unsigned
			return fmt.Sprintf("%d", v)
		}
	case uint64:
		switch elemType {
		case "hex":
			return fmt.Sprintf("0x%016X", v)
		case "binary":
			return fmt.Sprintf("0b%064b", v)
		case "signed":
			return fmt.Sprintf("%d", int64(v))
		default: // unsigned
			return fmt.Sprintf("%d", v)
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}

func EventListShow(srep *SimReport, srepOld *SimReport, vm *VM, oldVm *VM) (SimTickShow, error) {
	result := make(SimTickShow)

	for event, pointers := range srep.EventShow {
		switch event.event {
		case EVENTONVALID:
			// This is a case where we need to unwrap boolean pointers
			iposv := pointers[1]
			NewValisRef := (*srep.EventData[iposv]).(*bool)
			OldValisRef := (*srepOld.EventData[iposv]).(*bool)
			newValid := *NewValisRef
			oldValid := *OldValisRef
			if newValid != oldValid && newValid {
				ipos := pointers[0]
				result[ipos] = struct{}{}
			}
		case EVENTONCHANGE:
		case EVENTONRECV:
		}
	}

	return result, nil
}
