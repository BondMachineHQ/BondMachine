package bondgo

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/procbuilder"
)

func (bg *BondgoCheck) Visit(n ast.Node) ast.Visitor {
	//bgclone := &BondgoCheck{bg.BondgoResults, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg.Outer, bg.Clean, bg.Vars, bg.Returns, bg.CurrentLoop, bg.CurrentRoutine}

	if bg.Clean != nil {
		for vari, cell := range bg.Clean.Vars {
			if bg.In_debug() {
				fmt.Println("Cleaning " + vari)
			}
			bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell}
			if (<-bg.Answers).AnsType != ANS_OK {
				bg.Set_faulty("Resource clean failed")
				return nil
			}
		}
		bg.Clean = nil
	}

	switch x := n.(type) {
	case *ast.GenDecl:
		if bg.In_debug() {
			fmt.Printf("%p - Declaration\n", bg)
		}

		genDecl := n.(*ast.GenDecl)

		switch genDecl.Tok {
		case token.VAR:
			for _, s := range genDecl.Specs {
				spec := s.(*ast.ValueSpec)

				switch vtype := spec.Type.(type) {
				case *ast.SelectorExpr:
					x := vtype.X.(*ast.Ident)
					if x.Name == "bondgo" {
						if vtype.Sel.Name == "Input" {
							gent, _ := Type_from_string(bg.Basic_type)
							for _, vari := range spec.Names {
								if _, ok := bg.Vars[vari.Name]; ok {
									bg.Set_faulty(vari.Name + ": name already used")
									return nil
								} else {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, INPUT, 0, 0, 0, 0, 0, 0}}
									resp := <-bg.Answers
									if resp.AnsType == ANS_OK {

										cell := resp.Cell
										bg.Vars[vari.Name] = cell

										if bg.In_debug() {
											fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
										}
									} else {
										bg.Set_faulty("Resource reservation failed")
										return nil
									}
								}
							}
						} else if vtype.Sel.Name == "Output" {
							gent, _ := Type_from_string(bg.Basic_type)
							for _, vari := range spec.Names {
								if _, ok := bg.Vars[vari.Name]; ok {
									bg.Set_faulty(vari.Name + ": name already used")
									return nil
								} else {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, OUTPUT, 0, 0, 0, 0, 0, 0}}
									resp := <-bg.Answers
									if resp.AnsType == ANS_OK {

										cell := resp.Cell
										bg.Vars[vari.Name] = cell

										if bg.In_debug() {
											fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
										}
									} else {
										bg.Set_faulty("Resource reservation failed")
										return nil
									}
								}
							}
						} else {
							bg.Set_faulty("Unknown selector " + vtype.Sel.Name)
							return nil
						}
					} else {
						bg.Set_faulty("Unknown package " + x.Name)
						return nil
					}
				default:

					newt, _ := Type_from_ast(spec.Type)

					if gent, _ := Type_from_string(bg.Basic_type); Same_Type(newt, gent) {

						for _, vari := range spec.Names {
							if _, ok := bg.Vars[vari.Name]; ok {
								bg.Set_faulty(vari.Name + ": name already used")
								return nil
							} else {
								if len(vari.Name) > 4 && vari.Name[:4] == "reg_" {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
								} else {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, MEMORY, 0, 0, 0, 0, 0, 0}}
								}
								resp := <-bg.Answers
								if resp.AnsType == ANS_OK {

									cell := resp.Cell
									bg.Vars[vari.Name] = cell

									switch cell.Procobjtype {
									case REGISTER:
										regname := procbuilder.Get_register_name(cell.Id)
										bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
									case MEMORY:
										bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
										newresp := <-bg.Answers
										if newresp.AnsType == ANS_OK {
											newregcell := newresp.Cell

											regname := procbuilder.Get_register_name(newregcell.Id)

											bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
											bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
											bg.WriteLine(bg.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(cell.Id))
											bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "r2m", I_NIL}

											bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, newregcell}
											if (<-bg.Answers).AnsType != ANS_OK {
												bg.Set_faulty("Resource clean failed")
												return nil
											}
										} else {
											bg.Set_faulty("Resource reservation failed")
											return nil
										}
									}
									if bg.In_debug() {
										fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
									}
								} else {
									bg.Set_faulty("Resource reservation failed")
									return nil
								}
							}
						}
					} else if gent, _ := Type_from_string("bool"); Same_Type(newt, gent) {

						for _, vari := range spec.Names {
							if _, ok := bg.Vars[vari.Name]; ok {
								bg.Set_faulty(vari.Name + ": name already used")
								return nil
							} else {
								if len(vari.Name) > 4 && vari.Name[:4] == "reg_" {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
								} else {
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, MEMORY, 0, 0, 0, 0, 0, 0}}
								}
								resp := <-bg.Answers
								if resp.AnsType == ANS_OK {

									cell := resp.Cell
									bg.Vars[vari.Name] = cell

									switch cell.Procobjtype {
									case REGISTER:
										regname := procbuilder.Get_register_name(cell.Id)
										bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
									case MEMORY:
										bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
										newresp := <-bg.Answers
										if newresp.AnsType == ANS_OK {
											newregcell := newresp.Cell

											regname := procbuilder.Get_register_name(newregcell.Id)

											bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
											bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
											bg.WriteLine(bg.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(cell.Id))
											bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "r2m", I_NIL}

											bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, newregcell}
											if (<-bg.Answers).AnsType != ANS_OK {
												bg.Set_faulty("Resource clean failed")
												return nil
											}
										} else {
											bg.Set_faulty("Resource reservation failed")
											return nil
										}
									}
									if bg.In_debug() {
										fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
									}
								} else {
									bg.Set_faulty("Resource reservation failed")
									return nil
								}
							}
						}
					} else if gent, _ := Type_from_string(bg.Basic_chantype); Same_Type(newt, gent) {
						for _, vari := range spec.Names {
							if _, ok := bg.Vars[vari.Name]; ok {
								bg.Set_faulty(vari.Name + ": name already used")
								return nil
							} else {
								bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, CHANNEL, 0, 0, 0, 0, 0, 0}}
								resp := <-bg.Answers
								if resp.AnsType == ANS_OK {

									cell := resp.Cell
									bg.Vars[vari.Name] = cell

									if bg.In_debug() {
										fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
									}
								} else {
									bg.Set_faulty("Resource reservation failed")
									return nil
								}
							}
						}
					} else if gent, _ := Type_from_string("chan bool"); Same_Type(newt, gent) {
						for _, vari := range spec.Names {
							if _, ok := bg.Vars[vari.Name]; ok {
								bg.Set_faulty(vari.Name + ": name already used")
								return nil
							} else {
								bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, CHANNEL, 0, 0, 0, 0, 0, 0}}
								resp := <-bg.Answers
								if resp.AnsType == ANS_OK {

									cell := resp.Cell
									bg.Vars[vari.Name] = cell

									if bg.In_debug() {
										fmt.Println("\t\tAllocated to " + vari.Name + " the cell " + bg.Vars[vari.Name].String())
									}
								} else {
									bg.Set_faulty("Resource reservation failed")
									return nil
								}
							}
						}
					} else {

					}
				}
			}
		}

	case *ast.AssignStmt:

		// TODO Improvemente needed: type checks and channel *think*
		// TODO assignment and declaration

		if bg.In_debug() {
			fmt.Printf("%p - %s\n", bg, "Assignment Statement")
		}

		assignStmt := n.(*ast.AssignStmt)
		switch assignStmt.Tok {
		case token.ASSIGN:
			if len(assignStmt.Lhs) == len(assignStmt.Rhs) {
				destinations := make([]VarCell, len(assignStmt.Lhs))
				sources := make([]VarCell, len(assignStmt.Lhs))

				for assindex, _ := range assignStmt.Lhs {
					lhs := assignStmt.Lhs[assindex].(*ast.Ident)
					vari := lhs.Name

					varexist := false
					for scope := bg; scope != nil; scope = scope.Outer {
						if _, ok := scope.Vars[vari]; ok {
							varexist = true
							destinations[assindex] = scope.Vars[vari]
							break
						}
					}

					if !varexist {
						bg.Set_faulty("Unitialized variable " + vari)
						return nil
					}

					rhs := assignStmt.Rhs[assindex]

					if newcell, ok := bg.Expr_eval(rhs); ok {
						sources[assindex] = newcell[0]
					} else {
						bg.Set_faulty("Assignment RHS evaluation failed")
						return nil
					}
				}

				for assindex, cell := range destinations {
					newcell := sources[assindex]
					switch cell.Procobjtype {
					case REGISTER:
						regname := procbuilder.Get_register_name(cell.Id)
						newregname := procbuilder.Get_register_name(newcell.Id)
						bg.WriteLine(bg.CurrentRoutine, "cpy "+regname+" "+newregname)
						bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "cpy", I_NIL}
					case MEMORY:
						newregname := procbuilder.Get_register_name(newcell.Id)
						bg.WriteLine(bg.CurrentRoutine, "r2m "+newregname+" "+strconv.Itoa(cell.Id))
						bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "r2m", I_NIL}
					case INPUT:
						// Assing an input to another means to use the same global index
						cell.Global_id = newcell.Global_id
						cell.Start_globalid = newcell.Start_globalid
						cell.End_globalid = newcell.End_globalid

						lhs := assignStmt.Lhs[assindex].(*ast.Ident)
						vari := lhs.Name

						for scope := bg; scope != nil; scope = scope.Outer {
							if _, ok := scope.Vars[vari]; ok {
								scope.Vars[vari] = cell
								break
							}
						}

					case OUTPUT:
						// Assing an output to another means to use the same global index
						cell.Global_id = newcell.Global_id
						cell.Start_globalid = newcell.Start_globalid
						cell.End_globalid = newcell.End_globalid

						lhs := assignStmt.Lhs[assindex].(*ast.Ident)
						vari := lhs.Name

						for scope := bg; scope != nil; scope = scope.Outer {
							if _, ok := scope.Vars[vari]; ok {
								scope.Vars[vari] = cell
								break
							}
						}
					case CHANNEL:
					}

					bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, newcell}
					if (<-bg.Answers).AnsType != ANS_OK {
						bg.Set_faulty("Resource clean failed")
						return nil
					}
				}
			} else {
				bg.Set_faulty("Different arguments lenghts")
				return nil
			}
		case token.DEFINE:
			// TODO Finish
			if len(assignStmt.Lhs) == len(assignStmt.Rhs) {
				sources := make([]VarCell, len(assignStmt.Lhs))

				for assindex, _ := range assignStmt.Lhs {
					lhs := assignStmt.Lhs[assindex].(*ast.Ident)
					vari := lhs.Name

					varexist := false
					if _, ok := bg.Vars[vari]; ok {
						varexist = true
						break
					}

					if varexist {
						bg.Set_faulty("Already defined variable " + vari)
						return nil
					}

					rhs := assignStmt.Rhs[assindex]

					if newcell, ok := bg.Expr_eval(rhs); ok {
						sources[assindex] = newcell[0]
					} else {
						bg.Set_faulty("Assignment RHS evaluation failed")
						return nil
					}
				}

				for assindex, cell := range sources {
					lhs := assignStmt.Lhs[assindex].(*ast.Ident)
					vari := lhs.Name

					switch cell.Procobjtype {
					case REGISTER:
						gent := cell.Vtype
						if len(vari) > 4 && vari[:4] == "reg_" {
							bg.Vars[vari] = cell
						} else {
							bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, MEMORY, 0, 0, 0, 0, 0, 0}}
							resp := <-bg.Answers
							if resp.AnsType == ANS_OK {
								regname := procbuilder.Get_register_name(cell.Id)
								bg.WriteLine(bg.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(resp.Cell.Id))
								bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "r2m", I_NIL}
								bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell}
								if (<-bg.Answers).AnsType != ANS_OK {
									bg.Set_faulty("Resource clean failed")
									return nil
								}
							} else {
								bg.Set_faulty("Output request failed")
								return nil
							}
						}

					case CHANNEL:
						bg.Set_faulty("Channel reassign prohibited")
						return nil
					}

				}
			} else {
				bg.Set_faulty("Different arguments lenghts")
				return nil
			}
		default:
			bg.Set_faulty("Unknown assignment operation")
			return nil

		}

		// Do not reprocess assignment leaves
		return nil

	case *ast.IncDecStmt:

		incDecStmt := n.(*ast.IncDecStmt)

		vari := incDecStmt.X.(*ast.Ident).Name

		varexist := false
		for scope := bg; scope != nil; scope = scope.Outer {
			if _, ok := scope.Vars[vari]; ok {
				varexist = true
				cell := scope.Vars[vari]
				switch cell.Procobjtype {
				case REGISTER:
					regname := procbuilder.Get_register_name(cell.Id)
					switch incDecStmt.Tok {
					case token.INC:
						scope.WriteLine(scope.CurrentRoutine, "inc "+regname)
						scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "inc", I_NIL}
					case token.DEC:
						scope.WriteLine(scope.CurrentRoutine, "dec "+regname)
						scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "dec", I_NIL}
					}
				case MEMORY:
					gent, _ := Type_from_string(bg.Basic_type)
					scope.Reqs <- VarReq{REQ_NEW, scope.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {

						newregcell := resp.Cell

						regname := procbuilder.Get_register_name(newregcell.Id)

						scope.WriteLine(scope.CurrentRoutine, "m2r "+regname+" "+strconv.Itoa(cell.Id))
						scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "m2r", I_NIL}

						switch incDecStmt.Tok {
						case token.INC:
							scope.WriteLine(scope.CurrentRoutine, "inc "+regname)
							scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "inc", I_NIL}
						case token.DEC:
							scope.WriteLine(scope.CurrentRoutine, "dec "+regname)
							scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "dec", I_NIL}
						}

						scope.WriteLine(scope.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(cell.Id))
						scope.Used <- UsageNotify{TR_PROC, scope.CurrentRoutine, C_OPCODE, "r2m", I_NIL}
						scope.Reqs <- VarReq{REQ_REMOVE, scope.CurrentRoutine, newregcell}
						if (<-scope.Answers).AnsType != ANS_OK {
							bg.Set_faulty("Resource clean failed")
							return nil
						}
					} else {
						bg.Set_faulty("Resource allocation failed")
						return nil
					}
				}
				break
			}
		}
		if !varexist {
			bg.Set_faulty("Unitialized variable " + vari)
			return nil
		}

	case *ast.BlockStmt:
		{
			if bg.In_debug() {
				fmt.Printf("%p - %s ", bg, "Entering the new block")
			}

			vars := make(map[string]VarCell)

			bgleaf := &BondgoCheck{bg.BondgoResults, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, vars, bg.Returns, bg.CurrentLoop, bg.CurrentSwitch, bg.CurrentDevice, bg.CurrentRoutine}
			bg.Clean = bgleaf

			if bg.In_debug() {
				fmt.Printf("%p\n", bgleaf)
			}

			return bgleaf
		}
	case *ast.IfStmt:
		if bg.In_debug() {
			fmt.Println("Conditional statement")
		}

		// Create a new BondgoCheck for the if statement with an empty program
		results := new(BondgoResults) // Results go in here
		results.Init_Results(bg.BondgoConfig)

		bgif := &BondgoCheck{results, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, bg.Vars, bg.Returns, bg.CurrentLoop, bg.CurrentSwitch, bg.CurrentDevice, bg.CurrentRoutine}

		starting_point := bg.CountLines(bg.CurrentRoutine)

		if x.Init != nil {
			// Launch Walk on the init body using the new BondgoCheck
			ast.Walk(bgif, x.Init)
		}

		if x.Cond != nil {
			// Now eval the condition
			if newcell, ok := bgif.Expr_eval(x.Cond); ok {
				condevalued := newcell[0]

				gent, _ := Type_from_string("bool")
				if Same_Type(condevalued.Vtype, gent) {

					newregname := procbuilder.Get_register_name(condevalued.Id)
					bgif.WriteLine(bgif.CurrentRoutine, "jz "+newregname+" <<"+fmt.Sprintf("%p", bgif)+"FALSECONDITION>>")
					bgif.Used <- UsageNotify{TR_PROC, bgif.CurrentRoutine, C_OPCODE, "jz", I_NIL}

					bgif.Reqs <- VarReq{REQ_REMOVE, bgif.CurrentRoutine, condevalued}
					if (<-bgif.Answers).AnsType != ANS_OK {
						bg.Set_faulty("Resource clean failed")
						return nil
					}
				} else {
					bg.Set_faulty("If need a boolean condition")
					return nil
				}
			} else {
				bg.Set_faulty("If condition evaluation failed")
				return nil
			}
		}

		if x.Body != nil {
			ast.Walk(bgif, x.Body)
		}

		if x.Else != nil {
			bgif.WriteLine(bgif.CurrentRoutine, "j <<"+fmt.Sprintf("%p", bgif)+"IFEND>>")
			bgif.Used <- UsageNotify{TR_PROC, bgif.CurrentRoutine, C_OPCODE, "j", I_NIL}
		}

		prod_lines_body := bgif.CountLines(bgif.CurrentRoutine)

		if x.Else != nil {
			ast.Walk(bgif, x.Else)
		}

		prod_lines_total := bgif.CountLines(bgif.CurrentRoutine)

		bgif.Replacer(bgif.CurrentRoutine, "<<"+fmt.Sprintf("%p", bgif)+"FALSECONDITION>>", "<<"+strconv.Itoa(prod_lines_body)+">>")

		bgif.Replacer(bgif.CurrentRoutine, "<<"+fmt.Sprintf("%p", bgif)+"IFEND>>", "<<"+strconv.Itoa(prod_lines_total)+">>")

		// Shift eventually created reference to line number within the code
		bgif.Shift_program_location(bgif.CurrentRoutine, starting_point)

		for _, line := range bgif.GetProgram(bgif.CurrentRoutine) {
			bg.WriteLine(bg.CurrentRoutine, line)
		}

		// The node has already visited.
		return nil

	case *ast.ForStmt:
		if bg.In_debug() {
			fmt.Printf("%p - Entering The for loop ", bg)
		}

		// Create a new BondgoCheck for the loop
		results := new(BondgoResults) // Results go in here
		results.Init_Results(bg.BondgoConfig)

		bgfor := &BondgoCheck{results, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, bg.Vars, bg.Returns, "", bg.CurrentSwitch, bg.CurrentDevice, bg.CurrentRoutine}

		if bg.In_debug() {
			fmt.Printf("%p\n", bgfor)
		}

		starting_point := bg.CountLines(bg.CurrentRoutine)

		bgfor.CurrentLoop = fmt.Sprintf("%p", bgfor)

		if x.Init != nil {
			// Launch Walk on the init body using the new BondgoCheck
			ast.Walk(bgfor, x.Init)
		}

		// The for starting point
		loop_start := bgfor.CountLines(bgfor.CurrentRoutine)

		if x.Cond != nil {
			// Now eval the condition
			if newcell, ok := bgfor.Expr_eval(x.Cond); ok {
				condevalued := newcell[0]

				gent, _ := Type_from_string("bool")
				if Same_Type(condevalued.Vtype, gent) {

					newregname := procbuilder.Get_register_name(condevalued.Id)
					bgfor.WriteLine(bgfor.CurrentRoutine, "jz "+newregname+" <<"+bgfor.CurrentLoop+"ENDFOR>>")
					bgfor.Used <- UsageNotify{TR_PROC, bgfor.CurrentRoutine, C_OPCODE, "jz", I_NIL}

					bgfor.Reqs <- VarReq{REQ_REMOVE, bgfor.CurrentRoutine, condevalued}
					if (<-bgfor.Answers).AnsType != ANS_OK {
						bg.Set_faulty("Resource clean failed")
						return nil
					}
				} else {
					bg.Set_faulty("For need a boolean condition")
					return nil
				}
			} else {
				bg.Set_faulty("For condition evaluation failed")
				return nil
			}
		}

		if x.Body != nil {
			ast.Walk(bgfor, x.Body)
		}

		continue_point := bgfor.CountLines(bgfor.CurrentRoutine)

		if x.Post != nil {
			ast.Walk(bgfor, x.Post)
		}

		bgfor.WriteLine(bgfor.CurrentRoutine, "j <<"+bgfor.CurrentLoop+"STARTFOR>>")
		bgfor.Used <- UsageNotify{TR_PROC, bgfor.CurrentRoutine, C_OPCODE, "j", I_NIL}

		prod_lines_total := bgfor.CountLines(bgfor.CurrentRoutine)

		bgfor.Replacer(bgfor.CurrentRoutine, "<<"+bgfor.CurrentLoop+"STARTFOR>>", "<<"+strconv.Itoa(loop_start)+">>")
		bgfor.Replacer(bgfor.CurrentRoutine, "<<"+bgfor.CurrentLoop+"ENDFOR>>", "<<"+strconv.Itoa(prod_lines_total)+">>")
		bgfor.Replacer(bgfor.CurrentRoutine, "<<"+bgfor.CurrentLoop+"CONTINUEFOR>>", "<<"+strconv.Itoa(continue_point)+">>")

		// Shift eventually created reference to line number within the code
		bgfor.Shift_program_location(bgfor.CurrentRoutine, starting_point)

		for _, line := range bgfor.GetProgram(bgfor.CurrentRoutine) {
			bg.WriteLine(bg.CurrentRoutine, line)
		}

		// The node has already been visited.
		return nil

	case *ast.SelectStmt:
		if bg.In_debug() {
			fmt.Printf("%p - Select statement", bg)
		}

		// Create a new BondgoCheck for the select statement with an empty program
		results := new(BondgoResults) // Results go in here
		results.Init_Results(bg.BondgoConfig)

		bgsel := &BondgoCheck{results, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, bg.Vars, bg.Returns, bg.CurrentLoop, "", bg.CurrentDevice, bg.CurrentRoutine}

		starting_point := bg.CountLines(bg.CurrentRoutine)

		bgsel.CurrentSwitch = fmt.Sprintf("SEL%p", bgsel)

		selectbody := x.Body

		if selectbody.List != nil {
			// There is a default case or not
			defaulted := false
			var default_point int

			cases := len(selectbody.List)
			starting_points := make([]int, cases)

			comm_registers := make([]VarCell, cases)
			comm_direction := make([]bool, cases) // true = send
			comm_channel := make([]VarCell, cases)

			// The first loop set all the jumps toward cases
			ii := 0
			for _, clause := range selectbody.List {
				switch cc := clause.(type) {
				case *ast.CommClause:
					if cc.Comm != nil {
						switch commst := cc.Comm.(type) {
						case *ast.SendStmt:
							// This is a send operation
							if newcell, ok := bgsel.Expr_eval(commst.Value); ok {
								// TODO Optimization, this eval here may be a problem, everytime the expression is evaluted even if the channel is not used.
								comm_registers[ii] = newcell[0]
								comm_direction[ii] = true
							} else {
								bgsel.Set_faulty("Select condition evaluation failed")
								return nil
							}

							chn := commst.Chan.(*ast.Ident)

							for scope := bgsel; scope != nil; scope = scope.Outer {
								if _, ok := scope.Vars[chn.Name]; ok {
									comm_channel[ii] = scope.Vars[chn.Name]
									break
								}
							}
						case *ast.AssignStmt:
							// This is a receive operation, no matter if it's a assign or define space will be reserved.
							if len(commst.Lhs) == len(commst.Rhs) && len(commst.Lhs) == 1 {
								switch commst.Rhs[0].(type) {
								case *ast.UnaryExpr:
									rhs := commst.Rhs[0].(*ast.UnaryExpr)
									if rhs.Op == token.ARROW {
										switch rhs.X.(type) {
										case *ast.Ident:
											chn := rhs.X.(*ast.Ident)
											for scope := bgsel; scope != nil; scope = scope.Outer {
												if _, ok := scope.Vars[chn.Name]; ok {
													comm_channel[ii] = scope.Vars[chn.Name]
													if gent, _ := Type_from_string(bg.Basic_type); Same_Type(comm_channel[ii].Vtype.Values[0], gent) {
														bgsel.Reqs <- VarReq{REQ_NEW, bgsel.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
														resp := <-bgsel.Answers
														if resp.AnsType == ANS_OK {
															comm_registers[ii] = resp.Cell
														} else {
															bgsel.Set_faulty("Resource reservation failed")
															return nil
														}
													} else {
														bgsel.Set_faulty("Only allowed basic type")
														return nil
													}
													break
												}
											}
										default:
											bgsel.Set_faulty("Operation not valid in select assigment")
											return nil
										}
									} else {
										bgsel.Set_faulty("Operation not valid in select assigment")
										return nil
									}
								default:
									bgsel.Set_faulty("Operation not valid in select assigment")
									return nil
								}
							} else {
								bgsel.Set_faulty("Multivalue not allowed in receive operations")
								return nil
							}
						}
						ii++
					} else {
						// This is the default case
						if defaulted {
							bg.Set_faulty("Duplicate default case")
							return nil
						}
						defaulted = true
					}

				default:
					bgsel.Set_faulty("Wrong case")
					return nil
				}
			}

			if defaulted {
				comm_direction = comm_direction[:cases-1]
				comm_channel = comm_channel[:cases-1]
				comm_registers = comm_registers[:cases-1]
				starting_points = starting_points[:cases-1]
			}

			for i, dir := range comm_direction {
				if dir {
					bgsel.WriteLine(bgsel.CurrentRoutine, "wwr "+procbuilder.Get_register_name(comm_registers[i].Id)+" "+procbuilder.Get_channel_name(comm_channel[i].Id))
					bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "wwr", I_NIL}
				} else {
					bgsel.WriteLine(bgsel.CurrentRoutine, "wrd "+procbuilder.Get_register_name(comm_registers[i].Id)+" "+procbuilder.Get_channel_name(comm_channel[i].Id))
					bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "wrd", I_NIL}
				}
			}

			var eventreg VarCell

			gent, _ := Type_from_string(bg.Basic_type)
			bgsel.Reqs <- VarReq{REQ_NEW, bgsel.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
			resp := <-bgsel.Answers
			if resp.AnsType == ANS_OK {
				eventreg = resp.Cell
			} else {
				bgsel.Set_faulty("Resource reservation failed")
				return nil
			}

			if defaulted { // Jump to the end of the switch

				var occurreg VarCell

				bgsel.Reqs <- VarReq{REQ_NEW, bgsel.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
				resp := <-bgsel.Answers
				if resp.AnsType == ANS_OK {
					occurreg = resp.Cell
				} else {
					bgsel.Set_faulty("Resource reservation failed")
					return nil
				}

				bgsel.WriteLine(bgsel.CurrentRoutine, "chc "+procbuilder.Get_register_name(occurreg.Id)+" "+procbuilder.Get_register_name(eventreg.Id))
				bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "chc", I_NIL}

				bgsel.WriteLine(bgsel.CurrentRoutine, "jz "+procbuilder.Get_register_name(occurreg.Id)+" <<"+bgsel.CurrentSwitch+"CASEDEFAULT>>")
				bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "jz", I_NIL}

				bgsel.Reqs <- VarReq{REQ_REMOVE, bgsel.CurrentRoutine, occurreg}
				if (<-bgsel.Answers).AnsType != ANS_OK {
					bgsel.Set_faulty("Resource clean failed")
					return nil
				}
			} else {
				bgsel.WriteLine(bgsel.CurrentRoutine, "chw "+procbuilder.Get_register_name(eventreg.Id))
				bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "chw", I_NIL}
			}

			for i, _ := range comm_registers {
				bgsel.WriteLine(bgsel.CurrentRoutine, "jz "+procbuilder.Get_register_name(eventreg.Id)+" <<"+bgsel.CurrentSwitch+"CASE"+strconv.Itoa(i)+">>")
				bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "jz", I_NIL}
				if i < len(comm_registers)-1 {
					bgsel.WriteLine(bgsel.CurrentRoutine, "dec "+procbuilder.Get_register_name(eventreg.Id))
					bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "dec", I_NIL}
				}
			}

			bgsel.Reqs <- VarReq{REQ_REMOVE, bgsel.CurrentRoutine, eventreg}
			if (<-bgsel.Answers).AnsType != ANS_OK {
				bgsel.Set_faulty("Resource clean failed")
				return nil
			}

			// The second loop create the cases
			ii = 0
			for _, clause := range selectbody.List {
				switch cc := clause.(type) {
				case *ast.CommClause:

					if cc.Comm != nil {
						// This is a case
						starting_points[ii] = bgsel.CountLines(bgsel.CurrentRoutine)

						// It's a receive
						if !comm_direction[ii] {
							vari := cc.Comm.(*ast.AssignStmt).Lhs[0].(*ast.Ident).Name
							switch cc.Comm.(*ast.AssignStmt).Tok {
							case token.ASSIGN:

								varexist := false
								for scope := bgsel; scope != nil; scope = scope.Outer {
									if _, ok := scope.Vars[vari]; ok {
										varexist = true

										cell := scope.Vars[vari]
										newcell := comm_registers[ii]

										switch cell.Procobjtype {
										case REGISTER:
											regname := procbuilder.Get_register_name(cell.Id)
											newregname := procbuilder.Get_register_name(newcell.Id)
											bgsel.WriteLine(bgsel.CurrentRoutine, "cpy "+regname+" "+newregname)
											bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "cpy", I_NIL}
										case MEMORY:
											newregname := procbuilder.Get_register_name(newcell.Id)
											bgsel.WriteLine(bgsel.CurrentRoutine, "r2m "+newregname+" "+strconv.Itoa(cell.Id))
											bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "r2m", I_NIL}
										case INPUT:
											bgsel.Set_faulty("An input cannot be written")
											return nil
										case OUTPUT:
											newregname := procbuilder.Get_register_name(newcell.Id)
											bgsel.WriteLine(bgsel.CurrentRoutine, "r2o "+newregname+" "+strconv.Itoa(cell.Id))
											bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "r2o", I_NIL}
										case CHANNEL:
											bgsel.Set_faulty("A channel cannot be written")
											return nil
										}
										break
									}
								}

								if !varexist {
									bgsel.Set_faulty("Unitialized variable " + vari)
									return nil
								}

							case token.DEFINE:

								gent, _ := Type_from_string(bgsel.Basic_type)
								bgsel.Reqs <- VarReq{REQ_NEW, bgsel.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
								resp := <-bgsel.Answers
								if resp.AnsType == ANS_OK {

									newcell := resp.Cell
									cell := comm_registers[ii]

									regname := procbuilder.Get_register_name(cell.Id)
									newregname := procbuilder.Get_register_name(newcell.Id)
									bgsel.WriteLine(bgsel.CurrentRoutine, "cpy "+newregname+" "+regname)
									bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "cpy", I_NIL}

								} else {
									bgsel.Set_faulty("Resource reservation failed")
									return nil
								}
							default:
								bgsel.Set_faulty("Wrong assignment")
								return nil
							}
						}

					} else {
						// This is the default case
						default_point = bgsel.CountLines(bgsel.CurrentRoutine)
					}

					body := clause.(*ast.CommClause).Body

					for _, cl := range body {
						ast.Walk(bgsel, cl)
					}

				default:
					bgsel.Set_faulty("Wrong case")
					return nil
				}

				bgsel.WriteLine(bg.CurrentRoutine, "j <<"+bgsel.CurrentSwitch+"SELEND>>")
				bgsel.Used <- UsageNotify{TR_PROC, bgsel.CurrentRoutine, C_OPCODE, "j", I_NIL}

				ii++
			}

			for _, jreg := range comm_registers {
				bgsel.Reqs <- VarReq{REQ_REMOVE, bgsel.CurrentRoutine, jreg}
				if (<-bgsel.Answers).AnsType != ANS_OK {
					bgsel.Set_faulty("Resource clean failed")
					return nil
				}
			}

			select_end := bgsel.CountLines(bgsel.CurrentRoutine)

			for i, _ := range starting_points {
				bgsel.Replacer(bgsel.CurrentRoutine, "<<"+bgsel.CurrentSwitch+"CASE"+strconv.Itoa(i)+">>", "<<"+strconv.Itoa(starting_points[i])+">>")
			}
			bgsel.Replacer(bgsel.CurrentRoutine, "<<"+bgsel.CurrentSwitch+"CASEDEFAULT>>", "<<"+strconv.Itoa(default_point)+">>")
			bgsel.Replacer(bgsel.CurrentRoutine, "<<"+bgsel.CurrentSwitch+"SELEND>>", "<<"+strconv.Itoa(select_end)+">>")

			// Shift eventually created reference to line number within the code
			bgsel.Shift_program_location(bgsel.CurrentRoutine, starting_point)

			for _, line := range bgsel.GetProgram(bgsel.CurrentRoutine) {
				bg.WriteLine(bg.CurrentRoutine, line)
			}

		}
		return nil

	case *ast.SwitchStmt:
		if bg.In_debug() {
			fmt.Printf("%p - Switch statement", bg)
		}

		var tagexpr VarCell

		// Create a new BondgoCheck for the switch statement with an empty program
		results := new(BondgoResults) // Results go in here
		results.Init_Results(bg.BondgoConfig)

		bgsw := &BondgoCheck{results, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, bg.Vars, bg.Returns, bg.CurrentLoop, "", bg.CurrentDevice, bg.CurrentRoutine}

		starting_point := bg.CountLines(bg.CurrentRoutine)

		bgsw.CurrentSwitch = fmt.Sprintf("%p", bgsw)

		// I there is tag ok, otherwise use  a bool settet to true (not 0)
		if x.Tag != nil {
			if newcell, ok := bgsw.Expr_eval(x.Tag); ok {
				tagexpr = newcell[0]
			} else {
				bgsw.Set_faulty("Switch condition evaluation failed")
				return nil
			}
		} else {
			gent, _ := Type_from_string("bool")
			bgsw.Reqs <- VarReq{REQ_NEW, bgsw.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
			resp := <-bg.Answers
			if resp.AnsType == ANS_OK {
				tagexpr = resp.Cell

				regname := procbuilder.Get_register_name(tagexpr.Id)

				bgsw.WriteLine(bg.CurrentRoutine, "clr "+regname)
				bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "clr", I_NIL}
				bgsw.WriteLine(bg.CurrentRoutine, "rset "+regname+" 1")
				bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "rset", I_NIL}

			} else {
				bgsw.Set_faulty("Allocation failed")
				return nil

			}

		}

		switchbody := x.Body

		if switchbody.List != nil {
			// There is a default case or not
			defaulted := false
			var default_point int

			cases := len(switchbody.List)
			starting_points := make([]int, cases)

			// The first loop set all the jumps toward cases
			for i, clause := range switchbody.List {
				switch cc := clause.(type) {
				case *ast.CaseClause:

					if cc.List != nil {
						// This is a case
						for _, caseexpr := range cc.List {
							// Now eval the condition
							if newcell, ok := bgsw.Expr_eval(caseexpr); ok {

								condevalued := newcell[0]

								if Same_Type(tagexpr.Vtype, condevalued.Vtype) {

									gent_bool, _ := Type_from_string("bool")
									gent_basic_type, _ := Type_from_string(bg.Basic_type)

									if Same_Type(tagexpr.Vtype, gent_bool) || Same_Type(tagexpr.Vtype, gent_basic_type) {

										regname := procbuilder.Get_register_name(tagexpr.Id)
										newregname := procbuilder.Get_register_name(condevalued.Id)

										bgsw.WriteLine(bgsw.CurrentRoutine, "je "+regname+" "+newregname+" <<"+bgsw.CurrentSwitch+"CASE"+strconv.Itoa(i)+">>")
										bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "je", I_NIL}

										bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, condevalued}
										if (<-bgsw.Answers).AnsType != ANS_OK {
											bg.Set_faulty("Resource clean failed")
											return nil
										}
									} else {
										bg.Set_faulty("Variable type not supported")
										return nil
									}

								} else {
									bg.Set_faulty("Type mismatch in case clause")
									return nil
								}
							} else {
								bg.Set_faulty("Case condition evaluation failed")
								return nil
							}

						}
					} else {
						// This is the default case
						if defaulted {
							bg.Set_faulty("Duplicate default case")
							return nil
						}
						defaulted = true
					}

				default:
					bgsw.Set_faulty("Wrong case")
					return nil
				}
			}

			if defaulted { // Jump to the end of the switch
				bgsw.WriteLine(bg.CurrentRoutine, "j <<"+bgsw.CurrentSwitch+"CASEDEFAULT>>")
				bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "j", I_NIL}
			} else {
				bgsw.WriteLine(bg.CurrentRoutine, "j <<"+bgsw.CurrentSwitch+"SWEND>>")
				bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "j", I_NIL}
			}

			// The second loop create the cases
			for i, clause := range switchbody.List {
				switch cc := clause.(type) {
				case *ast.CaseClause:

					if cc.List != nil {
						// This is a case
						starting_points[i] = bgsw.CountLines(bgsw.CurrentRoutine)
					} else {
						// This is the default case
						default_point = bgsw.CountLines(bgsw.CurrentRoutine)
					}

					ast.Walk(bgsw, clause)

				default:
					bgsw.Set_faulty("Wrong case")
					return nil
				}

				if i < len(switchbody.List)-1 {
					bgsw.Replacer(bgsw.CurrentRoutine, "<<"+bgsw.CurrentSwitch+"FALLTHROUGH>>", "<<"+bgsw.CurrentSwitch+"CASE"+strconv.Itoa(i+1)+">>")
				} else {
					if bgsw.Checker(bgsw.CurrentRoutine, "<<"+bgsw.CurrentSwitch+"FALLTHROUGH>>") {
						bgsw.Set_faulty("Fallthrought on the last case is not permitted")
						return nil
					}
				}

				bgsw.WriteLine(bg.CurrentRoutine, "j <<"+bgsw.CurrentSwitch+"SWEND>>")
				bgsw.Used <- UsageNotify{TR_PROC, bgsw.CurrentRoutine, C_OPCODE, "j", I_NIL}
			}

			switch_end := bgsw.CountLines(bgsw.CurrentRoutine)

			for i, _ := range switchbody.List {
				bgsw.Replacer(bgsw.CurrentRoutine, "<<"+bgsw.CurrentSwitch+"CASE"+strconv.Itoa(i)+">>", "<<"+strconv.Itoa(starting_points[i])+">>")
			}
			bgsw.Replacer(bgsw.CurrentRoutine, "<<"+bgsw.CurrentSwitch+"CASEDEFAULT>>", "<<"+strconv.Itoa(default_point)+">>")
			bgsw.Replacer(bgsw.CurrentRoutine, "<<"+bgsw.CurrentSwitch+"SWEND>>", "<<"+strconv.Itoa(switch_end)+">>")

			// Shift eventually created reference to line number within the code
			bgsw.Shift_program_location(bgsw.CurrentRoutine, starting_point)

			for _, line := range bgsw.GetProgram(bgsw.CurrentRoutine) {
				bg.WriteLine(bg.CurrentRoutine, line)
			}

		}
		return nil

	case *ast.LabeledStmt:
		{
			if bg.In_debug() {
				fmt.Printf("%p - Entering the new label ", bg)
			}

			vars := make(map[string]VarCell)

			obj := x.Label.Obj
			label := obj.Name

			bglabel := &BondgoCheck{bg.BondgoResults, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, bg, nil, vars, bg.Returns, bg.CurrentLoop, bg.CurrentSwitch, label, bg.CurrentRoutine}

			if bg.In_debug() {
				fmt.Printf("%p\n", bglabel)
			}

			ast.Walk(bglabel, x.Stmt)

			return nil
		}

	case *ast.GoStmt:
		if bg.In_debug() {
			fmt.Printf("%p - Go statement\n", bg)
		}

		callExpr := x.Call

		// This is a function call
		var functcell FunctCell

		// TODO Some check here
		switch callExpr.Fun.(type) {
		case *ast.Ident:

			if bg.CurrentLoop != "" {
				bg.Set_faulty("Goroutines cannot be lauched from within a loop in bondgo")
				return nil
			}

			fun := callExpr.Fun.(*ast.Ident)
			funname := fun.Name

			if _, ok := bg.Functions[funname]; ok {
				functcell = bg.Functions[funname]
			} else {
				bg.Set_faulty("Undefined function " + funname)
				return nil
			}

			// Allocate variables for the function arguments and fill them with evalued values of the arguments
			vars := make(map[string]VarCell)

			if len(functcell.Inputs) == len(callExpr.Args) {
				for i, arg := range functcell.Inputs {
					argname := arg.Argname
					if _, ok := vars[argname]; ok {
						bg.Set_faulty("Already defined variable")
						return nil
					} else {
						if cell, ok := bg.Expr_eval(callExpr.Args[i]); !ok {
							bg.Set_faulty("Wrong evaluation")
							return nil
						} else {
							vars[argname] = cell[0]
						}
					}
				}
			} else {
				bg.Set_faulty("Call with wrong number of parameters")
				return nil
			}

			// Evaluate if a channel is needed, this happen if some data has to passed to the other processor before it start computing
			// Create also the variables for the other processor
			newvars := make(map[string]VarCell)
			newreturns := make([]VarCell, 0)

			next_currentroutine := len(bg.Program)

			proccode := new(BondgoRoutine)
			proccode.Lines = make([]string, 0)
			bg.Program[next_currentroutine] = proccode

			// Prepare the new bondgocheck for the other processor
			bggoroutine := &BondgoCheck{bg.BondgoResults, bg.BondgoConfig, bg.BondgoRequirements, bg.BondgoRuninfo, bg.BondgoMessages, bg.BondgoFunctions, bg.Used, bg.Reqs, bg.Answers, nil, nil, newvars, newreturns, "", "", bg.CurrentDevice, next_currentroutine}

			// Establish the Device parameter for the next goroutine
			bggoroutine.Used <- UsageNotify{TR_PROC, next_currentroutine, C_DEVICE, bg.CurrentDevice, I_NIL}

			// fmt.Println(len(bg.Program))

			needchan := false
			for varname, cell := range vars {
				if gent, _ := Type_from_string(bg.Basic_type); Same_Type(cell.Vtype, gent) {
					needchan = true
					bggoroutine.Reqs <- VarReq{REQ_NEW, bggoroutine.CurrentRoutine, cell}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newcell := resp.Cell

						newvars[varname] = newcell

						switch newcell.Procobjtype {
						case REGISTER:
							regname := procbuilder.Get_register_name(newcell.Id)
							bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "clr "+regname)
							bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "clr", I_NIL}
						case MEMORY:
							bggoroutine.Reqs <- VarReq{REQ_NEW, bggoroutine.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
							newresp := <-bggoroutine.Answers
							if newresp.AnsType == ANS_OK {
								newregcell := newresp.Cell

								regname := procbuilder.Get_register_name(newregcell.Id)

								bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "clr "+regname)
								bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "clr", I_NIL}
								bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(newcell.Id))
								bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "r2m", I_NIL}

								bggoroutine.Reqs <- VarReq{REQ_REMOVE, bggoroutine.CurrentRoutine, newregcell}
								if (<-bggoroutine.Answers).AnsType != ANS_OK {
									bggoroutine.Set_faulty("Resource clean failed")
									return nil
								}
							} else {
								bggoroutine.Set_faulty("Resource reservation failed")
								return nil
							}
						}
					} else {
						bg.Set_faulty("Allocation failed")
						return nil
					}
				} else if gent, _ := Type_from_string("bool"); Same_Type(cell.Vtype, gent) {
					needchan = true
					bggoroutine.Reqs <- VarReq{REQ_NEW, bggoroutine.CurrentRoutine, cell}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newcell := resp.Cell

						newvars[varname] = newcell

						switch newcell.Procobjtype {
						case REGISTER:
							regname := procbuilder.Get_register_name(newcell.Id)
							bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "clr "+regname)
							bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "clr", I_NIL}
						case MEMORY:
							bggoroutine.Reqs <- VarReq{REQ_NEW, bggoroutine.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
							newresp := <-bggoroutine.Answers
							if newresp.AnsType == ANS_OK {
								newregcell := newresp.Cell

								regname := procbuilder.Get_register_name(newregcell.Id)

								bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "clr "+regname)
								bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "clr", I_NIL}
								bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(newcell.Id))
								bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "r2m", I_NIL}

								bggoroutine.Reqs <- VarReq{REQ_REMOVE, bggoroutine.CurrentRoutine, newregcell}
								if (<-bggoroutine.Answers).AnsType != ANS_OK {
									bggoroutine.Set_faulty("Resource clean failed")
									return nil
								}
							} else {
								bggoroutine.Set_faulty("Resource reservation failed")
								return nil
							}
						}
					} else {
						bg.Set_faulty("Allocation failed")
						return nil
					}
				} else if gent, _ := Type_from_string(bg.Basic_chantype); Same_Type(cell.Vtype, gent) {
					bggoroutine.Reqs <- VarReq{REQ_ATTACH, bggoroutine.CurrentRoutine, cell}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newcell := resp.Cell

						newvars[varname] = newcell
					} else {
						bggoroutine.Set_faulty("Channel attach failed")
						return nil
					}
				} else if gent, _ := Type_from_string("chan bool"); Same_Type(cell.Vtype, gent) {
					bggoroutine.Reqs <- VarReq{REQ_ATTACH, bggoroutine.CurrentRoutine, cell}
					resp := <-bg.Answers
					if resp.AnsType == ANS_OK {
						newcell := resp.Cell

						newvars[varname] = newcell
					} else {
						bggoroutine.Set_faulty("Channel attach failed")
						return nil
					}
				} else {
					bg.Set_faulty("Unsupported type")
					return nil
				}
			}

			if needchan {
				// Prepare the channel on this side (if there is data to be passed)
				gent, _ := Type_from_string(bg.Basic_chantype)
				bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, CHANNEL, 0, 0, 0, 0, 0, 0}}
				resp := <-bg.Answers
				if resp.AnsType == ANS_OK {
					cell := resp.Cell

					// Prepare the channel on the other side (if)
					bggoroutine.Reqs <- VarReq{REQ_ATTACH, bggoroutine.CurrentRoutine, cell}
					oresp := <-bggoroutine.Answers
					if oresp.AnsType == ANS_OK {
						ocell := oresp.Cell

						// These channel is implicit and does not result as variable

						// Send the passed by value data to the channel
						channame := procbuilder.Get_channel_name(cell.Id)

						for _, cell := range vars {
							gent1, _ := Type_from_string(bg.Basic_type)
							gent2, _ := Type_from_string("bool")
							if Same_Type(cell.Vtype, gent1) || Same_Type(cell.Vtype, gent2) {

								switch cell.Procobjtype {
								case REGISTER:
									regname := procbuilder.Get_register_name(cell.Id)
									bg.WriteLine(bg.CurrentRoutine, "wwr "+regname+" "+channame)
									bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "wwr", I_NIL}
									bg.WriteLine(bg.CurrentRoutine, "chw")
									bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "chw", I_NIL}
								case MEMORY:
									bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
									newresp := <-bg.Answers
									if newresp.AnsType == ANS_OK {
										newregcell := newresp.Cell

										regname := procbuilder.Get_register_name(newregcell.Id)

										bg.WriteLine(bg.CurrentRoutine, "clr "+regname)
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "clr", I_NIL}
										bg.WriteLine(bg.CurrentRoutine, "m2r "+regname+" "+strconv.Itoa(cell.Id))
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "m2r", I_NIL}
										bg.WriteLine(bg.CurrentRoutine, "wwr "+regname+" "+channame)
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "wwr", I_NIL}
										bg.WriteLine(bg.CurrentRoutine, "chw")
										bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "chw", I_NIL}

										bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, newregcell}
										if (<-bg.Answers).AnsType != ANS_OK {
											bg.Set_faulty("Resource clean failed")
											return nil
										}
									} else {
										bg.Set_faulty("Resource reservation failed")
										return nil
									}
								}
							}
						}

						// Get the data passed by value from the channel on the other side
						ochanname := procbuilder.Get_channel_name(ocell.Id)

						for _, cell := range newvars {
							gent1, _ := Type_from_string(bg.Basic_type)
							gent2, _ := Type_from_string("bool")
							if Same_Type(cell.Vtype, gent1) || Same_Type(cell.Vtype, gent2) {

								switch cell.Procobjtype {
								case REGISTER:
									regname := procbuilder.Get_register_name(cell.Id)
									bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "wrd "+regname+" "+ochanname)
									bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "wrd", I_NIL}
									bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "chw")
									bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "chw", I_NIL}
								case MEMORY:
									bggoroutine.Reqs <- VarReq{REQ_NEW, bggoroutine.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
									newresp := <-bggoroutine.Answers
									if newresp.AnsType == ANS_OK {
										newregcell := newresp.Cell

										regname := procbuilder.Get_register_name(newregcell.Id)

										bggoroutine.WriteLine(bggoroutine.CurrentRoutine, "clr "+regname)
										bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "clr", I_NIL}
										bggoroutine.WriteLine(bg.CurrentRoutine, "wrd "+regname+" "+ochanname)
										bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "wrd", I_NIL}
										bggoroutine.WriteLine(bg.CurrentRoutine, "chw")
										bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "chw", I_NIL}
										bggoroutine.WriteLine(bg.CurrentRoutine, "r2m "+regname+" "+strconv.Itoa(cell.Id))
										bggoroutine.Used <- UsageNotify{TR_PROC, bggoroutine.CurrentRoutine, C_OPCODE, "m2r", I_NIL}

										bggoroutine.Reqs <- VarReq{REQ_REMOVE, bggoroutine.CurrentRoutine, newregcell}
										if (<-bggoroutine.Answers).AnsType != ANS_OK {
											bggoroutine.Set_faulty("Resource clean failed")
											return nil
										}
									} else {
										bggoroutine.Set_faulty("Resource reservation failed")
										return nil
									}
								}
							}
						}

					} else {
						bg.Set_faulty("Channel creation failed")
						return nil
					}
				} else {
					bg.Set_faulty("Channel creation failed")
					return nil
				}

			}

			// Compile the goroutine for the other processor
			ast.Walk(bggoroutine, functcell.Body)

			return nil

		default:
			bg.Set_faulty("Unknown function type")
			return nil
		}

	case *ast.BranchStmt:
		if bg.In_debug() {
			fmt.Println("Branch")
		}
		switch x.Tok {
		case token.BREAK:
			if bg.CurrentLoop != "" {
				bg.WriteLine(bg.CurrentRoutine, "j <<"+bg.CurrentLoop+"ENDFOR>>")
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "j", I_NIL}
			} else {
				bg.Set_faulty("break outside a loop")
				return nil
			}
		case token.CONTINUE:
			if bg.CurrentLoop != "" {
				bg.WriteLine(bg.CurrentRoutine, "j <<"+bg.CurrentLoop+"CONTINUEFOR>>")
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "j", I_NIL}
			} else {
				bg.Set_faulty("continue outside a loop")
				return nil
			}
		case token.FALLTHROUGH:
			if bg.CurrentSwitch != "" && !strings.HasPrefix(bg.CurrentSwitch, "SEL") {
				bg.WriteLine(bg.CurrentRoutine, "j <<"+bg.CurrentSwitch+"FALLTHROUGH>>")
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "j", I_NIL}
			} else {
				bg.Set_faulty("fallthrough outside a switch statement")
				return nil
			}
		}
	case *ast.ReturnStmt:
		// TODO this has to change, on the Check fill Returns with empty varcell of the right type and here check the type and substitute the append with an assignment
		if x.Results != nil {
			for _, resul := range x.Results {
				if newcell, ok := bg.Expr_eval(resul); ok {
					for scope := bg; ; scope = scope.Outer {
						if scope.Outer == nil {
							scope.Returns = append(scope.Returns, newcell[0])
							break
						}
					}
				} else {
					bg.Set_faulty("Return evaluation failed")
					return nil
				}
			}
		}
		bg.WriteLine(bg.CurrentRoutine, "j <<LASTN>>")
	case *ast.SendStmt:
		if bg.In_debug() {
			fmt.Println("Send Statement")
		}
		channelname := x.Chan.(*ast.Ident).Name

		var destchan VarCell

		chanexist := false
		for scope := bg; scope != nil; scope = scope.Outer {
			if _, ok := scope.Vars[channelname]; ok {
				chanexist = true
				destchan = scope.Vars[channelname]
				break
			}
		}

		if !chanexist {
			bg.Set_faulty("Unitialized channel " + channelname)
			return nil
		}

		value := x.Value

		if newcell, ok := bg.Expr_eval(value); ok {
			// TODO Missing types checks
			gent, _ := Type_from_string(bg.Basic_type)
			bg.Reqs <- VarReq{REQ_NEW, bg.CurrentRoutine, VarCell{gent, REGISTER, 0, 0, 0, 0, 0, 0}}
			resp := <-bg.Answers
			if resp.AnsType == ANS_OK {

				cell := newcell[0]
				cell2 := resp.Cell
				regname := procbuilder.Get_register_name(cell.Id)
				channame := procbuilder.Get_channel_name(destchan.Id)
				waitregname := procbuilder.Get_register_name(cell2.Id)

				bg.WriteLine(bg.CurrentRoutine, "wwr "+regname+" "+channame)
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "wwr", I_NIL}
				bg.WriteLine(bg.CurrentRoutine, "chw "+waitregname)
				bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "chw", I_NIL}

				bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell}
				if (<-bg.Answers).AnsType != ANS_OK {
					bg.Set_faulty("Resource clean failed")
					return nil
				}
				bg.Reqs <- VarReq{REQ_REMOVE, bg.CurrentRoutine, cell2}
				if (<-bg.Answers).AnsType != ANS_OK {
					bg.Set_faulty("Resource clean failed")
					return nil
				}
			} else {
				bg.Set_faulty("Resource reservation failed")
				return nil
			}

		} else {
			bg.Set_faulty("Send evaluation failed")
			return nil
		}

	case *ast.CallExpr:

		if bg.In_debug() {
			fmt.Printf("%p - %s\n", bg, "Call Statement")
		}

		switch fun := x.Fun.(type) {
		case (*ast.SelectorExpr):
			// This is the case of a "Method"
			xf := fun.X.(*ast.Ident)
			sel := fun.Sel

			if xf.Name == "bondgo" {
				switch sel.Name {
				case "Void":
				case "IOWrite":

					args := x.Args
					if len(args) == 2 {
						// TODO Finish
						switch arg := args[0].(type) {
						case *ast.Ident:
							identname := arg.Name

							varexist := false
							for scope := bg; scope != nil; scope = scope.Outer {
								if cell, ok := scope.Vars[identname]; ok {
									varexist = true
									if bg.In_debug() {
										fmt.Printf("%p - Found variable %s in scope %p\n", bg, identname, scope)
									}
									switch cell.Procobjtype {
									case OUTPUT:
										if cell.Global_id == 0 {
											bg.Set_faulty("Output not allocated")
											return nil
										} else {
											if bg.In_debug() {
												fmt.Printf("%p - IO Write call\n", bg)
											}
											if ncell, ok := bg.Expr_eval(args[1]); ok {
												newcell := ncell[0]
												newregname := procbuilder.Get_register_name(newcell.Id)
												newoutname := procbuilder.Get_output_name(cell.Id)
												bg.WriteLine(bg.CurrentRoutine, "r2o "+newregname+" "+newoutname)
												bg.Used <- UsageNotify{TR_PROC, bg.CurrentRoutine, C_OPCODE, "r2o", I_NIL}

											} else {
												bg.Set_faulty("Write evaluation failed")
												return nil
											}

										}
									default:
										bg.Set_faulty("Write can only be used on output registers")
										return nil
									}
									break
								}
							}
							if !varexist {
								bg.Set_faulty("Variable " + identname + " not defined")
								return nil
							}
						default:
							bg.Set_faulty("The fiers argument has to be an output register")
							return nil
						}

					} else {
						bg.Set_faulty("Two arguments expected")
						return nil
					}

				default:
					bg.Set_faulty("Unknown function " + sel.Name)
					return nil
				}
			} else {
				bg.Set_faulty("Unknown module " + xf.Name)
				return nil
			}

		case (*ast.Ident):
			// This id the case of a function with no receiver
			// TODO Finish
		}
		return nil
	}

	return bg
}
