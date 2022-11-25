package bmnumbers

import "errors"

func (n *BMNumber) ExportSingleUint() (uint64, error) {
	if n != nil && n.number != nil && len(n.number) != 1 {
		return 0, errors.New("number is not a single uint64")
	}
	return n.number[0], nil
}
