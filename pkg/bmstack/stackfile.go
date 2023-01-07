package bmstack

const (
	stack = `
module {{ .ModuleName }}(clk,
    {{- if .Senders }}
    {{- range .Senders }}
    {{ . }}Data,
    {{ . }}Write,
    {{ . }}Ack,
    {{- end }}
    {{- end }}
    {{- if .Receivers }}
    {{- range .Receivers }}
    {{ . }}Data,
    {{ . }}Read,
    {{ . }}Ack,
    {{- end }}
    {{- end }}
    reset,
    empty,
    full
);
    input clk;
    input reset;
    output empty;
    output full;
    {{- if .Senders }}
    {{- range .Senders }}
    input [{{ dec $.DataSize }}:0] {{ . }}Data;
    input {{ . }}Write;
    output {{ . }}Ack;
    {{- end }}
    {{- end }}
    {{- if .Receivers }}
    {{- range .Receivers }}
    output reg [{{ dec $.DataSize }}:0] {{ . }}Data;
    input {{ . }}Read;
    output {{ . }}Ack;
    {{- end }}
    {{- end }}
endmodule



    {{- $smindex:= 0 }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    assign {{ . }} = slv_reg{{ $smindex }}[31:0];
    {{- $smindex = inc $smindex }}
    {{- end }}
    {{- end }}
    assign DVDR_PS2PL = slv_reg{{ $smindex }}[31:0];
    {{- $smindex = inc $smindex }}
    assign states = slv_reg{{ $smindex }}[31:0];
    {{- $smindex = inc $smindex }}

    always @( posedge S_AXI_ACLK )
    begin
    {{- if .Outputs }}
    {{- range .Outputs }}
        slv_reg{{ $smindex }} <= {{ . }}[31:0];
        {{- $smindex = inc $smindex }}
    {{- end }}
    {{- end }}   
        slv_reg{{ $smindex }} <= DVDR_PL2PS[31:0];
    {{- $smindex = inc $smindex }}
        slv_reg{{ $smindex }} <= changes[31:0];
    end  

    
`
)
