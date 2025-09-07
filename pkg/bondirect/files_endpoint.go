package bondirect

const (
	bdEndpoint = `
-- The endpoint module is responsible for managing the communication between the BM in a node
-- and all the line interfaces towards other FPGAs (nodes)
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bd_endpoint IS
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
END bd_endpoint;

ARCHITECTURE Behavioral OF bd_endpoint IS
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
BEGIN

	-- Instantiations

    	-- Instantiate the lines
	{{- range $i := iter (len .Lines) }}
	{{- $lineName:= index $.Lines $i }}
	line_{{ $lineName }}_inst : ENTITY work.bd_line_1_1
	GENERIC MAP(
		message_length => {{ $.InnerMessLen }}
	)
	PORT MAP(
		clk => clk,
		reset => reset,
		{{- range $j := iter (len (index $.WiresOut $i)) }}
		tx_{{index (index $.WiresOut $i) $j}} => {{index $.Lines $i}}{{index $.TrOut $i}}{{index (index $.WiresOut $i) $j}},
		{{- end }}
	{{- end }}


            message => message_to_send,
            data_enable => send_data_enable,
            busy => send_busy,
            tx_clk => tx_clk,
            tx_out => tx_out
        );

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
