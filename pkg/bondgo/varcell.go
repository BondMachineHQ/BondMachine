package bondgo

import (
	"errors"
	"go/ast"
	"strconv"
	"strings"
)

const (
	REGISTER = uint8(0) + iota
	MEMORY
	INPUT
	OUTPUT // Used only in assignment
	CHANNEL
	SHAREDMEMORY
)

const (
	T_NAMED = uint8(0) + iota
	T_CHAN
	T_STAR
	T_STRUCT
)

type VarType struct {
	MainType uint8
	Name     string
	Values   []*VarType
}

// The type to hold variables in processor objects (mem or registers or whatever) and in bondmachine objects
type VarCell struct {
	Vtype          *VarType // The variable type
	Procobjtype    uint8    // 0 Register, 1 RAM, 2 output (notin vars) etc.
	Id             int      // Local (within the processor) reference to the object among objects of the same type
	Start_id       int      // Local (within the processor) reference to the object among objects of the same type (start data)
	End_id         int      // Local (within the processor) reference to the object among objects of the same type (end data)
	Global_id      int      // Global reference to the object among objects of the same type
	Start_globalid int      // Global reference to the object among objects of the same type (start data)
	End_globalid   int      // Global reference to the object among objects of the same type (end data)
}

// Vartype to String
func (t *VarType) String() string {
	result := ""
	switch t.MainType {
	case T_NAMED:
		result += t.Name
	case T_CHAN:
		result += "chan " + t.Values[0].String()
	case T_STAR:
		result += "* " + t.Values[0].String()
	case T_STRUCT:
		// TODO
	}
	return result
}

// Variable type equality
func Same_Type(t1 *VarType, t2 *VarType) bool {
	if t1 != nil && t2 != nil {
		if t1.MainType == t2.MainType {
			switch t1.MainType {
			case T_NAMED:
				if t1.Name == t2.Name {
					return true
				}
			case T_CHAN, T_STAR:
				if len(t1.Values) == 1 && len(t2.Values) == 1 {
					if Same_Type(t1.Values[0], t2.Values[0]) {
						return true
					}
				}
			case T_STRUCT:
				// TODO
			}
		}
	}
	return false
}

// Type from string
func Type_from_string(types string) (*VarType, error) {
	t := strings.TrimSpace(types)
	first_space := strings.Index(t, " ")
	if first_space == -1 {
		if t != "" {
			switch t {
			case "uint8", "uint16", "uint32", "uint64", "bool":
				newtype := new(VarType)
				newtype.MainType = T_NAMED
				newtype.Name = t
				newtype.Values = make([]*VarType, 0)
				return newtype, nil
			default:
				return nil, errors.New("Type " + t + " unsupported")
			}
		}
	} else {
		mod := t[:first_space]
		switch mod {
		case "*":
			if inner_type, err := Type_from_string(t[first_space:]); err == nil {
				newtype := new(VarType)
				newtype.MainType = T_STAR
				newtype.Name = ""
				newtype.Values = make([]*VarType, 1)
				newtype.Values[0] = inner_type
				return newtype, nil
			}
		case "chan":
			if inner_type, err := Type_from_string(t[first_space:]); err == nil {
				newtype := new(VarType)
				newtype.MainType = T_CHAN
				newtype.Name = ""
				newtype.Values = make([]*VarType, 1)
				newtype.Values[0] = inner_type
				return newtype, nil
			}
		case "struct":
			// TODO
		}
	}
	return nil, errors.New("Type unknown")
}

// Type from ast
func Type_from_ast(spec ast.Node) (*VarType, error) {
	switch vtype := spec.(type) {
	case *ast.Ident:
		newtype := new(VarType)
		newtype.MainType = T_NAMED
		newtype.Name = vtype.Name
		newtype.Values = make([]*VarType, 0)
		return newtype, nil
	case *ast.ChanType:
		if inner_type, err := Type_from_ast(vtype.Value); err == nil {
			newtype := new(VarType)
			newtype.MainType = T_CHAN
			newtype.Name = ""
			newtype.Values = make([]*VarType, 1)
			newtype.Values[0] = inner_type
			return newtype, nil
		}
	case *ast.StarExpr:
		if inner_type, err := Type_from_ast(vtype.X); err == nil {
			newtype := new(VarType)
			newtype.MainType = T_STAR
			newtype.Name = ""
			newtype.Values = make([]*VarType, 1)
			newtype.Values[0] = inner_type
			return newtype, nil
		}
	case *ast.StructType:
	}
	return nil, errors.New("Import failed")
}

// VarCell to string
func (m VarCell) String() string {
	result := "<"
	switch m.Procobjtype {
	case REGISTER:
		result += "reg "
	case INPUT:
		result += "input "
	case OUTPUT:
		result += "output "
	case MEMORY:
		result += "mem "
	case CHANNEL:
		result += "chan "
	}
	result += strconv.Itoa(m.Id) + ">"
	return result
}

func memused(a VarCell, list []VarCell) (int, bool) {
	for i, b := range list {
		if a == b {
			return i, true
		}
	}
	return 0, false
}
