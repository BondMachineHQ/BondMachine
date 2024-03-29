package bondmachine

const (
	uartSO = "\n`timescale 1ns / 1ps" + `
module {{.ModuleName}}(
    input clk,
    input rst,
    {{ range $i, $e := .Receivers }}
    input [7:0] {{ $e }}Data,
    input {{ $e }}Write,
    output {{ $e }}Ack,
    {{ end }}
    {{ range $i, $e := .Senders }}
    output [7:0] {{ $e }}Data,
    output {{ $e }}Read,
    input {{ $e }}Ack,
    {{ end }}
    input {{.ModuleName}}_rx,
    output {{.ModuleName}}_tx,     
    output rempty,
    output rfull,
    output wempty,
    output wfull
    );

    reg transmit;
    reg [7:0] tx_byte;

    wire is_receiving;
    wire is_transmitting;
    wire recv_error;
    wire [3:0] rx_samples;
    wire [3:0] rx_sample_countdown;

    reg tstart;
    reg [2:0] istrans;
    
    wire txactive;
    wire tx_ended;
    
    wire received;
    wire [7:0] rx_byte;
	
	wire [7:0] uartwriterData;
	reg uartwriterRead;
	wire uartwriterAck;

    reg [7:0] uartreaderData;
    reg uartreaderWrite;
    wire uartreaderAck;

{{.ModuleName}}rfifo {{.ModuleName}}rfifo_inst(.clk(clk),
    .reset(reset),
    {{ range $i, $e := .Senders }}
    .{{ $e }}Data({{ $e }}Data),
    .{{ $e }}Read({{ $e }}Read),
    .{{ $e }}Ack({{ $e }}Ack),
    {{ end }}
    .uartreaderData(uartreaderData),
    .uartreaderWrite(uartreaderWrite),
    .uartreaderAck(uartreaderAck),
    .empty(rempty),
    .full(rfull)
);    


{{.ModuleName}}wfifo {{.ModuleName}}wfifo_inst(.clk(clk),
    .reset(reset),
    {{ range $i, $e := .Receivers }}
    .{{ $e }}Data({{ $e }}Data),
    .{{ $e }}Write({{ $e }}Write),
    .{{ $e }}Ack({{ $e }}Ack),
    {{ end }}
    .uartwriterData(uartwriterData),
    .uartwriterRead(uartwriterRead),
    .uartwriterAck(uartwriterAck),
    .empty(wempty),
    .full(wfull)
);

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

reg [1:0] outSM;
 
localparam [1:0]     
    OUT_IDLE             = 2'd0,
    OUT_WAIT             = 2'd1,
    OUT_DONE             = 2'd2;
        
// Sending out to uart from the write FIFO
always @(posedge clk) begin
        if (reset) begin
            uartwriterRead <= 1'b0;
            transmit <= 1'b0;
        end
        else begin
            case (outSM)
            OUT_IDLE: begin
                if (!wempty) begin
                    if (uartwriterAck && uartwriterRead) begin
                        uartwriterRead <= 1'b0;
                        tx_byte[7:0] <= uartwriterData[7:0];
                        transmit <= 1'b1;
                        outSM <= OUT_WAIT;
                    end
                    else begin
                        uartwriterRead <= 1'b1;
                        transmit <= 1'b0;
                    end
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
            endcase
        end
end

reg [1:0] inSM;
 
localparam [1:0]     
    IN_IDLE             = 2'd0,
    IN_WAIT             = 2'd1,
    IN_DONE             = 2'd2;

// Reading the UART and pushing to the read FIFO
always @(posedge clk) begin
        if (reset) begin
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
    
endmodule    

`
)
