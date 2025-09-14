package bondmachine

import (
	//	"fmt"

	"fmt"
	"strconv"

	//	"strings"
	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bmstack"
	"github.com/BondMachineHQ/BondMachine/pkg/bondirect"
)

type Bondirect_extra struct {
	// Config  *bondirect.Config
	// Cluster *bmcluster.Cluster
	// Mesh    *bondirect.Mesh
	*bondirect.BondirectElement
	PeerID   uint32
	PeerName string
	Maps     *IOmap
	Flavor   string
}

func (sl *Bondirect_extra) Get_Name() string {
	return "bondirect"
}

func (sl *Bondirect_extra) Get_Params() *ExtraParams {
	result := new(ExtraParams)
	result.Params = make(map[string]string)
	result.Params["peer_id"] = strconv.Itoa(int(sl.PeerID))
	result.Params["cluster_id"] = strconv.Itoa(int(sl.Cluster.ClusterId))

	var mypeer bmcluster.Peer

	for _, peer := range sl.Cluster.Peers {
		if peer.PeerId == sl.PeerID {
			mypeer = peer
			result.Params["mesh_node_name"], _ = sl.AnyNameToMeshName(peer.PeerName)
			result.Params["node_name"], _ = sl.AnyNameToClusterName(peer.PeerName)
			break
		}

	}

	result.Params["input_ids"] = ""
	result.Params["inputs"] = ""
	result.Params["sources"] = ""
	// fmt.Println("mypeer", mypeer)
	// fmt.Println("cluster", sl.Cluster)

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

	result.Params["lines"] = ""
	for _, line := range sl.Lines {
		result.Params["lines"] += line + ","
	}
	if result.Params["lines"] != "" {
		result.Params["lines"] = result.Params["lines"][0 : len(result.Params["lines"])-1]
	}

	result.Params["insignals"] = ""
	result.Params["inports"] = ""
	result.Params["outsignals"] = ""
	result.Params["outports"] = ""

	for i, signals := range sl.TData.WiresIn {
		prefix := sl.TData.Lines[i] + sl.TData.TrIn[i]
		for _, sig := range signals {
			result.Params["insignals"] += prefix + sig + ","
		}
	}

	for _, ports := range sl.TData.WiresInNames {
		for _, port := range ports {
			result.Params["inports"] += port + ","
		}
	}

	for i, signals := range sl.TData.WiresOut {
		prefix := sl.TData.Lines[i] + sl.TData.TrOut[i]
		for _, sig := range signals {
			result.Params["outsignals"] += prefix + sig + ","
		}
	}

	for _, ports := range sl.TData.WiresOutNames {
		for _, port := range ports {
			result.Params["outports"] += port + ","
		}
	}
	if result.Params["insignals"] != "" {
		result.Params["insignals"] = result.Params["insignals"][0 : len(result.Params["insignals"])-1]
	}
	if result.Params["inports"] != "" {
		result.Params["inports"] = result.Params["inports"][0 : len(result.Params["inports"])-1]
	}
	if result.Params["outsignals"] != "" {
		result.Params["outsignals"] = result.Params["outsignals"][0 : len(result.Params["outsignals"])-1]
	}
	if result.Params["outports"] != "" {
		result.Params["outports"] = result.Params["outports"][0 : len(result.Params["outports"])-1]
	}

	// fmt.Println("lines:", result.Params)
	return result
}

func (sl *Bondirect_extra) Import(inp string) error {
	return nil
}

func (sl *Bondirect_extra) Export() string {
	return ""
}

func (sl *Bondirect_extra) Check(bmach *Bondmachine) error {
	return nil
}

func (sl *Bondirect_extra) Verilog_headers() string {
	result := "\n"
	return result
}
func (sl *Bondirect_extra) StaticVerilog() string {

	result := "\n"
	return result
}

func (sl *Bondirect_extra) ExtraFiles() ([]string, []string) {
	files := make([]string, 0)
	code := make([]string, 0)
	nodeName, _ := sl.BondirectElement.AnyNameToClusterName(sl.PeerName)

	// Generate extra files for the Bondirect module

	// Endpoint
	epCode, _ := sl.GenerateEndpoint("", nodeName)
	files = append(files, "bd_endpoint_"+sl.PeerName+".vhd")
	code = append(code, epCode)

	// Lines
	for _, line := range sl.Lines {
		lineCode, _ := sl.GenerateLine("", nodeName, line)
		files = append(files, "bd_line_"+sl.PeerName+"_"+line+".vhd")
		code = append(code, lineCode)
	}
	// Queues
	maps, _ := sl.BondirectElement.SolveMessages()
	peerName := sl.PeerName
	myMessages := maps[peerName]

	for _, line := range sl.Lines {
		fmt.Println("Generating queue for line", line)
		// Every line (let's call it wireB) has an input queue with several senders and one receiver
		// Senders are:
		// - One for every message coming from the local peer (bm out data and valid, bm in recv)
		//   These messages potentially can be concurrent, so they need to be queued by different senders
		// - One for every couple (wireA, wireB) if there are messages coming from wireA
		//   That has to be routed to the other wireB.
		//   Even if there are more than one message coming from wireA, they will be serialized by the
		//   endpoint, so only one sender is needed

		origins := *(myMessages.Origins)
		originsNextHopVia := *(myMessages.OriginsNextHopVia)
		routes := *(myMessages.Routes)
		routesNextHopVia := *(myMessages.RoutesNextHopVia)
		routesPrevHopVia := *(myMessages.RoutesPrevHopVia)

		s := bmstack.CreateBasicStack()
		s.ModuleName = "bond_queue_" + sl.PeerName + "_" + line
		s.DataSize = sl.InnerMessLen
		s.Depth = 8
		s.MemType = "FIFO"
		s.Receivers = make([]string, 1)
		s.Receivers[0] = line + "_queue_receiver"
		s.Senders = make([]string, 0)
		for i, msg := range origins {
			destWire := originsNextHopVia[i]
			if destWire == line {
				// This message is destined to this line, so it has a sender
				s.Senders = append(s.Senders, msg+"_to_"+destWire+"_sender")
			}
		}

		wire2wireSenders := make(map[string]struct{})

		for i, _ := range routes {
			prevWire := routesPrevHopVia[i]
			nextWire := routesNextHopVia[i]
			if nextWire == line {
				// This message is routed through this line, so it has a sender from the previous wire
				wire2wireSenders[prevWire+"_to_"+nextWire+"_sender"] = struct{}{}
			}
		}

		for w2wSender, _ := range wire2wireSenders {
			s.Senders = append(s.Senders, w2wSender)
		}

		if len(s.Senders) > 0 {
			queueCode, _ := s.WriteHDL()
			files = append(files, "bond_queue_"+sl.PeerName+"_"+line+".v")
			code = append(code, queueCode)
		} else {
			fmt.Println("No senders for line", line, "no queue generated")

		}
	}

	// Transceivers
	for _, line := range sl.Lines {
		trCodeIn, _ := sl.GenerateTransceiver("", nodeName, line, "in")
		trCodeOut, _ := sl.GenerateTransceiver("", nodeName, line, "out")
		files = append(files, "bond_tx_"+sl.PeerName+"_"+line+"_out.vhd")
		code = append(code, trCodeOut)
		files = append(files, "bond_rx_"+sl.PeerName+"_"+line+"_in.vhd")
		code = append(code, trCodeIn)
	}

	return files, code
}
