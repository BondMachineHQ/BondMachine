package bondgo

import (
	"go/ast"
	"go/token"
	"strconv"

	//	"strings"
	"fmt"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func (bg *BondgoCheck) Expr_eval(n ast.Expr) ([]VarCell, bool) {

	switch exptype := n.(type) {
	case *ast.BasicLit:
		if exptype.Kind == token.INT {
			value := exptype.Value

			gent, _ := Type_from_string(bg.Basic_type)

			bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}

			resp := <-bg.Answers
			if resp.AnsType == ANS_OK {

				newregcell := resp.Cell

				regname := procbuilder.Get_register_name(newregcell.Id)

				bg.WriteLine(bg.CurrentRoutine, "rset "+regname+" "+value)
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "rset", I_NIL}

				result := make([]VarCell, 1)
				result[0] = newregcell

				return result, true
			} else {
				bg.Set_faulty("Resource reservation failed")
				return []VarCell{}, false
			}
		}
	case *ast.Ident:
		identname := exptype.Name

		switch identname {
		case "true":
			gent, _ := Type_from_string("bool")
			bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
			resp := <-bg.Answers
			if resp.AnsType == ANS_OK {
				newregcell := resp.Cell
				regname := procbuilder.Get_register_name(newregcell.Id)
				bg.WriteLine(bg.CurrentRoutine, "rset "+regname+" 1")
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "rset", I_NIL}
				result := make([]VarCell, 1)
				result[0] = newregcell

				return result, true
			} else {
				bg.Set_faulty("Resource reservation failed")
				return []VarCell{}, false
			}
		case "false":
			gent, _ := Type_from_string("bool")
			bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
			resp := <-bg.Answers
			if resp.AnsType == ANS_OK {
				newregcell := resp.Cell
				regname := procbuilder.Get_register_name(newregcell.Id)
				bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
				result := make([]VarCell, 1)
				result[0] = newregcell

				return result, true
			} else {
				bg.Set_faulty("Resource reservation failed")
				return []VarCell{}, false
			}
		}

		varexist := false
		for scope := bg; scope != nil; scope = scope.Outer {
			if cell, ok := scope.Vars[identname]; ok {
				varexist = true

				var newregcell VarCell

				switch cell.Procobjtype {
				case REGISTER:
					gent, _ := Type_from_string(bg.Basic_type)
					bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newregcell = resp.Cell
						regname := procbuilder.Get_register_name(newregcell.Id)
						oldregname := procbuilder.Get_register_name(cell.Id)
						bg.WriteLine(bg.CurrentRoutine, "cpy "+regname+" "+oldregname)
						bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "cpy", I_NIL}
					} else {
						bg.Set_faulty("Resource reservation failed")
						return []VarCell{}, false
					}
				case MEMORY:
					gent, _ := Type_from_string(bg.Basic_type)
					bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newregcell = resp.Cell
						regname := procbuilder.Get_register_name(newregcell.Id)
						bg.WriteLine(bg.CurrentRoutine, "m2r "+regname+" "+strconv.Itoa(cell.Id))
						bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "m2r", I_NIL}
					} else {
						bg.Set_faulty("Resource reservation failed")
						return []VarCell{}, false
					}
				case INPUT:
					newregcell = cell
				case OUTPUT:
					newregcell = cell
				case CHANNEL:
					newregcell = cell
				}

				result := make([]VarCell, 1)
				result[0] = newregcell

				return result, true
			}
		}
		if !varexist {
			bg.Set_faulty("Variable " + identname + " not defined")
			return []VarCell{}, false
		}
	case *ast.UnaryExpr:
		x := exptype.X
		if cell, ok := bg.Expr_eval(x); !ok {
			bg.Set_faulty("Wrong expression")
			return []VarCell{}, false
		} else {
			if len(cell) == 1 {
				switch exptype.Op {
				case token.ARROW:
					gent, _ := Type_from_string(bg.Basic_chantype)
					if Same_Type(cell[0].Vtype, gent) {
						gent, _ := Type_from_string(bg.Basic_type)
						bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
						resp := <-bg.Answers

						if resp.AnsType == ANS_OK {
							bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
							resp2 := <-bg.Answers
							if resp2.AnsType == ANS_OK {
								newregcell := resp.Cell
								newregcell2 := resp2.Cell
								destname := procbuilder.Get_register_name(newregcell.Id)
								destname2 := procbuilder.Get_register_name(newregcell2.Id)
								channame := procbuilder.Get_channel_name(cell[0].Id)
								bg.WriteLine(bg.CurrentRoutine, "wrd "+destname+" "+channame)
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "wrd", I_NIL}
								bg.WriteLine(bg.CurrentRoutine, "chw "+destname2)
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "chw", I_NIL}

								bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, newregcell2}
								resp3 := <-bg.Answers
								if resp3.AnsType == ANS_OK {
									result := make([]VarCell, 1)
									result[0] = newregcell

									return result, true
								} else {
									bg.Set_faulty("Resource removal failed")
									return []VarCell{}, false
								}
							} else {
								bg.Set_faulty("Resource reservation failed")
								return []VarCell{}, false
							}
						} else {
							bg.Set_faulty("Resource reservation failed")
							return []VarCell{}, false
						}

					} else {
						bg.Set_faulty("Wrong type for the arrow operator")
						return []VarCell{}, false
					}
				default:
					bg.Set_faulty("Unsopported unary operation")
					return []VarCell{}, false
				}
			} else {
				bg.Set_faulty("Unary operations requires one returned value")
				return []VarCell{}, false
			}
		}

	case *ast.BinaryExpr:
		x := exptype.X
		y := exptype.Y
		if cell1, ok := bg.Expr_eval(x); !ok {
			bg.Set_faulty("Wrong expression")
			return []VarCell{}, false
		} else {
			if cell2, ok := bg.Expr_eval(y); !ok {
				bg.Set_faulty("Wrong expression")
				return []VarCell{}, false
			} else {
				if len(cell1) == 1 && len(cell2) == 1 {
					switch exptype.Op {
					case token.ADD:
						gent, _ := Type_from_string(bg.Basic_type)
						if Same_Type(cell1[0].Vtype, gent) && Same_Type(cell2[0].Vtype, gent) {
							destname := procbuilder.Get_register_name(cell1[0].Id)
							sourcename := procbuilder.Get_register_name(cell2[0].Id)
							bg.WriteLine(bg.CurrentRoutine, "add "+destname+" "+sourcename)
							bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "add", I_NIL}

							bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell2[0]}
							resp := <-bg.Answers
							if resp.AnsType == ANS_OK {
								result := make([]VarCell, 1)
								result[0] = cell1[0]

								return result, true
							} else {
								bg.Set_faulty("Resource removal failed")
								return []VarCell{}, false
							}

						} else {
							bg.Set_faulty("Variables cannot be added")
							return []VarCell{}, false
						}
					case token.MUL:
						gent, _ := Type_from_string(bg.Basic_type)
						if Same_Type(cell1[0].Vtype, gent) && Same_Type(cell2[0].Vtype, gent) {
							destname := procbuilder.Get_register_name(cell1[0].Id)
							sourcename := procbuilder.Get_register_name(cell2[0].Id)
							bg.WriteLine(bg.CurrentRoutine, "mult "+destname+" "+sourcename)
							bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "mult", I_NIL}

							bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell2[0]}
							resp := <-bg.Answers
							if resp.AnsType == ANS_OK {
								result := make([]VarCell, 1)
								result[0] = cell1[0]

								return result, true
							} else {
								bg.Set_faulty("Resource removal failed")
								return []VarCell{}, false
							}

						} else {
							bg.Set_faulty("Variables cannot be multiplied")
							return []VarCell{}, false
						}
					case token.EQL:
						gent_bool, _ := Type_from_string("bool")
						gent_basic_type, _ := Type_from_string(bg.Basic_type)

						if (Same_Type(cell1[0].Vtype, gent_bool) && Same_Type(cell2[0].Vtype, gent_bool)) || (Same_Type(cell1[0].Vtype, gent_basic_type) && Same_Type(cell2[0].Vtype, gent_basic_type)) {
							bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent_bool, REGISTER, 0, 0, 0, 0, 0, 0}}
							resp := <-bg.Answers
							if resp.AnsType == ANS_OK {
								compcell := resp.Cell

								starting_point := bg.CountLines(bg.CurrentRoutine)

								compname := procbuilder.Get_register_name(compcell.Id)
								destname := procbuilder.Get_register_name(cell1[0].Id)
								sourcename := procbuilder.Get_register_name(cell2[0].Id)
								bg.WriteLine(bg.CurrentRoutine, "je "+destname+" "+sourcename+" <<"+strconv.Itoa(starting_point+3)+">>") // 0
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "je", I_NIL}
								bg.WriteLine(bg.CurrentRoutine, "rset "+compname+" 0") // 1
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "rset", I_NIL}
								bg.WriteLine(bg.CurrentRoutine, "j <<"+strconv.Itoa(starting_point+4)+">>") // 2
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "j", I_NIL}
								bg.WriteLine(bg.CurrentRoutine, "rset "+compname+" 1") // 3
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "rset", I_NIL}

								bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell1[0]}
								resp := <-bg.Answers
								if resp.AnsType == ANS_OK {
									bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell2[0]}
									resp := <-bg.Answers
									if resp.AnsType == ANS_OK {
										result := make([]VarCell, 1)
										result[0] = compcell

										return result, true
									} else {
										bg.Set_faulty("Resource removal failed")
										return []VarCell{}, false
									}
								} else {
									bg.Set_faulty("Resource removal failed")
									return []VarCell{}, false
								}
							} else {
								bg.Set_faulty("Resource allocation failed")
								return []VarCell{}, false
							}
						} else {
							bg.Set_faulty("A Variable is not boolean")
							return []VarCell{}, false
						}
					default:
						bg.Set_faulty("Unsopported binary operation")
						return []VarCell{}, false
					}
				} else {
					bg.Set_faulty("Binary operations requires one returned value")
					return []VarCell{}, false
				}
			}
		}

	case *ast.CallExpr:
		// This is a function call
		var functcell FunctCell

		// TODO Some check here

		switch fun := exptype.Fun.(type) {
		case (*ast.SelectorExpr):
			// This is the case of a "Method"
			x := fun.X.(*ast.Ident)
			sel := fun.Sel

			if x.Name == "bondgo" {
				switch sel.Name {
				case "Make":
					args := exptype.Args
					if len(args) == 2 {
						switch maketype := args[0].(type) {
						case *ast.SelectorExpr:
							xt := maketype.X.(*ast.Ident)
							selt := maketype.Sel
							fullname := xt.Name + "." + selt.Name
							switch fullname {
							case "bondgo.Input":

								switch gidx := args[1].(type) {
								case *ast.BasicLit:
									if gidx.Kind == token.INT {
										value, _ := strconv.Atoi(gidx.Value)
										gent, _ := Type_from_string(bg.Basic_type)
										bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, INPUT, 0, 0, 0, value, value, value}}
										resp := <-bg.Answers
										if resp.AnsType == ANS_OK {

											cell := resp.Cell

											result := make([]VarCell, 1)
											result[0] = cell

											return result, true

										} else {
											bg.Set_faulty("Resource reservation failed")
											return []VarCell{}, false
										}

									}
								}

							case "bondgo.Output":

								switch gidx := args[1].(type) {
								case *ast.BasicLit:
									if gidx.Kind == token.INT {
										value, _ := strconv.Atoi(gidx.Value)
										gent, _ := Type_from_string(bg.Basic_type)
										bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, OUTPUT, 0, 0, 0, value, value, value}}
										resp := <-bg.Answers
										if resp.AnsType == ANS_OK {

											cell := resp.Cell

											result := make([]VarCell, 1)
											result[0] = cell

											return result, true

										} else {
											bg.Set_faulty("Resource reservation failed")
											return []VarCell{}, false
										}

									}
								}

							default:
								bg.Set_faulty("Unknown bondgo object type")
								return []VarCell{}, false
							}
						default:
							bg.Set_faulty("Unknown object type")
							return []VarCell{}, false
						}
					} else {
						bg.Set_faulty("Wrong argument number")
						return []VarCell{}, false
					}
				case "IORead":
					args := exptype.Args
					if len(args) == 1 {
						// TODO Finish
						switch arg := args[0].(type) {
						case *ast.Ident:
							identname := arg.Name

							varexist := false
							for scope := bg; scope != nil; scope = scope.Outer {
								if cell, ok := scope.Vars[identname]; ok {
									varexist = true
									var newregcell VarCell

									switch cell.Procobjtype {
									case INPUT:
										if cell.Global_id == 0 {
											bg.Set_faulty("Input not allocated")
											return []VarCell{}, false
										} else {
											gent, _ := Type_from_string(bg.Basic_type)
											bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
											resp := <-bg.Answers
											if resp.AnsType == ANS_OK {
												newregcell = resp.Cell
												regname := procbuilder.Get_register_name(newregcell.Id)
												newinname := procbuilder.Get_input_name(cell.Id)
												bg.WriteLine(bg.CurrentRoutine, "i2r "+regname+" "+newinname)
												bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "i2r", I_NIL}
												// TODO Insert here the notification for the input usage
											} else {
												bg.Set_faulty("Resource reservation failed")
												return []VarCell{}, false
											}
										}
									default:
										bg.Set_faulty("Read can only be used on input registers")
										return []VarCell{}, false
									}

									result := make([]VarCell, 1)
									result[0] = newregcell

									return result, true
								}
							}
							if !varexist {
								bg.Set_faulty("Variable " + identname + " not defined")
								return []VarCell{}, false
							}
						}

					} else {
						bg.Set_faulty("Only one input expected")
						return []VarCell{}, false
						// TODO maybe in the future
					}
				default:
					bg.Set_faulty("Unknown function " + sel.Name)
					return []VarCell{}, false
				}
			} else {
				bg.Set_faulty("Unknown module " + x.Name)
				return []VarCell{}, false
			}

		case (*ast.Ident):
			// This id the case of a function with no receiver
			funname := fun.Name

			if _, ok := bg.Functions[funname]; ok {
				functcell = bg.Functions[funname]
			} else {
				bg.Set_faulty("Undefined function " + funname)
				return []VarCell{}, false
			}

			// Allocate variables for the function arguments and fill them with evalued values of the arguments
			vars := make(map[string]VarCell)

			if len(functcell.Inputs) == len(exptype.Args) {
				for i, arg := range functcell.Inputs {
					argname := arg.Argname
					if _, ok := vars[argname]; ok {
						bg.Set_faulty("Already defined variable")
						return []VarCell{}, false
					} else {
						if cell, ok := bg.Expr_eval(exptype.Args[i]); !ok {
							bg.Set_faulty("Wrong evaluation")
							return []VarCell{}, false
						} else {
							vars[argname] = cell[0]
						}
					}
				}
			} else {
				bg.Set_faulty("Call with wrong number of parameters")
				return []VarCell{}, false
			}

			// Allocate variables for the returns and fill them with empty values
			returns := make([]VarCell, 0)

			// Create a new BondgoCheck with empty variable list e no parent but the same allocator
			results := new(BondgoResults) // Results go in here
			results.Init_Results(bg.BondgoConfig)

			bgfunct := &BondgoCheck{results, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, nil, nil, vars, returns, "", "", bg.CurrentDevice, bg.CurrentRoutine}

			// Launch Walk on the function body using the new BondgoCheck
			ast.Walk(bgfunct, functcell.Body)
			//fmt.Print("---\n", bgfunct.Write_assembly(), "\n----\n")

			// Get the generated code, count the lines and substitute the <<LASTN>> placeholders whit the actual line number
			starting_point := bg.CountLines(bg.CurrentRoutine)
			prod_lines := bgfunct.CountLines(bgfunct.CurrentRoutine)

			bgfunct.Replacer(bgfunct.CurrentRoutine, "<<LASTN>>", "<<"+strconv.Itoa(prod_lines)+">>")

			// Shift eventually created reference to line number within the code
			bgfunct.Shift_program_location(bgfunct.CurrentRoutine, starting_point)

			// Return the results
			for _, line := range bgfunct.GetProgram(bgfunct.CurrentRoutine) {
				bg.WriteLine(bg.CurrentRoutine, line)
			}

			result := make([]VarCell, len(bgfunct.Returns))
			for i, cell := range bgfunct.Returns {
				result[i] = cell
			}

			// Clean the inputs arguments VarCells
			for vari, cell := range bgfunct.Vars {
				if bg.In_debug() {
					fmt.Println("Cleaning " + vari)
				}

				bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell}
				resp := <-bg.Answers
				if resp.AnsType != ANS_OK {
					bg.Set_faulty("Resource removal failed")
					return []VarCell{}, false
				}
			}

			return result, true
		}
	}
	return []VarCell{}, false

}
