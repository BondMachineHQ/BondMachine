package bmline

import "strconv"

// Change the meta data of the BasmBody recursively adding the prefix to every existing key in prefix
func (body *BasmBody) PrefixMeta(prefix *BasmElement) {
	if body == nil {
		return
	}
	for k, v := range prefix.LoopMeta() {
		if body.BasmMeta.GetMeta(k) != "" {
			body.BasmMeta = body.BasmMeta.SetMeta(k, v+body.BasmMeta.GetMeta(k))
		}
	}

	for _, line := range body.Lines {

		for k, v := range prefix.LoopMeta() {
			if line.BasmMeta.GetMeta(k) != "" {
				line.BasmMeta = line.BasmMeta.SetMeta(k, v+line.BasmMeta.GetMeta(k))
			}
		}

		for k, v := range prefix.LoopMeta() {
			if line.Operation.BasmMeta.GetMeta(k) != "" {
				oldValue := line.Operation.GetMeta(k)
				line.Operation.BasmMeta = line.Operation.SetMeta(k, v+oldValue)
			}
		}

		for _, elem := range line.Elements {
			for k, v := range prefix.LoopMeta() {
				if elem.BasmMeta.GetMeta(k) != "" {
					oldValue := elem.GetMeta(k)
					elem.BasmMeta = elem.SetMeta(k, v+oldValue)
				}
			}
		}
	}
}

// Change the value of the BasmBody recursively adding the prefix to every existing matching value
func (body *BasmBody) PrefixValue(prefix *BasmElement) {
	if body == nil {
		return
	}

	for _, line := range body.Lines {

		if MatchMeta(prefix, line.Operation) {
			line.Operation.SetValue(prefix.GetValue() + line.Operation.GetValue())
		}

		for _, elem := range line.Elements {
			if MatchMeta(prefix, elem) {
				elem.SetValue(prefix.GetValue() + elem.GetValue())
			}
		}
	}
}

func (body *BasmBody) CheckArg(arg *BasmElement) bool {
	if body == nil {
		return false
	}
	for _, line := range body.Lines {
		for _, elem := range line.Elements {
			if MatchArg(arg, elem) {
				return true
			}
		}
	}
	return false
}

func (body *BasmBody) NextResource(res *BasmElement) *BasmElement {
	if body == nil {
		return nil
	}

	for i := 0; i < 10000; i++ {
		result := res.Copy()
		result.SetValue(res.GetValue() + strconv.Itoa(i))
		if !body.CheckArg(result) {
			return result
		}
	}

	return nil
}

func (body *BasmBody) ReplaceArg(s *BasmElement, t *BasmElement) {
	if body == nil {
		return
	}
	for _, line := range body.Lines {
		for i, elem := range line.Elements {
			if MatchArg(s, elem) {
				line.Elements[i] = t.Copy()
			}
		}
	}
}
