package bondgo

import (
	"fmt"
	"go/ast"
)

type BondgoFunctions struct {
	*BondgoConfig
	*BondgoMessages
	Functions map[string]FunctCell
}

type FunctArg struct {
	Argname string
	Argtype *VarType
}

// The type to hold functions
type FunctCell struct {
	Inputs  []FunctArg
	Outputs []FunctArg
	Body    *ast.BlockStmt
}

func (fc *FunctCell) String() string {
	result := "("
	for _, farg := range fc.Inputs {
		result += farg.Argname + " "
	}
	result += ")("
	for _, farg := range fc.Outputs {
		result += farg.Argname + " "
	}
	result += ")"
	return result
	// TODO format better
}
func (fn *BondgoFunctions) Init_Functions(cfg *BondgoConfig, ms *BondgoMessages) {
	fn.BondgoConfig = cfg
	fn.BondgoMessages = ms
	fn.Functions = make(map[string]FunctCell)
}

func (fn *BondgoFunctions) String() string {
	result := ""
	for fname, fcell := range fn.Functions {
		result += fname + " "
		result += fcell.String() + "\n"
	}
	return result
}

func (fn *BondgoFunctions) Visit(n ast.Node) ast.Visitor {

	switch n.(type) {
	case *ast.FuncDecl:
		funcDecl := n.(*ast.FuncDecl)
		fname := funcDecl.Name
		funcType := funcDecl.Type

		supported := []string{fn.Basic_type, "bool", fn.Basic_chantype, "chan bool"}

		if fn.In_debug() {
			fmt.Println("New function declaration:", fname)
		}

		inputs := make([]FunctArg, 0)
		outputs := make([]FunctArg, 0)

		//TODO Include here check for repeting names in input and outputs

		if funcType.Params != nil { // Not really necessary
			for _, param := range funcType.Params.List {
				argok := false
				for _, supp := range supported {
					gast, _ := Type_from_string(supp)
					gastc, _ := Type_from_ast(param.Type)
					if Same_Type(gast, gastc) {
						for _, vari := range param.Names {
							inputs = append(inputs, FunctArg{vari.Name, gast})
						}
						argok = true
						break
					}
				}
				if !argok {
					fn.Set_faulty("Function argument type not supported")
					return nil
				}
			}
		}

		if fn.In_debug() {
			fmt.Println("\tInputs:", inputs)
		}

		if funcType.Results != nil {
			for i, resul := range funcType.Results.List {
				argok := false
				for _, supp := range supported {
					gast, _ := Type_from_string(supp)
					gastc, _ := Type_from_ast(resul.Type)
					if Same_Type(gast, gastc) {
						if resul.Names == nil {
							outputs = append(outputs, FunctArg{"unspec" + fmt.Sprintf("%02d", i), gast})
						} else {
							for _, vari := range resul.Names {
								outputs = append(outputs, FunctArg{vari.Name, gast})
							}
						}
						argok = true
						break
					}
				}
				if !argok {
					fn.Set_faulty("Function argument type not supported")
					return nil
				}
			}
		}

		if fn.In_debug() {
			fmt.Println("\tOutputs:", outputs)
		}

		fcell := FunctCell{inputs, outputs, funcDecl.Body}
		fn.Functions[fname.Name] = fcell

		return nil
	}
	return fn
}
