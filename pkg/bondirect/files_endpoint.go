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
	`
)
