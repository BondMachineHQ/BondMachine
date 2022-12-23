package bondgo

import (
	"strconv"
	"strings"
)

type BondgoRoutine struct {
	Lines []string
}

func (p *BondgoRoutine) Shift_program_location(n int) {
	if p != nil {
		for i, line := range p.Lines {
			if strings.Contains(line, "<<") {
				fs := strings.Split(line, "<<")
				first := fs[0]
				ss := strings.Split(fs[1], ">>")
				last := ss[1]
				number_s := ss[0]
				if number, err := strconv.Atoi(number_s); err == nil {
					p.Lines[i] = first + "<<" + strconv.Itoa(number+n) + ">>" + last
				}
			}
		}
	}
}

func (p *BondgoRoutine) Remove_program_location() {
	if p != nil {
		for i, line := range p.Lines {
			if strings.Contains(line, "<<") {
				fs := strings.Split(line, "<<")
				first := fs[0]
				ss := strings.Split(fs[1], ">>")
				last := ss[1]
				number_s := ss[0]
				p.Lines[i] = first + number_s + last
			}
		}
	}
}

func (p *BondgoRoutine) Replacer(from string, to string) {
	if p != nil {
		for i, line := range p.Lines {
			p.Lines[i] = strings.Replace(line, from, to, -1)
		}
	}
}

func (p *BondgoRoutine) Checker(ck string) bool {
	if p != nil {
		for _, line := range p.Lines {
			if strings.Contains(line, ck) {
				return true
			}
		}
	}
	return false
}

func (p *BondgoRoutine) Append(line string) {
	if p != nil {
		p.Lines = append(p.Lines, line)
	}
}
