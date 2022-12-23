package bondmachine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/mmirko/mel"
)

func TestEvolutionaryBondmachine(t *testing.T) {
	rand.Seed(int64(time.Now().Unix()))
	ep := new(mel.EvolutionParameters)
	ep.Pars = make(map[string]string)
	fmt.Println(ep)
}

func TestEvolutionaryFitnell(t *testing.T) {
	rand.Seed(int64(time.Now().Unix()))
	ep := new(mel.EvolutionParameters)
	ep.Pars = make(map[string]string)
	fmt.Println(ep)
}
