package bondmachine

import (
	"encoding/json"
	"errors"
)

type Basys3Leds struct {
	MappedOutput string
}

func (sl *Basys3Leds) Get_Name() string {
	return "icefun_leds"
}

func (sl *Basys3Leds) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_output": sl.MappedOutput}
	return result
}
func (sl *Basys3Leds) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("unmarshalling failed")
	}
	return nil
}

func (sl *Basys3Leds) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *Basys3Leds) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *Basys3Leds) Verilog_headers() string {
	return ""
}

func (sl *Basys3Leds) StaticVerilog() string {
	return ""
}

func (sl *Basys3Leds) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
