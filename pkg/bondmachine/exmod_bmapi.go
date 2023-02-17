package bondmachine

import "strings"

//	"fmt"

//	"strings"

type BMAPIExtra struct {
	Maps          *IOmap
	Language      string
	Flavor        string // "uartusb" or "aximm" or "axist"
	FlavorVersion string
	Framework     string // "pynq" or ""
	LibOutDir     string
	ModOutDir     string
	AuxOutDir     string
	PackageName   string
	ModuleName    string
	Rsize         uint8
}

func (sl *BMAPIExtra) Get_Name() string {
	return "bmapi"
}

func (sl *BMAPIExtra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)

	result.Params["bmapi_language"] = sl.Language
	result.Params["bmapi_flavor"] = sl.Flavor
	result.Params["bmapi_flavor_version"] = sl.FlavorVersion
	result.Params["bmapi_framework"] = sl.Framework
	result.Params["bmapi_liboutdir"] = sl.LibOutDir
	result.Params["bmapi_auxoutdir"] = sl.AuxOutDir
	result.Params["bmapi_modoutdir"] = sl.ModOutDir
	result.Params["bmapi_packagename"] = sl.PackageName
	result.Params["bmapi_modulename"] = sl.ModuleName

	result.Params["inputs"] = ""
	result.Params["outputs"] = ""

	for bmport, apiport := range sl.Maps.Assoc {
		result.Params["assoc_"+bmport] = apiport
		if strings.HasPrefix(bmport, "i") {
			if result.Params["inputs"] != "" {
				result.Params["inputs"] += ","
			}
			result.Params["inputs"] += bmport
		} else {
			if result.Params["outputs"] != "" {
				result.Params["outputs"] += ","
			}
			result.Params["outputs"] += bmport
		}
	}

	return result
}

func (sl *BMAPIExtra) Import(inp string) error {
	return nil
}

func (sl *BMAPIExtra) Export() string {
	return ""
}

func (sl *BMAPIExtra) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *BMAPIExtra) Verilog_headers() string {
	result := "\n"
	return result
}
func (sl *BMAPIExtra) StaticVerilog() string {
	result := "\n"

	switch sl.Flavor {
	case "uartusb":
		result += verilogUART
	}
	return result
}

func (sl *BMAPIExtra) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
