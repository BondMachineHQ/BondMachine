package bondmachine

import (
	"encoding/json"
	"errors"
)

type Ps2KeyboardExtra struct {
	MappedInput string
}

func (sl *Ps2KeyboardExtra) Get_Name() string {
	return "ps2keyboard"
}

func (sl *Ps2KeyboardExtra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = map[string]string{"mapped_input": sl.MappedInput}
	return result
}

func (sl *Ps2KeyboardExtra) Import(inp string) error {
	if err := json.Unmarshal([]byte(inp), sl); err != nil {
		return errors.New("Unmarshalling failed")
	}
	return nil
}

func (sl *Ps2KeyboardExtra) Export() string {
	b, _ := json.Marshal(sl)
	return string(b)
}

func (sl *Ps2KeyboardExtra) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *Ps2KeyboardExtra) Verilog_headers() string {
	return ""
}

func (sl *Ps2KeyboardExtra) StaticVerilog() string {
	result := "\n`timescale 1ns / 1ps\n"
	result += `
module debouncer(
    input clk,
    input I,
    output reg O
    );
    parameter COUNT_MAX=255, COUNT_WIDTH=8;
    reg [COUNT_WIDTH-1:0] count;
    reg Iv=0;
    always@(posedge clk)
        if (I == Iv) begin
            if (count == COUNT_MAX)
                O <= I;
            else
                count <= count + 1'b1;
        end else begin
            count <= 'b0;
            Iv <= I;
        end
endmodule
	`

	result += "\n`timescale 1ns / 1ps\n"
	result += `
module PS2Receiver(
	input clk,
	input kclk,
	input kdata,
	output reg [15:0] keycode=0,
	output reg oflag
	);
		
	wire kclkf, kdataf;
	reg [7:0]datacur=0;
	reg [7:0]dataprev=0;
	reg [3:0]cnt=0;
	reg flag=0;
		
	debouncer #(
		.COUNT_MAX(19),
		.COUNT_WIDTH(5)
	) db_clk(
		.clk(clk),
		.I(kclk),
		.O(kclkf)
	);
	debouncer #(
	   .COUNT_MAX(19),
	   .COUNT_WIDTH(5)
	) db_data(
		.clk(clk),
		.I(kdata),
		.O(kdataf)
	);
		
	always@(negedge(kclkf))begin
		case(cnt)
		0:;//Start bit
		1:datacur[0]<=kdataf;
		2:datacur[1]<=kdataf;
		3:datacur[2]<=kdataf;
		4:datacur[3]<=kdataf;
		5:datacur[4]<=kdataf;
		6:datacur[5]<=kdataf;
		7:datacur[6]<=kdataf;
		8:datacur[7]<=kdataf;
		9:flag<=1'b1;
		10:flag<=1'b0;
		
		endcase
			if(cnt<=9) cnt<=cnt+1;
			else if(cnt==10) cnt<=0;
	end
	
	reg pflag;
	always@(posedge clk) begin
		if (flag == 1'b1 && pflag == 1'b0) begin
			keycode <= {dataprev, datacur};
			oflag <= 1'b1;
			dataprev <= datacur;
		end else
			oflag <= 'b0;
		pflag <= flag;
	end
	
endmodule
	`

	result += "\n`timescale 1ns / 1ps\n"
	result += `
	module bondkeydrv(
		input       clk,
		input       PS2Data,
		input       PS2Clk,
		output reg  [15:0] key,
		output reg  keycode_valid,
		input       keycode_recv
	);
	
		reg         start=0;
		reg         CLK50MHZ=0;
		wire [15:0] keycode;
		reg [15:0] keycodev;
		wire        flag;
		reg         cn=0;
	
		reg [5:0]	mods; // 0 Shift - 1 Ctrl - 3 Alt

		always @(posedge(clk))begin
			CLK50MHZ<=~CLK50MHZ;
		end
		
		PS2Receiver uut (
			.clk(CLK50MHZ),
			.kclk(PS2Clk),
			.kdata(PS2Data),
			.keycode(keycode),
			.oflag(flag)
		);
		
		always@(keycode)
			if (keycode[7:0] == 8'hf0) begin
				cn <= 1'b0;
			end else if (keycode[15:8] == 8'hf0) begin
				cn <= keycode != keycodev;
			end else begin
				cn <= keycode[7:0] != keycodev[7:0] || keycodev[15:8] == 8'hf0;
			end
		
		always @(posedge clk)
			if (flag == 1'b1 && cn == 1'b1) begin
				keycodev <= keycode;
				case (keycode[15:8])
					8'hf0: // Break Codes
						case (keycode[7:0])
							// Modificators
							8'h12 : begin // Left Shift
									mods[0] <= 1'b0;
									start <= 1'b0;
								end
							8'h59 : begin // Right Shift
									mods[0] <= 1'b0;
									start <= 1'b0;
								end
							default:
								start <= 1'b0;
						endcase
					default: // Make Codes
						case (keycode[7:0])
							// Modificators
							8'h12 : begin // Left Shift
									mods[0] <= 1'b1;
									start <= 1'b0;
								end
							8'h59 : begin // Right Shift
									mods[0] <= 1'b1;
									start <= 1'b0;
								end
							// Special
							8'h29 : begin
									key <= 16'h0020;  // Space
									start <= 1'b1;
								end
							8'h5a : begin
									key <= 16'h000d;  // Enter
									start <= 1'b1;
								end
							8'h76 : begin
									key <= 16'h001b;  // Escape
									start <= 1'b1;
								end
							8'h66 : begin
									key <= 16'h0008;  // Backspace
									start <= 1'b1;
								end
							8'h0d : begin
									key <= 16'h000b;  // Tab
									start <= 1'b1;
								end
							8'h52 : begin
									if (mods[0])
										key <= 16'h0022;  // "
									else
										key <= 16'h0027;  // '
									start <= 1'b1;
								end
							8'h4c : begin
									if (mods[0])
										key <= 16'h003a;  // :
									else
										key <= 16'h003b;  // ;
									start <= 1'b1;
								end
							8'h54 : begin
									if (mods[0])
										key <= 16'h007b;  // {
									else
										key <= 16'h005b;  // [
									start <= 1'b1;
								end
							8'h5b : begin
									if (mods[0])
										key <= 16'h007d;  // }
									else
										key <= 16'h005d;  // ]
									start <= 1'b1;
								end
							8'h0e : begin
									if (mods[0])
										key <= 16'h007e;  // ~
									else
										key <= 16'h0060;  // reverse '
									start <= 1'b1;
								end
							8'h5d : begin
									if (mods[0])
										key <= 16'h007c;  // |
									else
										key <= 16'h005c;  // \
									start <= 1'b1;
								end
							8'h4e : begin
									if (mods[0])
										key <= 16'h005f;  // _
									else
										key <= 16'h002d;  // -
									start <= 1'b1;
								end
							8'h55 : begin
									if (mods[0])
										key <= 16'h002b;  // +
									else
										key <= 16'h003d;  // =
									start <= 1'b1;
								end
							8'h4a : begin
									if (mods[0])
										key <= 16'h003f;  // ?
									else
										key <= 16'h002f;  // /
									start <= 1'b1;
								end
							8'h41 : begin
									if (mods[0])
										key <= 16'h003c;  // <
									else
										key <= 16'h002c;  // ,
									start <= 1'b1;
								end
							8'h49 : begin
									if (mods[0])
										key <= 16'h003e;  // >
									else
										key <= 16'h002e;  // .
									start <= 1'b1;
								end
							// Numbers 0-9
							8'h45 : begin
									if (mods[0])
										key <= 16'h0029;  // )
									else
										key <= 16'h0030;  // 0
									start <= 1'b1;
								end
							8'h16 : begin
									if (mods[0])
										key <= 16'h0021;  // !
									else
										key <= 16'h0031;  // 1
									start <= 1'b1;
								end
							8'h1e : begin
									if (mods[0])
										key <= 16'h0040;  // @
									else
										key <= 16'h0032;  // 2
									start <= 1'b1;
								end
							8'h26 : begin
									if (mods[0])
										key <= 16'h0023;  // #
									else
										key <= 16'h0033;  // 3
									start <= 1'b1;
								end
							8'h25 : begin
									if (mods[0])
										key <= 16'h0024;  // $
									else
										key <= 16'h0034;  // 4
									start <= 1'b1;
								end
							8'h2e : begin
									if (mods[0])
										key <= 16'h0025;  // %
									else
										key <= 16'h0035;  // 5
									start <= 1'b1;
								end
							8'h36 : begin
									if (mods[0])
										key <= 16'h005e;  // ^
									else
										key <= 16'h0036;  // 6
									start <= 1'b1;
								end
							8'h3d : begin
									if (mods[0])
										key <= 16'h0026;  // &
									else
										key <= 16'h0037;  // 7
									start <= 1'b1;
								end
							8'h3e : begin
									if (mods[0])
										key <= 16'h002a;  // *
									else
										key <= 16'h0038;  // 8
									start <= 1'b1;
								end
							8'h46 : begin
									if (mods[0])
										key <= 16'h0028;  // (
									else
										key <= 16'h0039;  // 9
									start <= 1'b1;
								end
							// Characters
							8'h15 : begin
									if (mods[0])
										key <= 16'h0051;  // Q
									else
										key <= 16'h0071;  // q 
									start <= 1'b1;
								end
							8'h1d : begin
									if (mods[0])
										key <= 16'h0057;  // W
									else
										key <= 16'h0077;  // w
									start <= 1'b1;
								end
							8'h24 : begin
									if (mods[0])
										key <= 16'h0045;  // E
									else
										key <= 16'h0065;  // e
									start <= 1'b1;
								end
							8'h2d : begin
									if (mods[0])
										key <= 16'h0052;  // R
									else
										key <= 16'h0072;  // r
									start <= 1'b1;
								end
							8'h2c : begin
									if (mods[0])
										key <= 16'h0054;  // T
									else
										key <= 16'h0074;  // t
									start <= 1'b1;
								end
							8'h35 : begin
									if (mods[0])
										key <= 16'h0059;  // Y
									else
										key <= 16'h0079;  // y
									start <= 1'b1;
								end
							8'h3c : begin
									if (mods[0])
										key <= 16'h0055;  // U
									else
										key <= 16'h0075;  // u
									start <= 1'b1;
								end
							8'h43 : begin
									if (mods[0])
										key <= 16'h0049;  // I
									else
										key <= 16'h0069;  // i
									start <= 1'b1;
								end
							8'h44 : begin
									if (mods[0])
										key <= 16'h004f;  // O
									else
										key <= 16'h006f;  // o
									start <= 1'b1;
								end
							8'h4d : begin
									if (mods[0])
										key <= 16'h0050;  // P
									else
										key <= 16'h0070;  // p
									start <= 1'b1;
								end
							8'h1c : begin
									if (mods[0])
										key <= 16'h0041;  // A
									else
										key <= 16'h0061;  // a
									start <= 1'b1;
								end
							8'h1b : begin
									if (mods[0])
										key <= 16'h0053;  // S
									else
										key <= 16'h0073;  // s
									start <= 1'b1;
								end
							8'h23 : begin
									if (mods[0])
										key <= 16'h0044;  // D
									else
										key <= 16'h0064;  // d
									start <= 1'b1;
								end
							8'h2b : begin
									if (mods[0])
										key <= 16'h0046;  // F
									else
										key <= 16'h0066;  // f
									start <= 1'b1;
								end
							8'h34 : begin
									if (mods[0])
										key <= 16'h0047;  // G
									else
										key <= 16'h0067;  // g
									start <= 1'b1;
								end
							8'h33 : begin
									if (mods[0])
										key <= 16'h0048;  // H
									else
										key <= 16'h0068;  // h 
									start <= 1'b1;
								end
							8'h3b : begin
									if (mods[0])
										key <= 16'h004a;  // J
									else
										key <= 16'h006a;  // j
									start <= 1'b1;
								end
							8'h42 : begin
									if (mods[0])
										key <= 16'h004b;  // K
									else
										key <= 16'h006b;  // k
									start <= 1'b1;
								end
							8'h4b : begin
									if (mods[0])
										key <= 16'h004c;  // L
									else
										key <= 16'h006c;  // l
									start <= 1'b1;
								end
							8'h1a : begin
									if (mods[0])
										key <= 16'h005a;  // Z
									else
										key <= 16'h007a;  // z
									start <= 1'b1;
								end
							8'h22 : begin
									if (mods[0])
										key <= 16'h0058;  // X
									else
										key <= 16'h0078;  // x
									start <= 1'b1;
								end
							8'h21 : begin
									if (mods[0])
										key <= 16'h0043;  // C
									else
										key <= 16'h0063;  // c
									start <= 1'b1;
								end
							8'h2a : begin
									if (mods[0])
										key <= 16'h0056;  // V
									else
										key <= 16'h0076;  // v
									start <= 1'b1;
								end
							8'h32 : begin
									if (mods[0])
										key <= 16'h0042;  // B
									else
										key <= 16'h0062;  // b
									start <= 1'b1;
								end
							8'h31 : begin
									if (mods[0])
										key <= 16'h004e;  // N
									else
										key <= 16'h006e;  // n
									start <= 1'b1;
								end
							8'h3a : begin
									if (mods[0])
										key <= 16'h004d;  // M
									else
										key <= 16'h006d;  // m
									start <= 1'b1;
								end
							default:
								start <= 1'b0;
						endcase
				endcase
			end else
				start <= 1'b0;
	 
		always @(posedge clk) begin
			if (start == 1'b1) begin
				keycode_valid <= 1'b1;
			end
			else begin 
				if (keycode_recv == 1'b1) begin
					keycode_valid <= 1'b0;
				end
			end
		end
		
	endmodule
	`
	return result
}

func (sl *Ps2KeyboardExtra) ExtraFiles() ([]string, []string) {
	return []string{}, []string{}
}
