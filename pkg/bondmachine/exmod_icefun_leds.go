package bondmachine

import (
	"encoding/json"
	"errors"
)

type IceFunLeds struct {
	MappedOutput string
}

func (sl *IceFunLeds) Get_Name() string {
	return "icefun_leds"
}

func (sl *IceFunLeds) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_output": sl.MappedOutput}
	return result
}
func (sl *IceFunLeds) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("unmarshalling failed")
	}
	return nil
}

func (sl *IceFunLeds) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *IceFunLeds) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *IceFunLeds) Verilog_headers() string {
	return ""
}

func (sl *IceFunLeds) StaticVerilog() string {
	return ""
}

func (sl *IceFunLeds) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
