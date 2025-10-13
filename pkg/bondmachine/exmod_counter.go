package bondmachine

import (
	"encoding/json"
	"errors"
	"strconv"
)

type CounterExtra struct {
	SlowFactor  string
	MappedInput string
}

func (sl *CounterExtra) Get_Name() string {
	return "slow"
}

func (sl *CounterExtra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"slow_factor": sl.SlowFactor}
	result.Params = map[string]string{"mapped_input": sl.MappedInput}
	return result
}
func (sl *CounterExtra) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("unmarshalling failed")
	}
	return nil
}

func (sl *CounterExtra) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *CounterExtra) Check(bmach *Bondmachine) error {
	if numeric, err := strconv.Atoi(sl.SlowFactor); err != nil {
		return errors.New("conversion failed")
	} else {
		if numeric < 0 || numeric > 31 {
			return errors.New("slow factor outside the 0-31 range")
		}
	}
	return nil
}

func (sl *CounterExtra) Verilog_headers() string {
	return ""
}

func (sl *CounterExtra) StaticVerilog() string {
	return ""
}

func (sl *CounterExtra) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
