package bondirect

import (
	"fmt"
	"strconv"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
)

// ShowMessages displays the message flow between nodes in the cluster
func ShowMessages(c *Config, mesh *Mesh, cluster *bmcluster.Cluster) {
	rawMess := cluster.GetMessages()

	for _, rawMes := range rawMess {
		from := rawMes.From
		to := rawMes.To

		bmFrom := strconv.Itoa(from.BmId)
		idxFrom := strconv.Itoa(from.Index)

		bmTo := strconv.Itoa(to.BmId)
		idxTo := strconv.Itoa(to.Index)

		fmt.Printf("bm%sidx%stobm%sidx%sdata\n", bmFrom, idxFrom, bmTo, idxTo)
		fmt.Printf("bm%sidx%stobm%sidx%svalid\n", bmFrom, idxFrom, bmTo, idxTo)
		fmt.Printf("bm%sidx%stobm%sidx%srecv\n", bmTo, idxTo, bmFrom, idxFrom)
	}
}
