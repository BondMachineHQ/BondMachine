package bondmachine

import (
	"encoding/json"
	"errors"
)

type UartExtra struct {
	Maps *IOmap
}

func (sl *UartExtra) Get_Name() string {
	return "uart"
}

func (sl *UartExtra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)
	for n, v := range sl.Maps.Assoc {
		result.Params[n] = v
	}
	return result
}

func (sl *UartExtra) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("Unmarshalling failed")
	}
	return nil
}

func (sl *UartExtra) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *UartExtra) Check(bmach *Bondmachine) error {
	// TODO check if the pins are valid
	return nil
}

func (sl *UartExtra) Verilog_headers() string {
	return ""
}

func (sl *UartExtra) StaticVerilog() string {
	return ""
}

func (sl *UartExtra) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
