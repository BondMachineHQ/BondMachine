package bmline

import "regexp"

func FilterMatcher(m *BasmLine, filterType string) bool {

	// The operand is not checked. This can change if necessary

	for _, elem := range m.Elements {
		if elem.GetMeta("type") == filterType {
			return true
		}
	}

	return false
}

func MatchMatcher(m *BasmLine, l *BasmLine) bool {
	if MatchArg(m.Operation, l.Operation) {
		if len(m.Elements) == len(l.Elements) {
			for i, arg := range l.Elements {
				if !MatchArg(m.Elements[i], arg) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func MatchArg(m *BasmElement, l *BasmElement) bool {
	for k, v := range m.LoopMeta() {
		if v != l.GetMeta(k) {
			return false
		}
	}
	re := regexp.MustCompile("^" + m.string + "$")
	return re.MatchString(l.string)
}

// Match if two second argument has a subset of metadata of the first
func MatchMeta(m *BasmElement, l *BasmElement) bool {
	for k, v := range m.LoopMeta() {
		if v != l.GetMeta(k) {
			return false
		}
	}
	return true
}
