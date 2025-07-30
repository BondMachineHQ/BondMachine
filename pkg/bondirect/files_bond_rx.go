package bondirect

const (
	bondRx = `
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bond_rx IS
    GENERIC (
        message_length : INTEGER := 8;
        counters_length : INTEGER := 32
    );
    PORT (
        clk : IN STD_LOGIC;
        reset : IN STD_LOGIC;
        rx_clk : IN STD_LOGIC;
        rx_in : IN STD_LOGIC;
        message : OUT STD_LOGIC_VECTOR (message_length-1 DOWNTO 0) := (OTHERS => '0');
        data_ready : OUT STD_LOGIC := '0';
        busy : OUT STD_LOGIC := '0';
        failed : OUT STD_LOGIC := '0'
    );
END bond_rx;

ARCHITECTURE Behavioral OF bond_rx IS
    TYPE state_type IS (IDLE, RECV, DONE, FAIL);
    SIGNAL current_state : state_type := IDLE;
    SIGNAL int_clk : STD_LOGIC := '0';
    SIGNAL int_clk_prev : STD_LOGIC := '0';
    CONSTANT clk_grace_period : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(5, counters_length);
    CONSTANT timeout : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(1000, counters_length);
    CONSTANT ones : STD_LOGIC_VECTOR(message_length-2 DOWNTO 0) := (OTHERS => '1');
    SIGNAL timeout_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL busy_sr : STD_LOGIC_VECTOR(message_length-1 DOWNTO 0) := (OTHERS => '1');
    SIGNAL counter : unsigned(counters_length-1 DOWNTO 0) := clk_grace_period;
    SIGNAL failed_tr : STD_LOGIC := '0';
    SIGNAL failure : STD_LOGIC := '0';
    SIGNAL message_read : STD_LOGIC_VECTOR (message_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL receiving : STD_LOGIC := '0';
BEGIN
    -- Overall failure signals
    failure <= failed_tr;
    failed <= failure;
    busy <= receiving;

    int_clk_proc : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            int_clk <= '0';
            counter <= clk_grace_period;
        ELSIF rising_edge(clk) THEN
            IF int_clk = '1' THEN
                IF counter = 0 THEN -- The external clock has been stable for the grace period
                    int_clk <= '0';
                    counter <= clk_grace_period;
                ELSE
                    IF rx_clk = '1' THEN
                        counter <= clk_grace_period; -- The external clock is not stable, reset the counter
                    ELSE
                        counter <= counter - 1;
                    END IF;
                END IF;
            ELSE
                IF counter = 0 THEN
                    int_clk <= '1'; -- The external clock is stable, set the internal clock
                    counter <= clk_grace_period;
                ELSE
                    IF rx_clk = '0' THEN
                        counter <= clk_grace_period; -- The external clock is not stable, reset the counter
                    ELSE
                        counter <= counter - 1;
                    END IF;
                END IF;
            END IF;
        END IF;
    END PROCESS;

    reading_proc : PROCESS (int_clk)
    BEGIN
        IF rising_edge(int_clk) THEN
            IF failure = '0' THEN
                message_read <= rx_in & message_read(message_length-1 DOWNTO 1); -- Shift in the received bit
            END IF;
        END IF;
    END PROCESS reading_proc;

    main_sm : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            current_state <= IDLE;
            data_ready <= '0';
            receiving <= '0';
            failed_tr <= '0';
            timeout_counter <= timeout;
        ELSIF rising_edge(clk) THEN
            int_clk_prev <= int_clk;
            CASE current_state IS
                WHEN IDLE =>
                    IF int_clk = '1' AND int_clk_prev = '0' THEN
                        current_state <= RECV;
                        receiving <= '1';
                        failed_tr <= '0';
                        data_ready <= '0';
                        timeout_counter <= timeout;
                        busy_sr <= '0' & ones;
                    END IF;
               WHEN RECV =>
                    IF int_clk = '1' AND int_clk_prev = '0' THEN
                        IF busy_sr(1) = '0' THEN
                            message <= message_read;
                            data_ready <= '1';
                            current_state <= DONE;
                        ELSE
                            busy_sr <= '0' & busy_sr(busy_sr'high DOWNTO 1);
                        END IF;
                    ELSE
                        IF timeout_counter = 0 THEN
                            failed_tr <= '1';
                            receiving <= '0';
                            data_ready <= '0';
                            current_state <= FAIL;
                        ELSE
                            timeout_counter <= timeout_counter - 1;
                        END IF;
                    END IF;
                WHEN FAIL =>
                    current_state <= IDLE;
                WHEN DONE =>
                    receiving <= '0';
                    current_state <= IDLE;
            END CASE;
        END IF;
    END PROCESS main_sm;

END Behavioral;
`
)
