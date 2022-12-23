package bondmachine

const (
	auxfilesAXIPatch = `
    wire [31:0] states;
    wire [31:0] changes;
    wire [31:0] DVDR_PS2PL;
    wire [31:0] DVDR_PL2PS;

    {{- if .Outputs }}
    {{- range .Outputs }}
    wire [31:0] {{ . }};
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    wire [31:0] {{ . }};
    {{- end }}
    {{- end }}

    bondmachine_main bondmachine_inst(
        .clk(S_AXI_ACLK),
        .btnC(btnC),
    //    .led(led),
        .A_DVDR_PS2PL(DVDR_PS2PL),
        .A_DVDR_PL2PS(DVDR_PL2PS),
        .A_changes(changes),
        .A_states(states),
	{{- if .Outputs }}
	{{- range .Outputs }}
	.A_{{ . }}({{ . }}),
	{{- end }}
	{{- end }}
	{{- if .Inputs }}
	{{- range .Inputs }}
	.A_{{ . }}({{ . }}),
	{{- end }}
	{{- end }}
        .interrupt(interrupt)
    );

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
