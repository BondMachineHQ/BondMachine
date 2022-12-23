package bondmachine

const (
	bmapiaximmtransceiver = "\n`timescale 1ns / 1ps\n" + `

module bmapiaximmtransceiver(
    input clk,
    input reset,
    input wire [31:0] A_DVDR_PS2PL,
    output reg [31:0] A_DVDR_PL2PS,
    output reg [31:0] A_changes,
    input wire [31:0] A_states,
    {{- if .Outputs }}
    {{- range .Outputs }}
    output wire [31:0] A_{{ . }},
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    input [31:0] A_{{ . }},
    {{- end }}
    {{- end }}
    {{- if .Outputs }}
    {{- range .Outputs }}
    input {{ $.Buswidth }} {{ . }},
    input {{ . }}_valid,
    output wire {{ . }}_recv,
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    output wire {{ $.Buswidth }} {{ . }},
    output wire {{ . }}_valid,
    input {{ . }}_recv,
    {{- end }}
    {{- end }}
    output interrupt
);

// Reg local to the transceiver
    {{- if .Outputs }}
    {{- range .Outputs }}
    reg [31:0] A_{{ . }}_local;
    reg A_{{ . }}_valid_local;
    reg A_{{ . }}_recv_local;
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    reg [31:0] A_{{ . }}_local;
    reg A_{{ . }}_valid_local;
    reg A_{{ . }}_recv_local;
    {{- end }}
    {{- end }}

    reg interrupt_local;
    wire interrupt_condition;
    wire interrupt_verified;

    wire kernel_exec;
    wire kernel_done;

    reg [1:0] INTR_SM;
    localparam INTR_IDLE=2'b00,
            INTR_SENT=2'b01,
            INTR_VER=2'b10,
            INTR_RESET=2'b11;

// Local replicas init
initial begin
    {{- if .Outputs }}
    {{- range  .Outputs }}
    A_{{ . }}_local{{ $.Buswidth }} = {{ $.Rsize }}'d0;
    A_{{ . }}_valid_local=1'b0;
    A_{{ . }}_recv_local=1'b0;
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range $key, $value := .Inputs }}
    A_{{ . }}_local{{ $.Buswidth }} = {{ $.Rsize }}'d0;
    A_{{ . }}_valid_local=1'b0;
    A_{{ . }}_recv_local=1'b0;
    {{- end }}
    {{- end }}
end

assign kernel_exec = A_states[0];
assign kernel_done = A_states[1];

// Store the signals from the BM

{{- $chindex:= 31 }}
{{- if .Inputs }}
{{- range .Inputs }}

reg {{ . }}_recv_need;

always @(posedge clk) begin
    if (INTR_SM[1:0] == INTR_RESET) begin
            {{ . }}_recv_need <= 1'b0;
            A_changes[{{ $chindex }}] <= 1'b0;
    end
    else begin
        if (INTR_SM[1:0] == INTR_IDLE) begin
            A_{{ . }}_recv_local <= {{ . }}_recv;
            if (A_{{ . }}_recv_local != {{ . }}_recv) begin
                {{ . }}_recv_need <= 1'b1;
                A_changes[{{ $chindex }}] <= 1'b1;
            end
        end
    end
end

{{- $chindex = dec $chindex }}

{{- end }}
{{- end }}

{{- if .Outputs }}
{{- range .Outputs }}

reg {{ . }}_need;

always @(posedge clk) begin
    if (INTR_SM[1:0] == INTR_RESET) begin
            {{ . }}_need <= 1'b0;
            A_changes[{{ $chindex }}] <= 1'b0;
    end
    else begin
        if (INTR_SM[1:0] == INTR_IDLE) begin
            A_{{ . }}_local{{ $.Buswidth }} <= {{ . }}{{ $.Buswidth }};
            if (A_{{ . }}_local{{ $.Buswidth }} != {{ . }}{{ $.Buswidth }}) begin
                {{ . }}_need <= 1'b1;
                A_changes[{{ $chindex }}] <= 1'b1;
            end
        end
    end
end

{{- $chindex = dec $chindex }}

reg {{ . }}_valid_need;

always @(posedge clk) begin
    if (INTR_SM[1:0] == INTR_RESET) begin
            {{ . }}_valid_need <= 1'b0;
            A_changes[{{ $chindex }}] <= 1'b0;
    end
    else begin
        if (INTR_SM[1:0] == INTR_IDLE) begin
            A_{{ . }}_valid_local <= {{ . }}_valid;
            if (A_{{ . }}_valid_local != {{ . }}_valid) begin
                {{ . }}_valid_need <= 1'b1;
                A_changes[{{ $chindex }}] <= 1'b1;
            end
        end
    end
end

{{- $chindex = dec $chindex }}

{{- end }}
{{- end }}



assign interrupt_condition = ((~ kernel_exec ) & ( 1'b0
{{- if .Inputs }}
{{- range .Inputs }}
        | ( {{ . }}_recv_need )
{{- end }}
{{- end }}
{{- if .Outputs }}
{{- range .Outputs }}
        | ( {{ . }}_need )
        | ( {{ . }}_valid_need )
{{- end }}
{{- end }}));

reg [31:0] interrupt_delay;

assign interrupt_verified = (kernel_exec & kernel_done);
assign interrupt_received = (~ kernel_exec);

// Interrupt state machine
always @(posedge clk) begin
    case (INTR_SM[1:0])
    INTR_IDLE: begin
        if (interrupt_condition) begin
            interrupt_delay <= 32'b00000001_00000000_00000000_00000000;
            interrupt_local <= 1'b1;
            INTR_SM[1:0] <= INTR_SENT;
        end
    end
    INTR_SENT: begin
        interrupt_local <= 1'b0;
        interrupt_delay <= interrupt_delay - 1;
        if (interrupt_verified) begin
            INTR_SM[1:0] <= INTR_VER;
        end
        else begin
            if (interrupt_delay == 32'd0) begin
                INTR_SM[1:0] <= INTR_IDLE;
            end
        end
    end
    INTR_VER: begin
        interrupt_delay <= interrupt_delay - 1;
        if (interrupt_received) begin
            INTR_SM[1:0] <= INTR_RESET;
        end
        else begin
            if (interrupt_delay == 32'd0) begin
                INTR_SM[1:0] <= INTR_IDLE;
            end
        end
    end   
    INTR_RESET: begin
        INTR_SM[1:0] <= INTR_IDLE;
        interrupt_delay <= 32'd0;
    end
    endcase
end

assign interrupt = interrupt_local;

// Assign of the signals meant to reach the PS
assign A_DVDR_PL2PS = {
{{- $windex:= 32 }}
{{- if .Inputs }}
{{- range .Inputs }}
        A_{{ . }}_recv_local, 
        {{- $windex = dec $windex }}
{{- end }}
{{- end }}
{{- if .Outputs }}
{{- range .Outputs }}
        A_{{ . }}_valid_local,
        {{- $windex = dec $windex }}
{{- end }}
{{- end }}
        {{ $windex }}'b0 };

{{- if .Outputs }}
{{- range .Outputs }}
assign A_{{ . }}[31:0] = A_{{ . }}_local[31:0];
{{- end }}
{{- end }}

// Store the signals from the PS
always @(posedge clk) begin
{{- $windex:= 31 }}
{{- if .Inputs }}
{{- range .Inputs }}
        A_{{ . }}_valid_local <= A_DVDR_PS2PL[{{ $windex }}];
        A_{{ . }}_local[31:0] <= A_{{ . }}[31:0];
        {{- $windex = dec $windex }}
{{- end }}
{{- end }}
{{- if .Outputs }}
{{- range .Outputs }}
        A_{{ . }}_recv_local <= A_DVDR_PS2PL[{{ $windex }}];
        {{- $windex = dec $windex }}
{{- end }}
{{- end }}
end

// Assign of outgoing signals (to BM)
{{- if .Inputs }}
{{- range $key, $value := .Inputs }}
assign {{ . }}{{ $.Buswidth }} = A_{{ . }}_local{{ $.Buswidth }};
assign {{ . }}_valid = A_{{ . }}_valid_local;
{{- end }}
{{- end }}
{{- if .Outputs }}
{{- range .Outputs }}
assign {{ . }}_recv = A_{{ . }}_recv_local;
{{- end }}
{{- end }}

endmodule
`
)
