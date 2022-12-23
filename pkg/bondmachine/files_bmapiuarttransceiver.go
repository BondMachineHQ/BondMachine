package bondmachine

const (
	bmapiuarttransceiver = "\n`timescale 1ns / 1ps\n" + `

module bmapiuarttransceiver(
    input clk,
    input reset,
    output TxD,
    input RxD,
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
    output [2:0] transconnected
);

// Possible states of the state machine for the send operations
reg {{ .SendSM.Buswidth }} send_state_machine;{{- printf "\n" -}}
        {{- $first := true -}}
        {{- range $key, $value := .SendSM.Names -}}
            {{- if $first -}}
                localparam{{- printf " " -}}
                {{- $first = false -}}
            {{- else -}},{{- printf "\n\t" -}}
            {{- end -}}
            {{$value}}={{ $.SendSM.Bits }}'b{{ index $.SendSM.Binary $key }}
        {{- end -}};

// Possible states of the state machine for the receive operations
reg [1:0] recv_state_machine;
localparam  RCVSTM_COMM=2'b00,
            RCVSTM_VAL=2'b01,
            RCVSTM_DUP=2'b10,
            RCVSTM_HANDSH=2'b11;

// Every trasmitted event is a 3 bit command followed by a 5 bit IO register (so the transceiver allows up to 2^5 Inputs and outputs)
// List of commands
localparam  CMD_NEWVAL=3'b000,
            CMD_DVALIDH=3'b001,
            CMD_DVALIDL=3'b010,
            CMD_DRECVH=3'b011,
            CMD_DRECVL=3'b100,
            CMD_HANDSH=3'b101,
            CMD_KEEP=3'b110;

// List of codes associated to input registers
{{- if .Outputs }}
{{- range $key, $value := .Outputs }}
localparam CODE_{{ $value }}=5'b{{ index $.OutputsBins $key }};
{{- end }}
{{- end }}
{{- if .Inputs }}
{{- range $key, $value := .Inputs }}
localparam CODE_{{ $value }}=5'b{{ index $.InputsBins $key }};
{{- end }}
{{- end }}

integer i;
reg step;

// Local replicas
{{- if .Outputs }}
{{- range  .Outputs }}
reg {{ $.Buswidth }} local_{{ . }};
reg local_{{ . }}_valid;
reg local_{{ . }}_recv;
{{- end }}
{{- end }}
{{- if .Inputs }}
{{- range $key, $value := .Inputs }}
reg {{ $.Buswidth }} local_{{ . }};
reg local_{{ . }}_valid;
reg local_{{ . }}_recv;
{{- end }}
{{- end }}

// Local replicas init
initial begin
    {{- if .Outputs }}
    {{- range  .Outputs }}
    local_{{ . }}{{ $.Buswidth }} = {{ $.Rsize }}'d0;
    local_{{ . }}_valid=1'b0;
    local_{{ . }}_recv=1'b0;
    {{- end }}
    {{- end }}
    {{- if .Inputs }}
    {{- range $key, $value := .Inputs }}
    local_{{ . }}{{ $.Buswidth }} = {{ $.Rsize }}'d0;
    local_{{ . }}_valid=1'b0;
    local_{{ . }}_recv=1'b0;
    {{- end }}
    {{- end }}
end

reg [7:0] handshakedata;

// Assign of outgoing (to BM) signals
{{- if .Outputs }}
{{- range  .Outputs }}
assign {{ . }}_recv = local_{{ . }}_recv;
{{- end }}
{{- end }}
{{- if .Inputs }}
{{- range $key, $value := .Inputs }}
assign {{ . }}{{ $.Buswidth }} = local_{{ . }}{{ $.Buswidth }};
assign {{ . }}_valid = local_{{ . }}_valid;
{{- end }}
{{- end }}

localparam HSMASK=8'b10101010;

// Receive buffer
reg {{ $.Buswidth }} receivebuffer;

// UART Transmitter logic
reg tstart;
reg [7:0] tx_byte;
reg {{ $.SendSM.Buswidth }} istrans;

wire txactive;
reg old_txactive;
wire tx_ended;

wire received;
wire [7:0] rx_byte;

wire recv_error;

uart UART_TX_INST
  (.clk(clk),
   .transmit(tstart),
   .received(received),
   .rx_byte(rx_byte[7:0]),
   .tx_byte(tx_byte[7:0]),
   .tx(TxD),
   .rx(RxD),
   .is_transmitting(txactive),
   .recv_error(recv_error)
   );


always @(posedge clk) begin
    old_txactive <= txactive;
end

assign tx_ended = old_txactive && !txactive;

// Inital values
initial begin
    send_state_machine{{ $.SendSM.Buswidth }} = 2'd0;
    tstart = 0;
    tx_byte[7:0] = 8'd0;
    istrans{{ $.SendSM.Buswidth }} = 2'd0;
end

// Main state machine, it store the overall status
localparam STATE_WAIT=3'b000,
           STATE_HSRECV=3'b001,
           STATE_MASKSENT=3'b010,
           STATE_ACK=3'b011,
           STATE_CONNECT=3'b100;

reg [27:0] resetcounter;
reg [2:0] main_state_machine;

// Overall connection status
assign transconnected[2:0] = main_state_machine[2:0];

initial begin
    resetcounter = {28{1'b1}};
    main_state_machine[2:0] = STATE_WAIT;
end

always @(posedge clk) begin
    if (main_state_machine[2:0] == STATE_CONNECT) begin
        if ((recv_state_machine[1:0] == RCVSTM_COMM) && received)
            if (rx_byte[7:5] == CMD_KEEP)
                resetcounter <= {28{1'b1}};
            else
                resetcounter <= resetcounter - 1;
        else
            resetcounter <= resetcounter - 1;
    end
    else
        resetcounter <= {28{1'b1}};
end

reg keep_need;
reg keep_reset;
reg [25:0] keep_counter;

initial begin
    keep_counter = {26{1'b1}};
    keep_need = 1'b0;
    keep_reset = 1'b0;
end

always @(posedge clk) begin
    if (keep_reset == 1'b1) begin
        keep_need <= 1'b0;
        keep_counter <= {26{1'b1}};
    end
    else begin
        if (keep_counter == 26'b0) begin
            keep_need <= 1'b1;
        end
        else begin
            keep_counter <= keep_counter - 1;
        end
    end
end

reg sending_mask;
reg sending_ack;

initial begin
    sending_mask = 1'b0;
    sending_ack = 1'b0;
end

always @(posedge clk) begin
    case (main_state_machine[2:0])
    STATE_WAIT: begin
        if ((recv_state_machine[1:0] == RCVSTM_COMM) && received) begin
            if (rx_byte[7:5] == CMD_HANDSH) begin
                main_state_machine[2:0] <= STATE_HSRECV;
            end
        end
    end
    STATE_HSRECV: begin
        if (tx_ended && sending_mask) begin
            main_state_machine[2:0] <= STATE_MASKSENT;
        end
    end
    STATE_MASKSENT: begin
        if (received) begin
            if (rx_byte[7:0] == (HSMASK & handshakedata)) begin
                main_state_machine[2:0] <= STATE_ACK;
            end
        end
    end
    STATE_ACK: begin
        if (tx_ended && sending_ack) begin
            main_state_machine[2:0] <= STATE_CONNECT;
        end
    end
    STATE_CONNECT: begin
       //if (recv_error || (resetcounter == {28{1'b0}})) begin
        if (resetcounter == {28{1'b0}}) begin
            main_state_machine[2:0] <= STATE_WAIT;
        end
    end
    endcase
end


// Output replication and transmission needs

{{- if .Outputs }}
{{- range .Outputs }}

reg {{ . }}_need;
reg {{ . }}_need_reset;

always @(posedge clk) begin
    if ({{ . }}_need_reset == 1'b1) begin
            {{ . }}_need <= 1'b0;
    end
    else begin
        if (istrans{{ $.SendSM.Buswidth }} != SENDSTM_{{ . }}) begin
            local_{{ . }}{{ $.Buswidth }} <= {{ . }}{{ $.Buswidth }};
            if ((local_{{ . }}{{ $.Buswidth }} != {{ . }}{{ $.Buswidth }})||(main_state_machine[2:0] != STATE_CONNECT)) begin
                {{ . }}_need <= 1'b1;
            end
        end
    end
end

reg {{ . }}_valid_need;
reg {{ . }}_valid_need_reset;

always @(posedge clk) begin
    if ({{ . }}_valid_need_reset == 1'b1) begin
            {{ . }}_valid_need <= 1'b0;
    end
    else begin
        if (istrans{{ $.SendSM.Buswidth }} != SENDSTM_{{ . }}_valid) begin
            local_{{ . }}_valid <= {{ . }}_valid;
            if ((local_{{ . }}_valid != {{ . }}_valid)||(main_state_machine[2:0] != STATE_CONNECT)) begin
                {{ . }}_valid_need <= 1'b1;
            end
        end
    end
end
 
{{- end }}
{{- end }}


{{- if .Inputs }}
{{- range .Inputs }}

reg {{ . }}_recv_need;
reg {{ . }}_recv_need_reset;

always @(posedge clk) begin
    if ({{ . }}_recv_need_reset == 1'b1) begin
            {{ . }}_recv_need <= 1'b0;
    end
    else begin
        if (istrans{{ $.SendSM.Buswidth }} != SENDSTM_{{ . }}_recv) begin
            local_{{ . }}_recv <= {{ . }}_recv;
            if ((local_{{ . }}_recv != {{ . }}_recv)||(main_state_machine[2:0] != STATE_CONNECT)) begin
                {{ . }}_recv_need <= 1'b1;
            end
        end
    end
end
 
{{- end }}
{{- end }}

// Main send State machine
{{- $smindex:= 0 }}
{{- $nextsmindex:= next $smindex }}
//{{ $smindex }}{{ $nextsmindex }}
always @(posedge clk) begin
    case (send_state_machine{{ $.SendSM.Buswidth }})
        {{ index $.SendSM.Names $smindex }}: begin
            case (main_state_machine[2:0])
            STATE_HSRECV: begin
                if (sending_mask == 1'b0) begin
                    tx_byte[7:0] <= HSMASK;
                    tstart <= 1'b1;
                    sending_mask <= 1'b1;
                end
                else begin
                    if (tx_ended) begin
                        sending_mask <= 1'b0;
                    end
                    else tstart <= 1'b0;
                end
            end
            STATE_ACK: begin
               if (sending_ack == 1'b0) begin
                    tx_byte[7:0] <= HSMASK & handshakedata;
                    tstart <= 1'b1;
                    sending_ack <= 1'b1;
                end
                else begin
                    if (tx_ended) begin
                        sending_ack <= 1'b0;
                    end
                    else tstart <= 1'b0;
                end                    
            end
            STATE_CONNECT: begin
                send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
            end
            endcase
        end
{{- if .Inputs }}
{{- range .Inputs }}
{{- $smindex = inc $smindex }}
{{- $nextsmindex = next $smindex }}
//{{ $smindex }}{{ $nextsmindex }}
        {{ index $.SendSM.Names $smindex }}: begin
            if (main_state_machine[2:0] == STATE_CONNECT) begin
                if (istrans{{ $.SendSM.Buswidth }} != {{ index $.SendSM.Names $smindex }}) begin
                    if ({{ . }}_recv_need) begin
                        if (local_{{ . }}_recv == 1'b1 ) 
                            tx_byte[7:0] <= {CMD_DRECVH ,CODE_{{ . }}};
                        else
                            tx_byte[7:0] <= {CMD_DRECVL ,CODE_{{ . }}};
                        tstart <= 1'b1;
                        istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $smindex }};
                        step <= 1'b0;
                        {{ . }}_recv_need_reset <= 1'b0;
                    end 
                    else begin
                        send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                        tstart <= 1'b0;
                        {{ . }}_recv_need_reset <= 1'b0;
                    end
                end
                else begin
                    if (tx_ended) begin
                        case (step)
                        1'b0: begin
                            send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                            tstart <= 1'b0;
                            istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
                            {{ . }}_recv_need_reset <= 1'b1;
                        end
                        endcase
                    end
                    else begin
                        tstart <= 1'b0;
                        {{ . }}_recv_need_reset <= 1'b0;
                    end
                end
            end
            else begin
                send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
            end
        end        {{- end }}
{{- end }}
{{- if .Outputs }}
{{- range .Outputs }}
{{- $smindex = inc $smindex }}
{{- $nextsmindex = next $smindex }}
//{{ $smindex }}{{ $nextsmindex }}
        {{ index $.SendSM.Names $smindex }}: begin
            if (main_state_machine[2:0] == STATE_CONNECT) begin
                if (istrans{{ $.SendSM.Buswidth }} != {{ index $.SendSM.Names $smindex }}) begin
                    if ({{ . }}_need) begin
                        tx_byte[7:0] <= {CMD_NEWVAL ,CODE_{{ . }}};
                        tstart <= 1'b1;
                        istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $smindex }};
                        step <= 1'b0;
                        {{ . }}_need_reset <= 1'b0;
                    end 
                    else begin
                        send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                        tstart <= 1'b0;
                        {{ . }}_need_reset <= 1'b0;
                    end
                end
                else begin
                    if (tx_ended) begin
                        case (step)
                        1'b0: begin
                            tx_byte[7:0] <= local_{{ . }}[7:0];
                            step <= 1'b1;
                            tstart <= 1'b1;
                            {{ . }}_need_reset <= 1'b0;
                        end
                        1'b1: begin
                            send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                            tstart <= 1'b0;
                            istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
                            {{ . }}_need_reset <= 1'b1;
                        end
                        endcase
                    end
                    else begin
                        tstart <= 1'b0;
                        {{ . }}_need_reset <= 1'b0;
                    end
                end
            end
            else begin
                send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
            end
        end
{{- $smindex = inc $smindex }}
{{- $nextsmindex = next $smindex }}       
//{{ $smindex }}{{ $nextsmindex }}
        {{ index $.SendSM.Names $smindex }}: begin
            if (main_state_machine[2:0] == STATE_CONNECT) begin
                if (istrans{{ $.SendSM.Buswidth }} != {{ index $.SendSM.Names $smindex }}) begin
                    if (({{ . }}_valid_need) && (!{{ . }}_need)) begin
                        if (local_{{ . }}_valid == 1'b1 ) 
                            tx_byte[7:0] <= {CMD_DVALIDH ,CODE_{{ . }}};
                        else
                            tx_byte[7:0] <= {CMD_DVALIDL ,CODE_{{ . }}};
                        tstart <= 1'b1;
                        istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $smindex }};
                        step <= 1'b0;
                        {{ . }}_valid_need_reset <= 1'b0;
                    end 
                    else begin
                        send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                        tstart <= 1'b0;
                        {{ . }}_valid_need_reset <= 1'b0;
                    end
                end
                else begin
                    if (tx_ended) begin
                        case (step)
                        1'b0: begin
                            send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                            tstart <= 1'b0;
                            istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
                            {{ . }}_valid_need_reset <= 1'b1;
                        end
                        endcase
                    end
                    else begin
                        tstart <= 1'b0;
                        {{ . }}_valid_need_reset <= 1'b0;
                    end
                end
            end
            else begin
                send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
            end
        end        
{{- end }}
{{- end }}     
{{- $smindex = inc $smindex }}
{{- $nextsmindex = next $smindex }}       
//{{ $smindex }}{{ $nextsmindex }}
        {{ index $.SendSM.Names $smindex }}: begin
            if (main_state_machine[2:0] == STATE_CONNECT) begin
                if (istrans{{ $.SendSM.Buswidth }} != {{ index $.SendSM.Names $smindex }}) begin
                    if (keep_need) begin
                        tx_byte[7:0] <= {CMD_KEEP ,5'b0};
                        istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $smindex }};
                        tstart <= 1'b1;
                        keep_reset <= 1'b0;
                    end
                    else begin
                        send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                        tstart <= 1'b0;
                        keep_reset <= 1'b0;
                    end      
                end
                else begin
                    if (tx_ended) begin
                        send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names $nextsmindex }};
                        tstart <= 1'b0;
                        istrans{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
                        keep_reset <= 1'b1;
                    end
                    else begin
                        tstart <= 1'b0;
                        keep_reset <= 1'b0;
                    end
                end
            end
            else begin
                send_state_machine{{ $.SendSM.Buswidth }} <= {{ index $.SendSM.Names 0 }};
                keep_reset <= 1'b0;
            end
        end
    endcase
end

reg [4:0] recv_reg; 
reg [1:0] recv_reg_count;

initial begin
    recv_state_machine[1:0] = RCVSTM_COMM;
    recv_reg_count = 0;
end

// Receive process
always @(posedge clk) begin
    case (recv_state_machine[1:0])
    RCVSTM_HANDSH: begin
        case (main_state_machine[2:0])
        STATE_CONNECT: begin
            recv_state_machine[1:0] <= RCVSTM_COMM;
        end
        endcase
    end
    RCVSTM_COMM: begin
        if (received) begin
            case (rx_byte[7:5])
            CMD_HANDSH: begin
                recv_state_machine[1:0] <= RCVSTM_HANDSH;
                handshakedata <= rx_byte[7:0];
            end
            CMD_NEWVAL: begin
                recv_reg[4:0] <= rx_byte[4:0];
                recv_state_machine[1:0] <= RCVSTM_VAL;
            end
            CMD_DVALIDH: begin
                case (rx_byte[4:0])
{{- if .Inputs }}
{{- range .Inputs }}
                CODE_{{ . }}: begin
                    local_{{ . }}_valid <= 1'b1;
                    recv_state_machine[1:0] <= RCVSTM_COMM;
                end
{{- end }}
{{- end }}
                endcase
            end
            CMD_DVALIDL: begin
                case (rx_byte[4:0])
{{- if .Inputs }}
{{- range .Inputs }}
                CODE_{{ . }}: begin
                    local_{{ . }}_valid <= 1'b0;
                    recv_state_machine[1:0] <= RCVSTM_COMM;
                end
{{- end }}
{{- end }}
                endcase
            end
            CMD_DRECVH: begin
                case (rx_byte[4:0])
{{- if .Outputs }}
{{- range .Outputs }}
                CODE_{{ . }}: begin
                    local_{{ . }}_recv <= 1'b1;
                    recv_state_machine[1:0] <= RCVSTM_COMM;
                end
{{- end }}
{{- end }}
                endcase
            end
            CMD_DRECVL: begin
                case (rx_byte[4:0])
{{- if .Outputs }}
{{- range .Outputs }}
                CODE_{{ . }}: begin
                    local_{{ . }}_recv <= 1'b0;
                    recv_state_machine[1:0] <= RCVSTM_COMM;
                end
{{- end }}
{{- end }}
                endcase
            end
            endcase
        end
    end
    RCVSTM_VAL: begin
        if (received) begin
            receivebuffer[7:0] <= {receivebuffer[7:0] << 8 , rx_byte[7:0]};
            recv_reg_count <= recv_reg_count + 1;
            if (recv_reg_count == 0) begin
                recv_state_machine[1:0] <= RCVSTM_DUP;
                recv_reg_count <= 0;
            end
        end
    end
    RCVSTM_DUP: begin
        case (recv_reg[4:0])
{{- if .Inputs }}
{{- range .Inputs }}
        CODE_{{ . }}: begin
            local_{{ . }}{{ $.Buswidth }} <= receivebuffer{{ $.Buswidth }};
            recv_state_machine[1:0] <= RCVSTM_COMM;
        end
{{- end }}
{{- end }}
        endcase
    end
    endcase
end

endmodule
`
)
