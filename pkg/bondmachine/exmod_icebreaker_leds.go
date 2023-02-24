package bondmachine

import (
	"encoding/json"
	"errors"
)

type IcebreakerLeds struct {
	MappedOutput string
}

func (sl *IcebreakerLeds) Get_Name() string {
	return "icebreaker_leds"
}

func (sl *IcebreakerLeds) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_output": sl.MappedOutput}
	return result
}
func (sl *IcebreakerLeds) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("unmarshalling failed")
	}
	return nil
}

func (sl *IcebreakerLeds) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *IcebreakerLeds) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *IcebreakerLeds) Verilog_headers() string {
	return ""
}

func (sl *IcebreakerLeds) StaticVerilog() string {
	return ""
}

func (sl *IcebreakerLeds) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
