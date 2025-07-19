package bondirect

const (
	bdLine = `
-- Thee bonddirect line transmitter is the component responsible for
-- transmitting data from two FPGAs. It contains a bond_tx and a bond_rx
-- component, which are used to send and receive data.
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

ENTITY bd_line IS
    PORT (
        clk : IN STD_LOGIC; -- Clock signal for the component
        reset : IN STD_LOGIC; -- Reset signal to initialize the component
        s_message : IN STD_LOGIC_VECTOR (7 DOWNTO 0); -- Message to be sent to the other FPGA
        s_valid : IN STD_LOGIC; -- Signal indicating that the message is valid
        s_busy : OUT STD_LOGIC; -- Signal indicating that the component is busy while transmitting
        s_ok : OUT STD_LOGIC; -- Signal indicating that the outgoing transmission was successful
        s_error : OUT STD_LOGIC; -- Signal indicating that an error occurred during transmission
        tx_clk : OUT STD_LOGIC; -- Clock signal to be used for transmission. Goes to the bond_tx component
        tx_out : OUT STD_LOGIC; -- Output signal for the transmitted data. Goes to the bond_tx component
        r_message : OUT STD_LOGIC_VECTOR (7 DOWNTO 0); -- Message received from the other FPGA
        r_busy : OUT STD_LOGIC; -- Signal indicating that the component is busy while receiving
        r_valid : OUT STD_LOGIC; -- Signal indicating that the received message is valid
        rx_clk : IN STD_LOGIC; -- Clock signal for receiving data. Comes from the bond_rx component
        rx_in : IN STD_LOGIC -- Input signal for the received data. Comes from the bond_rx component
    );
END bond_tx;

ARCHITECTURE Behavioral OF bd_line IS
    SIGNAL message_to_send : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL send_data_enable : STD_LOGIC := '0';
    SIGNAL send_busy : STD_LOGIC := '0';
    SIGNAL message_received : STD_LOGIC_VECTOR (7 DOWNTO 0) := (OTHERS => '0');
    SIGNAL message_valid : STD_LOGIC := '0';
    SIGNAL receive_busy : STD_LOGIC := '0';
    TYPE send_sm IS (IDLE, SEND);
    SIGNAL send_state : send_sm := IDLE;
BEGIN

    -- Instantiate the bond_tx component
    bond_tx_inst: ENTITY work.bond_tx
        PORT MAP (
            clk => clk,
            reset => reset,
            message => message_to_send,
            data_enable => send_data_enable,
            busy => send_busy,
            tx_clk => tx_clk,
            tx_out => tx_out
        );

    -- Instantiate the bond_rx component
    bond_rx_inst: ENTITY work.bond_rx
        PORT MAP (
            clk => clk,
            reset => reset,
            rx_clk => rx_clk,
            rx_in => rx_in,
            message => message_received,
            valid => message_valid,
            busy => receive_busy
        );

    -- The main process for sending data
    send_process: PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            send_state <= IDLE;
            message_to_send <= (OTHERS => '0');
            send_data_enable <= '0';
            s_busy <= '0';
            s_ok <= '0';
            s_error <= '0';
        ELSIF rising_edge(clk) THEN
            CASE send_state IS
                WHEN IDLE => 
                    -- When in IDLE state, check if there is a replay message to send
                    -- It take precedence over the message to send
                    -- TODO Finish this part
                    IF s_valid = '1' THEN
                        message_to_send <= s_message; -- Load the message to send
                        s_busy <= '1'; -- Indicate that the component is busy


END Behavioral;
`
)
