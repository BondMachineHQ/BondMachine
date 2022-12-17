package bmline

import "fmt"

func purple(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[35m" + ins + "\033[0m"
}
func green(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[32m" + ins + "\033[0m"
}
func yellow(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[33m" + ins + "\033[0m"
}
func blue(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[34m" + ins + "\033[0m"
}
func red(in ...interface{}) string {
	ins := fmt.Sprint(in...)
	return "\033[31m" + ins + "\033[0m"
}

func (body *BasmBody) String() string {
	result := ""
	if body == nil {
		result += blue("[]\n")
	} else {
		result += blue("[" + body.ListMeta() + "]\n")
		for lineN, line := range body.Lines {
			result += "\t\t\t" + fmt.Sprint(lineN) + blue("["+line.ListMeta()+"]") + ":" + line.String() + "\n"
		}
	}
	return result
}

func (line *BasmLine) String() string {
	result := line.Operation.String()
	for _, arg := range line.Elements {
		result += arg.String()
	}
	return result
}

func (el *BasmElement) String() string {
	result := ""
	if el != nil {
		result += " " + green(el.string)
		if el.BasmMeta == nil {
			result += blue("[]")
		} else {
			result += blue("[" + el.ListMeta() + "]")
		}
	}
	return result
}
