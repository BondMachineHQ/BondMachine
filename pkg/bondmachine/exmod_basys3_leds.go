package bondmachine

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Basys3Leds struct {
	LedName      string
	MappedOutput string
	Width        string
}

func (sl *Basys3Leds) Get_Name() string {
	return "basys3_leds"
}

func (sl *Basys3Leds) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)
	result.Params["led_name"] = sl.LedName
	result.Params["mapped_output"] = sl.MappedOutput
	result.Params["width"] = sl.Width
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
	files := make([]string, 0)
	code := make([]string, 0)

	width, _ := strconv.Atoi(sl.Width)
	widthS := strconv.Itoa(width - 1)
	source := "module basys3leds(clk, reset, mapped_output, mapped_output_valid, mapped_output_received, " + sl.LedName + ");\n"
	source += "input clk;\n"
	source += "input reset;\n"
	source += "input [" + widthS + ":0] mapped_output;\n"
	source += "input mapped_output_valid;\n"
	source += "output reg mapped_output_received;\n"
	source += "output reg [" + widthS + ":0] " + sl.LedName + ";\n"

	source += "always @(posedge clk) begin\n"
	source += "  if (reset) begin\n"
	source += "    " + sl.LedName + " <= 0;\n"
	source += "    mapped_output_received <= 0;\n"
	source += "  end else begin\n"
	source += "    if (mapped_output_valid) begin\n"
	source += "      " + sl.LedName + " <= mapped_output;\n"
	source += "      mapped_output_received <= 1;\n"
	source += "    end else begin\n"
	source += "      mapped_output_received <= 0;\n"
	source += "    end\n"
	source += "  end\n"
	source += "end\n"
	source += "endmodule\n"

	files = append(files, "basys3leds.v")
	code = append(code, source)

	return files, code
}
