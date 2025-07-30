package bondirect

const (
	bondRxTb = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bond_rx_tb IS
END bond_rx_tb;

ARCHITECTURE Behavioral OF bond_rx_tb IS
    SIGNAL clk : STD_LOGIC := '0';
    SIGNAL message : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL data_ready : STD_LOGIC := '0';
    SIGNAL busy : STD_LOGIC;
    SIGNAL rx_clk : STD_LOGIC;
    SIGNAL rx_in : STD_LOGIC;

    CONSTANT clk_period : TIME := 10 ns;
BEGIN
    uut: ENTITY work.bond_rx
        PORT MAP (
            clk => clk,
            reset => '0',  -- Assuming no reset for the test
            message => message,
            data_ready => data_ready,
            receiving => busy,
            rx_clk => rx_clk,
            rx_in => rx_in
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
        -- message <= "11110000"; -- Example message
        -- data_ready <= '1';
        -- WAIT FOR 100 ns;
        -- data_ready <= '0';

        -- Disable data enable to stop sending
        -- data_ready <= '0';
        -- WAIT FOR 100 ns;

        -- Test case 2: Send another message
        -- message <= "110011001"; -- Another example message
        -- data_ready <= '1';
        -- WAIT FOR 100 ns;

        -- Finish simulation
        WAIT;
    END PROCESS;

END Behavioral;
`
)
