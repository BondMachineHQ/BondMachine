package basm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
)

func memComposer(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tConnecting CP code and ROM:"))
	}

	// Loop over the CP code and data and connect them to the ROM
	for _, cp := range bi.cps {

		romAlternatives := make([]string, 0)
		ramAlternatives := make([]string, 0)
		romDataAlternatives := make([]string, 0)
		ramDataAlternatives := make([]string, 0)

		message := "\t\tConnecting "
		cpNewSectionName := ""
		dataNames := ""
		romCodeName := cp.GetMeta("romcode")
		if romCodeName != "" {
			message += "romcode " + romCodeName + ", "
			if _, ok := bi.sections[romCodeName]; !ok {
				return errors.New("ROM code section " + romCodeName + " not found")
			}
		}
		romDataName := cp.GetMeta("romdata")
		if romDataName != "" {
			message += "romdata " + romDataName + ", "
			dataNames += "_od_" + romDataName
			if _, ok := bi.sections[romDataName]; !ok {
				return errors.New("ROM data section " + romDataName + " not found")
			}
		}
		ramCodeName := cp.GetMeta("ramcode")
		if ramCodeName != "" {
			message += "ramcode " + ramCodeName + ", "
			if _, ok := bi.sections[ramCodeName]; !ok {
				return errors.New("RAM code section " + ramCodeName + " not found")
			}
		}
		ramDataName := cp.GetMeta("ramdata")
		if ramDataName != "" {
			message += "ramdata " + ramDataName + ", "
			dataNames += "_ad_" + ramDataName
			if _, ok := bi.sections[ramDataName]; !ok {
				return errors.New("RAM data section " + ramDataName + " not found")
			}
		}

		romSections := make([]string, 0)

		if romCodeName != "" {
			alts := bi.sections[romCodeName].sectionBody.GetMeta("alternatives")
			if alts != "" {
				splitAlts := strings.Split(alts, ":")
				romSections = append(romSections, splitAlts...)
			} else {
				romSections = append(romSections, romCodeName)
			}
		}

		ramSections := make([]string, 0)

		if ramCodeName != "" {
			alts := bi.sections[ramCodeName].sectionBody.GetMeta("alternatives")
			if alts != "" {
				splitAlts := strings.Split(alts, ":")
				ramSections = append(ramSections, splitAlts...)
			} else {
				ramSections = append(ramSections, ramCodeName)
			}
		}

		var romData *BasmSection
		var ramData *BasmSection

		if romDataName != "" {
			romData = bi.sections[romDataName]
		}
		if ramDataName != "" {
			ramData = bi.sections[ramDataName]
		}

		if len(romSections) == 0 {
			romSections = append(romSections, "")
		}
		if len(ramSections) == 0 {
			ramSections = append(ramSections, "")
		}

		for _, romSection := range romSections {
			for _, ramSection := range ramSections {
				cpNewSectionName = ""
				if bi.debug {
					fmt.Println(green(message))
				}

				var romCode *BasmSection
				var ramCode *BasmSection

				if romSection != "" {
					if _, ok := bi.sections[romSection]; !ok {
						return errors.New("ROM section " + romSection + " not found")
					} else {
						cpNewSectionName += "_ot_" + romSection
						romCode = bi.sections[romSection]
					}
				}

				if ramSection != "" {
					if _, ok := bi.sections[ramSection]; !ok {
						return errors.New("RAM section " + ramSection + " not found")
					} else {
						cpNewSectionName += "_at_" + ramSection
						ramCode = bi.sections[ramSection]
					}
				}

				cpNewSectionName += dataNames

				var romSectionLength int = 0
				var ramSectionLength int = 0
				var romDataLength int = 0
				var ramDataLength int = 0
				var romFinalLength int = 0
				var ramFinalLength int = 0

				if romCode != nil {
					romSectionLength = len(romCode.sectionBody.Lines)
					romFinalLength = romSectionLength
				}
				if ramCode != nil {
					ramSectionLength = len(ramCode.sectionBody.Lines)
					ramFinalLength = ramSectionLength
				}

				if bi.debug {
					fmt.Println(green("\t\t\tCode rom section length: " + fmt.Sprintf("%d", romSectionLength)))
					fmt.Println(green("\t\t\tCode ram section length: " + fmt.Sprintf("%d", ramSectionLength)))
				}

				// Resolving Symbols for the composed section
				if romData != nil {
					newSection := new(BasmSection)
					newSection.sectionName = "romdata" + cpNewSectionName
					newSection.sectionType = sectRomData
					newSection.sectionBody = romData.sectionBody.Copy()
					bi.resolveSymbols(newSection, cpNewSectionName)
					bi.sections[newSection.sectionName] = newSection
					romDataAlternatives = append(romDataAlternatives, newSection.sectionName)

					for _, symbol := range newSection.sectionBody.Lines {
						symbolName := symbol.Operation.GetValue()
						offset := symbol.Operation.GetMeta("offset")
						location, _ := strconv.Atoi(offset)
						location += romSectionLength
						bi.symbols["romdata.romdata"+cpNewSectionName+"."+symbolName] = int64(location)
					}
					resp := bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romdatas/sections:" + romDataName, Name: "datalength", Op: bmreqs.OpGet})
					if resp.Error != nil {
						return resp.Error
					}
					romDataLength, _ = strconv.Atoi(resp.Value)
					romFinalLength += romDataLength
				}

				if ramData != nil {
					newSection := new(BasmSection)
					newSection.sectionName = "ramdata" + cpNewSectionName
					newSection.sectionType = sectRamData
					newSection.sectionBody = ramData.sectionBody.Copy()
					bi.resolveSymbols(newSection, cpNewSectionName)
					bi.sections[newSection.sectionName] = newSection
					ramDataAlternatives = append(ramDataAlternatives, newSection.sectionName)

					for _, symbol := range ramData.sectionBody.Lines {
						symbolName := symbol.Operation.GetValue()
						offset := symbol.Operation.GetMeta("offset")
						location, _ := strconv.Atoi(offset)
						location += ramSectionLength
						bi.symbols["ramdata.ramdata"+cpNewSectionName+"."+symbolName] = int64(location)
					}
					resp := bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramdatas/sections:" + ramDataName, Name: "datalength", Op: bmreqs.OpGet})
					if resp.Error != nil {
						return resp.Error
					}
					ramDataLength, _ = strconv.Atoi(resp.Value)
					ramFinalLength += ramDataLength
				}

				// If the original section is not empty, we need to create a new section from it
				if romCode != nil {
					newSection := new(BasmSection)
					newSection.sectionName = "romcode" + cpNewSectionName
					newSection.sectionType = sectRomText
					newSection.sectionBody = romCode.sectionBody.Copy()
					bi.resolveSymbols(newSection, cpNewSectionName)
					bi.sections[newSection.sectionName] = newSection
					romAlternatives = append(romAlternatives, newSection.sectionName)
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: newSection.sectionName, Op: bmreqs.OpAdd})
					bi.rg.Clone("/code:romtexts/sections:"+romSection, "/code:romtexts/sections:"+newSection.sectionName)
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "codelength", Value: strconv.Itoa(romSectionLength), Op: bmreqs.OpAdd})
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "datalength", Value: strconv.Itoa(romDataLength), Op: bmreqs.OpAdd})
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:romtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "finallength", Value: strconv.Itoa(romFinalLength), Op: bmreqs.OpAdd})
				}

				// Same for RAM code
				if ramCode != nil {
					newSection := new(BasmSection)
					newSection.sectionName = "ramcode" + cpNewSectionName
					newSection.sectionType = sectRamText
					newSection.sectionBody = ramCode.sectionBody.Copy()
					bi.resolveSymbols(newSection, cpNewSectionName)
					bi.sections[newSection.sectionName] = newSection
					ramAlternatives = append(ramAlternatives, newSection.sectionName)
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramtexts", T: bmreqs.ObjectSet, Name: "sections", Value: newSection.sectionName, Op: bmreqs.OpAdd})
					bi.rg.Clone("/code:ramtexts/sections:"+ramSection, "/code:ramtexts/sections:"+newSection.sectionName)
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "codelength", Value: strconv.Itoa(ramSectionLength), Op: bmreqs.OpAdd})
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "datalength", Value: strconv.Itoa(ramDataLength), Op: bmreqs.OpAdd})
					bi.rg.Requirement(bmreqs.ReqRequest{Node: "code:ramtexts/sections:" + newSection.sectionName, T: bmreqs.ObjectMax, Name: "finallength", Value: strconv.Itoa(ramFinalLength), Op: bmreqs.OpAdd})
				}
			}
		}

		switch len(romAlternatives) {
		case 0:
		case 1:
			cp.BasmMeta = cp.BasmMeta.SetMeta("romcode", romAlternatives[0])
		default:
			cp.BasmMeta = cp.BasmMeta.SetMeta("romalternatives", strings.Join(romAlternatives, ":"))
			cp.BasmMeta.RmMeta("romcode")
		}

		switch len(ramAlternatives) {
		case 0:
		case 1:
			cp.BasmMeta = cp.BasmMeta.SetMeta("ramcode", ramAlternatives[0])
		default:
			cp.BasmMeta = cp.BasmMeta.SetMeta("ramalternatives", strings.Join(ramAlternatives, ":"))
			cp.BasmMeta.RmMeta("ramcode")
		}

		switch len(romDataAlternatives) {
		case 0:
		case 1:
			cp.BasmMeta = cp.BasmMeta.SetMeta("romdata", romDataAlternatives[0])
		default:
			cp.BasmMeta = cp.BasmMeta.SetMeta("romdataalternatives", strings.Join(romDataAlternatives, ":"))
			cp.BasmMeta.RmMeta("romdata")
		}

		switch len(ramDataAlternatives) {
		case 0:
		case 1:
			cp.BasmMeta = cp.BasmMeta.SetMeta("ramdata", ramDataAlternatives[0])
		default:
			cp.BasmMeta = cp.BasmMeta.SetMeta("ramdataalternatives", strings.Join(ramDataAlternatives, ":"))
			cp.BasmMeta.RmMeta("ramdata")
		}
	}

	return nil

}
