package bmmatrix

import (
	"github.com/mmirko/mel/pkg/mel"
)

func (mo *MatrixOperations) RunEngine() error {

	a := new(M3BasmMatrix)
	var env interface{} = newExporterEnv()
	a.Mel3Object.Environment = env

	var ep *mel.EvolutionParameters
	c := new(mel.MelConfig)
	c.Debug = true
	a.MelInit(c, ep)

	iString := mo.Expression

	a.MelStringImport(iString)
	if err := a.Compute(); err != nil {
		return err
	}

	mo.Result = *(a.Mel3Object.Environment.(exporterEnv).basmCode)

	return nil
}
