package bondirect

import (
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bmcluster"
)

func ShowMessages(c *Config, mesh *Mesh, cluster *bmcluster.Cluster) {
	fmt.Println(cluster.GetMessages())
}
