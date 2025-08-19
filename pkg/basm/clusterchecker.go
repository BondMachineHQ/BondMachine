package basm

import (
	"fmt"
	"strconv"

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
		devIdS := ""
		devName := "default"

		if cp.GetMeta("device") != "" {
			devName = cp.GetMeta("device")
		}
		if cp.GetMeta("devid") != "" {
			devIdS = cp.GetMeta("devid")
			if value, err := strconv.Atoi(devIdS); err != nil {
				return fmt.Errorf("invalid device id %q: %w", devIdS, err)
			} else {
				devId = value
			}
		}

		if _, ok := bi.clusteredNames[devName]; ok {
			if bi.clusteredNames[devName] != devId {
				return fmt.Errorf("device name %q is already associated with a different id: %d", devName, bi.clusteredNames[devName])
			}
		} else {
			bi.clusteredNames[devName] = devId

			if bi.debug {
				fmt.Println(green("\t\t\tAdding new device:"), red(devName), green("id"), blue(devId))
			}
		}
	}

	bi.clusteredBondMachines = make([]string, len(bi.clusteredNames))

	for i := 0; i < len(bi.clusteredNames); i++ {
		found := false
		for _, devId := range bi.clusteredNames {
			if i == devId {
				bi.clusteredBondMachines[i] = ""
				newMap := new(bondmachine.IOmap)
				newMap.Assoc = make(map[string]string)
				bi.clusteredMaps = append(bi.clusteredMaps, newMap)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("device id %d not found in clustered names", i)
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
