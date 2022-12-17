package bmline

func (body *BasmBody) Copy() *BasmBody {
	if body == nil {
		return nil
	}
	result := new(BasmBody)
	for k, v := range body.BasmMeta.LoopMeta() {
		result.BasmMeta = result.BasmMeta.SetMeta(k, v)
	}
	result.Lines = make([]*BasmLine, len(body.Lines))
	for i, line := range body.Lines {
		result.Lines[i] = line.Copy()
	}
	return result
}

func (line *BasmLine) Copy() *BasmLine {
	if line == nil {
		return nil
	}
	result := new(BasmLine)
	for k, v := range line.BasmMeta.LoopMeta() {
		result.BasmMeta = result.BasmMeta.SetMeta(k, v)
	}
	result.Operation = line.Operation.Copy()
	result.Elements = make([]*BasmElement, len(line.Elements))
	for i, el := range line.Elements {
		result.Elements[i] = el.Copy()
	}
	return result
}

// Create a copy of the BasmElement
func (el *BasmElement) Copy() *BasmElement {
	if el == nil {
		return nil
	}
	result := new(BasmElement)
	result.string = el.string
	for k, v := range el.BasmMeta.LoopMeta() {
		result.BasmMeta = result.BasmMeta.SetMeta(k, v)
	}
	return result
}
