package bondmachine

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/udpbond"
)

type NetParameters map[string]string

type Udpbond_extra struct {
	Config    *udpbond.Config
	Cluster   *udpbond.Cluster
	Ips       *udpbond.Ips
	PeerID    uint32
	Maps      *IOmap
	Flavor    string
	Ip        string
	Broadcast string
	Netmask   string
	Port      string
	NetParams *NetParameters
}

type FirmwareCommand struct {
	PrimaryState   string // SM state for the command
	SecondaryState string // state for eventually internal State machines
	Command        string // ASCII command
	Description    string
	VHDLRapp       string // VHDL hex
	Signal         string // Eventually associated signal
	Starting       int    // Position of start within the memory
	OmitReturn     bool
	Payload        []string // Payload to transmit
	Payload_relpos []int    // Pauload relative position
}

func Ascii2Hex(in string) string {
	encoded := ""
	remaining_string := in

	for remaining_string != "" {

		lidx := strings.Index(remaining_string, "<<<")

		if lidx == -1 {
			tenc := []byte(remaining_string)
			encoded += hex.EncodeToString(tenc)
			remaining_string = ""
		} else {
			if lidx != 0 {
				tenc := []byte(remaining_string[0:lidx])
				encoded += hex.EncodeToString(tenc)
			}

			ridx := strings.Index(remaining_string, ">>>")
			if lidx != -1 {
				encoded += remaining_string[lidx+3 : ridx]
				remaining_string = remaining_string[ridx+3:]
			}
		}
	}
	return encoded
}

func CompleteCommands(fc []FirmwareCommand) ([]FirmwareCommand, int) {
	result := make([]FirmwareCommand, len(fc))
	poscounter := 0

	statelen := make(map[string]int)

	for i, com := range fc {
		result[i].PrimaryState = com.PrimaryState
		result[i].SecondaryState = com.SecondaryState
		result[i].Command = com.Command
		result[i].Description = com.Description
		result[i].Signal = com.Signal
		result[i].OmitReturn = com.OmitReturn
		result[i].Starting = poscounter
		result[i].Payload = append([]string(nil), com.Payload...)
		result[i].Payload_relpos = append([]int(nil), com.Payload_relpos...)

		result[i].VHDLRapp = ""

		encoded := ""
		remaining_string := com.Command

		for remaining_string != "" {

			lidx := strings.Index(remaining_string, "<<<")

			if lidx == -1 {
				tenc := []byte(remaining_string)
				encoded += hex.EncodeToString(tenc)
				remaining_string = ""
			} else {
				if lidx != 0 {
					tenc := []byte(remaining_string[0:lidx])
					encoded += hex.EncodeToString(tenc)
				}

				ridx := strings.Index(remaining_string, ">>>")
				if lidx != -1 {
					encoded += remaining_string[lidx+3 : ridx]
					remaining_string = remaining_string[ridx+3:]
				}
			}
		}

		if !com.OmitReturn {
			encoded += "0d0a"
		}

		encoded += "ff"

		for j := 0; j < len(encoded); j = j + 2 {
			result[i].VHDLRapp += "x\"" + encoded[j:j+2] + "\","
		}

		poscounter += len(encoded) / 2

		pmstate := com.PrimaryState
		if sl, ok := statelen[pmstate]; ok {
			statelen[pmstate] = sl + 1
		} else {
			statelen[pmstate] = 1
		}

	}

	lencounter := make(map[string]int)

	for i, com := range fc {

		pmstate := com.PrimaryState
		nbits := Needed_bits(statelen[pmstate])

		if st, ok := lencounter[pmstate]; !ok {
			result[i].SecondaryState = zeros_prefix(nbits, get_binary(0))
			lencounter[pmstate] = 1
		} else {
			result[i].SecondaryState = zeros_prefix(nbits, get_binary(st))
			lencounter[pmstate] = st + 1
		}
	}
	return result, poscounter
}

func LocateCommand(fc []FirmwareCommand, ps string, ss string) int {
	for _, com := range fc {
		if com.PrimaryState == ps {
			if com.SecondaryState == ss {
				return com.Starting
			}
		}
	}
	return -1
}

func LocateCommandbyIndex(fc []FirmwareCommand, ps string, index int) (FirmwareCommand, bool) {
	i := 0
	for _, com := range fc {
		if com.PrimaryState == ps {
			if i == index {
				return com, true
			}
			i++
		}
	}
	return FirmwareCommand{}, false
}

func (sl *Udpbond_extra) Get_Name() string {
	return "udpbond"
}

func (sl *Udpbond_extra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)
	result.Params["peer_id"] = strconv.Itoa(int(sl.PeerID))
	result.Params["cluster_id"] = strconv.Itoa(int(sl.Cluster.ClusterId))
	result.Params["ip"] = sl.Ip
	result.Params["broadcast"] = sl.Broadcast
	result.Params["port"] = sl.Port

	netparams := *sl.NetParams

	if ssid, ok := netparams["ssid"]; ok {
		result.Params["ssid"] = ssid
	}
	if netmask, ok := netparams["netmask"]; ok {
		result.Params["netmask"] = netmask
	}
	if gateway, ok := netparams["gateway"]; ok {
		result.Params["gateway"] = gateway
	}

	var mypeer udpbond.Peer

	for _, peer := range sl.Cluster.Peers {
		if peer.PeerId == sl.PeerID {
			mypeer = peer
		}

		if sl.Ips != nil {
			peerstr := strconv.Itoa(int(peer.PeerId))
			if ipaddr, ok := sl.Ips.Assoc["peer_"+peerstr]; ok {
				if ipaddr == "auto" {
					result.Params["peer_"+peerstr+"_ip"] = "auto"
				} else if ipaddr == "adv" {
					result.Params["peer_"+peerstr+"_ip"] = "auto"
				} else {
					result.Params["peer_"+peerstr+"_ip"] = ipaddr
				}
			} else {
				result.Params["peer_"+peerstr+"_ip"] = "auto"
			}
		}

	}

	result.Params["input_ids"] = ""
	result.Params["inputs"] = ""
	result.Params["sources"] = ""

	for _, inp := range mypeer.Inputs {
		for iname, ival := range sl.Maps.Assoc {
			if iname[0] == 'i' && ival == strconv.Itoa(int(inp)) {
				result.Params["input_ids"] += "," + ival
				result.Params["inputs"] += "," + iname

				ressource := ""
				for _, opeer := range sl.Cluster.Peers {
					for _, oout := range opeer.Outputs {
						if strconv.Itoa(int(oout)) == ival {
							ressource = strconv.Itoa(int(opeer.PeerId))
							break
						}
					}
				}
				if ressource != "" {
					result.Params["sources"] += "," + ressource
				}

			}
		}
	}

	if result.Params["input_ids"] != "" {
		result.Params["input_ids"] = result.Params["input_ids"][1:len(result.Params["input_ids"])]
		result.Params["inputs"] = result.Params["inputs"][1:len(result.Params["inputs"])]
		result.Params["sources"] = result.Params["sources"][1:len(result.Params["sources"])]
	}

	result.Params["output_ids"] = ""
	result.Params["outputs"] = ""
	// Comma separated and - separated list of peer ids
	result.Params["destinations"] = ""

	for _, outp := range mypeer.Outputs {
		for oname, oval := range sl.Maps.Assoc {
			if oname[0] == 'o' && oval == strconv.Itoa(int(outp)) {
				result.Params["output_ids"] += "," + oval
				result.Params["outputs"] += "," + oname

				resdest := ""
				for _, ipeer := range sl.Cluster.Peers {
					for _, iin := range ipeer.Inputs {
						//fmt.Println(ipeer.PeerId, iin, oval, strconv.Itoa(int(iin)))
						if strconv.Itoa(int(iin)) == oval {
							resdest += "-" + strconv.Itoa(int(ipeer.PeerId))
						}
					}
				}
				//fmt.Println("resdest", resdest)
				if resdest != "" {
					result.Params["destinations"] += "," + resdest[1:len(resdest)]
				}

			}
		}
	}

	if result.Params["output_ids"] != "" {
		result.Params["output_ids"] = result.Params["output_ids"][1:len(result.Params["output_ids"])]
		result.Params["outputs"] = result.Params["outputs"][1:len(result.Params["outputs"])]
		result.Params["destinations"] = result.Params["destinations"][1:len(result.Params["destinations"])]
	}

	return result
}

func (sl *Udpbond_extra) Import(inp string) error {
	return nil
}

func (sl *Udpbond_extra) Export() string {
	return ""
}

func (sl *Udpbond_extra) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *Udpbond_extra) Verilog_headers() string {
	result := "\n"
	return result
}
func (sl *Udpbond_extra) StaticVerilog() string {

	result := "\n"
	return result
}

func (sl *Udpbond_extra) ExtraFiles() ([]string, []string) {
	rsize := int(sl.Config.Rsize)
	udpbond_params := sl.Get_Params().Params
	fmt.Println(udpbond_params)

	payload_fractions := rsize / 8

	if rsize%8 != 0 {
		payload_fractions++
	}

	intclusid, _ := strconv.Atoi(udpbond_params["cluster_id"])
	hexclusid := fmt.Sprintf("%08x", intclusid)

	intpeerid, _ := strconv.Atoi(udpbond_params["peer_id"])
	hexpeerid := fmt.Sprintf("%08x", intpeerid)

	result := ""

	result += "------------------ VHDL component for the esc8266 wireless chip -----------\n"
	result += "library IEEE;\n"
	result += "use IEEE.STD_LOGIC_1164.ALL;\n"
	result += "use IEEE.NUMERIC_STD.ALL;\n"
	result += "\n"
	// Entry points to the bondmachine
	result += "entity udpbond_main is\n"
	result += "     Port ( clk100         : in  STD_LOGIC;\n"
	result += "            reset          : in  STD_LOGIC;\n"
	result += "            wifi_enable    : out STD_LOGIC;\n"
	result += "            wifi_rx        : in  STD_LOGIC;\n"
	result += "            wifi_tx        : out STD_LOGIC;\n"
	for _, iname := range strings.Split(udpbond_params["inputs"], ",") {
		if iname != "" {
			result += "            input_" + iname + "       : out STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := (others => '0');\n"
			result += "            input_" + iname + "_dv    : out STD_LOGIC;\n"
			result += "            input_" + iname + "_recv  : in  STD_LOGIC;\n"
		}
	}
	for _, oname := range strings.Split(udpbond_params["outputs"], ",") {
		if oname != "" {
			result += "            output_" + oname + "      : in  STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := (others => '0');\n"
			result += "            output_" + oname + "_dv   : in  STD_LOGIC;\n"
			result += "            output_" + oname + "_recv : out STD_LOGIC;\n"
		}
	}
	result = result[0 : len(result)-2]
	result += "           );\n"
	result += "end udpbond_main;\n"
	result += "\n"
	result += "architecture Behavioral of udpbond_main is\n"
	result += "\n"
	result += "    component esp8266_driver is\n"
	result += "    Port ( clk100           : in  STD_LOGIC;\n"
	//result += "           -- roba da lasciare per il debug --\n"
	result += "           powerdown        : in  STD_LOGIC;\n"
	//result += "           otherpayload     : in  STD_LOGIC;\n"
	//result += "           status_active    : out STD_LOGIC;\n"
	//result += "           status_wifi_up   : out STD_LOGIC;\n"
	//result += "           status_connected : out STD_LOGIC;\n"
	//result += "           status_sending   : out STD_LOGIC;\n"
	//result += "           status_receiving : out STD_LOGIC;\n"
	//result += "           status_error     : out STD_LOGIC;\n"
	//counter := 0
	for _, iname := range strings.Split(udpbond_params["inputs"], ",") {
		if iname != "" {
			result += "           payload_" + iname + "       : out STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := x\"00\";\n"
			result += "           payload_" + iname + "_dv    : out STD_LOGIC;\n"
			result += "           payload_" + iname + "_recv  : in  STD_LOGIC;\n"
			//counter++
		}
	}
	for _, oname := range strings.Split(udpbond_params["outputs"], ",") {
		if oname != "" {
			result += "           payload_" + oname + "       : in  STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := x\"00\";\n"
			result += "           payload_" + oname + "_dv    : in  STD_LOGIC;\n"
			result += "           payload_" + oname + "_recv  : out STD_LOGIC;\n"
			//counter++
		}
	}
	result += "           wifi_enable      : out STD_LOGIC;\n"
	result += "           wifi_rx          : in  STD_LOGIC;\n"
	result += "           wifi_tx          : out STD_LOGIC);\n"
	result += "    end component;\n"
	result += "\n"
	//result += "    -- questi char non ho nai capito a cosa servono e vengono solo menzionati in questa parte di codice --\n"
	//result += "    signal char0 : std_logic_vector(7 downto 0) := x\"00\";\n"
	//result += "    signal char1 : std_logic_vector(7 downto 0) := x\"00\";\n"
	//result += "    signal char2 : std_logic_vector(7 downto 0) := x\"00\";\n"
	//result += "    signal char3 : std_logic_vector(7 downto 0) := x\"00\";\n"
	result += "begin\n"
	result += "\n"
	result += "i_esp8226: esp8266_driver Port map (\n"
	result += "           clk100           => clk100,\n"
	result += "           powerdown        => reset,\n"
	//result += "         --  status_active    => led(0),\n"
	//result += "         --  status_wifi_up   => led(1),\n"
	//result += "         --  status_connected => led(2),\n"
	//result += "         --  status_sending   => led(3),\n"
	//result += "         --  status_receiving => led(4),\n"
	//result += "         --  status_error     => led(5),\n"
	for _, iname := range strings.Split(udpbond_params["inputs"], ",") {
		if iname != "" {
			result += "           payload_" + iname + "       => input_" + iname + ",\n"
			result += "           payload_" + iname + "_dv    => input_" + iname + "_dv,\n"
			result += "           payload_" + iname + "_recv  => input_" + iname + "_recv,\n"
		}
	}
	for _, oname := range strings.Split(udpbond_params["outputs"], ",") {
		if oname != "" {
			result += "           payload_" + oname + "       => output_" + oname + ",\n"
			result += "           payload_" + oname + "_dv    => output_" + oname + "_dv,\n"
			result += "           payload_" + oname + "_recv  => output_" + oname + "_recv,\n"
		}
	}
	result += "           wifi_enable      => wifi_enable,\n"
	result += "           wifi_rx          => wifi_rx,\n"
	result += "           wifi_tx          => wifi_tx);\n"
	result += "\n"
	result += "end Behavioral;\n"
	result += "\n"
	result += "--------------------------------------------------\n"
	result += "-- esp8266_driver - Session setup and sending\n"
	result += "--                  packets of data using ESP8266\n"
	result += "--\n"
	result += "-- Author: Mike Field <hamster@snap.net.nz>\n"
	result += "--\n"
	result += "-- NOTE: You will need to edit the constants to put\n"
	result += "-- your own SSID & password, and the IP address\n"
	result += "-- and destination port number\n"
	result += "--\n"
	result += "-- This also has a watchdog, that resets the state\n"
	result += "-- of the design if no state change has occurred \n"
	result += "-- in the last 10 seconds.\n"
	result += "------------------------------------------------\n"
	result += "\n"
	result += "library IEEE;\n"
	result += "use IEEE.STD_LOGIC_1164.ALL;\n"
	result += "use IEEE.NUMERIC_STD.ALL;\n"
	result += "\n"
	result += "entity esp8266_driver is\n"
	result += "    Port ( clk100           : in  STD_LOGIC;\n"
	result += "           powerdown        : in  STD_LOGIC;\n"
	//result += "           status_active    : out STD_LOGIC := '0';\n"
	//result += "           status_wifi_up   : out STD_LOGIC := '0';\n"
	//result += "           status_connected : out STD_LOGIC := '0';\n"
	//result += "           status_sending   : out STD_LOGIC := '0';\n"
	//result += "           status_receiving : out STD_LOGIC := '0';\n"
	//result += "           status_error     : out STD_LOGIC := '0';\n"
	//result += "           \n"
	for _, iname := range strings.Split(udpbond_params["inputs"], ",") {
		if iname != "" {
			result += "           payload_" + iname + "       : out STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := x\"00\";\n"
			result += "           payload_" + iname + "_dv    : out STD_LOGIC;\n"
			result += "           payload_" + iname + "_recv  : in  STD_LOGIC;\n"
		}
	}
	for _, oname := range strings.Split(udpbond_params["outputs"], ",") {
		if oname != "" {
			result += "           payload_" + oname + "       : in  STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := x\"00\";\n"
			result += "           payload_" + oname + "_dv    : in  STD_LOGIC;\n"
			result += "           payload_" + oname + "_recv  : out STD_LOGIC;\n"
		}
	}
	result += "           wifi_enable      : out STD_LOGIC;\n"
	result += "           wifi_rx          : in  STD_LOGIC;\n"
	result += "           wifi_tx          : out STD_LOGIC);\n"
	result += "end esp8266_driver;\n"
	result += "\n"

	// Now lets build the firmare commands
	commands := make([]FirmwareCommand, 0)
	var command FirmwareCommand

	// command = FirmwareCommand{PrimaryState: "", SecondaryState: "", Command: "", VHDLRapp: "", Signal: "", Starting: -1}

	command = FirmwareCommand{PrimaryState: "00011", Command: "AT+CWMODE=3", Starting: -1, Description: "change mode"}
	commands = append(commands, command)

	command = FirmwareCommand{PrimaryState: "00100", Command: "AT+CIPMUX=1", Starting: -1, Description: "single socket mode"}
	commands = append(commands, command)

	command = FirmwareCommand{PrimaryState: "00101", Command: "AT+CWJAP=\"" + udpbond_params["ssid"] + "\",\"" + udpbond_params["password"] + "\"", Starting: -1, Description: "Connect to WIFI"}
	commands = append(commands, command)

	command = FirmwareCommand{PrimaryState: "11100", Command: "AT+CIPSTA=\"" + udpbond_params["ip"] + "\",\"" + udpbond_params["gateway"] + "\",\"" + udpbond_params["netmask"] + "\"", Starting: -1, Description: "Assing IP address"}
	commands = append(commands, command)

	// Start of the connection composing

	// Receiving socket
	command = FirmwareCommand{PrimaryState: "01110", Command: "AT+CIPSTART=0,\"UDP\",\"0.0.0.0\"," + udpbond_params["port"] + "," + udpbond_params["port"], Starting: -1, Description: "Open UDP receiving socket"}
	commands = append(commands, command)

	// Multicast socket
	command = FirmwareCommand{PrimaryState: "01101", Command: "AT+CIPSTART=1,\"UDP\",\"" + udpbond_params["broadcast"] + "\"," + udpbond_params["port"], Starting: -1, Description: "Multicast socket"}
	commands = append(commands, command)

	// Keep track to witch connection map every destination
	dest2conn := make(map[string]string)

	{
		done := make(map[string]bool)
		conn := 2

		// Find the unique destinations and set a CISSTART for each one
		if udpbond_params["destinations"] != "" {
			for _, destlist := range strings.Split(udpbond_params["destinations"], ",") {
				for _, dest := range strings.Split(destlist, "-") {
					if _, ok := done[dest]; !ok {
						done[dest] = true

						var peerid string
						var peerip string
						var peerport string

						if pid, ok := udpbond_params["peer_"+dest+"_ip"]; ok {
							peerid = dest
							peerip = strings.Split(pid, "/")[0]
							peerport = strings.Split(pid, ":")[1]
						} else {
							break
						}

						command = FirmwareCommand{PrimaryState: "00110", Command: "AT+CIPSTART=" + strconv.Itoa(conn) + ",\"UDP\",\"" + peerip + "\"," + peerport, Starting: -1, Description: "Start an UDP connection to peer " + peerid}
						commands = append(commands, command)

						dest2conn[dest] = strconv.Itoa(conn)

						conn++
					}
				}
			}
		}

		// Find the received destination and set a CISSTART for them (if not already settend from the previous
		if udpbond_params["sources"] != "" {
			for _, source := range strings.Split(udpbond_params["sources"], ",") {
				if _, ok := done[source]; !ok {
					done[source] = true

					var peerid string
					var peerip string
					var peerport string

					if pid, ok := udpbond_params["peer_"+source+"_ip"]; ok {
						peerid = source
						peerip = strings.Split(pid, "/")[0]
						peerport = strings.Split(pid, ":")[1]
					} else {
						break
					}

					command = FirmwareCommand{PrimaryState: "00110", Command: "AT+CIPSTART=" + strconv.Itoa(conn) + ",\"UDP\",\"" + peerip + "\"," + peerport, Starting: -1, Description: "Start an UDP connection to peer " + peerid}
					commands = append(commands, command)

					dest2conn[source] = strconv.Itoa(conn)

					conn++
				}
			}
		}
	}

	// Send packets mapped to 2 state with the same SecondaryState

	// Processing Broadcasts
	command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=1,13", Starting: -1, Description: "Broadcast ADV_CLU cipsend", Signal: "broadcast_ready_adv_clu"}
	commands = append(commands, command)

	command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888801" + hexclusid + hexpeerid + ">>>", Starting: -1, Description: "Broadcast ADV_CLU message", Signal: "broadcast_ready_adv_clu"}
	commands = append(commands, command)

	if udpbond_params["input_ids"] != "" {
		for i, resid := range strings.Split(udpbond_params["input_ids"], ",") {

			residint, _ := strconv.Atoi(resid)
			residhex := fmt.Sprintf("%08x", residint)

			command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=1,17", Starting: -1, Description: "Broadcast ADV_IN cipsend " + strconv.Itoa(i), Signal: "broadcast_ready_adv_in_" + strconv.Itoa(i)}
			commands = append(commands, command)

			command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888803" + hexclusid + hexpeerid + residhex + ">>>", Starting: -1, Description: "Broadcast ADV_IN message " + strconv.Itoa(i), Signal: "broadcast_ready_adv_in_" + strconv.Itoa(i)}
			commands = append(commands, command)
		}
	}

	if udpbond_params["output_ids"] != "" {
		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			residint, _ := strconv.Atoi(resid)
			residhex := fmt.Sprintf("%08x", residint)

			command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=1,17", Starting: -1, Description: "Broadcast ADV_OUT cipsend " + strconv.Itoa(i), Signal: "broadcast_ready_adv_out_" + strconv.Itoa(i)}
			commands = append(commands, command)

			command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888804" + hexclusid + hexpeerid + residhex + ">>>", Starting: -1, Description: "Broadcast ADV_OUT message " + strconv.Itoa(i), Signal: "broadcast_ready_adv_out_" + strconv.Itoa(i)}
			commands = append(commands, command)
		}
	}

	// Processing IO transfers
	if udpbond_params["output_ids"] != "" {

		dests := strings.Split(udpbond_params["destinations"], ",")
		outs := strings.Split(udpbond_params["outputs"], ",")

		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			oname := outs[i]
			destlist := dests[i]
			for _, dest := range strings.Split(destlist, "-") {

				residint, _ := strconv.Atoi(resid)
				residhex := fmt.Sprintf("%08x", residint)

				conn := dest2conn[dest]

				command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=" + conn + "," + strconv.Itoa(19+payload_fractions), Starting: -1, Description: "Send Payload of " + strconv.Itoa(residint), Signal: "io_tr_ready_" + strconv.Itoa(residint) + "_on_" + conn}
				commands = append(commands, command)

				payplaceholder := strings.Repeat("X", payload_fractions)

				command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888805>>>XXXX<<<" + hexclusid + hexpeerid + residhex + ">>>" + payplaceholder, Payload: []string{"tag", "payload_" + oname}, Payload_relpos: []int{3, 19}, Starting: -1, Description: "Send Payload of " + strconv.Itoa(residint), Signal: "io_tr_ready_" + strconv.Itoa(residint) + "_on_" + conn, OmitReturn: true}
				commands = append(commands, command)
			}
		}
	}

	// Processing ACKs
	// TODO

	// Processing data valid
	if udpbond_params["output_ids"] != "" {

		dests := strings.Split(udpbond_params["destinations"], ",")
		outs := strings.Split(udpbond_params["outputs"], ",")

		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			oname := outs[i]
			destlist := dests[i]
			for _, dest := range strings.Split(destlist, "-") {

				residint, _ := strconv.Atoi(resid)
				residhex := fmt.Sprintf("%08x", residint)

				conn := dest2conn[dest]

				command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=" + conn + "," + strconv.Itoa(20), Starting: -1, Description: "Send Payload data valid of " + strconv.Itoa(residint), Signal: "io_tr_dv_" + strconv.Itoa(residint) + "_on_" + conn}
				commands = append(commands, command)

				command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888806>>>XXXX<<<" + hexclusid + hexpeerid + residhex + ">>>X", Payload: []string{"tag", "payload_" + oname + "_dv"}, Payload_relpos: []int{3, 19}, Starting: -1, Description: "Send Payload data valid of " + strconv.Itoa(residint), Signal: "io_tr_dv_" + strconv.Itoa(residint) + "_on_" + conn, OmitReturn: true}
				commands = append(commands, command)
			}
		}
	}

	// Processing data received
	if udpbond_params["input_ids"] != "" {

		ins := strings.Split(udpbond_params["inputs"], ",")
		sources := strings.Split(udpbond_params["sources"], ",")

		for i, resid := range strings.Split(udpbond_params["input_ids"], ",") {

			iname := ins[i]
			dest := sources[i]

			residint, _ := strconv.Atoi(resid)
			residhex := fmt.Sprintf("%08x", residint)

			conn := dest2conn[dest]

			command = FirmwareCommand{PrimaryState: "00111", Command: "AT+CIPSEND=" + conn + "," + strconv.Itoa(20), Starting: -1, Description: "Send Payload data received of " + strconv.Itoa(residint), Signal: "io_tr_recv_" + strconv.Itoa(residint) + "_from_" + conn}
			commands = append(commands, command)

			command = FirmwareCommand{PrimaryState: "01000", Command: "<<<888807>>>XXXX<<<" + hexclusid + hexpeerid + residhex + ">>>X", Payload: []string{"tag", "payload_" + iname + "_recv"}, Payload_relpos: []int{3, 19}, Starting: -1, Description: "Send Payload data received of " + strconv.Itoa(residint), Signal: "io_tr_recv_" + strconv.Itoa(residint) + "_from_" + conn, OmitReturn: true}
			commands = append(commands, command)
		}
	}

	completedcommands, totall := CompleteCommands(commands)

	membits := Needed_bits(totall)

	//fmt.Println(commands)
	//fmt.Println(completedcommands, totall)

	//sockets_bits := Needed_bits(len("TODO"))
	var cipstart_loop_bits int
	cipstart_loop_present := true
	if cipstart_example, ok := LocateCommandbyIndex(completedcommands, "00110", 0); ok {
		cipstart_loop_bits = len(cipstart_example.SecondaryState)
	} else {
		cipstart_loop_present = false
	}

	var message_loop_bits int

	if send_example, ok := LocateCommandbyIndex(completedcommands, "00111", 0); ok {
		message_loop_bits = len(send_example.SecondaryState)
	}

	result += "architecture Behavioral of esp8266_driver is\n"
	result += "    type a_mem is array(0 to " + strconv.Itoa(totall-1) + ") of std_logic_vector(7 downto 0);\n"
	result += "    signal memory : a_mem := \n"
	result += "   (\n"
	for _, command = range completedcommands {
		result += "-- " + command.Description + " - Starting from " + strconv.Itoa(command.Starting) + " - Command \"" + command.Command + "\"\n"
		result += "\t" + command.VHDLRapp + "\n"
	}
	result = result[0:len(result)-2] + ");\n"

	result += "\n"

	// Old state for signal that has to go out via wireless: payloads, data valids, data received of incoming payloads

	for _, oname := range strings.Split(udpbond_params["outputs"], ",") {
		if oname != "" {
			//for i := 0; i < payload_fractions; i++ {
			//	result += "    signal payload_" + oname + "_f_" + strconv.Itoa(i) + "         : std_logic_vector(7 downto 0) := x\"00\";\n"
			//}
			result += "    signal payload_" + oname + "_old         : STD_LOGIC_VECTOR(" + strconv.Itoa(rsize-1) + " downto 0) := x\"00\";\n"
			result += "    signal payload_" + oname + "_dv_old      : STD_LOGIC;\n"
		}
	}

	result += "\n"

	for _, iname := range strings.Split(udpbond_params["inputs"], ",") {
		if iname != "" {
			result += "    signal payload_" + iname + "_recv_old      : STD_LOGIC;\n"
		}
	}

	result += "\n"

	// These are the received signals, one each destination
	if udpbond_params["output_ids"] != "" {

		dests := strings.Split(udpbond_params["destinations"], ",")

		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			destlist := dests[i]
			for _, dest := range strings.Split(destlist, "-") {

				conn := dest2conn[dest]

				residint, _ := strconv.Atoi(resid)

				cname := "io_tr_recv_" + strconv.Itoa(residint) + "_on_" + conn
				result += "    signal " + cname + "      : STD_LOGIC;\n"
			}
		}
	}

	result += "\n"
	result += "    signal current_char  : unsigned(" + strconv.Itoa(membits-1) + " downto 0)  := (others => '0');\n"
	result += "    signal delay_counter : unsigned(26 downto 0) := (others => '0');\n"
	result += "    signal in_delay      : std_logic             := '0';\n"
	result += "\n"
	result += "    component tx is\n"
	result += "        Port ( clk         : in  STD_LOGIC;\n"
	result += "               data        : in  STD_LOGIC_VECTOR (7 downto 0);\n"
	result += "               data_enable : in  STD_LOGIC;\n"
	result += "               busy        : out STD_LOGIC;\n"
	result += "               tx_out      : out STD_LOGIC);\n"
	result += "    end component;\n"
	result += "\n"
	result += "    component rx is\n"
	result += "        Port ( clk         : in  STD_LOGIC;\n"
	result += "               data        : out STD_LOGIC_VECTOR (7 downto 0);\n"
	result += "               data_enable : out STD_LOGIC;\n"
	result += "               rx_in       : in  STD_LOGIC);\n"
	result += "    end component;\n"
	result += "\n"
	result += "    signal tx_data        : std_logic_vector(7 downto 0 ) := (others => '0');\n"
	result += "    signal tx_busy        : std_logic                     := '0';\n"
	result += "    signal tx_data_enable : std_logic                     := '0';\n"
	result += "\n"
	result += "    signal rx_data        : std_logic_vector(7 downto 0 ) := (others => '0');\n"
	result += "    signal rx_data_enable : std_logic                     := '0';\n"
	result += "    signal sending        : std_logic                     := '0';\n"
	result += "    signal receiving      : std_logic                     := '0';\n"
	result += "    signal state          : std_logic_vector(4 downto 0 ) := (others => '0');\n"
	result += "    signal state_last     : std_logic_vector(4 downto 0 ) := (others => '0');\n"
	result += "    signal last_rx_chars  : std_logic_vector(" + strconv.Itoa(231+rsize) + " downto 0 ) := (others => '0');\n"
	result += "\n"
	result += "    signal rx_seeing_ok     : std_logic                     := '0';\n"
	result += "    signal rx_seeing_recv   : std_logic                     := '0';\n"
	result += "    signal rx_seeing_ready  : std_logic                     := '0';\n"
	result += "    signal rx_seeing_change : std_logic                     := '0';\n"
	result += "    signal rx_seeing_prompt : std_logic                     := '0';\n"
	result += "\n"

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "00111", currindex)

		if currcom.Signal != "" {

			result += "    signal " + currcom.Signal + "  : std_logic                     := '0';\n"
			result += "    signal reset_" + currcom.Signal + "  : std_logic                     := '0';\n"
		}

		if !currstat {
			break
		}
	}

	result += "\n"
	result += "    signal tag     : std_logic_vector(31 downto 0 ) := (others => '0');\n"
	result += "\n"

	result += "    signal message_loop     : std_logic_vector(" + strconv.Itoa(message_loop_bits-1) + " downto 0 ) := (others => '0');\n"

	if cipstart_loop_present {
		result += "    signal cipstart_loop    : std_logic_vector(" + strconv.Itoa(cipstart_loop_bits-1) + " downto 0 ) := (others => '0');\n"
	}
	result += "\n"
	result += "    -- Watchdog timer for recovery.\n"
	result += "    -- This has to count for the number of clock cycles in 1ms\n"
	result += "    signal watchdog_low      : unsigned(16 downto 0) := (others => '0');\n"
	result += "    -- This is high for one cycle every millisecond\n"
	result += "    signal inc_wd_high       : std_logic := '0';\n"
	result += "    -- This has to count to 9,9999, for the 10,000 ms timout\n"
	result += "    signal watchdog_high     : unsigned(13 downto 0) := (others => '0');\n"
	result += "    signal counter           : unsigned(12 downto 0) := (others => '0');\n"
	result += "    signal bcounter          : unsigned(11 downto 0) := (others => '0');\n"
	result += "begin\n"
	result += "\n"
	result += "i_tx : tx port map (\n"
	result += "	clk         => clk100,\n"
	result += "	data        => tx_data,\n"
	result += "	data_enable => tx_data_enable,\n"
	result += "	busy        => tx_busy,\n"
	result += "	tx_out      => wifi_tx);\n"
	result += "\n"
	result += "i_rx : rx  port map (\n"
	result += "	clk         => clk100,\n"
	result += "	data        => rx_data,\n"
	result += "	data_enable => rx_data_enable,\n"
	result += "	rx_in       => wifi_rx);\n"
	result += "\n"

	// These are the received signals, one each destination merged in the output received
	if udpbond_params["output_ids"] != "" {

		dests := strings.Split(udpbond_params["destinations"], ",")
		outs := strings.Split(udpbond_params["outputs"], ",")

		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			oname := outs[i]
			destlist := dests[i]

			result += "merge_payload_" + oname + ": process(clk100)\n"
			result += "    begin\n"
			result += "    if rising_edge(clk100) then\n"
			result += "        payload_" + oname + "_recv <= "

			for j, dest := range strings.Split(destlist, "-") {

				conn := dest2conn[dest]

				residint, _ := strconv.Atoi(resid)

				cname := "io_tr_recv_" + strconv.Itoa(residint) + "_on_" + conn
				if j == 0 {
					result += cname
				} else {
					result += " and " + cname
				}
			}

			result += ";\n"

			result += "    end if;\n"
			result += "    end process;\n"
			result += "\n"

		}
	}

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "00111", currindex)

		// Broadcasts triggerred avery 5s
		if strings.HasPrefix(currcom.Signal, "broadcast_ready_adv") {
			result += "send_" + currcom.Signal + ": process(clk100)\n"
			result += "    begin\n"
			result += "    if rising_edge(clk100) then\n"
			result += "        if reset_" + currcom.Signal + " = '1' then\n"
			result += "            " + currcom.Signal + " <= '0'; \n"
			result += "        else\n"
			result += "            if counter = 0 then\n"
			result += "                " + currcom.Signal + " <= '1';\n"
			result += "            end if;\n"
			result += "        end if;\n"
			result += "    end if;\n"
			result += "    end process;\n"
			result += "\n"
		}

		if !currstat {
			break
		}
	}

	payloads_done := make(map[string]bool)

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "01000", currindex)

		if strings.HasPrefix(currcom.Signal, "io_tr_ready_") {
			asspay := ""
			for _, pay := range currcom.Payload {
				if strings.HasPrefix(pay, "payload_") {
					asspay = pay
					break
				}
			}
			result += "send_" + currcom.Signal + ": process(clk100)\n"
			result += "    begin\n"
			result += "    if rising_edge(clk100) then\n"
			result += "        if reset_" + currcom.Signal + " = '1' then\n"
			result += "            " + currcom.Signal + " <= '0'; \n"
			result += "        else\n"
			result += "            if bcounter = 0 then\n"
			result += "                " + currcom.Signal + " <= '1';\n"
			result += "            else\n"
			result += "                if " + asspay + "_old /= " + asspay + " then\n"
			result += "                     " + currcom.Signal + " <= '1';\n"
			result += "                end if;\n"
			result += "            end if;\n"
			result += "        end if;\n"
			result += "    end if;\n"
			result += "    end process;\n"
			result += "    \n"

			if _, ok := payloads_done[asspay]; !ok {
				payloads_done[asspay] = true

				result += "check_" + asspay + ": process(clk100)\n"
				result += "    begin\n"
				result += "    if rising_edge(clk100) then\n"
				result += "             " + asspay + "_old <= " + asspay + ";\n"
				result += "    end if;\n"
				result += "    end process;\n"
				result += "    \n"
			}
		}
		if !currstat {
			break
		}
	}

	payloads_done = make(map[string]bool)

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "01000", currindex)

		if strings.HasPrefix(currcom.Signal, "io_tr_dv_") {
			asspay := ""
			for _, pay := range currcom.Payload {
				if strings.HasPrefix(pay, "payload_") {
					asspay = pay
					break
				}
			}
			result += "send_" + currcom.Signal + ": process(clk100)\n"
			result += "    begin\n"
			result += "    if rising_edge(clk100) then\n"
			result += "        if reset_" + currcom.Signal + " = '1' then\n"
			result += "            " + currcom.Signal + " <= '0'; \n"
			result += "        else\n"
			result += "            if bcounter = 0 then\n"
			result += "                " + currcom.Signal + " <= '1';\n"
			result += "            else\n"
			result += "                if " + asspay + "_old /= " + asspay + " then\n"
			result += "                     " + currcom.Signal + " <= '1';\n"
			result += "                end if;\n"
			result += "            end if;\n"
			result += "        end if;\n"
			result += "    end if;\n"
			result += "    end process;\n"
			result += "    \n"

			if _, ok := payloads_done[asspay]; !ok {
				payloads_done[asspay] = true

				result += "check_" + asspay + ": process(clk100)\n"
				result += "    begin\n"
				result += "    if rising_edge(clk100) then\n"
				result += "             " + asspay + "_old <= " + asspay + ";\n"
				result += "    end if;\n"
				result += "    end process;\n"
				result += "    \n"
			}
		}
		if !currstat {
			break
		}
	}

	payloads_done = make(map[string]bool)

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "01000", currindex)

		if strings.HasPrefix(currcom.Signal, "io_tr_recv_") {
			asspay := ""
			for _, pay := range currcom.Payload {
				if strings.HasPrefix(pay, "payload_") {
					asspay = pay
					break
				}
			}
			result += "send_" + currcom.Signal + ": process(clk100)\n"
			result += "    begin\n"
			result += "    if rising_edge(clk100) then\n"
			result += "        if reset_" + currcom.Signal + " = '1' then\n"
			result += "            " + currcom.Signal + " <= '0'; \n"
			result += "        else\n"
			result += "            if bcounter = 0 then\n"
			result += "                " + currcom.Signal + " <= '1';\n"
			result += "            else\n"
			result += "                if " + asspay + "_old /= " + asspay + " then\n"
			result += "                     " + currcom.Signal + " <= '1';\n"
			result += "                end if;\n"
			result += "            end if;\n"
			result += "        end if;\n"
			result += "    end if;\n"
			result += "    end process;\n"
			result += "    \n"

			if _, ok := payloads_done[asspay]; !ok {
				payloads_done[asspay] = true

				result += "check_" + asspay + ": process(clk100)\n"
				result += "    begin\n"
				result += "    if rising_edge(clk100) then\n"
				result += "             " + asspay + "_old <= " + asspay + ";\n"
				result += "    end if;\n"
				result += "    end process;\n"
				result += "    \n"
			}
		}
		if !currstat {
			break
		}
	}

	result += "send_chars: process(clk100)\n"
	result += "    begin\n"
	result += "       if rising_edge(clk100) then\n"
	result += "            tx_data_enable <= '0'; -- Default to no character being sent\n"
	result += "                if sending = '1' then\n"
	result += "                if memory(to_integer(current_char)) = x\"FF\" then\n"
	result += "                    sending <= '0';\n"
	result += "                elsif tx_busy ='0' then\n"
	result += "                    -- Send the next character\n"

	firstchar := true
	for _, command = range completedcommands {
		starpos := command.Starting
		for i, paystart := range command.Payload_relpos {
			pay := command.Payload[i]
			if pay == "tag" {
				for j := 0; j < 4; j++ {
					fractstart := 31 - j*8
					if firstchar {
						firstchar = false
						result += "                    if current_char = " + strconv.Itoa(starpos+paystart+j) + " then\n"
					} else {
						result += "                    elsif current_char = " + strconv.Itoa(starpos+paystart+j) + " then\n"
					}
					result += "                        tx_data <= tag(" + strconv.Itoa(fractstart) + " downto " + strconv.Itoa(fractstart-7) + ");\n"
				}
			} else if strings.HasPrefix(pay, "payload_") {
				if strings.HasSuffix(pay, "_dv") || strings.HasSuffix(pay, "_recv") {
					if firstchar {
						firstchar = false
						result += "                    if current_char = " + strconv.Itoa(starpos+paystart) + " then\n"
					} else {
						result += "                    elsif current_char = " + strconv.Itoa(starpos+paystart) + " then\n"
					}
					result += "                        tx_data <= " + pay + strings.Repeat(" & '0' ", 7) + ";\n"

				} else {
					for j := 0; j < payload_fractions; j++ {
						fractstart := rsize - 1 - j*8
						if firstchar {
							firstchar = false
							result += "                    if current_char = " + strconv.Itoa(starpos+paystart+j) + " then\n"
						} else {
							result += "                    elsif current_char = " + strconv.Itoa(starpos+paystart+j) + " then\n"
						}
						if fractstart-7 >= 0 {
							result += "                        tx_data <= " + pay + "(" + strconv.Itoa(fractstart) + " downto " + strconv.Itoa(fractstart-7) + ");\n"
						} else {
							result += "                        tx_data <= " + pay + "(" + strconv.Itoa(fractstart) + " downto 0) & x\"" + strings.Repeat("0", 7-fractstart) + "\";\n"
						}
					}
				}
			}

		}
	}

	if firstchar {
		result += "                    tx_data        <= memory(to_integer(current_char));\n"
	} else {
		result += "                    else\n"
		result += "                        tx_data        <= memory(to_integer(current_char));\n"
		result += "                    end if;\n"
	}
	result += "                    tx_data_enable <= '1';\n"
	result += "                    current_char   <= current_char + 1;\n"
	result += "                end if;\n"
	result += "            elsif delay_counter /= 0 then                \n"
	result += "                delay_counter <= delay_counter - 1;\n"
	result += "            else\n"
	result += "                case state is                  \n"
	result += "                    when \"00000\" => \n"
	result += "                        -- Power down the module and delay \n"
	//result += "                        status_connected <= '0';\n"
	//result += "                        status_sending   <= '0';\n"
	//result += "                        status_receiving <= '0';\n"
	//result += "                        status_active    <= '0';\n"
	//result += "                        status_wifi_up   <= '0';\n"
	//result += "                        status_error     <= '0';\n"
	//result += "                        wifi_enable      <= '0'; \n"

	if cipstart_loop_present {
		firstcom, _ := LocateCommandbyIndex(completedcommands, "00110", 0)
		result += "                                   cipstart_loop <= \"" + firstcom.SecondaryState + "\";\n"
	}

	//result += "                        if powerdown = '0' then\n"
	result += "                            state <= \"00001\";\n"
	//result += "                        end if;\n"
	result += "                        delay_counter <= (others => '1');\n"
	result += "                    when \"00001\" => \n"
	result += "                        -- Power up the module and delay \n"
	//result += "                        status_active  <= '1';\n"
	result += "                        wifi_enable    <= '1'; \n"
	result += "                        delay_counter <= (others => '1');\n"
	result += "                        state <= \"00010\";\n"
	result += "                        \n"
	result += "                    when \"00010\" =>\n"
	result += "                        -- Power up the module and delay \n"
	result += "                        delay_counter <= (others => '1');\n"
	result += "                        state <= \"00011\";\n"
	result += "                        \n"
	result += "                    when \"00011\" =>\n"
	result += "                        -- Should be waiting for \"ready\"\n"
	result += "                        if rx_seeing_ready = '1' then                    \n"
	result += "                            -- Set wifi mode\n"
	result += "                            current_char <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "00011", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                            sending <= '1'; \n"
	result += "                            state <= \"00100\";\n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                        end if;                    \n"
	result += "                    when \"00100\" =>\n"
	result += "                        -- Should be waiting \"OK\" or \"no change\"\n"
	result += "                        if rx_seeing_ok = '1' or rx_seeing_change = '1' then                    \n"
	result += "                            -- Set to multiple connection mode\n"
	result += "                            current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "00100", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                            sending       <= '1'; \n"
	result += "                            state <= \"00101\";\n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                        end if;\n"
	result += "                        \n"
	result += "                    when \"00101\" =>\n"
	result += "                        -- Should be waiting \"OK\" or \"no change\"\n"
	result += "                        if rx_seeing_ok = '1' then                                        \n"
	result += "                            -- Connect to the Wifi network\n"
	result += "                            current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "00101", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                            sending       <= '1'; \n"
	result += "                            state <= \"11100\";\n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                        end if;\n"
	result += "                        \n"
	result += "                    when \"11100\" =>\n"
	result += "                       --  Should be waiting \"OK\" or \"no change\"\n"
	result += "                        if rx_seeing_ok = '1' then                                        \n"
	result += "                             current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "11100", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                             sending       <= '1'; \n"
	result += "                             state <= \"01110\";\n"
	result += "                             last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                       end if;                        \n"
	result += "                        \n"
	result += "                    when \"01110\" =>\n"
	result += "                        -- Should be waiting \"OK\"\n"
	result += "                        if rx_seeing_ok = '1' then                                            \n"
	result += "                            -- Open the receiving socket\n"
	//result += "                            status_wifi_up   <= '1';\n"
	result += "                            current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "01110", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                            sending       <= '1'; \n"
	result += "                            state <= \"01101\"; \n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                        end if;\n"
	result += "                        \n"
	result += "                    when \"01101\" =>\n"
	result += "                        -- Should be waiting \"OK\"\n"
	result += "                        if rx_seeing_ok = '1' then                                            \n"
	result += "                            -- Open the broadcast socket\n"
	//result += "                            status_wifi_up   <= '1';\n"
	result += "                            current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "01101", "0")) + "," + strconv.Itoa(membits) + ");\n"
	result += "                            sending       <= '1'; \n"
	result += "                            state <= \"00110\"; \n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                        end if;\n"
	result += "                        \n"
	result += "                    when \"00110\" =>\n"
	if cipstart_loop_present {
		result += "                            case cipstart_loop is \n"
		var firstcom FirmwareCommand
		var oldcom FirmwareCommand
		for currindex := 0; ; currindex++ {
			currcom, currstat := LocateCommandbyIndex(completedcommands, "00110", currindex)

			if currindex != 0 {
				result += "                            when \"" + oldcom.SecondaryState + "\" =>\n"

				result += "                               -- Should be waiting \"OK\"\n"
				result += "                               if rx_seeing_ok = '1' then                                            \n"
				result += "                                   -- Open the sending sockets\n"
				//result += "                                   status_wifi_up   <= '1';\n"
				result += "                                   current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "00110", oldcom.SecondaryState)) + "," + strconv.Itoa(membits) + ");\n"
				result += "                                   sending       <= '1'; \n"
				result += "                                   receiving     <= '1';\n"
				if currstat {
					result += "                                   cipstart_loop <= \"" + currcom.SecondaryState + "\";\n"
				} else {
					result += "                                   cipstart_loop <= \"" + firstcom.SecondaryState + "\";\n"
					result += "                                   state <= \"00111\";\n"
				}
				result += "                                   last_rx_chars(4 downto 0) <= (others => '0');\n"
				result += "                               end if; \n"
				result += "                        \n"
			} else {
				firstcom = currcom
			}

			if currstat {
				oldcom = currcom
			} else {
				break
			}
		}
		result += "                            when others =>\n"
		result += "                                cipstart_loop <= \"" + firstcom.SecondaryState + "\";\n"
		result += "                             \n"
		result += "                            end case;\n"
	} else {
		result += "                            state <= \"00111\";\n"
	}

	result += "                   when \"00111\" =>\n"

	{
		result += "                            case message_loop is \n"
		var firstcom FirmwareCommand
		var oldcom FirmwareCommand
		for currindex := 0; ; currindex++ {
			currcom, currstat := LocateCommandbyIndex(completedcommands, "00111", currindex)

			if currindex != 0 {
				result += "                            when \"" + oldcom.SecondaryState + "\" =>\n"
				result += "                                if " + oldcom.Signal + " = '1' then\n"
				result += "                                    -- Should be waiting \"linked\" and \"OK\"\n"
				result += "                                    if rx_seeing_ok = '1' or rx_seeing_recv = '1' then\n"
				result += "                                        current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "00111", oldcom.SecondaryState)) + "," + strconv.Itoa(membits) + ");\n"
				result += "                                        sending       <= '1'; \n"
				result += "                                        state <= \"01000\";\n"
				result += "                                        if rx_seeing_ok = '1' then \n"
				result += "                                            last_rx_chars(4 downto 0) <= (others => '0');\n"
				result += "                                        end if; \n"
				result += "                                    end if; \n"
				result += "                                else\n"
				if currstat {
					result += "                                    message_loop <= \"" + currcom.SecondaryState + "\";\n"
				} else {
					result += "                                    message_loop <= \"" + firstcom.SecondaryState + "\";\n"
				}
				result += "                                end if;\n"
			} else {
				firstcom = currcom
			}

			if currstat {
				oldcom = currcom
			} else {
				break
			}
		}
		result += "                            when others =>\n"
		result += "                                message_loop <= \"" + firstcom.SecondaryState + "\";\n"
		result += "                             \n"
		result += "                            end case;\n"
	}

	result += "                    when \"01000\" =>\n"

	{
		result += "                            case message_loop is \n"
		var firstcom FirmwareCommand
		var oldcom FirmwareCommand
		for currindex := 0; ; currindex++ {
			currcom, currstat := LocateCommandbyIndex(completedcommands, "01000", currindex)

			if currindex != 0 {
				result += "                            when \"" + oldcom.SecondaryState + "\" =>\n"
				result += "                                if " + oldcom.Signal + " = '1' then\n"
				result += "                                    -- Should be waiting \"linked\" and \"OK\"\n"
				result += "                                    if rx_seeing_prompt = '1' then                                                                                            \n"
				result += "                                        current_char  <= to_unsigned(" + strconv.Itoa(LocateCommand(completedcommands, "01000", oldcom.SecondaryState)) + "," + strconv.Itoa(membits) + ");\n"
				result += "                                        sending       <= '1'; \n"
				result += "                                        state <= \"01001\";\n"
				if currstat {
					result += "                                        message_loop <= \"" + currcom.SecondaryState + "\";\n"
				} else {
					result += "                                        message_loop <= \"" + firstcom.SecondaryState + "\";\n"
				}
				result += "                                        reset_" + oldcom.Signal + " <= '1';\n"
				result += "                                        last_rx_chars(4 downto 0) <= (others => '0');\n"
				result += "                                    end if; \n"
				result += "                                else\n"
				result += "                                        state <= \"00111\";\n"
				if currstat {
					result += "                                        message_loop <= \"" + currcom.SecondaryState + "\";\n"
				} else {
					result += "                                        message_loop <= \"" + firstcom.SecondaryState + "\";\n"
				}
				result += "                                end if;\n"
			} else {
				firstcom = currcom
			}

			if currstat {
				oldcom = currcom
			} else {
				break
			}
		}
		result += "                            when others =>\n"
		result += "                                state <= \"00111\";\n"
		result += "                                message_loop <= \"" + firstcom.SecondaryState + "\";\n"
		result += "                             \n"
		result += "                            end case;\n"

	}

	result += "\n"
	result += "                    when \"01001\" =>\n"
	result += "                        if rx_seeing_ok = '1' then -- Should really be looking for \"SEND OK\", but...                                                                    \n"
	result += "                            -- Close the socket \n"
	result += "                            if rx_seeing_ok = '1' then\n"
	//result += "                                last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "                            \n"
	result += "                                if powerdown = '1' then     \n"
	result += "                                    -- Jump to the shutdown state                                                               \n"
	result += "                                    state         <= \"01100\";\n"
	// TODO Original to close the connections	result += "                                    state         <= \"01011\";\n"
	result += "                                else\n"

	for currindex := 0; ; currindex++ {
		currcom, currstat := LocateCommandbyIndex(completedcommands, "00111", currindex)

		if currcom.Signal != "" {

			result += "                                    reset_" + currcom.Signal + " <= '0';\n"
		}

		if !currstat {
			break
		}
	}

	//result += "                                --    delay_counter <= (others => '1');\n"
	result += "                                    state         <= \"00111\";\n"
	result += "                                end if;\n"
	result += "                            end if;\n"
	result += "                        end if;\n"

	// TODO The connections are not closed
	//	result += "                    when \"01011\" =>\n"
	//result += "                        status_sending <= '0';                        \n"
	//	result += "                        current_char   <= to_unsigned(66," + strconv.Itoa(membits) + ");\n"
	//	result += "                        sending        <= '1'; \n"
	//	result += "                        state          <= \"01100\";\n"
	//	result += "                    \n"

	result += "                    when \"01100\" =>  -- Power down the module.\n"
	result += "                        if rx_seeing_ok = '1' then           \n"
	//result += "                             status_connected <= '0';\n"
	result += "                            -- Wait a while before power down the module.                                                         \n"
	result += "                            delay_counter <= (others => '1');\n"
	result += "                            state <= \"00000\";\n"
	result += "                            last_rx_chars(4 downto 0) <= (others => '0');\n"
	result += "\n"
	result += "                        end if;\n"
	result += "\n"
	result += "                    when \"01111\" =>  -- Error state\n"
	//result += "                        status_connected <= '0';\n"
	//result += "                        status_sending   <= '0';\n"
	//result += "                        status_active    <= '0';\n"
	//result += "                        status_wifi_up   <= '0';\n"
	//result += "                        status_error     <= '1';                        \n"
	result += "                        delay_counter <= (others => '1');\n"
	result += "                        -- Power down and hang here.\n"
	//result += "                        wifi_enable      <= '0';\n"
	result += "                        state <= \"00000\"; -- restart\n"
	result += "                         \n"
	result += "                    when others =>\n"
	result += "                        state <= \"00000\"; -- restart\n"
	result += "                end case;\n"
	result += "            end if;\n"
	result += "\n"
	result += "            --==================================================================\n"
	result += "            -- Sort of a watchdog  \n"
	result += "            -- inc_wd_high is '1' every one millisecond.\n"
	result += "            -- so if we don't see a state change for 10 seconds, then\n"
	result += "            -- trigger the watchdog to reset everything\n"
	result += "            --==================================================================\n"
	result += "            if inc_wd_high = '1' then\n"
	result += "                if watchdog_high = 10000 then\n"
	result += "                    state <= \"01111\"; -- Flag error and restart\n"
	result += "                end if;\n"
	result += "                watchdog_high <= watchdog_high + 1;\n"
	result += "            end if;\n"
	result += "\n"
	result += "            -- reset the watchdog if the state changes\n"
	result += "            if state_last /= state or state = \"00000\" then\n"
	result += "               watchdog_high <= (others => '0');\n"
	result += "            end if;\n"
	result += "            state_last <= state;\n"
	result += "\n"
	result += "            if watchdog_low = 99999 then\n"
	result += "                watchdog_low <= (others => '0');\n"
	result += "                inc_wd_high <= '1';\n"
	result += "            else \n"
	result += "                watchdog_low <= watchdog_low + 1;\n"
	result += "                inc_wd_high <= '0';\n"
	result += "            end if;\n"
	result += "\n"
	result += "            --==================================================================\n"
	result += "            -- Broadcast counter \n"
	result += "            --==================================================================\n"
	result += "            if inc_wd_high = '1' then\n"
	//result += "                if counter = 5001 then\n"
	//result += "                   counter <=  (others => '0');\n"
	//result += "                else\n"
	result += "    	              counter <= counter + 1;\n"
	result += "    	              bcounter <= bcounter + 1;\n"
	//result += "                end if;\n"
	result += "            end if;\n"
	result += "            --==================================================================\n"
	result += "            -- Processing the received bytes of data \n"
	result += "            --==================================================================             \n"
	result += "            if rx_data_enable = '1' then\n"
	result += "                last_rx_chars <= last_rx_chars(last_rx_chars'high-8 downto 0) & rx_data;\n"
	result += "            end if;\n"
	result += "\n"
	result += "            if last_rx_chars(63 downto 0) = x\"6368616e67650d0a\" then -- ASCII for \"change\\r\\n\"\n"
	result += "                 rx_seeing_change <= '1';\n"
	result += "             else\n"
	result += "                 rx_seeing_change <= '0';\n"
	result += "             end if;\n"
	result += "\n"
	result += "            if last_rx_chars(55 downto 0) = x\"72656164790d0a\" or last_rx_chars(55 downto 0) = x\"4f542049500d0a\" then -- ASCII for \"ready\\r\\n\"\n"
	result += "                rx_seeing_ready <= '1';\n"
	result += "            else\n"
	result += "                rx_seeing_ready <= '0';\n"
	result += "            end if;\n"
	result += "\n"
	result += "            if last_rx_chars(15 downto 0) = x\"3e20\" then -- ASCII for \"\"> \" prompt\n"
	result += "                rx_seeing_prompt <= '1';\n"
	result += "            else\n"
	result += "                rx_seeing_prompt <= '0';\n"
	result += "            end if;\n"
	result += "\n"

	result += "            if last_rx_chars(31 downto 0) = x\"4f4b0d0a\" then -- ASCII for \"OK\\r\\n\"\n"
	result += "                rx_seeing_ok <= '1';\n"
	result += "            else\n"
	result += "                rx_seeing_ok <= '0';\n"
	result += "            end if;\n"
	result += "\n"

	firstreveicer := false

	// Receiving the real input payloads
	if udpbond_params["input_ids"] != "" {

		sources := strings.Split(udpbond_params["sources"], ",")
		ins := strings.Split(udpbond_params["inputs"], ",")

		for i, resid := range strings.Split(udpbond_params["input_ids"], ",") {

			iname := ins[i]
			source := sources[i]

			intsourceid, _ := strconv.Atoi(source)
			hexsourceid := fmt.Sprintf("%08x", intsourceid)

			residint, _ := strconv.Atoi(resid)
			residhex := fmt.Sprintf("%08x", residint)

			expected_prefix := Ascii2Hex("+IPD,0," + strconv.Itoa(19+payload_fractions) + ":<<<888805>>>")
			expected_other := Ascii2Hex("<<<" + hexclusid + hexsourceid + residhex + ">>>")

			expected_prefix_len := len(expected_prefix) / 2
			expected_other_len := len(expected_other) / 2
			if !firstreveicer {
				firstreveicer = true
				result += "            if last_rx_chars(" + strconv.Itoa(rsize+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(rsize+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(rsize+(expected_other_len*8)-1) + " downto " + strconv.Itoa(rsize) + ") = x\"" + expected_other + "\" then\n"
			} else {
				result += "            elsif last_rx_chars(" + strconv.Itoa(rsize+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(rsize+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(rsize+(expected_other_len*8)-1) + " downto " + strconv.Itoa(rsize) + ") = x\"" + expected_other + "\" then\n"
			}
			result += "                --tagtoack <= last_rx_chars(" + strconv.Itoa(rsize+expected_other_len*8+32-1) + " downto " + strconv.Itoa(rsize+expected_other_len*8) + ") ;\n"
			result += "                --optoack <= x\"05\" ;\n"
			result += "                --peertoack <= x\"" + hexsourceid + "\" ;\n"
			result += "                payload_" + iname + " <= last_rx_chars(" + strconv.Itoa(rsize-1) + " downto 0) ;\n"
			result += "                rx_seeing_recv <= '1';\n"

		}
	}

	// Receiving the input payloads data valid
	if udpbond_params["input_ids"] != "" {

		sources := strings.Split(udpbond_params["sources"], ",")
		ins := strings.Split(udpbond_params["inputs"], ",")

		for i, resid := range strings.Split(udpbond_params["input_ids"], ",") {

			iname := ins[i]
			source := sources[i]

			intsourceid, _ := strconv.Atoi(source)
			hexsourceid := fmt.Sprintf("%08x", intsourceid)

			residint, _ := strconv.Atoi(resid)
			residhex := fmt.Sprintf("%08x", residint)

			cname := "payload_" + iname + "_dv"

			expected_prefix := Ascii2Hex("+IPD,0," + strconv.Itoa(20) + ":<<<888806>>>")
			expected_other := Ascii2Hex("<<<" + hexclusid + hexsourceid + residhex + ">>>")

			expected_prefix_len := len(expected_prefix) / 2
			expected_other_len := len(expected_other) / 2
			if !firstreveicer {
				firstreveicer = true
				result += "            if last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(8+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)-1) + " downto " + strconv.Itoa(8) + ") = x\"" + expected_other + "\" then\n"
			} else {
				result += "            elsif last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(8+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)-1) + " downto " + strconv.Itoa(8) + ") = x\"" + expected_other + "\" then\n"
			}
			result += "                --tagtoack <= last_rx_chars(" + strconv.Itoa(8+expected_other_len*8+32-1) + " downto " + strconv.Itoa(8+expected_other_len*8) + ") ;\n"
			result += "                --optoack <= x\"06\" ;\n"
			result += "                --peertoack <= x\"" + hexsourceid + "\" ;\n"
			result += "                    if last_rx_chars(7 downto 0) = x\"00\" then \n"
			result += "                        " + cname + " <= '0';\n"
			result += "                    else \n"
			result += "                        " + cname + " <= '1';\n"
			result += "                    end if;\n"
			result += "                rx_seeing_recv <= '1';\n"

		}
	}

	// These are the received signals, one each destination
	if udpbond_params["output_ids"] != "" {

		dests := strings.Split(udpbond_params["destinations"], ",")

		for i, resid := range strings.Split(udpbond_params["output_ids"], ",") {

			destlist := dests[i]
			for _, dest := range strings.Split(destlist, "-") {

				conn := dest2conn[dest]

				intsourceid, _ := strconv.Atoi(conn)
				hexsourceid := fmt.Sprintf("%08x", intsourceid)

				residint, _ := strconv.Atoi(resid)
				residhex := fmt.Sprintf("%08x", residint)

				cname := "io_tr_recv_" + strconv.Itoa(residint) + "_on_" + conn

				expected_prefix := Ascii2Hex("+IPD,0," + strconv.Itoa(20) + ":<<<888807>>>")
				expected_other := Ascii2Hex("<<<" + hexclusid + hexsourceid + residhex + ">>>")

				expected_prefix_len := len(expected_prefix) / 2
				expected_other_len := len(expected_other) / 2

				if !firstreveicer {
					firstreveicer = true
					result += "            if last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(8+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)-1) + " downto " + strconv.Itoa(8) + ") = x\"" + expected_other + "\" then\n"
				} else {
					result += "            elsif last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)+(4*8)+(expected_prefix_len*8)-1) + " downto " + strconv.Itoa(8+(expected_other_len*8)+(4*8)) + ") = x\"" + expected_prefix + "\" and last_rx_chars(" + strconv.Itoa(8+(expected_other_len*8)-1) + " downto " + strconv.Itoa(8) + ") = x\"" + expected_other + "\" then\n"
				}
				result += "                --tagtoack <= last_rx_chars(" + strconv.Itoa(8+expected_other_len*8+32-1) + " downto " + strconv.Itoa(8+expected_other_len*8) + ") ;\n"
				result += "                --optoack <= x\"07\" ;\n"
				result += "                --peertoack <= x\"" + hexsourceid + "\" ;\n"
				result += "                    if last_rx_chars(7 downto 0) = x\"00\" then \n"
				result += "                        " + cname + " <= '0';\n"
				result += "                    else \n"
				result += "                        " + cname + " <= '1';\n"
				result += "                    end if;\n"
				result += "                rx_seeing_recv <= '1';\n"

			}
		}
	}

	// Close the cases whenever it was opened
	if firstreveicer {
		result += "            else\n"
		result += "                rx_seeing_recv <= '0';\n"
		result += "            end if;\n"
		result += "\n"
	}

	result += "        end if;\n"
	result += "    end process;\n"
	result += "\n"
	result += "end Behavioral;\n"
	result += "\n"
	result += "-----------------------------------------\n"
	result += "-- tx.vhd - Transmit data to an ESP8266\n"
	result += "--\n"
	result += "-- Author: Mike Field <hamster@snap.net.nz>\n"
	result += "--\n"
	result += "-- Designed for 9600 baud and 100MHz clock\n"
	result += "--\n"
	result += "------------------------------------------------\n"
	result += "library IEEE;\n"
	result += "use IEEE.STD_LOGIC_1164.ALL;\n"
	result += "use IEEE.NUMERIC_STD.ALL;\n"
	result += "\n"
	result += "entity tx is\n"
	result += "    Port ( clk         : in  STD_LOGIC;\n"
	result += "           data        : in  STD_LOGIC_VECTOR (7 downto 0);\n"
	result += "           data_enable : in  STD_LOGIC;\n"
	result += "           busy        : out STD_LOGIC;\n"
	result += "           tx_out      : out STD_LOGIC);\n"
	result += "end tx;\n"
	result += "\n"
	result += "architecture Behavioral of tx is\n"
	result += "    signal baud_count       : unsigned(13 downto 0) := (others => '0');\n"
	result += "    constant baud_count_max : unsigned(13 downto 0) := to_unsigned(100000000/115200, 14);\n"
	result += "    signal busy_sr          : std_logic_vector(9 downto 0) := (others => '0');\n"
	result += "\n"
	result += "    signal sending          : std_logic_vector(9 downto 0) := (others => '0');\n"
	result += "begin\n"
	result += "    busy <= busy_sr(0) or data_enable;\n"
	result += "    \n"
	result += "clk_proc: process(clk)\n"
	result += "    begin\n"
	result += "        if rising_edge(clk) then\n"
	result += "            if baud_count = 0 then\n"
	result += "                baud_count <= baud_count_max;\n"
	result += "                tx_out     <= sending(0);\n"
	result += "                sending    <= '1' & sending(sending'high downto 1);\n"
	result += "                busy_sr    <= '0' & busy_sr(busy_sr'high downto 1);\n"
	result += "            else\n"
	result += "                baud_count  <= baud_count - 1;         \n"
	result += "            end if;\n"
	result += "\n"
	result += "            if busy_sr(0) = '0' and data_enable = '1' then\n"
	result += "                baud_count <= baud_count_max;\n"
	result += "                sending    <= \"1\" & data & \"0\";\n"
	result += "                busy_sr    <= (others =>'1');\n"
	result += "            end if;\n"
	result += "            \n"
	result += "        end if;\n"
	result += "    end process;\n"
	result += "\n"
	result += "end Behavioral;\n"
	result += "\n"
	result += "-----------------------------------------\n"
	result += "-- rx.vhd - Receive serial data from an ESP8266\n"
	result += "--\n"
	result += "-- Author: Mike Field <hamster@snap.net.nz>\n"
	result += "--\n"
	result += "-- Designed for 9600 baud and 100MHz clock\n"
	result += "--\n"
	result += "------------------------------------------------\n"
	result += "library IEEE;\n"
	result += "use IEEE.STD_LOGIC_1164.ALL;\n"
	result += "use IEEE.NUMERIC_STD.ALL;\n"
	result += "\n"
	result += "entity rx is\n"
	result += "    Port ( clk : in STD_LOGIC;\n"
	result += "           rx_in : in STD_LOGIC;\n"
	result += "           data : out STD_LOGIC_VECTOR (7 downto 0);\n"
	result += "           data_enable : out STD_LOGIC);\n"
	result += "end rx;\n"
	result += "\n"
	result += "architecture Behavioral of rx is\n"
	result += "    signal baud_count       : unsigned(13 downto 0) := (others => '0');\n"
	result += "    constant baud_count_max : unsigned(13 downto 0) := to_unsigned(100000000/115200, 14);\n"
	result += "    signal busy                : std_logic:= '0';\n"
	result += "    signal receiving           : std_logic_vector(7 downto 0) := (others => '0');\n"
	result += "    signal rx_in_last          : std_logic:= '1';\n"
	result += "    signal rx_in_synced        : std_logic:= '1';\n"
	result += "    signal rx_in_almost_synced : std_logic:= '1';\n"
	result += "\n"
	result += "    signal bit_count           : unsigned(3 downto 0) := (others => '0');\n"
	result += "begin\n"
	result += "    \n"
	result += "process(clk)\n"
	result += "    begin\n"
	result += "        if rising_edge(clk) then\n"
	result += "            data_enable <= '0';\n"
	result += "            if busy = '1' then\n"
	result += "            \n"
	result += "                if baud_count = 0 then\n"
	result += "                    if bit_count = 9 then\n"
	result += "                        -- We've got all the bits we need\n"
	result += "                        busy        <= '0';\n"
	result += "                        data        <= receiving(7 downto 0);\n"
	result += "                        data_enable <= '1';\n"
	result += "                    end if;\n"
	result += "                    \n"
	result += "                    -- receive this bit\n"
	result += "                    receiving  <= rx_in_synced & receiving(7 downto 1);\n"
	result += "                    -- Set timer for the next bit\n"
	result += "                    bit_count  <= bit_count + 1;        \n"
	result += "                    baud_count <= baud_count_max;\n"
	result += "                else\n"
	result += "                    baud_count <= baud_count-1;\n"
	result += "                end if; \n"
	result += "            else\n"
	result += "                -- Is this the falling edge of the start bit?\n"
	result += "               if rx_in_last = '1' and rx_in_synced = '0' then\n"
	result += "                    -- Load it up with half the count so we sample in the middle of the bit\n"
	result += "                    baud_count <= '0' & baud_count_max(13 downto 1);\n"
	result += "                    bit_count  <= (others => '0');\n"
	result += "                    busy       <= '1';\n"
	result += "               end if;   \n"
	result += "            end if;\n"
	result += "\n"
	result += "            rx_in_last   <= rx_in_synced;\n"
	result += "            -- Synchronise the RX signal\n"
	result += "            rx_in_synced        <= rx_in_almost_synced;\n"
	result += "            rx_in_almost_synced <= rx_in;\n"
	result += "        end if;\n"
	result += "    end process;\n"
	result += "end Behavioral;\n"
	result += "\n"

	return []string{"udpbond.vhd"}, []string{result}
}
