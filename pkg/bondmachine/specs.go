package bondmachine

import "strconv"

func (b *Bondmachine) Specs() string {
	// TODO - implement this
	result := "Register size: " + strconv.Itoa(int(b.Rsize)) + "\n"
	result += "Processors:\n"

	for i, dom_id := range b.Processors {
		result += "  " + strconv.Itoa(i) + ":\n"
		result += "    Domain ID: " + strconv.Itoa(dom_id) + "\n"
		result += b.Domains[dom_id].Specs()
	}
	return result
}
