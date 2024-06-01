package bondmachine

import (
	"encoding/json"
	"errors"
)

type Ice40Lp1kLeds struct {
	MappedOutput string
}

func (sl *Ice40Lp1kLeds) Get_Name() string {
	return "ice40lp1k_leds"
}

func (sl *Ice40Lp1kLeds) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_output": sl.MappedOutput}
	return result
}
func (sl *Ice40Lp1kLeds) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("unmarshalling failed")
	}
	return nil
}

func (sl *Ice40Lp1kLeds) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *Ice40Lp1kLeds) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *Ice40Lp1kLeds) Verilog_headers() string {
	return ""
}

func (sl *Ice40Lp1kLeds) StaticVerilog() string {
	return ""
}

func (sl *Ice40Lp1kLeds) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
