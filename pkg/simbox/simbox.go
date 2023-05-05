package simbox

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	TIMEC_ABS  uint8 = 0 + iota // Rules in absolute time
	TIMEC_NONE                  // Rules with no temporal components
	TIMEC_REL                   // Relative time aka periodic
)

const (
	ACTION_SET uint8 = 0 + iota
	ACTION_GET
	ACTION_SHOW
	ACTION_CONFIG
)

type Prerror struct {
	string
}

func (e Prerror) Error() string {
	return e.string
}

type Rule struct {
	Timec  uint8  // Time contraint type:
	Tick   uint64 // Tick (if appliable)
	Action uint8
	Object string
	Extra  string
}

type Simbox struct {
	Rules []Rule
}

type Report struct {
}

func (rule Rule) String() string {
	switch rule.Timec {
	case TIMEC_ABS:
		switch rule.Action {
		case ACTION_SET:
			return "absolute:" + strconv.Itoa(int(rule.Tick)) + ":set:" + rule.Object + ":" + rule.Extra
		case ACTION_GET:
			return "absolute:" + strconv.Itoa(int(rule.Tick)) + ":get:" + rule.Object + ":" + rule.Extra
		case ACTION_SHOW:
			return "absolute:" + strconv.Itoa(int(rule.Tick)) + ":show:" + rule.Object + ":" + rule.Extra
		}
	case TIMEC_NONE:
		switch rule.Action {
		case ACTION_CONFIG:
			switch rule.Object {
			case "get_all":
				return "config:get_all:" + rule.Extra
			case "get_all_internal":
				return "config:get_all_internal:" + rule.Extra
			default:
				return "config:" + rule.Object
			}
		}
	case TIMEC_REL:
		switch rule.Action {
		case ACTION_SET:
			return "relative:" + strconv.Itoa(int(rule.Tick)) + ":set:" + rule.Object + ":" + rule.Extra
		case ACTION_GET:
			return "relative:" + strconv.Itoa(int(rule.Tick)) + ":get:" + rule.Object + ":" + rule.Extra
		case ACTION_SHOW:
			return "relative:" + strconv.Itoa(int(rule.Tick)) + ":show:" + rule.Object + ":" + rule.Extra
		}
	}
	return ""
}

func (r *Simbox) Print() string {
	result := ""
	for i, rule := range r.Rules {
		result = result + fmt.Sprintf("%03d - ", i) + rule.String() + "\n"
	}
	return result
}

func (r *Simbox) Del(idx int) error {
	if idx < len(r.Rules) {
		r.Rules = append(r.Rules[:idx], r.Rules[idx+1:]...)
		return nil
	}
	return Prerror{"Wrong rule index"}
}

func (r *Simbox) Add(adds string) error {
	words := strings.Split(adds, ":")
	if len(words) == 5 {
		if words[0] == "absolute" && words[2] == "set" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_ABS, uint64(tick), ACTION_SET, words[3], words[4]})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "set" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_REL, uint64(every), ACTION_SET, words[3], words[4]})
				return nil
			}
		}
		if words[0] == "absolute" && words[2] == "get" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_ABS, uint64(tick), ACTION_GET, words[3], words[4]})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "get" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_REL, uint64(every), ACTION_GET, words[3], words[4]})
				return nil
			}
		}
		if words[0] == "absolute" && words[2] == "show" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_ABS, uint64(tick), ACTION_SHOW, words[3], words[4]})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "show" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_REL, uint64(every), ACTION_SHOW, words[3], words[4]})
				return nil
			}
		}
	} else if len(words) == 4 {
		if words[0] == "absolute" && words[2] == "get" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_ABS, uint64(tick), ACTION_GET, words[3], "unsigned"})
				return nil
			}
		} else if words[0] == "absolute" && words[2] == "show" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_ABS, uint64(tick), ACTION_SHOW, words[3], "unsigned"})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "get" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_REL, uint64(every), ACTION_GET, words[3], "unsigned"})
				return nil
			}
		} else if words[0] == "relative" && words[2] == "show" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{TIMEC_REL, uint64(every), ACTION_SHOW, words[3], "unsigned"})
				return nil
			}
		}
	} else if len(words) == 3 {
		if words[0] == "config" {
			switch words[1] {
			case "get_all":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "get_all", words[2]})
				return nil
			case "get_all_internal":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "get_all_internal", words[2]})
				return nil
			}
		}
	} else if len(words) == 2 {
		if words[0] == "config" {
			switch words[1] {
			case "show_pc":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_pc", ""})
				return nil
			case "show_instruction":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_instruction", ""})
				return nil
			case "show_disasm":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_disasm", ""})
				return nil
			case "show_ticks":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_ticks", ""})
				return nil
			case "get_ticks":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "get_ticks", ""})
				return nil
			case "show_proc_regs_pre":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_proc_regs_pre", ""})
				return nil
			case "show_proc_regs_post":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_proc_regs_post", ""})
				return nil
			case "show_proc_io_pre":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_proc_io_pre", ""})
				return nil
			case "show_proc_io_post":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_proc_io_post", ""})
				return nil
			case "show_io_pre":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_io_pre", ""})
				return nil
			case "show_io_post":
				r.Rules = append(r.Rules, Rule{TIMEC_NONE, uint64(0), ACTION_CONFIG, "show_io_post", ""})
				return nil
			}
		}
	}
	return Prerror{"Rule cannot be decoded"}
}
