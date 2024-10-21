package bmmatrix

import (
	"fmt"
	"testing"

	"github.com/mmirko/mel/pkg/mel"
)

func TestM3numberImporter(t *testing.T) {

	fmt.Println("---- Test: M3number importer ----")

	a := new(M3BasmMatrix)
	var ep *mel.EvolutionParameters
	a.MelInit(nil, ep)

	istrings := []string{
		`
m(54)

`,
		`
mult(
	m(3),
	m(1)
)

`}

	for i := 0; i < len(istrings); i++ {
		fmt.Println(">>>")
		fmt.Println(istrings[i])
		err := a.MelStringImport(istrings[i])
		fmt.Println("---")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			a.MelDump(nil)
		}
		fmt.Println("<<<")
	}

	fmt.Println("---- End test ----")

}
