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

	deserializer = `
module {{ .ModuleName }}(
	   input clk,
	   input reset,
	   input impulse,
	   input [31:0] data,
	   output reg ready,
	   
	   output wire [31:0] i0,
	   output i0_valid,
	   input i0_recv,
	   
	   output wire [31:0] i1,
	   output i1_valid,
	   input i1_recv,
	   
	   output wire [31:0] i2,
	   output i2_valid,
	   input i2_recv,
	   
	   output wire [31:0] i3,
	   output i3_valid,
	   input i3_recv	   	   
	);
	
	reg [2:0] input_index;
    reg [0:0] SM;
    
    reg [31:0] localdata;

    reg [3:0] valids;

    wire [3:0] recvs;
    
    localparam SMIDLE=1'b0,
                SMBM=1'b1;
	
	always @( posedge clk) begin
	   if (reset) begin
	       ready <= 1'b0;
	       input_index <= 3'b000;
	       SM<=SMIDLE;
	       localdata[31:0] <= 32'b0;
	       valids[3:0] <= 4'b0000;
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

//     empty,
//     full
// );
//     input clk;
//     input reset;
//     output empty;
//     output full;
//     {{- if .Senders }}
//     {{- range .Senders }}
//     input [{{ dec $.DataSize }}:0] {{ . }}Data;
//     input {{ . }}Write;
//     output reg {{ . }}Ack;
//     {{- end }}
//     {{- end }}
//     {{- if .Receivers }}
//     {{- range .Receivers }}
//     output reg [{{ dec $.DataSize }}:0] {{ . }}Data;
//     input {{ . }}Read;
//     output reg {{ . }}Ack;
//     {{- end }}
//     {{- end }}

//     reg [{{ dec $.DataSize }}:0] memory[{{ dec $.Depth }}:0];
//     reg [{{ dec (bits (inc $.Depth)) }}:0] sp;
//     {{- if eq $.MemType "FIFO" }}
//     reg [{{ dec (bits (inc $.Depth)) }}:0] readsp;
//     reg [{{ dec (bits (inc $.Depth)) }}:0] writesp;
//     {{- end }}

//     assign empty = (sp==0)? 1'b1:1'b0;
//     assign full = (sp=={{ $.Depth }})? 1'b1:1'b0;

//     wire readneed;
//     wire writeneed;

//     assign writeneed = ( 1'b0
//     {{- if .Senders }}
//     {{- range .Senders }}
//             | {{ . }}Write
//     {{- end }}
//     {{- end }} );

//     assign readneed = ( 1'b0
//     {{- if .Receivers }}
//     {{- range .Receivers }}
//             | {{ . }}Read
//     {{- end }}
//     {{- end }} );

//     reg [{{ dec ( bits (len .Senders) ) }}:0] sendSM;
//     //{{- range $key, $value := .Senders }}
//     //localparam sendSM{{ . }} = {{ bits (len $.Senders) }}'d{{ $key }};
//     //{{- end }}

//     reg [{{ dec ( bits (len .Receivers) ) }}:0] recvSM;
//     //{{- range $key, $value := .Receivers }}
//     //localparam recvSM{{ . }} = {{ bits (len $.Receivers) }}'d{{ $key }};
//     //{{- end }}

//     integer i;

//     always @(posedge clk) begin
//         if (reset) begin
//             sp <= {{ bits (inc $.Depth) }}'d0;
//             {{- if eq $.MemType "FIFO" }}
//             readsp <= {{ bits (inc $.Depth) }}'d0;
//             writesp <= {{ bits (inc $.Depth) }}'d0;
//             {{- end }}
//             {{- range $key, $value := .Receivers }}
//             {{ $value }}Data <= {{ $.DataSize}}'d0;
//             {{ $value }}Ack <= 1'b0;
//             {{- end }}
//             {{- range $key, $value := .Senders }}
//             {{ $value }}Ack <= 1'b0;
//             {{- end }}
//             sendSM <= {{bits (len .Receivers)}}'d0;
//             recvSM <= {{bits (len .Receivers)}}'d0;
//             for (i=0;i<{{ .Depth }};i=i+1) begin
//                 memory[i]<={{ $.DataSize}}'d0;
//             end
//         end
//         else begin
//             // Read state machine part
//             if (readneed && !empty) begin
//                 case (recvSM)
//                 {{- range $key, $value := .Receivers }}
//                 {{ bits (len $.Receivers) }}'d{{ $key }}: begin
//                     if ({{ $value }}Read && !{{ $value }}Ack) begin
//                         {{- if eq $.MemType "LIFO" }}
//                         {{ $value }}Data[{{ dec $.DataSize }}:0] <= memory[sp-1];
//                         sp <= sp - 1;
//                         {{- end }}
//                         {{- if eq $.MemType "FIFO" }}
//                         {{ $value }}Data[{{ dec $.DataSize }}:0] <= memory[readsp];
//                         if (readsp=={{ dec $.Depth }}) begin
//                             readsp <= 0;
//                             sp <=  writesp;
//                         end
//                         else begin
//                             readsp <= readsp + 1;
//                             if (writesp < readsp + 1) begin
//                                 sp <= {{ $.Depth }} - readsp -1 + writesp;
//                             end
//                             else begin
//                                 sp <= writesp - readsp - 1;
//                             end
//                         end
//                         {{- end }}
//                     end
//                     recvSM <= {{ bits (len $.Receivers) }}'d{{ next $key (len $.Receivers) }};
//                 end
//                 {{- end }}
//                 endcase
//             end
//             // Write state machine part
//             else if (writeneed && !full) begin
//                 case (sendSM)
//                 {{- range $key, $value := .Senders }}
//                 {{ bits (len $.Senders) }}'d{{ $key }}: begin
//                     if ({{ $value }}Write && !{{ $value }}Ack) begin
//                         {{- if eq $.MemType "LIFO" }}
//                         memory[sp] <= {{ $value }}Data[{{ dec $.DataSize }}:0];
//                         sp <= sp + 1;
//                         {{- end }}
//                         {{- if eq $.MemType "FIFO" }}
//                         memory[writesp] <= {{ $value }}Data[{{ dec $.DataSize }}:0];
//                         if (writesp=={{ dec $.Depth }}) begin
//                             writesp <= 0;
//                             sp <= {{ $.Depth }} - readsp;
//                         end
//                         else begin
//                             writesp <= writesp + 1;
//                             if (writesp + 1 > readsp) begin
//                                 sp <= writesp - readsp + 1;
//                             end
//                             else begin
//                                 sp <= {{ $.Depth }} - readsp + writesp + 1;
//                             end
//                         end
//                         {{- end }}
//                     end
//                     sendSM <= {{ bits (len $.Senders) }}'d{{ next $key (len $.Senders) }};
//                 end
//                 {{- end }}
//                 endcase
//             end

//             // Read ack process
//             {{- range $key, $value := .Receivers }}
//             if ({{ $value }}Read && !{{ $value }}Ack && recvSM=={{ bits (len $.Senders) }}'d{{ $key }} && !empty) begin
//                 {{ $value }}Ack <= 1'b1;
//             end
//             else begin
//                 if (!{{ $value }}Read) begin
//                     {{ $value }}Ack <= 1'b0;
//                 end
//             end
//             {{- end }}

//             // Write ack process
//             {{- range $key, $value := .Senders }}
//             if ({{ $value }}Write && !{{ $value }}Ack && sendSM=={{ bits (len $.Receivers) }}'d{{ $key }} && !full) begin
//                 {{ $value }}Ack <= 1'b1;
//             end
//             else begin
//                 if (!{{ $value }}Write) begin
//                     {{ $value }}Ack <= 1'b0;
//                 end
//             end
//             {{- end }}
//         end
//     end
// endmodule
// `

// {{- $smindex:= 0 }}
// {{- if .Inputs }}
// {{- range .Inputs }}
// assign {{ . }} = slv_reg{{ $smindex }}[31:0];
// {{- $smindex = inc $smindex }}
// {{- end }}
// {{- end }}
// assign DVDR_PS2PL = slv_reg{{ $smindex }}[31:0];
// {{- $smindex = inc $smindex }}
// assign states = slv_reg{{ $smindex }}[31:0];
// {{- $smindex = inc $smindex }}

// always @( posedge S_AXI_ACLK )
// begin
// {{- if .Outputs }}
// {{- range .Outputs }}
//     slv_reg{{ $smindex }} <= {{ . }}[31:0];
//     {{- $smindex = inc $smindex }}
// {{- end }}
// {{- end }}
//     slv_reg{{ $smindex }} <= DVDR_PL2PS[31:0];
// {{- $smindex = inc $smindex }}
//     slv_reg{{ $smindex }} <= changes[31:0];
// end

)
