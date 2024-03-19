package bondmachine

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

type AXIsTemplateData struct {
	*templateData
	Samples      int
	FifoDepth    int
	CountersBits int
}

func (bmach *Bondmachine) WriteBMAPI(conf *Config, flavor string, iomaps *IOmap, extramods []ExtraModule, sbox *simbox.Simbox) error {

	var bmapiFlavor string
	var bmapiFlavorVersion string
	var bmapiLanguage string
	var bmapiFramework string
	var bmapiLibOutDir string
	var bmapiModOutDir string
	var bmapiAuxOutDir string
	var bmapiPackageName string
	var bmapiModuleName string
	var bmapiGenerateExample string
	var bmapiDataType string

	var bmapiParams map[string]string

	// Extracting and check of BMAPI params
	for _, mod := range extramods {
		if mod.Get_Name() == "bmapi" {
			bmapiParams = mod.Get_Params().Params

			if val, ok := bmapiParams["bmapi_flavor"]; ok {
				bmapiFlavor = val
			} else {
				return errors.New("Missing bmapi flavor")
			}

			if val, ok := bmapiParams["bmapi_flavor_version"]; ok {
				bmapiFlavorVersion = val
			} else {
				return errors.New("Missing bmapi flavor version")
			}

			if val, ok := bmapiParams["bmapi_language"]; ok {
				bmapiLanguage = val
			} else {
				return errors.New("Missing bmapi language")
			}

			if val, ok := bmapiParams["bmapi_framework"]; ok {
				bmapiFramework = val
			} else {
				return errors.New("Missing bmapi framework")
			}

			if val, ok := bmapiParams["bmapi_liboutdir"]; ok {
				bmapiLibOutDir = val
			} else {
				return errors.New("Missing bmapi liboutdir")
			}

			if val, ok := bmapiParams["bmapi_modoutdir"]; ok {
				bmapiModOutDir = val
			} else {
				return errors.New("Missing bmapi modoutdir")
			}

			if val, ok := bmapiParams["bmapi_auxoutdir"]; ok {
				bmapiAuxOutDir = val
			} else {
				return errors.New("Missing bmapi auxoutdir")
			}

			if val, ok := bmapiParams["bmapi_packagename"]; ok {
				bmapiPackageName = val
			} else {
				return errors.New("Missing bmapi packagename")
			}

			if val, ok := bmapiParams["bmapi_modulename"]; ok {
				bmapiModuleName = val
			} else {
				return errors.New("Missing bmapi modulename")
			}

			if val, ok := bmapiParams["bmapi_generate_example"]; ok {
				bmapiGenerateExample = val
			} else {
				return errors.New("Missing bmapi generate example")
			}

			if val, ok := bmapiParams["bmapi_datatype"]; ok {
				bmapiDataType = val
			} else {
				return errors.New("Missing bmapi datatype")
			}

			break
		}
	}

	switch bmapiFlavor {
	case "axist":

		// This is the generation of the AXI Stream interface

		axiStData := new(AXIsTemplateData)
		axiStData.templateData = bmach.createBasicTemplateData()
		axiStData.ModuleName = "bmaccelerator"
		axiStData.InputNum = bmach.Inputs
		axiStData.OutputNum = bmach.Outputs
		axiStData.funcmap = template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
			"dec": func(i int) int {
				return i - 1
			},
			"add": func(i, j int) int {
				return i + j
			},
			"sub": func(i, j int) int {
				return i - j
			},
		}

		// This fields are temporarely hardcoded, in the future could be get from the command line
		axiStData.Samples = 16
		axiStData.FifoDepth = 256

		if axiStData.InputNum > axiStData.OutputNum {
			axiStData.CountersBits = Needed_bits(axiStData.InputNum * axiStData.FifoDepth)
		} else {
			axiStData.CountersBits = Needed_bits(axiStData.OutputNum * axiStData.FifoDepth)
		}

		axiStData.Inputs = make([]string, 0)
		axiStData.Outputs = make([]string, 0)

		sortedKeys := make([]string, 0)
		for param, _ := range bmapiParams {
			sortedKeys = append(sortedKeys, param)
		}

		sort.Slice(sortedKeys, func(i, j int) bool {
			first := sortedKeys[i]
			second := sortedKeys[j]
			for {
				if len(first) == 0 || len(second) == 0 {
					return first < second
				} else {
					if first[0] != second[0] {
						return first < second
					} else {
						first = first[1:]
						second = second[1:]

						if numA, err := strconv.Atoi(first); err == nil {
							if numB, err := strconv.Atoi(second); err == nil {
								return numA < numB
							}
						}
					}
				}
			}
		})

		for _, param := range sortedKeys {
			if strings.HasPrefix(param, "assoc") {
				bmport := strings.Split(param, "_")[1]
				if strings.HasPrefix(bmport, "o") {
					axiStData.Outputs = append(axiStData.Outputs, bmport)
				} else if strings.HasPrefix(bmport, "i") {
					axiStData.Inputs = append(axiStData.Inputs, bmport)
				}
			}
		}

		switch flavor {
		case "alveou50", "alveou55c":
			// This is the generation of the directory for the kernel project
			if _, err := os.Stat(bmapiModOutDir); os.IsNotExist(err) {
				os.Mkdir(bmapiModOutDir, 0700)
				os.Mkdir(bmapiModOutDir+"/src", 0700)
				os.Mkdir(bmapiModOutDir+"/src/krnl_bondmachine", 0700)
				os.Mkdir(bmapiModOutDir+"/src/krnl_bondmachine/hdl", 0700)
			} else {
				return errors.New("BMAPI modoutdir already exists")
			}
			// Adding the files to the project
			vFiles := make(map[string]string)
			vFiles["krnl_bondmachine_rtl.v"] = krnlBondmachineRTL
			vFiles["krnl_bondmachine_rtl_axi_read_master.sv"] = krnlBondmachineRTLAxiReadMaster
			vFiles["krnl_bondmachine_rtl_axi_write_master.sv"] = krnlBondmachineRTLAxiWriteMaster
			vFiles["krnl_bondmachine_rtl_caller.sv"] = krnlBondmachineRTLCaller
			vFiles["krnl_bondmachine_rtl_control_s_axi.v"] = krnlBondmachineRTLControlSAxi
			vFiles["krnl_bondmachine_rtl_counter.sv"] = krnlBondmachineRTLCounter
			vFiles["krnl_bondmachine_rtl_int.sv"] = krnlBondmachineRTLInt
			if bmapiFlavorVersion == "basic" {
				vFiles["bmaccelerator_v1_0.v"] = basicAXIStream
			} else if bmapiFlavorVersion == "optimized" {
				vFiles["bmaccelerator_v1_0.v"] = optimizedAXIStream
			} else {
				return errors.New("unknown AXI Stream flavor version")
			}

			for file, temp := range vFiles {
				t, err := template.New(file).Funcs(axiStData.funcmap).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(bmapiModOutDir + "/src/krnl_bondmachine/hdl/" + file)
				if err != nil {
					return err
				}

				err = t.Execute(f, axiStData)
				if err != nil {
					return err
				}

				f.Close()
			}

		case "zedboard", "ebaz4205", "zc702":
			vFiles := make(map[string]string)
			if bmapiFlavorVersion == "basic" {
				vFiles["axistream.v"] = basicAXIStream
			} else if bmapiFlavorVersion == "optimized" {
				vFiles["axistream.v"] = optimizedAXIStream
			} else {
				return errors.New("unknown AXI Stream flavor version")
			}

			for file, temp := range vFiles {
				t, err := template.New(file).Funcs(axiStData.funcmap).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(file)
				if err != nil {
					return err
				}

				err = t.Execute(f, axiStData)
				if err != nil {
					return err
				}

				f.Close()
			}
		default:
			return errors.New("unknown board")
		}

		switch bmapiLanguage {
		case "python":
			switch bmapiFramework {
			case "pynq":
				if bmapiGenerateExample != "" {
					if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
						os.Mkdir(bmapiLibOutDir, 0700)
					} else {
						return errors.New("BMAPI liboutdir already exists")
					}

					// Compiling the data for the templates
					bmapiExample := bmach.createBasicTemplateData()

					exFiles := make(map[string]string)
					exFiles[bmapiGenerateExample] = axistPynqExample

					for file, temp := range exFiles {
						t, err := template.New(file).Parse(temp)
						if err != nil {
							return err
						}

						f, err := os.Create(bmapiLibOutDir + "/" + file)
						if err != nil {
							return err
						}

						err = t.Execute(f, bmapiExample)
						if err != nil {
							return err
						}

						f.Close()
					}
				}
			}
		}

	case "aximm":

		// This is the generation of the Linux kernel module
		if _, err := os.Stat(bmapiModOutDir); os.IsNotExist(err) {
			os.Mkdir(bmapiModOutDir, 0700)
		} else {
			return errors.New("BMAPI modoutdir already exists")
		}

		kmoddata := bmach.createBasicTemplateData()

		modFiles := make(map[string]string)
		modFiles["bm.c"] = moduleFilesBm

		for file, temp := range modFiles {
			t, err := template.New(file).Parse(temp)
			if err != nil {
				return err
			}

			f, err := os.Create(bmapiModOutDir + "/" + file)
			if err != nil {
				return err
			}

			err = t.Execute(f, kmoddata)
			if err != nil {
				return err
			}

			f.Close()
		}

		// This is the generation of the AXI auxiliary files
		if _, err := os.Stat(bmapiAuxOutDir); os.IsNotExist(err) {
			os.Mkdir(bmapiAuxOutDir, 0700)
		} else {
			return errors.New("BMAPI auxoutdir already exists")
		}

		auxdata := bmach.createBasicTemplateData()

		auxdata.Inputs = make([]string, 0)
		auxdata.Outputs = make([]string, 0)

		sortedKeys := make([]string, 0)
		for param, _ := range bmapiParams {
			sortedKeys = append(sortedKeys, param)
		}

		sort.Slice(sortedKeys, func(i, j int) bool {
			first := sortedKeys[i]
			second := sortedKeys[j]
			for {
				if len(first) == 0 || len(second) == 0 {
					return first < second
				} else {
					if first[0] != second[0] {
						return first < second
					} else {
						first = first[1:]
						second = second[1:]

						if numA, err := strconv.Atoi(first); err == nil {
							if numB, err := strconv.Atoi(second); err == nil {
								return numA < numB
							}
						}
					}
				}
			}
		})

		for _, param := range sortedKeys {
			if strings.HasPrefix(param, "assoc") {
				bmport := strings.Split(param, "_")[1]
				if strings.HasPrefix(bmport, "o") {
					auxdata.Outputs = append(auxdata.Outputs, "port_"+bmport)
				} else if strings.HasPrefix(bmport, "i") {
					auxdata.Inputs = append(auxdata.Inputs, "port_"+bmport)
				}
			}
		}

		auxFiles := make(map[string]string)
		auxFiles["axipatch.txt"] = auxfilesAXIPatch
		auxFiles["outregs.txt"] = auxfilesAXIOutRegs
		auxFiles["designexternal.txt"] = auxfilesDesignExternal
		auxFiles["designexternalinst.txt"] = auxfilesDesignExternalInst

		for file, temp := range auxFiles {
			t, err := template.New(file).Funcs(auxdata.funcmap).Parse(temp)
			if err != nil {
				return err
			}

			f, err := os.Create(bmapiAuxOutDir + "/" + file)
			if err != nil {
				return err
			}

			err = t.Execute(f, auxdata)
			if err != nil {
				return err
			}

			f.Close()
		}

		// Tivial aux files
		f, err := os.Create(bmapiAuxOutDir + "/axiregnum.txt")
		if err != nil {
			return err
		}
		f.Write([]byte(fmt.Sprintf("%d", len(auxdata.Inputs)+len(auxdata.Outputs)+4)))
		f.Close()

		switch bmapiLanguage {
		case "c":
			if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
				os.Mkdir(bmapiLibOutDir, 0700)
			} else {
				return errors.New("BMAPI liboutdir already exists")
			}

			// Compiling the data for the templates
			bmapidata := bmach.createBasicTemplateData()
			bmapidata.PackageName = bmapiPackageName
			var _ = bmapiModuleName // TODO TEMP
			cFiles := make(map[string]string)
			cFiles["Makefile"] = cFilesMakefile

			for file, temp := range cFiles {
				t, err := template.New(file).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(bmapiLibOutDir + "/" + file)
				if err != nil {
					return err
				}

				err = t.Execute(f, bmapidata)
				if err != nil {
					return err
				}

				f.Close()
			}
		case "go":
			if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
				os.Mkdir(bmapiLibOutDir, 0700)
			} else {
				return errors.New("BMAPI liboutdir already exists")
			}

			// Compiling the data for the templates
			bmapidata := bmach.createBasicTemplateData()
			bmapidata.PackageName = bmapiPackageName
			var _ = bmapiModuleName // TODO TEMP
			apiFiles := make(map[string]string)
			apiFiles["bmapi.go"] = bmapi
			apiFiles["encoder.go"] = bmapiEncoder
			apiFiles["decoder.go"] = bmapiDecoder
			apiFiles["commands.go"] = bmapiCommands
			apiFiles["functions.go"] = bmapiFunctions
			apiFiles["go.mod"] = bmapigomod

			for file, temp := range apiFiles {
				t, err := template.New(file).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(bmapiLibOutDir + "/" + file)
				if err != nil {
					return err
				}

				err = t.Execute(f, bmapidata)
				if err != nil {
					return err
				}

				f.Close()
			}

		case "python":
			switch bmapiFramework {
			case "pynq":
				if bmapiGenerateExample != "" {
					if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
						os.Mkdir(bmapiLibOutDir, 0700)
					} else {
						return errors.New("BMAPI liboutdir already exists")
					}

					// Compiling the data for the templates
					bmapiExample := bmach.createBasicTemplateData()
					bmapiExample.DataType = bmapiDataType

					exFiles := make(map[string]string)
					exFiles[bmapiGenerateExample] = aximmPynqExample

					for file, temp := range exFiles {
						t, err := template.New(file).Parse(temp)
						if err != nil {
							return err
						}

						f, err := os.Create(bmapiLibOutDir + "/" + file)
						if err != nil {
							return err
						}

						err = t.Execute(f, bmapiExample)
						if err != nil {
							return err
						}

						f.Close()
					}
				}
			}

		}
	case "uartusb":
		switch bmapiLanguage {
		case "c":
			if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
				os.Mkdir(bmapiLibOutDir, 0700)
			} else {
				return errors.New("BMAPI liboutdir already exists")
			}

			// Compiling the data for the templates
			bmapidata := bmach.createBasicTemplateData()
			bmapidata.PackageName = bmapiPackageName
			var _ = bmapiModuleName // TODO TEMP
			cFiles := make(map[string]string)
			cFiles["Makefile"] = cFilesMakefile

			for file, temp := range cFiles {
				t, err := template.New(file).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(bmapiLibOutDir + "/" + file)
				if err != nil {
					return err
				}

				err = t.Execute(f, bmapidata)
				if err != nil {
					return err
				}

				f.Close()
			}
		case "go":
			if _, err := os.Stat(bmapiLibOutDir); os.IsNotExist(err) {
				os.Mkdir(bmapiLibOutDir, 0700)
			} else {
				return errors.New("BMAPI liboutdir already exists")
			}

			// Compiling the data for the templates
			bmapidata := bmach.createBasicTemplateData()
			bmapidata.PackageName = bmapiPackageName
			var _ = bmapiModuleName // TODO TEMP
			apiFiles := make(map[string]string)
			apiFiles["bmapi.go"] = bmapi
			apiFiles["encoder.go"] = bmapiEncoder
			apiFiles["decoder.go"] = bmapiDecoder
			apiFiles["commands.go"] = bmapiCommands
			apiFiles["functions.go"] = bmapiFunctions
			apiFiles["go.mod"] = bmapigomod

			for file, temp := range apiFiles {
				t, err := template.New(file).Parse(temp)
				if err != nil {
					return err
				}

				f, err := os.Create(bmapiLibOutDir + "/" + file)
				if err != nil {
					return err
				}

				err = t.Execute(f, bmapidata)
				if err != nil {
					return err
				}

				f.Close()
			}

		default:
			return errors.New("unimplemented language")
		}
	default:
		return errors.New("unknown bmapi flavor")
	}

	return nil
}
