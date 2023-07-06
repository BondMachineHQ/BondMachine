package basm

import (
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

func (bi *BasmInstance) LLVMGetLocalID(inst value.Value) (int64, error) {
	switch inst := inst.(type) {
	case *ir.InstAlloca:
		return inst.LocalIdent.LocalID, nil
	case *ir.InstLoad:
		return inst.LocalIdent.LocalID, nil
	case *ir.InstAdd:
		return inst.LocalIdent.LocalID, nil
	case *ir.Param:
		return inst.LocalIdent.LocalID, nil
	}
	return 0, fmt.Errorf("Could not get local ID for %v", inst)
}

func (bi *BasmInstance) ParseAssemblyLLVM(filePath string) error {

	if bi.debug {
		fmt.Println(purple("Phase 0") + ": " + red("Reading LLVM IR file "+filePath))
	}

	m, err := asm.ParseFile(filePath)
	if err != nil {
		return err
	}

	for _, f := range m.Funcs {
		name := f.GlobalIdent.Ident()
		fName := strings.ReplaceAll(name, "@", "")
		if bi.debug {
			fmt.Println("\t" + red("Adding LLVM Function "+name+" as fragment "+fName))
		}

		assocMap := make(map[int64]string)

		newFragment := new(BasmFragment)
		newFragment.fragmentName = fName
		newFragment.fragmentBody = new(bmline.BasmBody)
		newFragment.fragmentBody.Lines = make([]*bmline.BasmLine, 0)

		resin := make([]string, len(f.Params))

		for i, param := range f.Params {
			assocMap[param.LocalID] = "r" + fmt.Sprint(param.LocalID)
			resin[i] = "r" + fmt.Sprint(param.LocalID)
		}

		newFragment.fragmentBody.BasmMeta = newFragment.fragmentBody.BasmMeta.SetMeta("template", "false")
		newFragment.fragmentBody.BasmMeta = newFragment.fragmentBody.BasmMeta.SetMeta("llvm", "true")
		if len(resin) > 0 {
			newFragment.fragmentBody.BasmMeta = newFragment.fragmentBody.BasmMeta.SetMeta("resin", strings.Join(resin, ":"))
		}

		for _, bb := range f.Blocks {
			for _, inst := range bb.Insts {
				switch inst := inst.(type) {
				case *ir.InstAlloca:
					dest := "r" + fmt.Sprint(inst.LocalIdent.LocalID)

					newLine := new(bmline.BasmLine)
					newLine.Operation = new(bmline.BasmElement)
					newLine.Operation.SetValue("clr")
					newLine.Elements = make([]*bmline.BasmElement, 1)
					newLine.Elements[0] = new(bmline.BasmElement)
					newLine.Elements[0].SetValue(dest)

					newFragment.fragmentBody.Lines = append(newFragment.fragmentBody.Lines, newLine)

				case *ir.InstLoad:

					dest := "r" + fmt.Sprint(inst.LocalIdent.LocalID)
					srcP := inst.Src
					id, err := bi.LLVMGetLocalID(srcP)
					if err != nil {
						return err
					}
					src := "r" + fmt.Sprint(id)
					newLine := new(bmline.BasmLine)
					newLine.Operation = new(bmline.BasmElement)
					newLine.Operation.SetValue("mov")
					newLine.Elements = make([]*bmline.BasmElement, 2)
					newLine.Elements[0] = new(bmline.BasmElement)
					newLine.Elements[0].SetValue(dest)
					newLine.Elements[1] = new(bmline.BasmElement)
					newLine.Elements[1].SetValue(src)

					newFragment.fragmentBody.Lines = append(newFragment.fragmentBody.Lines, newLine)

				case *ir.InstStore:
					srcP := inst.Src
					srcId, err := bi.LLVMGetLocalID(srcP)
					if err != nil {
						return err
					}
					src := "r" + fmt.Sprint(srcId)
					destP := inst.Dst
					destId, err := bi.LLVMGetLocalID(destP)
					if err != nil {
						return err
					}
					dest := "r" + fmt.Sprint(destId)
					newLine := new(bmline.BasmLine)
					newLine.Operation = new(bmline.BasmElement)
					newLine.Operation.SetValue("mov")
					newLine.Elements = make([]*bmline.BasmElement, 2)
					newLine.Elements[0] = new(bmline.BasmElement)
					newLine.Elements[0].SetValue(dest)
					newLine.Elements[1] = new(bmline.BasmElement)
					newLine.Elements[1].SetValue(src)

					newFragment.fragmentBody.Lines = append(newFragment.fragmentBody.Lines, newLine)

				case *ir.InstAdd:
					var dest string
					src1P := inst.X
					src1Id, err := bi.LLVMGetLocalID(src1P)
					if err != nil {
						return err
					}
					src1 := "r" + fmt.Sprint(src1Id)
					src2P := inst.Y
					src2Id, err := bi.LLVMGetLocalID(src2P)
					if err != nil {
						return err
					}
					src2 := "r" + fmt.Sprint(src2Id)
					destP := inst.LocalIdent
					dest = "r" + fmt.Sprint(destP.LocalID)

					newLine1 := new(bmline.BasmLine)
					newLine1.Operation = new(bmline.BasmElement)
					newLine1.Operation.SetValue("mov")
					newLine1.Elements = make([]*bmline.BasmElement, 2)
					newLine1.Elements[0] = new(bmline.BasmElement)
					newLine1.Elements[0].SetValue(dest)
					newLine1.Elements[1] = new(bmline.BasmElement)
					newLine1.Elements[1].SetValue(src1)

					newFragment.fragmentBody.Lines = append(newFragment.fragmentBody.Lines, newLine1)

					newLine := new(bmline.BasmLine)
					newLine.Operation = new(bmline.BasmElement)
					newLine.Operation.SetValue("add")
					newLine.Elements = make([]*bmline.BasmElement, 2)
					newLine.Elements[0] = new(bmline.BasmElement)
					newLine.Elements[0].SetValue(dest)
					newLine.Elements[1] = new(bmline.BasmElement)
					newLine.Elements[1].SetValue(src2)

					newFragment.fragmentBody.Lines = append(newFragment.fragmentBody.Lines, newLine)

				default:
					return fmt.Errorf("unknown inst type %T", inst)

				}

				var ret string
				retP := bb.Term
				switch retP := retP.(type) {
				case *ir.TermRet:
					X := retP.X
					XId, err := bi.LLVMGetLocalID(X)
					if err != nil {
						return err
					}
					ret = "r" + fmt.Sprint(XId)
				}
				newFragment.fragmentBody.BasmMeta = newFragment.fragmentBody.BasmMeta.SetMeta("resout", ret)

			}
		}
		//pretty.Println(f)

		bi.fragments[fName] = newFragment
	}

	return nil
}
