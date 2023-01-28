package bmserialize

const (
	testbench = "`" + `timescale 1ns / 1ps

module request(clk, reset, req, ack, impulse);
    input clk;
    input reset;
    output reg req;
    input ack;
    input impulse;

    reg state;

    initial begin
        state = 0;
        req = 0;
    end

    always @(posedge clk) begin
        if (reset) begin
            state <= 0;
        end else begin
            case (state)
                0: begin
                    req <= 0;
                    if (impulse) begin
                        state <= 1;
                    end
                end
                1: begin
                    req <= 1;
                    if (ack) begin
                        state <= 0;
                    end
                end
            endcase
        end
    end

endmodule

module {{ .ModuleName }}_tb;

    // Inputs and outputs
    reg clk;
    reg reset;
    wire empty;
    wire full;
    {{- if .Senders }}
    {{- range .Senders }}
    reg [{{ dec $.DataSize }}:0] {{ . }}Data;
    wire {{ . }}Write;
    wire {{ . }}Ack;
    reg {{ . }}Impulse;
    {{- end }}
    {{- end }}
    {{- if .Receivers }}
    {{- range .Receivers }}
    wire [{{ dec $.DataSize }}:0] {{ . }}Data;
    wire {{ . }}Read;
    wire {{ . }}Ack;
    reg {{ . }}Impulse;
    {{- end }}
    {{- end }}

    // Clock
    always #1 clk = ~clk;

    // Instantiate the Unit Under Test (UUT)
    {{ .ModuleName }} uut (
        .clk(clk),
        .reset(reset),
        {{- if .Senders }}
        {{- range .Senders }}
        .{{ . }}Data({{ . }}Data),
        .{{ . }}Write({{ . }}Write),
        .{{ . }}Ack({{ . }}Ack),
        {{- end }}
        {{- end }}
        {{- if .Receivers }}
        {{- range .Receivers }}
        .{{ . }}Data({{ . }}Data),
        .{{ . }}Read({{ . }}Read),
        .{{ . }}Ack({{ . }}Ack),
        {{- end }}
        {{- end }}
        .empty(empty),
        .full(full)
    );

    // Instantiate the stimulus process
    {{- if .Senders }}
    {{- range .Senders }}
    request {{ . }}_req(
        .clk(clk),
        .reset(reset),
        .req({{ . }}Write),
        .ack({{ . }}Ack),
        .impulse({{ . }}Impulse)
    );
    {{- end }}
    {{- end }}
    {{- if .Receivers }}
    {{- range .Receivers }}
    request {{ . }}_req(
        .clk(clk),
        .reset(reset),
        .req({{ . }}Read),
        .ack({{ . }}Ack),
        .impulse({{ . }}Impulse)
    );
    {{- end }}
    {{- end }}

    initial begin
		$dumpfile("{{ .ModuleName }}.vcd");
		$dumpvars;
	end

    initial begin
        clk = 0;
        reset = 1;
        {{- if .Senders }}
        {{- range .Senders }}
        {{ . }}Data = 0;
        {{ . }}Impulse = 0;
        {{- end }}
        {{- end }}
        {{- if .Receivers }}
        {{- range .Receivers }}
        {{ . }}Impulse = 0;
        {{- end }}
        {{- end }}

        // Wait 100 ns for reset
        #100;        
  
        // Release reset
        reset = 1'b0;

        // Start tests
        {{- range .TestSequence }}
        {{ . }}
        {{- end }}

        #1000;
        $finish;

    end
endmodule
`
)
