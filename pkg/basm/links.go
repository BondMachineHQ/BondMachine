package basm

import (
	"errors"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

const (
	IOLINK = uint8(0) + iota
	FILINK
)

type endpoint struct {
	seq      int
	name     string
	index    string
	typ      string
	linkName string
}

func (bi *BasmInstance) GetLinks(t uint8, name string, index string, dir string) ([]string, error) {
	result := make([]string, 0)
	switch t {
	case IOLINK:
		for _, link := range bi.ioAttach {
			if link.GetMeta("cp") == name && link.GetMeta("index") == index && link.GetMeta("type") == dir {
				result = append(result, link.GetValue())
			}
		}
	case FILINK:
		for _, fi := range bi.fiLinkAttach {
			if fi.GetMeta("fi") == name && fi.GetMeta("index") == index && fi.GetMeta("type") == dir {
				result = append(result, fi.GetValue())
			}
		}
	}
	return result, nil
}

func (bi *BasmInstance) GetEndpoints(t uint8, name string) (in endpoint, out endpoint, err error) {
	switch t {
	case FILINK:
		inok := false
		outok := false
		isExt := false
		in = endpoint{}
		out = endpoint{}
		for i, fi := range bi.fiLinkAttach {
			if fi.GetValue() == name {
				dir := fi.GetMeta("type")
				if fi.GetMeta("fi") == "ext" {
					isExt = true
					if dir == "input" {
						dir = "output"
					} else {
						dir = "input"
					}
				}
				switch dir {
				case "input":
					if inok {
						return in, out, errors.New("Multiple input endpoints found")
					}
					in.name = fi.GetMeta("fi")
					in.index = fi.GetMeta("index")
					in.linkName = fi.GetValue()
					if isExt {
						in.typ = "output"
					} else {
						in.typ = "input"
					}
					in.seq = i
					inok = true
				case "output":
					if outok {
						return in, out, errors.New("Multiple output endpoints found")
					}
					out.name = fi.GetMeta("fi")
					out.index = fi.GetMeta("index")
					out.linkName = fi.GetValue()
					if isExt {
						out.typ = "input"
					} else {
						out.typ = "output"
					}
					out.seq = i
					outok = true
				default:
					err = errors.New("Invalid direction")
					return in, out, err
				}
			}
		}
		if inok && outok {
			return in, out, nil
		} else {
			err = errors.New("Invalid endpoint")
			return in, out, err
		}
	}
	return in, out, errors.New("Invalid link type")
}

func (bi *BasmInstance) GetFI(finame string) (int, *bmline.BasmElement, error) {
	for i, fi := range bi.fis {
		if fi.GetValue() == finame {
			return i, fi, nil
		}
	}
	return -1, nil, errors.New("FI not found")
}
