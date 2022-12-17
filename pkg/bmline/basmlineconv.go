package bmline

import (
	"strings"
)

func Text2BasmLine(line string) (*BasmLine, error) {
	splitted := strings.Split(line, "::")
	argN := len(splitted)
	operand := splitted[0]

	argFlags := strings.Split(operand, "--")
	newElem := new(BasmElement)
	newElem.string = argFlags[0]

	if len(argFlags) > 1 {
		for _, meta := range argFlags[1:] {
			key, value := strings.Split(meta, "=")[0], strings.Split(meta, "=")[1]
			newElem.BasmMeta = newElem.SetMeta(key, value)
		}
	}

	newLine := new(BasmLine)
	newLine.Operation = newElem

	if argN > 1 {
		arguments := splitted[1:]
		newArgs := make([]*BasmElement, len(arguments))
		for i, arg := range arguments {
			argFlags := strings.Split(arg, "--")

			newArg := new(BasmElement)
			newArg.string = argFlags[0]

			if len(argFlags) > 1 {
				for _, meta := range argFlags[1:] {
					key, value := strings.Split(meta, "=")[0], strings.Split(meta, "=")[1]
					newArg.BasmMeta = newArg.SetMeta(key, value)
				}
			}

			newArgs[i] = newArg
		}
		newLine.Elements = newArgs
	}
	return newLine, nil
}

func BasmLine2Text(bline *BasmLine) (string, error) {

	result := ""

	result += bline.Operation.string
	for key, val := range bline.Operation.LoopMeta() {
		result += "--" + key + "=" + val
	}

	for _, arg := range bline.Elements {
		result += "::" + arg.string
		for key, val := range arg.LoopMeta() {
			result += "--" + key + "=" + val
		}

	}
	return result, nil
}
