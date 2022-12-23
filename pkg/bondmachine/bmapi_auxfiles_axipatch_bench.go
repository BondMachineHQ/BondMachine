package bondmachine

const (
	auxfilesAXIPatchBench = `
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

    reg [31:0] i0state;
    reg [31:0] o0state;
    reg ch_i0;
    reg ch_00;
    reg [31:0] measure;
    reg [31:0] measure_final;

    initial begin
        i0state = 0;
        o0state = 0;
        ch_i0 = 0;
        ch_00 = 0;
        measure = 0;
        measure_final = 0;
    end

    always @( posedge S_AXI_ACLK) begin
        i0state <= port_i0;
        o0state <= port_o0;
    end

    always @( posedge S_AXI_ACLK) begin
        if (i0state != port_i0) begin
            ch_i0 <= 1;
        end else begin
            ch_i0 <= 0;
        end
    end

    always @( posedge S_AXI_ACLK) begin
        if (o0state != port_o0) begin
            ch_00 <= 1;
        end else begin
            ch_00 <= 0;
        end
    end

    always @( posedge S_AXI_ACLK) begin
        if (ch_i0) begin
//            measure <= 0;
//        end else if (ch_00) begin
            measure_final[31:0] <= measure[31:0];
        end else begin
            measure <= measure + 1;
        end
    end

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
        slv_reg4 <= port_o0[31:0];
        slv_reg5 <= measure_final[31:0];
        slv_reg6 <= DVDR_PL2PS[31:0];
        slv_reg7 <= changes[31:0];
    end  

    
`
)
