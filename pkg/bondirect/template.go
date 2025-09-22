package bondirect

import (
	"fmt"
	"text/template"
)

const (
	ActionsNum = 3
)

type IOSpec struct {
	SignalName   string
	SignalType   string // "data", "valid" or "recv"
	AssociatedIO string // Associated IO signal name
	IOHeader     string // Header for data,valid and recv signals
}

type TData struct {
	// Define the fields for Tdata
	Prefix        string
	NodeName      string // Cluster name of the node
	MeshNodeName  string // Mesh name of the node
	EdgeName      string
	TransName     string
	Rsize         int               // Register size
	NodeNum       int               // Number of nodes
	NodeBits      int               // Bits needed for node addressing
	IONum         int               // Maximum number of inputs or outputs Among all nodes
	IOBits        int               // Bits needed for IO addressing
	InnerMessLen  int               // Length of inner messages
	Inputs        []string          // List of input signals
	Outputs       []string          // List of output signals
	Lines         []string          // List of line signals
	IOSenders     [][]IOSpec        // List of input signal senders, the first dimension is the line index, the second dimension is the input index
	RouteSenders  [][]IOSpec        // List of route signal senders, the first dimension is the line index, the second dimension is the route index
	IOReceivers   [][]IOSpec        // List of output signal receivers, the first dimension is the line index, the second dimension is the output index
	TrIn          []string          // List of transceivers for incoming signals
	TrOut         []string          // List of transceivers for outgoing signals
	WiresIn       [][]string        // List of wire signals, the first dimension is the line index, the second dimension is the incoming wire index (the 0 is the clock)
	WiresInNames  [][]string        // List of wire signal port names, the first dimension is the line index, the second dimension is the incoming wire index (the 0 is the clock)
	WiresOut      [][]string        // List of wire signals, the first dimension is the line index, the second dimension is the outgoing wire index (the 0 is the clock)
	WiresOutNames [][]string        // List of wire signal port names, the first dimension is the line index, the second dimension is the outgoing wire index (the 0 is the clock)
	NodeParams    map[string]string // Additional parameters for the specified node
	EdgeParams    map[string]string // Additional parameters for the specified edge
	TransParams   map[string]string // Additional parameters for the specified transceiver
}

func (be *BondirectElement) InitTData() {
	be.TData = &TData{
		Rsize:    int(be.Config.Rsize),
		NodeNum:  len(be.Cluster.Peers),
		NodeBits: NeededBits(len(be.Cluster.Peers)),
	}

	maxIO := 0
	for _, node := range be.Cluster.Peers {
		if len(node.Inputs) > maxIO {
			maxIO = len(node.Inputs)
		}
		if len(node.Outputs) > maxIO {
			maxIO = len(node.Outputs)
		}
	}
	be.TData.IONum = maxIO
	be.TData.IOBits = NeededBits(maxIO)
	be.TData.InnerMessLen = be.TData.NodeBits + be.TData.IOBits + NeededBits(ActionsNum) + be.TData.Rsize
}

func (be *BondirectElement) PopulateIOData(nodeName string) error {
	inputs := make([]string, 0)
	outputs := make([]string, 0)

	if clusterNodeName, err := be.AnyNameToClusterName(nodeName); err == nil {
		nodeName = clusterNodeName
	} else {
		return fmt.Errorf("failed to get cluster node name: %v", err)
	}

	found := false
	for _, node := range be.Cluster.Peers {
		if node.PeerName == nodeName {
			for _, inp := range node.Inputs {
				inName := fmt.Sprintf("input%d", inp)
				inputs = append(inputs, inName)
			}
			for _, out := range node.Outputs {
				outName := fmt.Sprintf("output%d", out)
				outputs = append(outputs, outName)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("node %s not found", nodeName)
	}

	be.TData.Inputs = inputs
	be.TData.Outputs = outputs

	return nil
}

func (be *BondirectElement) PopulateWireData(nodeName string) error {
	lines := make([]string, 0)
	trIn := make([]string, 0)
	trOut := make([]string, 0)
	wiresIn := make([][]string, 0)
	wiresOut := make([][]string, 0)
	wiresInNames := make([][]string, 0)
	wiresOutNames := make([][]string, 0)

	// Using cluster names to find the mesh node name (that can be different)
	if meshNodeName, err := be.AnyNameToMeshName(nodeName); err == nil {
		nodeName = meshNodeName
	} else {
		return fmt.Errorf("failed to get mesh node name: %v", err)
	}

	// fmt.Println("Populating wire data for node:", nodeName)

	for lineName, line := range be.Mesh.Edges {

		if line.NodeA == nodeName {
			// Check if NodeB is in the cluster
			if be.CheckMeshNodeName(line.NodeB) {
				lines = append(lines, lineName)
				incoming := line.FromBtoA.ATransceiver
				trIn = append(trIn, incoming)
				if signals, ports, err := be.GetTransceiverSignals(incoming); err == nil {
					wiresIn = append(wiresIn, signals)
					wiresInNames = append(wiresInNames, ports)
				} else {
					return fmt.Errorf("failed to get incoming transceiver signals: %v", err)
				}

				outgoing := line.FromAtoB.ATransceiver
				trOut = append(trOut, outgoing)
				if signals, ports, err := be.GetTransceiverSignals(outgoing); err == nil {
					wiresOut = append(wiresOut, signals)
					wiresOutNames = append(wiresOutNames, ports)
				} else {
					return fmt.Errorf("failed to get outgoing transceiver signals: %v", err)
				}

				continue
			}
		} else if line.NodeB == nodeName {
			// Check if NodeA is in the cluster
			if be.CheckMeshNodeName(line.NodeA) {
				lines = append(lines, lineName)

				incoming := line.FromAtoB.BTransceiver
				trIn = append(trIn, incoming)
				if signals, ports, err := be.GetTransceiverSignals(incoming); err == nil {
					wiresIn = append(wiresIn, signals)
					wiresInNames = append(wiresInNames, ports)
				} else {
					return fmt.Errorf("failed to get incoming transceiver signals: %v", err)
				}

				outgoing := line.FromBtoA.BTransceiver
				trOut = append(trOut, outgoing)
				if signals, ports, err := be.GetTransceiverSignals(outgoing); err == nil {
					wiresOut = append(wiresOut, signals)
					wiresOutNames = append(wiresOutNames, ports)
				} else {
					return fmt.Errorf("failed to get outgoing transceiver signals: %v", err)
				}

				continue
			}
		}
	}

	be.TData.Lines = lines
	be.TData.TrIn = trIn
	be.TData.TrOut = trOut
	be.TData.WiresIn = wiresIn
	be.TData.WiresOut = wiresOut
	be.TData.WiresInNames = wiresInNames
	be.TData.WiresOutNames = wiresOutNames

	// fmt.Println("Tdata:", be.TData)
	if maps, err := be.SolveMessages(); err == nil {

		myMessages := maps[nodeName]

		origins := *(myMessages.Origins)
		originsHeader := *(myMessages.OriginsHeader)
		originsType := *(myMessages.OriginsType)
		originsAssociatedIO := *(myMessages.OriginIO)
		originsNextHopVia := *(myMessages.OriginsNextHopVia)
		destinations := *(myMessages.Destinations)
		destinationsHeader := *(myMessages.DestinationsHeader)
		destinationsType := *(myMessages.DestinationsType)
		destinationsAssociatedIO := *(myMessages.DestinationIO)
		destinationsPrevHopVia := *(myMessages.DestinationsPrevHopVia)
		routes := *(myMessages.Routes)
		routesHeader := *(myMessages.RoutesHeader)
		routesNextHopVia := *(myMessages.RoutesNextHopVia)
		routesPrevHopVia := *(myMessages.RoutesPrevHopVia)

		for _, line := range be.Lines {

			ioSenders := make([]IOSpec, 0)
			for i, msg := range origins {
				destWire := originsNextHopVia[i]
				if destWire == line {
					newSender := IOSpec{SignalName: msg}
					newSender.SignalType = originsType[i]
					newSender.AssociatedIO = originsAssociatedIO[i]
					newSender.IOHeader = originsHeader[i]
					ioSenders = append(ioSenders, newSender)
				}
			}

			ioReceivers := make([]IOSpec, 0)
			for i, msg := range destinations {
				srcWire := destinationsPrevHopVia[i]
				if srcWire == line {
					newReceiver := IOSpec{SignalName: msg}
					newReceiver.SignalType = destinationsType[i]
					newReceiver.AssociatedIO = destinationsAssociatedIO[i]
					newReceiver.IOHeader = destinationsHeader[i]
					ioReceivers = append(ioReceivers, newReceiver)
				}
			}

			wire2wireSenders := make(map[string]map[string]struct{})
			for i := range routes {
				prevWire := routesPrevHopVia[i]
				nextWire := routesNextHopVia[i]
				if prevWire == line {
					if val, exists := wire2wireSenders[prevWire+"_to_"+nextWire+"_sender"]; !exists {
						newMap := make(map[string]struct{})
						newMap[routesHeader[i]] = struct{}{}
						wire2wireSenders[prevWire+"_to_"+nextWire+"_sender"] = newMap
					} else {
						val[routesHeader[i]] = struct{}{}
						wire2wireSenders[prevWire+"_to_"+nextWire+"_sender"] = val
					}
				}
			}

			w2wSenders := make([]IOSpec, 0)
			for w2wSender, header := range wire2wireSenders {
				newSender := IOSpec{SignalName: w2wSender}
				for h := range header {
					newSender.IOHeader += "\"" + h + "\"|"
				}
				if len(newSender.IOHeader) > 0 {
					newSender.IOHeader = newSender.IOHeader[:len(newSender.IOHeader)-1]
				}
				w2wSenders = append(w2wSenders, newSender)
			}

			be.TData.IOSenders = append(be.TData.IOSenders, ioSenders)
			be.TData.RouteSenders = append(be.TData.RouteSenders, w2wSenders)
			be.TData.IOReceivers = append(be.TData.IOReceivers, ioReceivers)
		}
	} else {
		return fmt.Errorf("failed to solve messages: %v", err)
	}
	return nil
}

func (be *BondirectElement) DumpTemplateData() string {
	result := ""
	result += fmt.Sprintf("Register Size: %d\n", be.TData.Rsize)
	result += fmt.Sprintf("Node Number: %d\n", be.TData.NodeNum)
	result += fmt.Sprintf("Node Bits: %d\n", be.TData.NodeBits)
	result += fmt.Sprintf("IO Number: %d\n", be.TData.IONum)
	result += fmt.Sprintf("IO Bits: %d\n", be.TData.IOBits)
	result += fmt.Sprintf("Inner Message Length: %d\n", be.TData.InnerMessLen)
	result += fmt.Sprintf("Inputs: %v\n", be.TData.Inputs)
	result += fmt.Sprintf("Outputs: %v\n", be.TData.Outputs)
	result += fmt.Sprintf("Lines: %v\n", be.TData.Lines)
	result += fmt.Sprintf("IO Senders: %v\n", be.TData.IOSenders)
	result += fmt.Sprintf("Route Senders: %v\n", be.TData.RouteSenders)
	result += fmt.Sprintf("IO Receivers: %v\n", be.TData.IOReceivers)
	result += fmt.Sprintf("Transceivers In: %v\n", be.TData.TrIn)
	result += fmt.Sprintf("Transceivers Out: %v\n", be.TData.TrOut)
	result += fmt.Sprintf("Wires In: %v\n", be.TData.WiresIn)
	result += fmt.Sprintf("Wires Out: %v\n", be.TData.WiresOut)
	result += fmt.Sprintf("Node Params: %v\n", be.TData.NodeParams)
	result += fmt.Sprintf("Edge Params: %v\n", be.TData.EdgeParams)
	result += fmt.Sprintf("Transceiver Params: %v\n", be.TData.TransParams)

	return result
}

var funcMap = template.FuncMap{
	"inc": func(i int) int {
		return i + 1
	},
	"dec": func(i int) int {
		return i - 1
	},
	"next": func(i int, max int) int {
		if i < max-1 {
			return i + 1
		} else {
			return 0
		}
	},
	"bits": func(i int) int {
		return NeededBits(i)
	},
	"len": func(s []string) int {
		return len(s)
	},
	"ios": func(s []IOSpec) int {
		return len(s)
	},
	"iter": func(n int) []int {
		result := make([]int, n)
		for i := 0; i < n; i++ {
			result[i] = i
		}
		return result
	},
	"str": func(i int) string {
		return fmt.Sprintf("%d", i)
	},
	"int": func(s string) int {
		var i int
		fmt.Sscanf(s, "%d", &i)
		return i
	},
	"add": func(a, b int) int {
		return a + b
	},
}

func NeededBits(num int) int {
	if num > 0 {
		for bits := 1; true; bits++ {
			if 1<<uint8(bits) >= num {
				return bits
			}
		}
	}
	return 0
}
