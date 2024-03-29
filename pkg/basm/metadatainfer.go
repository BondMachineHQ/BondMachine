package basm

import (
	"fmt"
	"regexp"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bmreqs"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func (bi *BasmInstance) bodyMetadataInfer(body *bmline.BasmBody, soShortNames []string) error {
	symbolFilteredMatchers := make([]*bmline.BasmLine, 0)
	for _, matcher := range bi.matchers {
		if bmline.FilterMatcher(matcher, "symbol") {
			symbolFilteredMatchers = append(symbolFilteredMatchers, matcher)
		}
	}

	revSymbolFilteredMatchers := make([]*bmline.BasmLine, 0)
	for _, matcher := range bi.matchers {
		if bmline.FilterMatcher(matcher, "revsymbol") {
			revSymbolFilteredMatchers = append(symbolFilteredMatchers, matcher)
		}
	}
	for _, line := range body.Lines {

		if bi.debug {
			fmt.Println(green("\t\t\tLine: ") + line.String())
		}

		for _, arg := range line.Elements {
			re := regexp.MustCompile("^r[0-9]+$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "reg")
				continue
			}
			re = regexp.MustCompile("^\\[(?P<reg>r[0-9]+)\\]$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "loc")
				arg.BasmMeta = arg.SetMeta("locaddressing", "register")
				regAddr := re.ReplaceAllString(arg.GetValue(), "${reg}")
				arg.BasmMeta = arg.SetMeta("locregister", regAddr)
				continue
			}
			// re = regexp.MustCompile("^[0-9]+$")
			// if re.MatchString(arg.GetValue()) {
			// 	arg.BasmMeta = arg.SetMeta("type", "number")
			// 	arg.BasmMeta = arg.SetMeta("numbertype", "decimal")
			// }
			// re = regexp.MustCompile("^0f[+-]?([0-9]*[.])?[0-9]+$")
			// if re.MatchString(arg.GetValue()) {
			// 	arg.BasmMeta = arg.SetMeta("type", "number")
			// 	arg.BasmMeta = arg.SetMeta("numbertype", "float")
			// }

			re = regexp.MustCompile("^rom:(?P<location>[0-9]+)$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "rom")
				arg.BasmMeta = arg.SetMeta("romaddressing", "immediate")
				location := re.ReplaceAllString(arg.GetValue(), "${location}")
				arg.BasmMeta = arg.SetMeta("location", location)
				continue
			}
			re = regexp.MustCompile("^rom:(?P<symb>[0-9a-zA-Z_\\.]+)$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "rom")
				arg.BasmMeta = arg.SetMeta("romaddressing", "symbol")
				symbol := re.ReplaceAllString(arg.GetValue(), "${symb}")
				arg.BasmMeta = arg.SetMeta("symbol", symbol)
				continue
			}
			re = regexp.MustCompile("^rom:\\[(?P<reg>r[0-9]+)\\]$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "rom")
				arg.BasmMeta = arg.SetMeta("romaddressing", "register")
				regAddr := re.ReplaceAllString(arg.GetValue(), "${reg}")
				arg.BasmMeta = arg.SetMeta("romregister", regAddr)
				continue
			}
			re = regexp.MustCompile("^ram:(?P<location>[0-9]+)$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "ram")
				arg.BasmMeta = arg.SetMeta("ramaddressing", "immediate")
				location := re.ReplaceAllString(arg.GetValue(), "${location}")
				arg.BasmMeta = arg.SetMeta("location", location)
				continue
			}
			re = regexp.MustCompile("^ram:(?P<symb>[0-9a-zA-Z_\\.]+)$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "ram")
				arg.BasmMeta = arg.SetMeta("ramaddressing", "symbol")
				symbol := re.ReplaceAllString(arg.GetValue(), "${symb}")
				arg.BasmMeta = arg.SetMeta("symbol", symbol)
				continue
			}
			re = regexp.MustCompile("^ram:\\[(?P<reg>r[0-9]+)\\]$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "ram")
				arg.BasmMeta = arg.SetMeta("ramaddressing", "register")
				regAddr := re.ReplaceAllString(arg.GetValue(), "${reg}")
				arg.BasmMeta = arg.SetMeta("ramregister", regAddr)
				continue
			}
			re = regexp.MustCompile("^(?P<input>i[0-9]+):?$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "input")
				input := re.ReplaceAllString(arg.GetValue(), "${input}")
				arg.SetValue(input)
				switch line.GetMeta("iomode") {
				case "async":
					line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
				case "sync":
					line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
				default:
					switch body.GetMeta("iomode") {
					case "async":
						line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
					case "sync":
						line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
					default:
						switch bi.global.GetMeta("iomode") {
						case "async":
							line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
						case "sync":
							line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
						}
					}
				}
				continue
			}
			re = regexp.MustCompile("^(?P<output>o[0-9]+):?$")
			if re.MatchString(arg.GetValue()) {
				arg.BasmMeta = arg.SetMeta("type", "output")
				output := re.ReplaceAllString(arg.GetValue(), "${output}")
				arg.SetValue(output)
				switch line.GetMeta("iomode") {
				case "async":
					line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
				case "sync":
					line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
				default:
					switch body.GetMeta("iomode") {
					case "async":
						line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
					case "sync":
						line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
					default:
						switch bi.global.GetMeta("iomode") {
						case "async":
							line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "async")
						case "sync":
							line.Operation.BasmMeta = line.Operation.SetMeta("iomode", "sync")
						}
					}
				}
				continue
			}

			for _, soname := range soShortNames {
				re = regexp.MustCompile("^" + soname + "(?P<index>[0-9]+)$")
				if re.MatchString(arg.GetValue()) {
					index := re.ReplaceAllString(arg.GetValue(), "${index}")
					arg.BasmMeta = arg.SetMeta("type", "somov")
					arg.BasmMeta = arg.SetMeta("sotype", soname)
					arg.BasmMeta = arg.SetMeta("soaddressing", "immediate")
					arg.BasmMeta = arg.SetMeta("soindex", index)
					arg.BasmMeta = arg.SetMeta("soport", "0")
				}
				re = regexp.MustCompile("^" + soname + "(?P<index>[0-9]+):(?P<port>[0-9]+)$")
				if re.MatchString(arg.GetValue()) {
					index := re.ReplaceAllString(arg.GetValue(), "${index}")
					port := re.ReplaceAllString(arg.GetValue(), "${port}")
					arg.BasmMeta = arg.SetMeta("type", "somov")
					arg.BasmMeta = arg.SetMeta("sotype", soname)
					arg.BasmMeta = arg.SetMeta("soaddressing", "immediate")
					arg.BasmMeta = arg.SetMeta("soindex", index)
					arg.BasmMeta = arg.SetMeta("soport", port)
				}
				re = regexp.MustCompile("^" + soname + "(?P<index>[0-9]+):\\[(?P<reg>r[0-9]+)\\]$")
				if re.MatchString(arg.GetValue()) {
					index := re.ReplaceAllString(arg.GetValue(), "${index}")
					reg := re.ReplaceAllString(arg.GetValue(), "${reg}")
					arg.BasmMeta = arg.SetMeta("type", "somov")
					arg.BasmMeta = arg.SetMeta("sotype", soname)
					arg.BasmMeta = arg.SetMeta("soaddressing", "register")
					arg.BasmMeta = arg.SetMeta("soindex", index)
					arg.BasmMeta = arg.SetMeta("soregister", reg)
				}
				re = regexp.MustCompile("^" + soname + "(?P<index>[0-9]+):\\[(?P<addr>[0-9]+)\\]$")
				if re.MatchString(arg.GetValue()) {
					index := re.ReplaceAllString(arg.GetValue(), "${index}")
					addr := re.ReplaceAllString(arg.GetValue(), "${addr}")
					arg.BasmMeta = arg.SetMeta("type", "somov")
					arg.BasmMeta = arg.SetMeta("sotype", soname)
					arg.BasmMeta = arg.SetMeta("soaddressing", "direct")
					arg.BasmMeta = arg.SetMeta("soindex", index)
					arg.BasmMeta = arg.SetMeta("soaddr", addr)
				}
			}

			if bmNumber, err := bmnumbers.ImportString(arg.GetValue()); err == nil {
				arg.BasmMeta = arg.SetMeta("type", "number")
				arg.BasmMeta = arg.SetMeta("numbertype", bmNumber.GetTypeName())
			}
		}

		// Every unknown arg is a symbol if the operation is a symbol matcher
		for _, matcher := range symbolFilteredMatchers {
			if bmline.MatchArg(matcher.Operation, line.Operation) {
				for _, arg := range line.Elements {
					if arg.GetMeta("type") == "" {
						arg.BasmMeta = arg.SetMeta("type", "symbol")
					}
				}
			}
		}

		// Every unknown arg is a revsymbol if the operation is a revsymbol matcher
		for _, matcher := range revSymbolFilteredMatchers {
			if bmline.MatchArg(matcher.Operation, line.Operation) {
				for _, arg := range line.Elements {
					if arg.GetMeta("type") == "" {
						arg.BasmMeta = arg.SetMeta("type", "revsymbol")
					}
				}
			}
		}
	}

	//TODO add more meta data
	return nil
}

func metadataInfer(bi *BasmInstance) error {
	// TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing sections:"))
	}

	soShortNames := make([]string, 0)

	for _, so := range procbuilder.Allshared {
		soShortNames = append(soShortNames, so.Shortname())
	}

	// Loop over the sections
	for sectName, section := range bi.sections {
		if section.sectionType == sectRomText || section.sectionType == sectRamText {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			if section.sectionType == sectRomText {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			} else {
				bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			}

			body := section.sectionBody

			if err := bi.bodyMetadataInfer(body, soShortNames); err != nil {
				return err
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	for sectName, section := range bi.sections {
		if section.sectionType == sectRomData || section.sectionType == sectRamData {
			if bi.debug {
				fmt.Println(green("\t\tSection: ") + sectName)
			}

			// if section.sectionType == sectRomText {
			// 	bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:romtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			// } else {
			// 	bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:ramtexts", T: bmreqs.ObjectSet, Name: "sections", Value: sectName, Op: bmreqs.OpAdd})
			// }

			body := section.sectionBody

			if err := bi.bodyMetadataInfer(body, soShortNames); err != nil {
				return err
			}
		} else {
			if bi.debug {
				fmt.Println(yellow("\t\tSection type not handled: ") + sectName)
			}
		}
	}

	for fragName, frag := range bi.fragments {
		if bi.debug {
			fmt.Println(green("\t\tFragment: ") + fragName)
		}

		bi.rg.Requirement(bmreqs.ReqRequest{Node: "/code:fragments", T: bmreqs.ObjectSet, Name: "fragments", Value: fragName, Op: bmreqs.OpAdd})

		body := frag.fragmentBody

		if err := bi.bodyMetadataInfer(body, soShortNames); err != nil {
			return err
		} else {
			if bi.debug {
				fmt.Println(green("\t\t\tFragment body metadata inferred"))
			}
		}
	}

	return nil
}
