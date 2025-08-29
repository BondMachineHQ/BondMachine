package bondirect

const (
	bondTx = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bond_tx IS
    GENERIC (
        message_length : INTEGER := {{.Rsize}};
        counters_length : INTEGER := 32
    );
    PORT (
        clk : IN STD_LOGIC;
        reset : IN STD_LOGIC;
        message : IN STD_LOGIC_VECTOR (message_length-1 DOWNTO 0);
        data_enable : IN STD_LOGIC;
        busy : OUT STD_LOGIC;
        tx_clk : OUT STD_LOGIC;
        tx_out : OUT STD_LOGIC
    );
END bond_tx;

ARCHITECTURE Behavioral OF bond_tx IS
    SIGNAL counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    CONSTANT out_clock_tick : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(10, counters_length);
    SIGNAL busy_sr : STD_LOGIC_VECTOR(message_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL sending : STD_LOGIC_VECTOR(message_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL int_clk : STD_LOGIC := '0';
    SIGNAL doing : STD_LOGIC := '0';
BEGIN
    busy <= doing;
    tx_clk <= int_clk;

    clk_proc : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            counter <= (OTHERS => '0');
            busy_sr <= (OTHERS => '0');
            sending <= (OTHERS => '0');
            int_clk <= '0';
            tx_out <= '0';
            doing <= '0';
        ELSIF rising_edge(clk) THEN
            IF doing = '1' AND busy_sr(0) /= '0' THEN
                IF counter = 0 THEN
                    counter <= out_clock_tick;
                    IF int_clk = '0' THEN
                        int_clk <= '1';
                        tx_out <= sending(0);
                        sending <= '1' & sending(sending'high DOWNTO 1);
                        busy_sr <= '0' & busy_sr(busy_sr'high DOWNTO 1);
                    ELSE
                        int_clk <= '0';
                    END IF;
                ELSE
                    counter <= counter - 1;
                END IF;
            ELSE
                IF counter = 0 THEN
                    int_clk <= '0';
                    tx_out <= '0';
                ELSE
                    counter <= counter - 1;
                END IF;
            END IF;

            IF busy_sr(0) = '0' THEN
                IF data_enable = '1' THEN
                    IF doing = '0' THEN
                        counter <= to_unsigned(0, counters_length);
                        sending <= message;
                        busy_sr <= (OTHERS => '1');
                        doing <= '1';
                    END IF;
                ELSE
                    doing <= '0';
                END IF;
            ELSE
                doing <= '1';
            END IF;
        END IF;
    END PROCESS;

END Behavioral;
`
)
