package bmstack

const (
	stack = `
module {{ .ModuleName }}(clk,
    reset,
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
    output reg {{ . }}Ack;
    {{- end }}
    {{- end }}
    {{- if .Receivers }}
    {{- range .Receivers }}
    output reg [{{ dec $.DataSize }}:0] {{ . }}Data;
    input {{ . }}Read;
    output reg {{ . }}Ack;
    {{- end }}
    {{- end }}

    reg [{{ dec $.DataSize }}:0] memory[{{ dec $.Depth }}:0];
    reg [{{ dec (bits (inc $.Depth)) }}:0] sp;
    {{- if eq $.MemType "FIFO" }}
    reg [{{ dec (bits (inc $.Depth)) }}:0] readsp;
    reg [{{ dec (bits (inc $.Depth)) }}:0] writesp;
    {{- end }}

    assign empty = (sp==0)? 1'b1:1'b0; 
    assign full = (sp=={{ $.Depth }})? 1'b1:1'b0;
    
    wire readneed;
    wire writeneed;

    assign writeneed = ( 1'b0
    {{- if .Senders }}
    {{- range .Senders }}
            | {{ . }}Write
    {{- end }}
    {{- end }} );

    assign readneed = ( 1'b0
    {{- if .Receivers }}
    {{- range .Receivers }}
            | {{ . }}Read
    {{- end }}
    {{- end }} );

    reg [{{ dec ( bits (len .Senders) ) }}:0] sendSM;
    //{{- range $key, $value := .Senders }}
    //localparam sendSM{{ . }} = {{ bits (len $.Senders) }}'d{{ $key }};
    //{{- end }}
    
    reg [{{ dec ( bits (len .Receivers) ) }}:0] recvSM;
    //{{- range $key, $value := .Receivers }}
    //localparam recvSM{{ . }} = {{ bits (len $.Receivers) }}'d{{ $key }};
    //{{- end }}

    integer i;

    always @(posedge clk) begin
        if (reset) begin
            sp <= {{ bits (inc $.Depth) }}'d0;
            {{- if eq $.MemType "FIFO" }}
            readsp <= {{ bits (inc $.Depth) }}'d0;
            writesp <= {{ bits (inc $.Depth) }}'d0;
            {{- end }}
            {{- range $key, $value := .Receivers }}
            {{ $value }}Data <= {{ $.DataSize}}'d0;
            {{ $value }}Ack <= 1'b0;
            {{- end }}
            {{- range $key, $value := .Senders }}
            {{ $value }}Ack <= 1'b0;
            {{- end }}
            sendSM <= {{bits (len .Receivers)}}'d0;
            recvSM <= {{bits (len .Receivers)}}'d0;
            for (i=0;i<{{ .Depth }};i=i+1) begin
                memory[i]<={{ $.DataSize}}'d0;
            end
        end
        else begin
            // Read state machine part
            if (readneed && !empty) begin
                case (recvSM)
                {{- range $key, $value := .Receivers }}
                {{ bits (len $.Receivers) }}'d{{ $key }}: begin
                    if ({{ $value }}Read && !{{ $value }}Ack) begin
                        {{- if eq $.MemType "LIFO" }}
                        {{ $value }}Data[{{ dec $.DataSize }}:0] <= memory[sp-1];
                        sp <= sp - 1;
                        {{- end }}
                        {{- if eq $.MemType "FIFO" }}
                        {{ $value }}Data[{{ dec $.DataSize }}:0] <= memory[readsp];
                        if (readsp=={{ dec $.Depth }}) begin
                            readsp <= 0;
                            sp <=  writesp;
                        end
                        else begin
                            readsp <= readsp + 1;
                            if (writesp < readsp + 1) begin
                                sp <= {{ $.Depth }} - readsp -1 + writesp;
                            end
                            else begin
                                sp <= writesp - readsp - 1;
                            end
                        end
                        {{- end }}
                    end
                    recvSM <= {{ bits (len $.Receivers) }}'d{{ next $key (len $.Receivers) }};
                end
                {{- end }}
                endcase
            end
            // Write state machine part
            else if (writeneed && !full) begin
                case (sendSM)
                {{- range $key, $value := .Senders }}
                {{ bits (len $.Senders) }}'d{{ $key }}: begin
                    if ({{ $value }}Write && !{{ $value }}Ack) begin
                        {{- if eq $.MemType "LIFO" }}
                        memory[sp] <= {{ $value }}Data[{{ dec $.DataSize }}:0];
                        sp <= sp + 1;
                        {{- end }}
                        {{- if eq $.MemType "FIFO" }}
                        memory[writesp] <= {{ $value }}Data[{{ dec $.DataSize }}:0];
                        if (writesp=={{ dec $.Depth }}) begin
                            writesp <= 0;
                            sp <= {{ $.Depth }} - readsp;
                        end
                        else begin
                            writesp <= writesp + 1;
                            if (writesp + 1 > readsp) begin
                                sp <= writesp - readsp + 1;
                            end
                            else begin
                                sp <= {{ $.Depth }} - readsp + writesp + 1;
                            end
                        end
                        {{- end }}
                    end
                    sendSM <= {{ bits (len $.Senders) }}'d{{ next $key (len $.Senders) }};
                end
                {{- end }}
                endcase
            end

            // Read ack process
            {{- range $key, $value := .Receivers }}
            if ({{ $value }}Read && !{{ $value }}Ack && recvSM=={{ bits (len $.Receivers) }}'d{{ $key }} && !empty) begin
                {{ $value }}Ack <= 1'b1;
            end
            else begin
                if (!{{ $value }}Read) begin
                    {{ $value }}Ack <= 1'b0;
                end
            end
            {{- end }}

            // Write ack process
            {{- range $key, $value := .Senders }}
            if (!(readneed && !empty) && {{ $value }}Write && !{{ $value }}Ack && sendSM=={{ bits (len $.Senders) }}'d{{ $key }} && !full) begin
                {{ $value }}Ack <= 1'b1;
            end
            else begin
                if (!{{ $value }}Write) begin
                    {{ $value }}Ack <= 1'b0;
                end
            end
            {{- end }}
        end
    end
endmodule
`

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
