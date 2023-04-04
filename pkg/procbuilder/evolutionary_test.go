package procbuilder

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mmirko/mel/pkg/mel"
)

func TestEvolutionaryMachine(t *testing.T) {
	rand.Seed(int64(time.Now().Unix()))
	ep := new(mel.EvolutionParameters)
	ep.Pars = make(map[string]string)
	ep.Pars["procbuilder:opcodes"] = "nop,rset"
	fmt.Println(Machine_Program_Generate(ep))
}
