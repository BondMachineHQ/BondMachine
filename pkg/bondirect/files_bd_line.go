package bondirect

// TODO make the number of wires configurable, currently hardcoded to 4, also the name will reflect
// the number of wires used in the design for example: bd_line_1_1
const (
	bdLine = `
-- The bondirect line transmitter is the component responsible for
-- transmitting data from two FPGAs. It contains a bond_tx and a bond_rx
-- component, which are used to send and receive data.
LIBRARY IEEE;
USE IEEE.STD_LOGIC_1164.ALL;
USE IEEE.NUMERIC_STD.ALL;

{{- $lineIdx := -1 }}
{{- range $i := iter (len .Lines) }}
{{- if eq $.EdgeName (index $.Lines $i) }}
{{- $lineIdx = $i }}
{{- end }}
{{- end }}
-- Line index: {{ $lineIdx }}

ENTITY {{.Prefix}}bd_line_{{.MeshNodeName}}_{{ .EdgeName }} IS
    GENERIC (
        rsize : INTEGER := {{.Rsize}}; -- Size of the register
        message_length : INTEGER := {{ .InnerMessLen }}; -- Length of the message to be sent, in this length is not included bits used by tx and rx
        counters_length : INTEGER := 32 -- Length of the counters used for timeouts
    );
    PORT (
        clk : IN STD_LOGIC; -- Clock signal for the component
        reset : IN STD_LOGIC; -- Reset signal to initialize the component
        tx_clk : OUT STD_LOGIC; -- Clock signal to be used for transmission. Goes to the bond_tx component
{{- range $j := iter (dec (len (index $.WiresOut $lineIdx))) }}
        tx_data{{$j}} : OUT STD_LOGIC; -- Output signal for the wire {{$j}}. Goes to the physical line
{{- end }}
        rx_clk : IN STD_LOGIC; -- Clock signal for receiving data. Comes from the bond_rx component
{{- range $j := iter (dec (len (index $.WiresIn $lineIdx))) }}
        rx_data{{$j}} : IN STD_LOGIC; -- Input signal for the wire {{$j}}. Comes from the physical line
{{- end }}
        s_message : IN STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0); -- Message to be sent to the other FPGA
        s_valid : IN STD_LOGIC; -- Signal indicating that the message is valid
        s_busy : OUT STD_LOGIC := '0'; -- Signal indicating that the component is busy while transmitting
        s_ok : OUT STD_LOGIC := '0'; -- Signal indicating that the outgoing transmission was successful
        s_error : OUT STD_LOGIC := '0'; -- Signal indicating that an error occurred during transmission
        r_message : OUT STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0) := (OTHERS => '0'); -- Message received from the other FPGA
        r_busy : OUT STD_LOGIC := '0'; -- Signal indicating that the component is busy while receiving
        r_valid : OUT STD_LOGIC := '0'; -- Signal indicating that the received message is valid
        r_error : OUT STD_LOGIC := '0' -- Signal indicating that an error occurred during reception
    );
END {{.Prefix}}bd_line_{{.MeshNodeName}}_{{ .EdgeName }};

ARCHITECTURE Behavioral OF {{.Prefix}}bd_line_{{.MeshNodeName}}_{{ .EdgeName }} IS
    SIGNAL message_to_send : STD_LOGIC_VECTOR (message_length + 1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL send_data_enable : STD_LOGIC := '0';
    SIGNAL send_busy : STD_LOGIC := '0';
    SIGNAL message_to_validate : STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL message_valid : STD_LOGIC := '0';
    SIGNAL start_wait_for_reply : STD_LOGIC := '0'; -- Signal to indicate if we are waiting for a reply
    SIGNAL start_wait_for_ack : STD_LOGIC := '0'; -- Signal to indicate if we are waiting for an ack
    SIGNAL stop_wait_for_reply : STD_LOGIC := '0'; -- Signal to stop waiting for a reply
    SIGNAL stop_wait_for_ack : STD_LOGIC := '0'; -- Signal to stop waiting for an ack
    SIGNAL wait_for_reply : STD_LOGIC := '0'; -- Signal to indicate if we are waiting for a reply
    SIGNAL wait_for_ack : STD_LOGIC := '0'; -- Signal to indicate if we are waiting for an ack
    CONSTANT timeout : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(1000000000, counters_length);
    SIGNAL timeout_reply_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL timeout_ack_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL timeout_ack : STD_LOGIC := '0';
    SIGNAL timeout_reply : STD_LOGIC := '0';
    CONSTANT timeout_send : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(1000000000, counters_length);
    SIGNAL timeout_send_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL message_waiting_for_reply : STD_LOGIC_VECTOR (message_length - 1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL receive_busy : STD_LOGIC := '0';
    SIGNAL receive_failed : STD_LOGIC := '0';
    SIGNAL ack_result : STD_LOGIC := '0'; -- Result of the ACK operation
    SIGNAL ack_send_needed : STD_LOGIC := '0'; -- Signal to indicate if an ACK is needed to be sent
    SIGNAL ack_send_needed_reset : STD_LOGIC := '0'; -- Signal to reset the ack_send_needed signal
    SIGNAL reply_send_needed : STD_LOGIC := '0'; -- Signal to indicate if a REPLY is needed to be sent
    SIGNAL reply_send_needed_reset : STD_LOGIC := '0'; -- Signal to reset the reply_send_needed signal
    SIGNAL receiving: STD_LOGIC := '0'; -- Signal to indicate if we are currently receiving a message
    SIGNAL sending: STD_LOGIC := '0'; -- Signal to indicate if we are currently sending a message
    SIGNAL recv_header : STD_LOGIC_VECTOR (1 DOWNTO 0) := (OTHERS => '0'); -- Header of the received message

    CONSTANT sending_limit : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(1000000000, counters_length);
    CONSTANT receiving_limit : unsigned(counters_length-1 DOWNTO 0) := to_unsigned(1000000000, counters_length);
    SIGNAL sending_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL receiving_counter : unsigned(counters_length-1 DOWNTO 0) := (OTHERS => '0');
    SIGNAL sending_timeout : STD_LOGIC := '0';
    SIGNAL receiving_timeout : STD_LOGIC := '0';

    CONSTANT MSG_DATA : STD_LOGIC_VECTOR(1 DOWNTO 0) := "00";
    CONSTANT MSG_REPLY : STD_LOGIC_VECTOR(1 DOWNTO 0) := "01";
    CONSTANT MSG_ACK : STD_LOGIC_VECTOR(1 DOWNTO 0) := "10";
    CONSTANT MSG_NACK : STD_LOGIC_VECTOR(1 DOWNTO 0) := "11";

    TYPE send_sm IS (IDLE, ACK, REPLY, SEND);
    SIGNAL send_state : send_sm := IDLE;

    TYPE check IS (ACK, REPLY, DATA);
    SIGNAL check_type : check := DATA;

    TYPE phase IS (P1, P2, P3, P4);
    SIGNAL current_phase : phase := P1;

    TYPE recv_sm IS (RECV_IDLE, RECV_WAIT, RECV_PROCESS, RECV_FAIL);
    SIGNAL recv_state : recv_sm := RECV_IDLE;

    SIGNAL message_received : STD_LOGIC_VECTOR (message_length + 1 DOWNTO 0) := (OTHERS => '0'); -- Message received from the bond_rx component
    SIGNAL message_to_examine : STD_LOGIC_VECTOR (message_length + 1 DOWNTO 0) := (OTHERS => '0'); -- Message to be examined after receiving

BEGIN
    r_busy <= receiving;
    s_busy <= sending;

    recv_header <= message_received(message_length + 1 DOWNTO message_length);

    -- Instantiate the bond_tx component
    {{.Prefix}}bond_tx_{{.MeshNodeName}}_{{.EdgeName}}_inst : ENTITY work.{{.Prefix}}bond_tx_{{.MeshNodeName}}_{{.EdgeName}}
        GENERIC MAP(
            message_length => message_length + 2, -- +2 for the message length plus bits used by tx and rx
            counters_length => 32
        )
        PORT MAP(
            clk => clk,
            reset => reset,
            message => message_to_send,
            data_enable => send_data_enable,
            busy => send_busy,
{{- range $j := iter (dec (len (index $.WiresOut $lineIdx))) }}
            tx_data{{$j}} => tx_data{{$j}},
{{- end }}
            tx_clk => tx_clk
        );

    -- Instantiate the bond_rx component
    {{.Prefix}}bond_rx_{{.MeshNodeName}}_{{.EdgeName}}_inst : ENTITY work.{{.Prefix}}bond_rx_{{.MeshNodeName}}_{{.EdgeName}}
        GENERIC MAP(
            message_length => message_length + 2, -- +2 for the message length plus bits used by tx and rx
            counters_length => 32
        )
        PORT MAP(
            clk => clk,
            reset => reset,
            rx_clk => rx_clk,
{{- range $j := iter (dec (len (index $.WiresIn $lineIdx))) }}
            rx_data{{$j}} => rx_data{{$j}},
{{- end }}
            message => message_received,
            data_ready => message_valid,
            busy => receive_busy,
            failed => receive_failed
        );

    -- Manage s_busy signal and timeout
    busy_s_process : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            sending_counter <= (OTHERS => '0');
            sending_timeout <= '0';
        ELSIF rising_edge(clk) THEN
            IF sending = '1' THEN
                IF sending_counter = sending_limit THEN
                    sending_timeout <= '1';
                ELSE
                    sending_timeout <= '0';
                    sending_counter <= sending_counter + 1;
                END IF;
            ELSE
                sending_counter <= (OTHERS => '0');
                sending_timeout <= '0';
            END IF;
        END IF;
    END PROCESS busy_s_process;

    -- Manage r_busy signal and timeout
    busy_r_process : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            receiving_counter <= (OTHERS => '0');
            receiving_timeout <= '0';
        ELSIF rising_edge(clk) THEN
            IF receiving = '1' THEN
                IF receiving_counter = receiving_limit THEN
                    receiving_timeout <= '1';
                ELSE
                    receiving_timeout <= '0';
                    receiving_counter <= receiving_counter + 1;
                END IF;
            ELSE
                receiving_counter <= (OTHERS => '0');
                receiving_timeout <= '0';
            END IF;
        END IF;
    END PROCESS busy_r_process;

    -- The timeout process for managing timeouts for replies and acks
    timeout_process : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            timeout_reply_counter <= (OTHERS => '0');
            timeout_ack_counter <= (OTHERS => '0');
            timeout_reply <= '0';
            timeout_ack <= '0';
            wait_for_ack <= '0';
            wait_for_reply <= '0';
        ELSIF rising_edge(clk) THEN
            -- Manage timeout for reply

            IF start_wait_for_reply = '1' THEN
                wait_for_reply <= '1';
            ELSIF stop_wait_for_reply = '1' THEN
                wait_for_reply <= '0';
            END IF;

            IF wait_for_reply = '1' THEN
                IF timeout_reply_counter = 0 THEN
                    timeout_reply <= '0';
                ELSE
                    timeout_reply_counter <= timeout_reply_counter - 1;
                END IF;
            ELSE
                timeout_reply_counter <= timeout;
                timeout_reply <= '0';
            END IF;
            -- Manage timeout for ack

            IF start_wait_for_ack = '1' THEN
                wait_for_ack <= '1';
            ELSIF stop_wait_for_ack = '1' THEN
                wait_for_ack <= '0';
            END IF;

            IF wait_for_ack = '1' THEN
                IF timeout_ack_counter = 0 THEN
                    timeout_ack <= '0';
                ELSE
                    timeout_ack_counter <= timeout_ack_counter - 1;
                END IF;
            ELSE
                timeout_ack_counter <= timeout;
                timeout_ack <= '0';
            END IF;
        END IF;
    END PROCESS timeout_process;

    -- The main process for sending data
    send_process : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            send_state <= IDLE;
            message_to_send <= (OTHERS => '0');
            send_data_enable <= '0';
            sending <= '0';
            s_ok <= '0';
            s_error <= '0';
            start_wait_for_ack <= '0';
            start_wait_for_reply <= '0';
            current_phase <= P1;
            check_type <= DATA;
        ELSIF rising_edge(clk) THEN
            reply_send_needed_reset <= '0';
            ack_send_needed_reset <= '0';

            IF timeout_reply = '1' OR sending_timeout = '1' THEN
                s_error <= '1'; -- Indicate that an error occurred during transmission due to reply timeout
                sending <= '0'; -- Reset the busy signal
                s_ok <= '0'; -- Reset the ok signal
                send_data_enable <= '0'; -- Disable data sending
                current_phase <= P1; -- Reset the phase to P1
                send_state <= IDLE; -- Go back to IDLE state
            ELSE
                CASE send_state IS
                    WHEN IDLE =>
                        -- When in IDLE state:
                        -- 1 - Check if there is a ack message to send and eventually go to the ACK state
                        -- 2 - Check if there is a reply message to send and eventually go to the REPLY state
                        -- 3 - Check if there is a message to send and eventually go to the SEND state
                        CASE check_type IS
                            WHEN ACK =>
                                IF ack_send_needed = '1' THEN
                                    send_state <= ACK; -- Move to ACK state if ACK is needed
                                    ack_send_needed_reset <= '1';
                                    current_phase <= P1; -- Reset the phase to P1
                                END IF;
                                check_type <= REPLY; -- After checking ACK, check for REPLY next
                            WHEN REPLY =>
                                IF reply_send_needed = '1' THEN
                                    send_state <= REPLY; -- Move to REPLY state if REPLY is needed
                                    reply_send_needed_reset <= '1';
                                    current_phase <= P1; -- Reset the phase to P1
                                END IF;
                                check_type <= DATA; -- After checking REPLY, check for DATA next
                            WHEN DATA =>
                                IF  s_valid = '1' THEN
                                    IF sending = '0' THEN
                                        sending <= '1'; -- Indicate that the component send channel is busy
                                        s_ok <= '0'; -- Reset the ok signal
                                        s_error <= '0'; -- Reset the error signal
                                        send_state <= SEND; -- Move to SEND state
                                        current_phase <= P1; -- Reset the phase to P1
                                    ELSE
                                        sending <= '0'; -- Reset the busy signal if no message to send
                                        s_ok <= '0'; -- Reset the ok signal
                                        s_error <= '1'; -- Reset the error signal
                                    END IF;
                                END IF;
                                check_type <= ACK; -- After checking DATA, check for ACK next
                        END CASE;
                    WHEN ACK =>
                        CASE current_phase IS
                            WHEN P1 => -- The phase1 of the ACK sending process is preparing the ACK message
                                -- Prepare the ACK message to send
                                IF ack_result = '1' THEN
                                    message_to_send <= MSG_ACK & (message_length - 1 DOWNTO 0 => '0'); -- ACK message
                                ELSE
                                    message_to_send <= MSG_NACK & (message_length - 1 DOWNTO 0 => '0'); -- NACK message
                                END IF;
                                send_data_enable <= '1'; -- Enable data sending
                                current_phase <= P2; -- Move to next phase
                                timeout_send_counter <= timeout_send; -- Reset the timeout counter
                            WHEN P2 => -- The phase2 of the ACK sending process is waiting for the bond_tx to be busy
                                IF send_busy = '1' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    current_phase <= P3; -- Move to next phase if bond_tx is busy
                                    timeout_send_counter <= timeout_send; -- Reset the timeout counter
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during ACK sending due to timeout
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P3 => -- The phase3 of the ACK sending process is waiting for the bond_tx to finish sending
                                IF send_busy = '0' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                    current_phase <= P1; -- Reset the phase to P1
                                    IF ack_result = '1' THEN
                                        s_ok <= '1'; -- Indicate that the ACK was sent successfully
                                        s_error <= '0'; -- Reset the error signal
                                    ELSE
                                        s_error <= '1'; -- Indicate that an error occurred during ACK sending
                                        s_ok <= '0'; -- Reset the ok signal
                                    END IF;
                                    sending <= '0'; -- Reset the busy signal
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during ACK sending due to timeout
                                        s_ok <= '0'; -- Reset the ok signal
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P4 =>
                                current_phase <= P1;
                                send_data_enable <= '0'; -- Disable data sending
                                send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                        END CASE;
                    WHEN REPLY =>
                        CASE current_phase IS
                            WHEN P1 => -- The phase1 of the REPLY sending process is preparing the REPLY message
                                -- Prepare the REPLY message to send
                                message_to_send <= MSG_REPLY & message_to_validate;
                                send_data_enable <= '1'; -- Enable data sending
                                current_phase <= P2; -- Move to next phase
                                timeout_send_counter <= timeout_send; -- Reset the timeout counter
                            WHEN P2 => -- The phase2 of the REPLY sending process is waiting for the bond_tx to be busy
                                IF send_busy = '1' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    current_phase <= P3; -- Move to next phase if bond_tx is busy
                                    timeout_send_counter <= timeout_send; -- Reset the timeout counter
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during REPLY sending due to timeout
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P3 => -- The phase3 of the REPLY sending process is waiting for the bond_tx to finish sending
                                IF send_busy = '0' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    start_wait_for_ack <= '1'; -- Reset wait_for_ack signal
                                    send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                    current_phase <= P4; -- Reset the phase to P4
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during REPLY sending due to timeout
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P4 =>
                                start_wait_for_ack <= '0';
                                current_phase <= P1;
                                send_data_enable <= '0';
                                send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                        END CASE;
                    WHEN SEND =>
                        -- In SEND state, we are sending the message
                        CASE current_phase IS
                            WHEN P1 => -- The phase1 of the SEND process is preparing the message to send
                                message_to_send <= MSG_DATA & s_message; -- Load the message to send
                                send_data_enable <= '1'; -- Enable data sending
                                current_phase <= P2; -- Move to next phase
                                timeout_send_counter <= timeout_send; -- Reset the timeout counter
                            WHEN P2 => -- The phase2 of the SEND process is waiting for the bond_tx to be busy
                                IF send_busy = '1' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    current_phase <= P3; -- Move to next phase if bond_tx is busy
                                    timeout_send_counter <= timeout_send; -- Reset the timeout counter
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during sending due to timeout
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P3 => -- The phase3 of the SEND process is waiting for the bond_tx to finish sending
                                IF send_busy = '0' THEN
                                    send_data_enable <= '0'; -- Disable data sending
                                    start_wait_for_reply <= '1'; -- Reset wait_for_reply signal
                                    message_waiting_for_reply <= s_message; -- Store the message that was sent
                                    current_phase <= P4; -- Reset the phase to P4
                                ELSE
                                    IF timeout_send_counter = 0 THEN
                                        send_data_enable <= '0'; -- Disable data sending
                                        s_error <= '1'; -- Indicate that an error occurred during sending due to timeout
                                        send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                        current_phase <= P1; -- Reset the phase to P1
                                        sending <= '0'; -- Reset the busy signal
                                    ELSE
                                        timeout_send_counter <= timeout_send_counter - 1;
                                    END IF;
                                END IF;
                            WHEN P4 =>
                                start_wait_for_reply <= '0';
                                send_data_enable <= '0';
                                send_state <= IDLE; -- Go back to IDLE state if bond_tx is not busy
                                current_phase <= P1;
                        END CASE;
                END CASE;
            END IF;
        END IF;
    END PROCESS send_process;

    -- The main process for receiving data
    receive_process : PROCESS (clk, reset)
    BEGIN
        IF reset = '1' THEN
            r_message <= (OTHERS => '0');
            receiving <= '0';
            r_valid <= '0';
            r_error <= '0';
            message_to_validate <= (OTHERS => '0');
            ack_send_needed <= '0';
            reply_send_needed <= '0';
            stop_wait_for_reply <= '0';
            stop_wait_for_ack <= '0';
        ELSIF rising_edge(clk) THEN
            IF timeout_ack = '1' OR receiving_timeout = '1' THEN
                r_error <= '1'; -- Indicate that an error occurred during transmission due to ack timeout
                receiving <= '0'; -- Reset the busy signal
                r_valid <= '0'; -- Reset the valid signal
                ack_send_needed <= '0'; -- Reset the ack_send_needed signal
                recv_state <= RECV_IDLE; -- Go back to IDLE state
                stop_wait_for_ack <= '1'; -- Stop waiting for an ack
            ELSIF timeout_reply = '1' THEN
                stop_wait_for_reply <= '1'; -- Stop waiting for a reply
            ELSIF reply_send_needed_reset = '1' THEN
                reply_send_needed <= '0'; -- Reset the reply_send_needed signal
            ELSIF ack_send_needed_reset = '1' THEN
                ack_send_needed <= '0'; -- Reset the ack_send_needed signal
            ELSE
                CASE recv_state IS
                    WHEN RECV_IDLE =>
                        stop_wait_for_reply <= '0'; -- Reset the stop_wait_for_reply signal
                        stop_wait_for_ack <= '0'; -- Reset the stop_wait_for_ack signal
                        IF receive_busy = '1' THEN
                            recv_state <= RECV_WAIT;
                        END IF;
                    WHEN RECV_WAIT =>
                        IF receive_busy = '0' THEN
                            IF message_valid = '1' THEN
                                message_to_examine <= message_received;
                                recv_state <= RECV_PROCESS;
                            ELSIF receive_failed = '1' THEN
                                recv_state <= RECV_FAIL;
                            END IF;
                        END IF;
                    WHEN RECV_FAIL =>
                        recv_state <= RECV_IDLE; -- Reset to IDLE state after failure
                    WHEN RECV_PROCESS =>
                        -- Process the received message
                        CASE recv_header IS
                            WHEN MSG_DATA =>
                                IF receiving = '1' THEN -- Already busy receiving a message - cannot accept new message this has to be an error
                                    r_error <= '1'; -- Set the r_error signal to indicate that an error occurred
                                    r_valid <= '0'; -- Reset the r_valid signal
                                    receiving <= '0'; -- Reset the receiving signal
                                ELSE
                                    message_to_validate <= message_to_examine(message_length - 1 DOWNTO 0); -- Extract the message part
                                    reply_send_needed <= '1'; -- Set the reply_send_needed signal to indicate that a reply is needed
                                    receiving <= '1'; -- Set the receiving signal to indicate that the component is busy receiving
                                    r_valid <= '0'; -- Reset the r_valid signal
                                    r_error <= '0'; -- Reset the r_error signal
                                END IF;
                            WHEN MSG_REPLY =>
                                IF message_to_examine(message_length - 1 DOWNTO 0) = message_waiting_for_reply THEN
                                    ack_result <= '1'; -- Set ack_result to indicate that the reply matches the sent message
                                ELSE
                                    ack_result <= '0'; -- Set ack_result to indicate that the reply does not match the sent message
                                END IF;
                                ack_send_needed <= '1'; -- Set the ack_send_needed signal to indicate that an ACK is needed
                                stop_wait_for_reply <= '1'; -- Stop waiting for a reply
                            WHEN MSG_ACK =>
                                if wait_for_ack = '0' THEN
                                    -- We are not waiting for an ack but we received one - this is an error
                                    r_error <= '1'; -- Set the r_error signal to indicate that an error occurred
                                    r_valid <= '0'; -- Reset the r_valid signal
                                    receiving <= '0'; -- Reset the receiving signal
                                ELSE
                                    r_message <= message_to_validate;
                                    r_valid <= '1'; -- Set the r_valid signal to indicate that a valid message has been received
                                    receiving <= '0'; -- Reset the receiving signal
                                    r_error <= '0'; -- Reset the r_error signal
                                    stop_wait_for_ack <= '1'; -- Stop waiting for an ack
                                END IF;
                            WHEN MSG_NACK =>
                                r_message <= (OTHERS => '0'); -- Reset the received message
                                r_valid <= '0'; -- Reset the r_valid signal
                                receiving <= '0'; -- Reset the receiving signal
                                r_error <= '1'; -- Set the r_error signal to indicate that an error occurred
                                stop_wait_for_ack <= '1'; -- Stop waiting for an ack
                            WHEN OTHERS =>
                                r_error <= '1'; -- Set the r_error signal to indicate that an error occurred due to unknown message type
                                r_valid <= '0'; -- Reset the r_valid signal
                                receiving <= '0'; -- Reset the receiving signal
                        END CASE;
                        recv_state <= RECV_IDLE; -- Reset to IDLE state after processing
                END CASE;
            END IF;
        END IF;
    END PROCESS receive_process;

END Behavioral;
`
)
