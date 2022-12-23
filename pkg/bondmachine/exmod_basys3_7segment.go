package bondmachine

import (
	"encoding/json"
	"errors"
)

type B37s struct {
	Mapped_output string
}

func (sl *B37s) Get_Name() string {
	return "basys3_7segment"
}

func (sl *B37s) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_output": sl.Mapped_output}
	return result
}
func (sl *B37s) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("Unmarshalling failed")
	}
	return nil
}

func (sl *B37s) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *B37s) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *B37s) Verilog_headers() string {
	return ""
}

func (sl *B37s) StaticVerilog() string {
	return ""
}

func (sl *B37s) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
