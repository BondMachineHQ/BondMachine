package bondirect

const (
    bondTb = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bond_tb IS
END bond_tb;

ARCHITECTURE Behavioral OF bond_tb IS
    SIGNAL clk : STD_LOGIC := '0';
    SIGNAL message : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL message_received : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL data_enable : STD_LOGIC := '0';
    SIGNAL busy : STD_LOGIC;
    SIGNAL tx_clk : STD_LOGIC;
    SIGNAL data_line : STD_LOGIC;
    SIGNAL receiving : STD_LOGIC;
    SIGNAL failed : STD_LOGIC;
    SIGNAL data_ready : STD_LOGIC;

    CONSTANT clk_period : TIME := 10 ns;
BEGIN
    uut1: ENTITY work.bond_tx
        PORT MAP (
            clk => clk,
            reset => '0',
            message => message,
            data_enable => data_enable,
            busy => busy,
            tx_clk => tx_clk,
            tx_out => data_line
        );
    
    uut2: ENTITY work.bond_rx
        PORT MAP (
            clk => clk,
            reset => '0',
            rx_clk => tx_clk,
            rx_in => data_line,
            message => message_received,
            data_ready => data_ready,
            receiving => receiving,
            failed => failed
        );

    clk_process: PROCESS
    BEGIN
        WHILE TRUE LOOP
            clk <= '0';
            WAIT FOR clk_period / 2;
            clk <= '1';
            WAIT FOR clk_period / 2;
        END LOOP;
    END PROCESS;

    stimulus_process: PROCESS
    BEGIN
        -- Test case 1: Send a message
        message <= "11110000"; -- Example message
        data_enable <= '1';
        WAIT FOR 100 ns;
        data_enable <= '0';

        -- Disable data enable to stop sending
        -- data_enable <= '0';
        -- WAIT FOR 100 ns;

        -- Test case 2: Send another message
        -- message <= "110011001"; -- Another example message
        -- data_enable <= '1';
        -- WAIT FOR 100 ns;

        -- Finish simulation
        WAIT;
    END PROCESS;

END Behavioral;
`
)