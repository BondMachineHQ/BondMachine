package bondirect

const (
	bdEndpoint = `
-- The endpoint module is responsible for managing the communication between the BM in a node
-- and all the line interfaces towards other FPGAs (nodes)
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY {{.Prefix}}bd_endpoint_{{.NodeName}} IS
    GENERIC (
	rsize : INTEGER := {{.Rsize}} -- Size of the register
        message_length : INTEGER := {{.InnerMessLen}} -- Length of the message to be sent, in this length is not included bits used by tx and rx
    );
    PORT (
        clk : IN STD_LOGIC; -- Clock signal for the component
        reset : IN STD_LOGIC; -- Reset signal to initialize the component
	-- Interface towards the BM
	-- BM Inputs (are my outputs)
	{{- range .Inputs }}
	{{ . }} : OUT STD_LOGIC_VECTOR(rsize-1 DOWNTO 0);
	{{ . }}_valid : OUT STD_LOGIC;
	{{ . }}_recv : IN STD_LOGIC;
	{{- end }}
	-- BM Outputs (are my inputs)
	{{- range .Outputs }}
	{{ . }} : IN STD_LOGIC_VECTOR(rsize-1 DOWNTO 0);
	{{ . }}_valid : IN STD_LOGIC;
	{{ . }}_recv : OUT STD_LOGIC;
	{{- end }}
	-- Interface towards the outside world (other nodes)
	{{- range $i := iter (len .Lines) }}
	-- Line: {{index $.Lines $i}}
	-- Input transceiver: {{index $.TrIn $i}}
	{{- range $j := iter (len (index $.WiresIn $i)) }}
	{{index $.Lines $i}}{{index $.TrIn $i}}{{index (index $.WiresIn $i) $j}} : IN STD_LOGIC;
	{{- end }}
	-- Output transceiver: {{index $.TrOut $i}}
	{{- range $j := iter (len (index $.WiresOut $i)) }}
	{{index $.Lines $i}}{{index $.TrOut $i}}{{index (index $.WiresOut $i) $j}} : OUT STD_LOGIC;
	{{- end }}
	{{- end }}

    );
END {{.Prefix}}bd_endpoint_{{.NodeName}};

ARCHITECTURE Behavioral OF {{.Prefix}}bd_endpoint_{{.NodeName}} IS
	-- BM Cache signals
		-- BM Inputs
		{{- range .Inputs }}
		{{ . }}_local: STD_LOGIC_VECTOR(rsize-1 DOWNTO 0);
		{{ . }}_valid_local: STD_LOGIC := '0';
		{{ . }}_recv_local: STD_LOGIC := '0';
		{{- end }}
		-- BM Outputs
		{{- range .Outputs }}
		{{ . }}_local: STD_LOGIC_VECTOR(rsize-1 DOWNTO 0);
		{{ . }}_valid_local: STD_LOGIC := '0';
		{{ . }}_recv_local: STD_LOGIC := '0';
		{{- end }}
	-- BM need signals
		-- BM Inputs
		{{- range .Inputs }}
		{{ . }}_recv_need: STD_LOGIC := '0';
		{{ . }}_recv_need_reset: STD_LOGIC := '0';
		{{- end }}
		-- BM Outputs
		{{- range .Outputs }}
		{{ . }}_need: STD_LOGIC := '0';
		{{ . }}_need_reset: STD_LOGIC := '0';
		{{ . }}_valid_need: STD_LOGIC := '0';
		{{ . }}_valid_need_reset: STD_LOGIC := '0';
		{{- end }}
	-- Signals towards the lines
	{{- range $i := iter (len .Lines) }}
	{{- $lineName:= index $.Lines $i }}
		-- Signals towards the line {{ $lineName }}
		{{$lineName}}_s_message : STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0); -- Message to be sent to the other FPGA
        	{{$lineName}}_s_valid : STD_LOGIC; -- Signal indicating that the message is valid
        	{{$lineName}}_s_busy : STD_LOGIC := '0'; -- Signal indicating that the component is busy while transmitting
        	{{$lineName}}_s_ok : STD_LOGIC := '0'; -- Signal indicating that the outgoing transmission was successful
        	{{$lineName}}_s_error : STD_LOGIC := '0'; -- Signal indicating that an error occurred during transmission
        	{{$lineName}}_r_message : STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0) := (OTHERS => '0'); -- Message received from the other FPGA
        	{{$lineName}}_r_busy : STD_LOGIC := '0'; -- Signal indicating that the component is busy while receiving
        	{{$lineName}}_r_valid : STD_LOGIC := '0'; -- Signal indicating that the received message is valid
        	{{$lineName}}_r_error : STD_LOGIC := '0' -- Signal indicating that an error occurred during reception
	{{- end }}
BEGIN

	-- Instantiations

    	-- Instantiate the lines
	{{- range $i := iter (len .Lines) }}
	{{- $lineName:= index $.Lines $i }}
	{{$.Prefix}}bd_line_{{$.NodeName}}_{{ $lineName }}_inst : ENTITY work.{{$.Prefix}}bd_line_{{$.NodeName}}_{{ $lineName }}
	GENERIC MAP(
		message_length => {{ $.InnerMessLen }},
	)
	PORT MAP(
		clk => clk,
		reset => reset,
		{{- range $j := iter (len (index $.WiresOut $i)) }}
		{{- if eq $j 0 }}
		tx_clk => {{index $.Lines $i}}{{index $.TrOut $i}}{{index (index $.WiresOut $i) $j}},
		{{- else }}
		tx_data{{ dec $j }} => {{index $.Lines $i}}{{index $.TrOut $i}}{{index (index $.WiresOut $i) $j}},
		{{- end }}
		{{- end }}
		{{- range $j := iter (len (index $.WiresIn $i)) }}
		{{- if eq $j 0 }}
		rx_clk => {{index $.Lines $i}}{{index $.TrIn $i}}{{index (index $.WiresIn $i) $j}},
		{{- else }}
		rx_data{{ dec $j }} => {{index $.Lines $i}}{{index $.TrIn $i}}{{index (index $.WiresIn $i) $j}},
		{{- end }}
		{{- end }}
		s_message => {{$lineName}}_s_message,
		s_valid => {{$lineName}}_s_valid,
		s_busy => {{$lineName}}_s_busy,
		s_ok => {{$lineName}}_s_ok,
		s_error => {{$lineName}}_s_error,
		r_message => {{$lineName}}_r_message,
		r_busy => {{$lineName}}_r_busy,
		r_valid => {{$lineName}}_r_valid,
		r_error => {{$lineName}}_r_error
	);
	{{- end }}

	-- Processes

	-- Towards BM wires
	-- BM Inputs
	{{- range .Inputs }}
	{{ . }} <= {{ . }}_local;
	{{ . }}_valid <= {{ . }}_valid_local;
	{{- end }}
	-- BM Outputs
	{{- range .Outputs }}
	{{ . }}_recv <= {{ . }}_recv_local;
	{{- end }}

	-- BM need signals processes
	-- BM Inputs
	{{- range .Inputs }}
	{{ . }}_recv_need_proc : PROCESS (clk, reset)
    	BEGIN
        IF reset = '1' THEN
		{{ . }}_recv_need <= '0';
        ELSIF rising_edge(clk) THEN
		IF {{ . }}_recv_need_reset = '1' THEN
			{{ . }}_recv_need <= '0';
		ELSE
			{{ . }}_recv_local <= {{ . }}_recv;
			IF {{ . }}_recv_local \= {{ . }}_recv THEN
				{{ . }}_recv_need <= '1';
			END IF;
		END IF;
	END PROCESS;
	{{- end }}
	-- BM Outputs
	{{- range .Outputs }}
	{{ . }}_valid_need_proc : PROCESS (clk, reset)
    	BEGIN
        IF reset = '1' THEN
		{{ . }}_valid_need <= '0';
        ELSIF rising_edge(clk) THEN
		IF {{ . }}_valid_need_reset = '1' THEN
			{{ . }}_valid_need <= '0';
		ELSE
			{{ . }}_valid_local <= {{ . }}_valid;
			IF {{ . }}_valid_local \= {{ . }}_valid THEN
				{{ . }}_valid_need <= '1';
			END IF;
		END IF;
	END PROCESS;

	{{ . }}_need_proc : PROCESS (clk, reset)
    	BEGIN
        IF reset = '1' THEN
		{{ . }}_need <= '0';
        ELSIF rising_edge(clk) THEN
		IF {{ . }}_need_reset = '1' THEN
			{{ . }}_need <= '0';
		ELSE
			{{ . }}_local <= {{ . }};
			IF {{ . }}_local \= {{ . }} THEN
				{{ . }}_need <= '1';
			END IF;
		END IF;
	END PROCESS;	

	{{- end }}
END Behavioral;
	`
)
