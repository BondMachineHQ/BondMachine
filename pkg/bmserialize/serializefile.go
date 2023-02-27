package bmserialize

const (
	serializer = `
module {{ .ModuleName }}(
	input clk,
	input reset,
	{{- range $y := intRange 0 .Terminals }}
	input [{{ dec $.TerminalDataSize }}:0] o{{ $y }},
	input o{{ $y }}_valid,
	output o{{ $y }}_recv,
	{{- end }}
	input ack,
	output [{{ dec .SerialDataSize }}:0] data,
	output reg ready,
	);
	
reg [{{ bits .Terminals }}:0] output_index;
reg [1:0] SM;
    
reg [{{ dec .SerialDataSize }}:0] localdata;

wire [{{ bits .Terminals }}:0] valids;
reg [{{ bits .Terminals }}:0] recvs;
    
wire [{{ dec .SerialDataSize }}:0] outputs[{{ bits .Terminals }}:0];
    
localparam SMIDLE=2'b00,
        SMRES=2'b01,
        SMBM=2'b10;
	
always @( posedge clk) begin
	if (reset) begin
		ready <= 1'b0;
		output_index <= {{ inc (bits .Terminals) }}'d0;
		SM<=SMIDLE;
		recvs[{{ bits .Terminals }}:0] <= {{ inc (bits .Terminals) }}'d0;
	end 
	else begin
		case (SM)
		SMIDLE: begin
			if (valids[output_index]) begin
				ready <= 1'b1;
				localdata[{{ dec .SerialDataSize }}:0] <= outputs[output_index][{{ dec .SerialDataSize }}:0];
				SM<=SMRES;
			end
			else begin
				ready <= 1'b0;
			end
		end
		SMRES: begin
	        	if (ack) begin
	        		ready <= 1'b0;
	        		SM<=SMBM;
	        	end   
	       	end
		SMBM: begin
			if (!valids[output_index]) begin
				if (output_index + 1 == {{ inc (bits .Terminals) }}'d{{ .Terminals }}) begin
					output_index <= 0;
				end
				else begin
					output_index <= output_index + 1;
				end
				recvs[output_index] <= 1'b0;
				SM<=SMIDLE;           
			end
			else begin
				recvs[output_index] <= 1'b1;
	        	end
		end
		endcase
	end
end
	
assign data[{{ dec .SerialDataSize }}:0] = localdata[{{ dec .SerialDataSize }}:0];
{{- range $y := intRange 0 .Terminals }}
assign outputs[{{ $y }}]=o{{ $y }}[{{ dec $.SerialDataSize }}:0];
assign o{{ $y }}_recv=recvs[{{ $y }}];
assign valids[{{ $y }}] = o{{ $y }}_valid;
{{- end }}
	
endmodule 
`
)
