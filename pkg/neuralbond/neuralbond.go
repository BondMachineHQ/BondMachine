package neuralbond

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/BondMachineHQ/BondMachine/pkg/bminfo"
)

const (
	ROMCODE = uint8(0) + iota
	FRAGMENT
)

const (
	ASYNC = uint8(0) + iota
	SYNC
)

type TrainedNet struct {
	Nodes         []Node
	Weights       []Weight
	Neurons       map[string]*Neuron
	NetConfig     *Config
	RegisterSize  int
	IOMode        uint8
	OperatingMode uint8
}

type Group []string
type Config struct {
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Pruned        []string
	Collapsed     []Group
	Debug         bool
	Verbose       bool
	NeuronLibPath string
	*bminfo.BMinfo
}

type Neuron struct {
	Params []string
}

type Weight struct {
	Layer        int
	PosCurrLayer int
	PosPrevLayer int
	RelPosDown   int
	RelPosUp     int
	Value        float32
}

type Node struct {
	Layer   int
	Pos     int
	Type    string
	Bias    float32
	Inputs  int
	Outputs int
}

func (n *TrainedNet) Init(config *Config) error {
	n.NetConfig = config
	n.Neurons = make(map[string]*Neuron)

	// List nb files in the neuron library path and load them
	neuronFiles, err := os.ReadDir(n.NetConfig.NeuronLibPath)
	if err != nil {
		return err
	}
	for _, f := range neuronFiles {
		if len(f.Name()) > 3 && f.Name()[len(f.Name())-3:] == ".nb" {
			neuronFile, err := os.ReadFile(n.NetConfig.NeuronLibPath + "/" + f.Name())
			if err != nil {
				return err
			}
			neuron := new(Neuron)
			if err := json.Unmarshal(neuronFile, neuron); err != nil {
				return err
			}
			n.Neurons[f.Name()[0:len(f.Name())-3]] = neuron
		}
	}

	return nil
}

func (n *TrainedNet) Normalize() {
	for w, weight := range n.Weights {
		downL, downP := weight.Layer-1, weight.PosPrevLayer
		upL, upP := weight.Layer, weight.PosCurrLayer

		for i, node := range n.Nodes {
			if node.Layer == downL && node.Pos == downP {
				n.Nodes[i].Outputs++
			}
			if node.Layer == upL && node.Pos == upP {
				n.Nodes[i].Inputs++
			}
		}

		sameDown := make([]int, 0)
		sameUp := make([]int, 0)

		for _, chWeight := range n.Weights {
			if chWeight.Layer-1 == downL && chWeight.PosPrevLayer == downP {
				sameDown = append(sameDown, chWeight.PosCurrLayer)
			}
			if chWeight.Layer == upL && chWeight.PosCurrLayer == upP {
				sameUp = append(sameUp, chWeight.PosPrevLayer)
			}
		}

		// Sort the sameDown and sameUp arrays
		sort.Ints(sameDown)
		sort.Ints(sameUp)

		// fmt.Println(w, sameDown, sameUp)

		for i, v := range sameDown {
			if weight.PosCurrLayer == v {
				n.Weights[w].RelPosDown = i
				break
			}
		}

		for i, v := range sameUp {
			if weight.PosPrevLayer == v {
				n.Weights[w].RelPosUp = i
				break
			}
		}

	}
}

func (n *TrainedNet) WriteBasm() (string, error) {
	c := n.NetConfig
	regSize := n.RegisterSize
	result := fmt.Sprintf("%%meta bmdef     global registersize:%d\n", regSize)
	switch n.IOMode {
	case ASYNC:
		result += fmt.Sprintf("%%meta bmdef     global iomode:async\n")
	case SYNC:
		result += fmt.Sprintf("%%meta bmdef     global iomode:sync\n")
	}
	switch n.OperatingMode {
	case ROMCODE:
		for _, node := range n.Nodes {
			if node.Type == "input" {
				result += fmt.Sprintf("%%meta cpdef node_0_%d romcode:terminal\n", node.Pos)
				c.List[fmt.Sprintf("node_0_%d", node.Pos)] = "input"
				result += fmt.Sprintf("%%meta iodef input_%d type:io\n", node.Pos)
				result += fmt.Sprintf("%%meta ioatt input_%d cp:node_0_%d, type:input, index:0\n", node.Pos, node.Pos)
				result += fmt.Sprintf("%%meta ioatt input_%d cp:bm, type:input, index:%d\n", node.Pos, node.Pos)
			} else if node.Type == "output" {
				result += fmt.Sprintf("%%meta cpdef node_%d_%d romcode:terminal\n", node.Layer, node.Pos)
				c.List[fmt.Sprintf("node_%d_%d", node.Layer, node.Pos)] = "output"
				result += fmt.Sprintf("%%meta iodef output_%d type:io\n", node.Pos)
				result += fmt.Sprintf("%%meta ioatt output_%d cp:node_%d_%d, type:output, index:0\n", node.Pos, node.Layer, node.Pos)
				result += fmt.Sprintf("%%meta ioatt output_%d cp:bm, type:output, index:%d\n", node.Pos, node.Pos)
			} else {
				if neuron, ok := n.Neurons[node.Type]; ok {
					result += fmt.Sprintf("%%meta cpdef node_%d_%d romcode:%s", node.Layer, node.Pos, node.Type)
					c.List[fmt.Sprintf("node_%d_%d", node.Layer, node.Pos)] = node.Type
					for _, param := range neuron.Params {
						switch param {
						case "inputs":
							result += fmt.Sprintf(", inputs:%d", node.Inputs)
						case "outputs":
							result += fmt.Sprintf(", outputs:%d", node.Outputs)
						case "bias":
							result += fmt.Sprintf(", bias:"+c.TypePrefix+"%f", node.Bias)
						case "pos":
							result += fmt.Sprintf(", pos:%d", node.Pos)
						default:
							if value, ok := c.Params[param]; ok {
								result += fmt.Sprintf(", %s:%s", param, value)
							} else {
								return "", errors.New("Unknown parameter " + param)
							}
						}
					}
					result += "\n"
				}
			}
		}

		for _, weight := range n.Weights {
			weightCP := fmt.Sprintf("weightcp_%d_%d__%d_%d", weight.Layer-1, weight.PosPrevLayer, weight.Layer, weight.PosCurrLayer)
			downNode := fmt.Sprintf("node_%d_%d", weight.Layer-1, weight.PosPrevLayer)
			upNode := fmt.Sprintf("node_%d_%d", weight.Layer, weight.PosCurrLayer)
			result += fmt.Sprintf("%%meta cpdef %s romcode:weight, weight:"+c.TypePrefix+"%f\n", weightCP, weight.Value)
			c.List[weightCP] = "weight"
			result += fmt.Sprintf("%%meta iodef up%s type:io\n", weightCP)
			result += fmt.Sprintf("%%meta iodef down%s type:io\n", weightCP)
			result += fmt.Sprintf("%%meta ioatt down%s cp:%s, type:input, index:0\n", weightCP, weightCP)
			result += fmt.Sprintf("%%meta ioatt down%s cp:%s, type:output, index:0\n", weightCP, downNode)
			result += fmt.Sprintf("%%meta ioatt up%s cp:%s, type:input, index:%d\n", weightCP, upNode, weight.RelPosUp)
			result += fmt.Sprintf("%%meta ioatt up%s cp:%s, type:output, index:0\n", weightCP, weightCP)
		}
	case FRAGMENT:
		for _, node := range n.Nodes {
			if node.Type == "input" {
				result += fmt.Sprintf("%%meta fidef node_0_%d fragment:terminal\n", node.Pos)
				c.List[fmt.Sprintf("node_0_%d", node.Pos)] = "input"
				result += fmt.Sprintf("%%meta filinkdef input_%d type:fl\n", node.Pos)
				result += fmt.Sprintf("%%meta filinkatt input_%d fi:node_0_%d, type:input, index:0\n", node.Pos, node.Pos)
				result += fmt.Sprintf("%%meta filinkatt input_%d fi:ext, type:input, index:%d\n", node.Pos, node.Pos)
			} else if node.Type == "output" {
				result += fmt.Sprintf("%%meta fidef node_%d_%d fragment:terminal\n", node.Layer, node.Pos)
				c.List[fmt.Sprintf("node_%d_%d", node.Layer, node.Pos)] = "output"
				result += fmt.Sprintf("%%meta filinkdef output_%d type:fl\n", node.Pos)
				result += fmt.Sprintf("%%meta filinkatt output_%d fi:node_%d_%d, type:output, index:0\n", node.Pos, node.Layer, node.Pos)
				result += fmt.Sprintf("%%meta filinkatt output_%d fi:ext, type:output, index:%d\n", node.Pos, node.Pos)
			} else {
				if neuron, ok := n.Neurons[node.Type]; ok {
					result += fmt.Sprintf("%%meta fidef node_%d_%d fragment:%s", node.Layer, node.Pos, node.Type)
					c.List[fmt.Sprintf("node_%d_%d", node.Layer, node.Pos)] = node.Type
					for _, param := range neuron.Params {
						switch param {
						case "inputs":
							result += fmt.Sprintf(", inputs:%d", node.Inputs)
						case "outputs":
							result += fmt.Sprintf(", outputs:%d", node.Outputs)
						case "bias":
							result += fmt.Sprintf(", bias:"+c.TypePrefix+"%f", node.Bias)
						case "pos":
							result += fmt.Sprintf(", pos:%d", node.Pos)
						default:
							if value, ok := c.Params[param]; ok {
								result += fmt.Sprintf(", %s:%s", param, value)
							} else {
								return "", errors.New("Unknown parameter " + param)
							}
						}
					}
					result += "\n"
				}
			}
		}

		for _, weight := range n.Weights {
			weightFI := fmt.Sprintf("weightfi_%d_%d__%d_%d", weight.Layer-1, weight.PosPrevLayer, weight.Layer, weight.PosCurrLayer)
			downNode := fmt.Sprintf("node_%d_%d", weight.Layer-1, weight.PosPrevLayer)
			upNode := fmt.Sprintf("node_%d_%d", weight.Layer, weight.PosCurrLayer)
			result += fmt.Sprintf("%%meta fidef %s fragment:weight", weightFI)
			for _, param := range n.Neurons["weight"].Params {
				switch param {
				case "weight":
					result += fmt.Sprintf(", weight:"+c.TypePrefix+"%f", weight.Value)
				default:
					if value, ok := c.Params[param]; ok {
						result += fmt.Sprintf(", %s:%s", param, value)
					} else {
						return "", errors.New("Unknown parameter " + param)
					}
					result += "\n"
				}
			}
			c.List[weightFI] = "weight"
			result += fmt.Sprintf("%%meta filinkdef up%s type:fi\n", weightFI)
			result += fmt.Sprintf("%%meta filinkdef down%s type:fi\n", weightFI)
			result += fmt.Sprintf("%%meta filinkatt down%s fi:%s, type:input, index:0\n", weightFI, weightFI)
			result += fmt.Sprintf("%%meta filinkatt down%s fi:%s, type:output, index:0\n", weightFI, downNode)
			result += fmt.Sprintf("%%meta filinkatt up%s fi:%s, type:input, index:%d\n", weightFI, upNode, weight.RelPosUp)
			result += fmt.Sprintf("%%meta filinkatt up%s fi:%s, type:output, index:0\n", weightFI, weightFI)
		}

		ProcessedNodes := make(map[string]struct{})
		for node, _ := range c.List {
			ProcessedNodes[node] = struct{}{}
		}

		// Processing pruned nodes
		for _, pNode := range c.Pruned {
			if _, ok := ProcessedNodes[pNode]; !ok {
				return "", errors.New("Pruned node " + pNode + " not found")
			} else {
				result += fmt.Sprintf("%%meta fidef %s pruned:true\n", pNode)
				delete(ProcessedNodes, pNode)
			}
		}

		// Processing collapsed nodes
		for _, cGroup := range c.Collapsed {
			cName := ""
			cCode := ""
			for _, cNode := range cGroup {
				if _, ok := ProcessedNodes[cNode]; !ok {
					return "", errors.New("Collapsed node " + cNode + " not found")
				} else {
					cName += cNode + "_"
					cCode += ":" + cNode
					delete(ProcessedNodes, cNode)
				}
			}
			result += fmt.Sprintf("%%meta cpdef %s fragcollapse%s\n", cName[:len(cName)-1], cCode)
		}

		// Processing remaining nodes
		for node := range ProcessedNodes {
			result += fmt.Sprintf("%%meta cpdef %s fragcollapse:%s\n", node, node)
		}

	default:
		return "", errors.New("unknown operating mode")
	}

	return result, nil
}
