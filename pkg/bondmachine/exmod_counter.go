package bondmachine

import (
	"encoding/json"
	"errors"
	"strconv"
)

type CounterExtra struct {
	SlowFactor  string
	MappedInput string
	Width       string
}

func (sl *CounterExtra) Get_Name() string {
	return "counter"
}

func (sl *CounterExtra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)
	result.Params["slow_factor"] = sl.SlowFactor
	result.Params["mapped_input"] = sl.MappedInput
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
	files := make([]string, 0)
	code := make([]string, 0)

	width, _ := strconv.Atoi(sl.Width)
	slowFactor, _ := strconv.Atoi(sl.SlowFactor)

	max := width + slowFactor + 1
	sel := slowFactor + 1
	maxS := strconv.Itoa(max)
	selS := strconv.Itoa(sel)

	source := "module counter(clk, reset, mapped_input, mapped_input_valid, mapped_input_received);\n"
	source += "input clk;\n"
	source += "input reset;\n"
	source += "output reg [" + sl.Width + ":0] mapped_input;\n"
	source += "output reg mapped_input_valid;\n"
	source += "input mapped_input_received;\n"
	source += "reg [" + maxS + ":0] counter_reg;\n"
	source += "reg checkbit;\n"

	source += "always @(posedge clk) begin\n"
	source += "  if (reset) begin\n"
	source += "    counter_reg <= 0;\n"
	source += "    checkbit <= 0;\n"
	source += "  end else begin\n"
	source += "    counter_reg <= counter_reg + 1;\n"
	source += "    if (! mapped_input_received) begin\n"
	source += "      checkbit <= counter_reg[" + sl.SlowFactor + "];\n"
	source += "    end\n"
	source += "  end\n"
	source += "end\n"

	source += "always @(posedge clk) begin\n"
	source += "    if (mapped_input_received) begin\n"
	source += "      mapped_input_valid <= 0;\n"
	source += "    end else begin\n"
	source += "    if (checkbit == 1 && counter_reg[" + sl.SlowFactor + "] == 0) begin\n"
	source += "      mapped_input <= counter_reg[" + maxS + ":" + selS + "];\n"
	source += "      mapped_input_valid <= 1;\n"
	source += "    end\n"
	source += "  end\n"
	source += "end\n"

	source += "endmodule\n"

	files = append(files, "counter.v")
	code = append(code, source)

	return files, code
}
