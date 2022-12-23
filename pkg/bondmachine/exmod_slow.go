package bondmachine

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Slow_extra struct {
	Slow_factor string
}

func (sl *Slow_extra) Get_Name() string {
	return "slow"
}

func (sl *Slow_extra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"slow_factor": sl.Slow_factor}
	return result
}
func (sl *Slow_extra) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("Unmarshalling failed")
	}
	return nil
}

func (sl *Slow_extra) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *Slow_extra) Check(bmach *Bondmachine) error {
	if numeric, err := strconv.Atoi(sl.Slow_factor); err != nil {
		return errors.New("Conversion failed")
	} else {
		if numeric < 0 || numeric > 31 {
			return errors.New("Slow factor outside the 0-31 range")
		}
	}
	return nil
}

func (sl *Slow_extra) Verilog_headers() string {
	return ""
}

func (sl *Slow_extra) StaticVerilog() string {
	return ""
}

func (sl *Slow_extra) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
