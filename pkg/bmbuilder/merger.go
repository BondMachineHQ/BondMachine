package bmbuilder

import (
	"errors"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

const (
	FinalBMInput = uint8(1) + iota
	FinalBMOutput
	FirstBMInput
	FirstBMOutput
	SecondBMInput
	SecondBMOutput
)

type BMEndPoint struct {
	BType uint8
	BNum  int
}

type BMLink struct {
	E1 BMEndPoint
	E2 BMEndPoint
}

type BMConnections struct {
	Links []BMLink
}

func (e BMEndPoint) String() string {
	switch e.BType {
	case FinalBMInput:
		return fmt.Sprintf("i%d", e.BNum)
	case FinalBMOutput:
		return fmt.Sprintf("o%d", e.BNum)
	case FirstBMInput:
		return fmt.Sprintf("bm1i%d", e.BNum)
	case FirstBMOutput:
		return fmt.Sprintf("bm1o%d", e.BNum)
	case SecondBMInput:
		return fmt.Sprintf("bm2i%d", e.BNum)
	case SecondBMOutput:
		return fmt.Sprintf("bm2o%d", e.BNum)
	}
	return fmt.Sprintf("Unknown(%d)", e.BNum)
}

func (c *BMConnections) String() string {
	str := ""
	for _, link := range c.Links {
		str += fmt.Sprintf("\t%v -> %v\n", link.E1, link.E2)
	}
	return str
}

type bmOrigin struct {
	originBM int
	originID int
}

type chMap map[int]bmOrigin

func (o chMap) getNewID(bm int, idx int) int {
	for i, or := range o {
		if or.originBM == bm && or.originID == idx {
			return i
		}
	}
	return -1
}

func followLink(bm1 *bondmachine.Bondmachine, bm2 *bondmachine.Bondmachine, result *bondmachine.Bondmachine, intOutMap chMap, endPoint BMEndPoint) (int, error) {
	// This function finds out the internal output that corresponds to the given endpoint

	switch endPoint.BType {
	case FinalBMInput:
		for i, bond := range result.Internal_outputs {
			if bond.Map_to == bondmachine.BMINPUT && bond.Res_id == endPoint.BNum {
				return i, nil
			}
		}
	case FirstBMOutput:
		for i, bond := range bm1.Internal_inputs {
			if bond.Map_to == bondmachine.BMOUTPUT && bond.Res_id == endPoint.BNum {
				if bm1.Links[i] == -1 {
					return -1, nil
				}
				return intOutMap.getNewID(1, bm1.Links[i]), nil
			}
		}
	case SecondBMOutput:
		for i, bond := range bm2.Internal_inputs {
			if bond.Map_to == bondmachine.BMOUTPUT && bond.Res_id == endPoint.BNum {
				if bm2.Links[i] == -1 {
					return -1, nil
				}
				return intOutMap.getNewID(2, bm2.Links[i]), nil
			}
		}
	}
	return 0, errors.New("endpoint not Found")
}

func (bld *BMBuilder) BMMerge(bm1 *bondmachine.Bondmachine, bm2 *bondmachine.Bondmachine, l *BMConnections) (*bondmachine.Bondmachine, error) {
	// Merge two bondmachines
	// bm1 and bm2 are the bondmachines to merge
	// Returns the merged bondmachine

	if bm1.Rsize != bm2.Rsize {
		return nil, errors.New("bondmachines have different register sizes")
	}

	if bld.debug {
		fmt.Println("Connections:\n", l)
	}

	result := new(bondmachine.Bondmachine)
	result.Rsize = bm1.Rsize

	result.Init()

	mergedInputs := 0
	mergedOutputs := 0
	for _, link := range l.Links {
		if link.E1.BType == FinalBMInput {
			if link.E1.BNum >= mergedInputs {
				mergedInputs = link.E1.BNum + 1
			}
		}
		if link.E1.BType == FinalBMOutput {
			if link.E1.BNum >= mergedOutputs {
				mergedOutputs = link.E1.BNum + 1
			}
		}

		if link.E2.BType == FinalBMInput {
			if link.E2.BNum >= mergedInputs {
				mergedInputs = link.E2.BNum + 1
			}
		}
		if link.E2.BType == FinalBMOutput {
			if link.E2.BNum >= mergedOutputs {
				mergedOutputs = link.E2.BNum + 1
			}
		}
	}

	result.Inputs = mergedInputs
	result.Outputs = mergedOutputs

	if bld.debug {
		fmt.Printf("Inputs of the merged bondmachine: %d\n", mergedInputs)
		fmt.Printf("Outputs of the merged bondmachine: %d\n", mergedOutputs)
	}

	// Processing domains
	domainsMap := make(chMap)

	result.Domains = make([]*procbuilder.Machine, len(bm1.Domains)+len(bm2.Domains))

	for i, dom := range bm1.Domains {
		result.Domains[i] = dom
		domainsMap[i] = bmOrigin{originBM: 1, originID: i}
	}
	for i, dom := range bm2.Domains {
		result.Domains[i+len(bm1.Domains)] = dom
		domainsMap[i+len(bm1.Domains)] = bmOrigin{originBM: 2, originID: i}
	}

	if bld.debug {
		fmt.Println("Domains map: ", domainsMap)
		fmt.Println("Domains: ", result.Domains)
	}

	// Processing processors
	procsMap := make(chMap)

	result.Processors = make([]int, len(bm1.Processors)+len(bm2.Processors))

	for i, domId := range bm1.Processors {
		result.Processors[i] = domainsMap.getNewID(1, i)
		procsMap[i] = bmOrigin{originBM: 1, originID: domId}
	}

	for i, domId := range bm2.Processors {
		result.Processors[i+len(bm1.Processors)] = domainsMap.getNewID(2, domId)
		procsMap[i+len(bm1.Processors)] = bmOrigin{originBM: 2, originID: i}
	}

	if bld.debug {
		fmt.Println("Processors map: ", procsMap)
		fmt.Println("Processors: ", result.Processors)
	}

	if bld.debug {
		fmt.Println("BM 1 internal inputs: ", bm1.Internal_inputs)
		fmt.Println("BM 1 internal outputs: ", bm1.Internal_outputs)
		fmt.Println("BM 2 internal inputs: ", bm2.Internal_inputs)
		fmt.Println("BM 2 internal outputs: ", bm2.Internal_outputs)
		fmt.Println("BM 1 links: ", bm1.Links)
		fmt.Println("BM 2 links: ", bm2.Links)
	}

	// Processing shared objects
	soMap := make(chMap)
	result.Shared_objects = make([]bondmachine.Shared_instance, len(bm1.Shared_objects)+len(bm2.Shared_objects))

	for i, so := range bm1.Shared_objects {
		result.Shared_objects[i] = so
		soMap[i] = bmOrigin{originBM: 1, originID: i}
	}

	for i, so := range bm2.Shared_objects {
		result.Shared_objects[i+len(bm1.Shared_objects)] = so
		soMap[i+len(bm1.Shared_objects)] = bmOrigin{originBM: 2, originID: i}
	}

	result.Shared_links = make([]bondmachine.Shared_instance_list, len(result.Processors))
	for i, _ := range result.Processors {
		soList := make(bondmachine.Shared_instance_list, 0)
		bm := procsMap[i].originBM
		origProc := procsMap[i].originID
		switch bm {
		case 1:
			for _, so := range bm1.Shared_links[origProc] {
				soList = append(soList, soMap.getNewID(1, so))
			}
		case 2:
			for _, so := range bm2.Shared_links[origProc] {
				soList = append(soList, soMap.getNewID(2, so))
			}
		}
		result.Shared_links[i] = soList
	}

	// Merging internal outputs

	intOutMap := make(chMap)
	result.Internal_outputs = make([]bondmachine.Bond, mergedInputs)

	for i := 0; i < mergedInputs; i++ {
		newBond := bondmachine.Bond{Map_to: bondmachine.BMINPUT, Res_id: i, Ext_id: 0}
		result.Internal_outputs[i] = newBond
	}

	no := mergedInputs
	for i, bond := range bm1.Internal_outputs {
		switch bond.Map_to {
		case bondmachine.BMINPUT:
			if bld.debug {
				fmt.Println("Ignoring internal output ", bond, " from BM 1")
			}
		case bondmachine.CPOUTPUT:
			newBond := bondmachine.Bond{Map_to: bondmachine.CPOUTPUT, Res_id: procsMap.getNewID(1, bond.Res_id), Ext_id: bond.Ext_id}
			result.Internal_outputs = append(result.Internal_outputs, newBond)
			intOutMap[no] = bmOrigin{originBM: 1, originID: i}
			no++
		}
	}

	for i, bond := range bm2.Internal_outputs {
		switch bond.Map_to {
		case bondmachine.BMINPUT:
			if bld.debug {
				fmt.Println("Ignoring internal output ", bond, " from BM 2")
			}
		case bondmachine.CPOUTPUT:
			newBond := bondmachine.Bond{Map_to: bondmachine.CPOUTPUT, Res_id: procsMap.getNewID(2, bond.Res_id), Ext_id: bond.Ext_id}
			result.Internal_outputs = append(result.Internal_outputs, newBond)
			intOutMap[no] = bmOrigin{originBM: 2, originID: i}
			no++
		}
	}

	if bld.debug {
		fmt.Println("Internal outputs: ", result.Internal_outputs)
		fmt.Println("\t", intOutMap)
	}

	// Merging internal inputs
	intInMap := make(chMap)
	result.Internal_inputs = make([]bondmachine.Bond, mergedOutputs)

	for i := 0; i < mergedOutputs; i++ {
		newBond := bondmachine.Bond{Map_to: bondmachine.BMOUTPUT, Res_id: i, Ext_id: 0}
		result.Internal_inputs[i] = newBond
	}

	ni := mergedOutputs
	for i, bond := range bm1.Internal_inputs {
		switch bond.Map_to {
		case bondmachine.BMOUTPUT:
			if bld.debug {
				fmt.Println("Ignoring internal input ", bond, " from BM 1")
			}

		case bondmachine.CPINPUT:
			newBond := bondmachine.Bond{Map_to: bondmachine.CPINPUT, Res_id: procsMap.getNewID(1, bond.Res_id), Ext_id: bond.Ext_id}
			result.Internal_inputs = append(result.Internal_inputs, newBond)
			intInMap[ni] = bmOrigin{originBM: 1, originID: i}
			ni++
		}
	}

	for i, bond := range bm2.Internal_inputs {
		switch bond.Map_to {
		case bondmachine.BMOUTPUT:
			if bld.debug {
				fmt.Println("Ignoring internal input ", bond, " from BM 2")
			}
		case bondmachine.CPINPUT:
			newBond := bondmachine.Bond{Map_to: bondmachine.CPINPUT, Res_id: procsMap.getNewID(2, bond.Res_id), Ext_id: bond.Ext_id}
			result.Internal_inputs = append(result.Internal_inputs, newBond)
			intInMap[ni] = bmOrigin{originBM: 2, originID: i}
			ni++
		}
	}

	if bld.debug {
		fmt.Println("Internal inputs: ", result.Internal_inputs)
	}

	// Merging links
	result.Links = make([]int, len(result.Internal_inputs))

	for i, bond := range result.Internal_inputs {
		if bld.debug {
			fmt.Println("Following internal input: ", bond)
		}

		switch bond.Map_to {
		case bondmachine.BMOUTPUT:
			for _, link := range l.Links {
				if link.E2.BType == FinalBMOutput && link.E2.BNum == bond.Res_id {
					if bld.debug {
						fmt.Println("Linking internal input", bond, "to", link.E1)
					}
					if io, err := followLink(bm1, bm2, result, intOutMap, link.E1); err == nil {
						result.Links[i] = io
					} else {
						return nil, err
					}

					break
				}

				if link.E1.BType == FinalBMOutput && link.E1.BNum == bond.Res_id {
					if bld.debug {
						fmt.Println("Linking internal input", bond, "to", link.E2)
					}
					if io, err := followLink(bm1, bm2, result, intOutMap, link.E2); err == nil {
						result.Links[i] = io
					} else {
						return nil, err
					}

					break
				}
			}
		case bondmachine.CPINPUT:
			origIntIn := intInMap[i]
			switch origIntIn.originBM {
			case 1:
				origIntOut := bm1.Links[origIntIn.originID]
				origBond := bm1.Internal_outputs[origIntOut]
				switch origBond.Map_to {
				case bondmachine.BMINPUT:
					for _, link := range l.Links {
						if link.E2.BType == FirstBMInput && link.E2.BNum == origBond.Res_id {
							if bld.debug {
								fmt.Println("Linking internal input", bond, "to", link.E1)
							}
							if io, err := followLink(bm1, bm2, result, intOutMap, link.E1); err == nil {
								result.Links[i] = io
							} else {
								return nil, err
							}

							break
						}

						if link.E1.BType == FirstBMInput && link.E1.BNum == origBond.Res_id {
							if bld.debug {
								fmt.Println("Linking internal input", bond, "to", link.E2)
							}
							if io, err := followLink(bm1, bm2, result, intOutMap, link.E2); err == nil {
								result.Links[i] = io
							} else {
								return nil, err
							}

							break
						}
					}
				case bondmachine.CPOUTPUT:
					result.Links[i] = intOutMap.getNewID(1, origIntOut)
				}
			case 2:
				origIntOut := bm2.Links[origIntIn.originID]
				origBond := bm2.Internal_outputs[origIntOut]
				switch origBond.Map_to {
				case bondmachine.BMINPUT:
					for _, link := range l.Links {
						if link.E2.BType == SecondBMInput && link.E2.BNum == origBond.Res_id {
							if bld.debug {
								fmt.Println("Linking internal input", bond, "to", link.E1)
							}
							if io, err := followLink(bm1, bm2, result, intOutMap, link.E1); err == nil {
								result.Links[i] = io
							} else {
								return nil, err
							}

							break
						}

						if link.E1.BType == SecondBMInput && link.E1.BNum == origBond.Res_id {
							if bld.debug {
								fmt.Println("Linking internal input", bond, "to", link.E2)
							}
							if io, err := followLink(bm1, bm2, result, intOutMap, link.E2); err == nil {
								result.Links[i] = io
							} else {
								return nil, err
							}

							break
						}
					}
				case bondmachine.CPOUTPUT:
					result.Links[i] = intOutMap.getNewID(2, origIntOut)
				}
			}
		}

	}

	return result, nil
}
