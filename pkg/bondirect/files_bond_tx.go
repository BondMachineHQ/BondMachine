package bondirect

const (
	bondTx = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY {{.Prefix}}bond_tx_{{.MeshNodeName}}_{{.EdgeName}} IS
    GENERIC (
        message_length : INTEGER := {{add .InnerMessLen 2}}; -- Length of the message to be sent including 2 extra bits
        num_wires : INTEGER := {{.TransParams.NumWires}}; -- Number of wires in the bond direct interface
        counters_length : INTEGER := {{.TransParams.CountersLen}}; -- Length of the counters used in the design
        out_clock_wait: INTEGER := {{.TransParams.OutClockWait}} -- Number of clock cycles to wait before toggling the output clock
    );
    PORT (
        clk : IN STD_LOGIC;
        message : IN STD_LOGIC_VECTOR (message_length-1 DOWNTO 0);
        data_enable : IN STD_LOGIC;
        busy : OUT STD_LOGIC;
        s_error : OUT STD_LOGIC;
        tx_clk : OUT STD_LOGIC;
{{- $iSeq := ""}}
{{- range $i := (iter (int .TransParams.NumWires )) }}
        tx_data{{ $i }} : OUT STD_LOGIC;
        {{- $iSeq = printf "%s'1' & " $iSeq }}
{{- end }}
        reset : IN STD_LOGIC
    );
END {{.Prefix}}bond_tx_{{.MeshNodeName}}_{{.EdgeName}};

ARCHITECTURE Behavioral OF {{.Prefix}}bond_tx_{{.MeshNodeName}}_{{.EdgeName}} IS
    SIGNAL counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    CONSTANT out_clock_tick : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(out_clock_wait, counters_length);
    CONSTANT adjusted_length : INTEGER := ((message_length + num_wires - 1) / num_wires) * num_wires; -- Adjusted message length to be a multiple of num_wires
    CONSTANT extra_bits : INTEGER := adjusted_length - message_length;
    CONSTANT readings: INTEGER := adjusted_length / num_wires;
    CONSTANT zeroes : STD_LOGIC_VECTOR(extra_bits-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL busy_sr : STD_LOGIC_VECTOR(readings-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL sending : STD_LOGIC_VECTOR(adjusted_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL int_clk : STD_LOGIC := '0';
    SIGNAL doing : STD_LOGIC := '0';
    SIGNAL error : STD_LOGIC := '0';
    SIGNAL data_enable_d : STD_LOGIC := '0';

    TYPE tx_sm IS (IDLE, SEND, FINALIZE, EXTRA);
    SIGNAL send_state : tx_sm := IDLE;
BEGIN
    busy <= doing;
    s_error <= error;
    tx_clk <= int_clk;

    clk_proc : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            counter <= (OTHERS => '0');
            busy_sr <= (OTHERS => '0');
            sending <= (OTHERS => '0');
            int_clk <= '0';
{{- range $i := (iter (int .TransParams.NumWires )) }}
            tx_data{{ $i }} <= '0';
{{- end }}
            doing <= '0';
            error <= '0';
        ELSIF rising_edge(clk) THEN
            data_enable_d <= data_enable;

            CASE send_state IS
                WHEN IDLE =>
                    doing <= '0';
                    IF data_enable = '1' AND data_enable_d = '0' THEN
                        busy_sr <= (OTHERS => '1');
                        sending <= message;
                        send_state <= SEND;
                        counter <= out_clock_tick;
                    END IF;
                WHEN SEND =>
                    doing <= '1';
                    IF busy_sr(0) /= '0' THEN
                        IF counter = 0 THEN
                            counter <= out_clock_tick;
                            IF int_clk = '0' THEN
                                int_clk <= '1';
{{- range $i := (iter (int .TransParams.NumWires )) }}
	                        tx_data{{ $i }} <= sending({{ $i }});
{{- end }}
                                sending <= {{$iSeq}} sending(sending'high DOWNTO {{.TransParams.NumWires}});
                                busy_sr <= '0' & busy_sr(busy_sr'high DOWNTO 1);
                            ELSE
                                int_clk <= '0';
                            END IF;
                        ELSE
                            counter <= counter - 1;
                        END IF;
                    ELSE
                        send_state <= FINALIZE;
                        counter <= out_clock_tick;
                    END IF;

                WHEN FINALIZE =>
                    doing <= '1';
                    if counter = 0 THEN
                        counter <= out_clock_tick;
                        send_state <= EXTRA;
                        int_clk <= '0';
{{- range $i := (iter (int .TransParams.NumWires )) }}
                        tx_data{{ $i }} <= '0';
{{- end }}
                    ELSE
                        counter <= counter - 1;
                    END IF;
                WHEN EXTRA =>
                    doing <= '1';
                    if counter = 0 THEN
                        counter <= out_clock_tick;
                        send_state <= IDLE;
                        int_clk <= '0';
{{- range $i := (iter (int .TransParams.NumWires )) }}
                        tx_data{{ $i }} <= '0';
{{- end }}
                    ELSE
                        counter <= counter - 1;
                    END IF;
                WHEN OTHERS =>
                    send_state <= IDLE;
            END CASE;
        END IF;
    END PROCESS;

END Behavioral;
`
)
