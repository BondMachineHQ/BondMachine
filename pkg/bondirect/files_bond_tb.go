package bondirect

const (
	bondTb = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bd_line_tb IS
END bd_line_tb;

ARCHITECTURE Behavioral OF bd_line_tb IS
    SIGNAL clk : STD_LOGIC := '0';
    SIGNAL reset : STD_LOGIC := '0';
    SIGNAL e0_tx_clk : STD_LOGIC;
    SIGNAL e0_tx_out : STD_LOGIC;
    SIGNAL e1_tx_clk : STD_LOGIC;
    SIGNAL e1_tx_out : STD_LOGIC;

    SIGNAL e0_s_message : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL e0_s_valid : STD_LOGIC := '0';
    SIGNAL e0_s_busy : STD_LOGIC;
    SIGNAL e0_s_ok : STD_LOGIC;
    SIGNAL e0_s_error : STD_LOGIC;
    SIGNAL e0_r_message : STD_LOGIC_VECTOR (7 DOWNTO 0);
    SIGNAL e0_r_busy : STD_LOGIC;
    SIGNAL e0_r_valid : STD_LOGIC;

    SIGNAL e1_s_message : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL e1_s_valid : STD_LOGIC := '0';
    SIGNAL e1_s_busy : STD_LOGIC;
    SIGNAL e1_s_ok : STD_LOGIC;
    SIGNAL e1_s_error : STD_LOGIC;
    SIGNAL e1_r_message : STD_LOGIC_VECTOR (7 DOWNTO 0);
    SIGNAL e1_r_busy : STD_LOGIC;
    SIGNAL e1_r_valid : STD_LOGIC;

    CONSTANT clk_period : TIME := 10 ns;
BEGIN
    endpoint0 : ENTITY work.bd_line
        PORT MAP(
            clk => clk,
            reset => reset,
            tx_clk => e0_tx_clk,
            tx_out => e0_tx_out,
            rx_clk => e1_tx_clk,
            rx_in => e1_tx_out,
            s_message => e0_s_message,
            s_valid => e0_s_valid,
            s_busy => e0_s_busy,
            s_ok => e0_s_ok,
            s_error => e0_s_error,
            r_message => e0_r_message,
            r_busy => e0_r_busy,
            r_valid => e0_r_valid
        );
    endpoint1 : ENTITY work.bd_line
        PORT MAP(
            clk => clk,
            reset => reset,
            tx_clk => e1_tx_clk,
            tx_out => e1_tx_out,
            rx_clk => e0_tx_clk,
            rx_in => e0_tx_out,
            s_message => e1_s_message,
            s_valid => e1_s_valid,
            s_busy => e1_s_busy,
            s_ok => e1_s_ok,
            s_error => e1_s_error,
            r_message => e1_r_message,
            r_busy => e1_r_busy,
            r_valid => e1_r_valid
        );

    clk_process : PROCESS
    BEGIN
        WHILE TRUE LOOP
            clk <= '0';
            WAIT FOR clk_period / 2;
            clk <= '1';
            WAIT FOR clk_period / 2;
        END LOOP;
    END PROCESS;

    stimulus_process : PROCESS
    BEGIN
        -- Test case 1: Send a message
        e0_s_message <= "10101010"; -- Example message
        e0_s_valid <= '1';
        WAIT UNTIL e0_s_busy = '1';
        e0_s_valid <= '0';

        WAIT FOR 1000 ns;

        e1_s_message <= "11100011"; -- Example message
        e1_s_valid <= '1';
        WAIT UNTIL e1_s_busy = '1';
        e1_s_valid <= '0';

        WAIT UNTIL e0_s_ok = '1' OR e0_s_error = '1';

        WAIT UNTIL e1_s_ok = '1' OR e1_s_error = '1';

        WAIT FOR 1000 ns; -- Wait for some time to observe the behavior

        e0_s_message <= "01010101"; -- Example message
        e0_s_valid <= '1';
        WAIT UNTIL e0_s_busy = '1';
        e0_s_valid <= '0';

        WAIT FOR 1000 ns;

        e1_s_message <= "11001100"; -- Example message
        e1_s_valid <= '1';
        WAIT UNTIL e1_s_busy = '1';
        e1_s_valid <= '0';

        WAIT UNTIL e0_s_ok = '1' OR e0_s_error = '1';

        WAIT UNTIL e1_s_ok = '1' OR e1_s_error = '1';

        -- Finish simulation
        WAIT;
    END PROCESS;

END Behavioral;
`
)
