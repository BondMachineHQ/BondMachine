package bondmachine

const (
	uartSO = "\n`timescale 1ns / 1ps" + `
module {{.ModuleName}}(
    input clk,
    input reset,
    {{- range .CpPorts }}
    {{- if eq .Dir "send" }}
    input [7:0] {{ .Name }}Data,
    input {{ .Name }}Write,
    output {{ .Name }}Ack,
    {{- else }}
    output [7:0] {{ .Name }}Data,
    output {{ .Name }}Read,
    input {{ .Name }}Ack,
    {{- end }}
    {{- end }}
    input {{.ModuleName}}_rx,
    output {{.ModuleName}}_tx,
    output rempty,
    output rfull,
    output wempty,
    output wfull
    );

    reg transmit = 1'b0;
    reg [7:0] tx_byte = 8'd0;

    wire is_receiving;
    wire is_transmitting;
    wire recv_error;
    wire [3:0] rx_samples;
    wire [3:0] rx_sample_countdown;

    wire received;
    wire [7:0] rx_byte;

{{- if .Senders }}

    // Read FIFO: filled by the UART receiver, read by the CPs (u2r)
    reg [7:0] uartreaderData = 8'd0;
    reg uartreaderWrite = 1'b0;
    wire uartreaderAck;

{{.ModuleName}}rfifo {{.ModuleName}}rfifo_inst(.clk(clk),
    .reset(reset),
    {{- range .Senders }}
    .{{ . }}Data({{ . }}Data),
    .{{ . }}Read({{ . }}Read),
    .{{ . }}Ack({{ . }}Ack),
    {{- end }}
    .uartreaderData(uartreaderData),
    .uartreaderWrite(uartreaderWrite),
    .uartreaderAck(uartreaderAck),
    .empty(rempty),
    .full(rfull)
);
{{- else }}

    assign rempty = 1'b1;
    assign rfull = 1'b0;
{{- end }}

{{- if .Receivers }}

    // Write FIFO: filled by the CPs (r2u), read by the UART transmitter
    wire [7:0] uartwriterData;
    reg uartwriterRead = 1'b0;
    wire uartwriterAck;

{{.ModuleName}}wfifo {{.ModuleName}}wfifo_inst(.clk(clk),
    .reset(reset),
    {{- range .Receivers }}
    .{{ . }}Data({{ . }}Data),
    .{{ . }}Write({{ . }}Write),
    .{{ . }}Ack({{ . }}Ack),
    {{- end }}
    .uartwriterData(uartwriterData),
    .uartwriterRead(uartwriterRead),
    .uartwriterAck(uartwriterAck),
    .empty(wempty),
    .full(wfull)
);
{{- else }}

    assign wempty = 1'b1;
    assign wfull = 1'b0;
{{- end }}

{{.ModuleName}}uart {{.ModuleName}}uart_inst(.clk(clk),
    .rst(reset),
    .rx({{.ModuleName}}_rx),
    .tx({{.ModuleName}}_tx),
    .transmit(transmit),
    .tx_byte(tx_byte),
    .received(received),
    .rx_byte(rx_byte),
    .is_receiving(is_receiving),
    .is_transmitting(is_transmitting),
    .recv_error(recv_error),
    .rx_samples(rx_samples),
    .rx_sample_countdown(rx_sample_countdown)
);

{{- if .Receivers }}

reg [1:0] outSM = 2'd0;

localparam [1:0]
    OUT_IDLE             = 2'd0,
    OUT_WAIT             = 2'd1,
    OUT_DONE             = 2'd2;

// Sending out to uart from the write FIFO
always @(posedge clk) begin
        if (reset) begin
            uartwriterRead <= 1'b0;
            transmit <= 1'b0;
            outSM <= OUT_IDLE;
        end
        else begin
            case (outSM)
            OUT_IDLE: begin
                // The FIFO may become empty as soon as the read handshake
                // starts, so a handshake in flight is always completed and
                // wempty only gates the start of a new one
                if (uartwriterAck && uartwriterRead) begin
                    uartwriterRead <= 1'b0;
                    tx_byte[7:0] <= uartwriterData[7:0];
                    transmit <= 1'b1;
                    outSM <= OUT_WAIT;
                end
                else if (!wempty && !uartwriterRead) begin
                    uartwriterRead <= 1'b1;
                    transmit <= 1'b0;
                end
            end
            OUT_WAIT: begin
                if (is_transmitting) begin
                    outSM <= OUT_DONE;
                    transmit <= 1'b0;
                end
            end
            OUT_DONE: begin
                if (!is_transmitting) begin
                    outSM <= OUT_IDLE;
                    transmit <= 1'b0;
                end
            end
            default: begin
                outSM <= OUT_IDLE;
            end
            endcase
        end
end
{{- end }}

{{- if .Senders }}

reg inSM = 1'b0;

localparam
    IN_IDLE             = 1'd0,
    IN_WAIT             = 1'd1;

// Reading the UART and pushing to the read FIFO
always @(posedge clk) begin
        if (reset) begin
            uartreaderWrite <= 1'b0;
            inSM <= IN_IDLE;
        end
        else begin
            case (inSM)
            IN_IDLE: begin
                if (received) begin
                    if (!uartreaderAck) begin
                        uartreaderData[7:0] <= rx_byte[7:0];
                        uartreaderWrite <= #1 1'b1;
                        inSM <= IN_WAIT;
                    end
                end
            end
            IN_WAIT: begin
                if (uartreaderAck) begin
                    uartreaderWrite <= #1 1'b0;
                    inSM <= IN_IDLE;
                end
            end
            endcase
        end
end
{{- end }}

endmodule
`
)
