package basm

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

func clusterChecker(bi *BasmInstance) error {

	if bi.debug {
		fmt.Println(green("\tProcessing cp:"))
	}

	bi.clusteredBondMachines = make([]string, 0)
	bi.clusteredNames = make(map[string]int)
	bi.clusteredMaps = make([]*bondmachine.IOmap, 0)
	bi.cluster = new(Cluster)
	bi.cluster.ClusterId = 0
	bi.cluster.Peers = make([]Peer, 0)

	for _, cp := range bi.cps {
		cpName := cp.GetValue()
		if bi.debug {
			fmt.Println(green("\t\tProcessing cp: "), cpName)
		}

		devId := 0
		devName := "default"

		if cp.GetMeta("device") != "" {
			devName = cp.GetMeta("device")
		}
		if _, ok := bi.clusteredNames[devName]; !ok {
			devId = len(bi.clusteredBondMachines)
			bi.clusteredBondMachines = append(bi.clusteredBondMachines, "")
			bi.clusteredNames[devName] = devId
			newMap := new(bondmachine.IOmap)
			newMap.Assoc = make(map[string]string)
			bi.clusteredMaps = append(bi.clusteredMaps, newMap)

			if bi.debug {
				fmt.Println(green("\t\t\tAdding new device:"), red(devName), green("id"), blue(devId))
			}
		}
	}

	if len(bi.clusteredBondMachines) > 1 {
		bi.isClustered = true
		if bi.debug {
			fmt.Println(green("\t\tClustered BondMachine detected"))
		}
	} else {
		bi.isClustered = false
		if bi.debug {
			fmt.Println(green("\t\tSingle BondMachine detected"))
		}
	}

	return nil
}
