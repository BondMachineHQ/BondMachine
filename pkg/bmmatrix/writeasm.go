package bmmatrix

import (
	"errors"
)

func (mo *MatrixOperations) WriteBasm() (string, error) {
	if mo == nil {
		return "", errors.New("MatrixOperations is nil")
	}

	return mo.Result, nil
}
