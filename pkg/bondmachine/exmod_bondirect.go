package bondmachine

import (
	//	"fmt"
	"fmt"
	"strconv"

	//	"strings"
	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
	"github.com/BondMachineHQ/BondMachine/pkg/bondirect"
)

type Bondirect_extra struct {
	// Config  *bondirect.Config
	// Cluster *bmcluster.Cluster
	// Mesh    *bondirect.Mesh
	*bondirect.BondirectElement
	PeerID uint32
	Maps   *IOmap
	Flavor string
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
			result.Params["peer_name"], _ = sl.GetMeshNodeName(peer.PeerName)
			break
		}

	}

	result.Params["input_ids"] = ""
	result.Params["inputs"] = ""
	result.Params["sources"] = ""
	fmt.Println("mypeer", mypeer)
	fmt.Println("cluster", sl.Cluster)
	fmt.Println("ppp", sl.Maps)
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
	return []string{}, []string{}
}
