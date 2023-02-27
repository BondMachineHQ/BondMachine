package bmserialize

const (
	deserializer = `
module {{ .ModuleName }}(
	input clk,
	input reset,
	{{- range $y := intRange 0 .Terminals }}
	output wire [{{ dec $.TerminalDataSize }}:0] i{{ $y }},
	output i{{ $y}}_valid,
	input i{{ $y }}_recv,
	{{- end }}  
	input impulse,
	input [{{ dec .SerialdataSize }}:0] data,
	output reg ready
	);

reg [{{ bits .Terminals }}:0] input_index;
reg [1:0] SM;
    
reg [{{ dec .SerialDataSize }}:0] localdata;

reg [{{ bits .Terminals }}:0] valids;
wire [{{ bits .Terminals }}:0] recvs;
    
localparam	SMIDLE=2'b0,
        	SMBM=2'b1;
	
always @( posedge clk) begin
	if (reset) begin
		ready <= 1'b0;
		input_index <= {{ inc (bits .Terminals) }}'d0;
		SM<=SMIDLE;
		localdata[ {{ dec .SerialDataSize }}:0] <= {{ .SerialDataSize }}'d0;
		valids[{{ dec (bits Terminals) }}:0] <= {{ bits Terminals }}'d0;
	end 
	else begin
		case (SM)
		SMIDLE: begin
	        	if (impulse) begin
				ready <= 1'b0;
				localdata[31:0] <= data[31:0];
				SM<=SMBM;
			end
			else begin
	                	ready <= 1'b1;
	        	end
	        end
		SMBM: begin
	           if (recvs[input_index] == 1'b0) begin
	               valids[input_index] <= 1'b1;
	               ready <= 1'b0;
	           end
	           else begin
	               valids[input_index] <= 1'b0;
	               if (input_index + 1 == 4'd4) begin
	                   input_index <= 0;
	               end
	               else begin
	                   input_index <= input_index + 1;
	               end
	               SM<=SMIDLE;
	               ready <= 1'b1;
	           end
		end
		endcase
	end
end
	
assign i0[31:0] = localdata[31:0];
assign i1[31:0] = localdata[31:0];
assign i2[31:0] = localdata[31:0];
assign i3[31:0] = localdata[31:0];
	
assign i0_valid=valids[0];
assign i1_valid=valids[1];
assign i2_valid=valids[2];
assign i3_valid=valids[3];
	
assign recvs[3:0] = {i3_recv, i2_recv, i1_recv, i0_recv};
	
endmodule
`
)
