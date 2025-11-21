package simbox

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	TIMEC_ABS      uint8 = 0 + iota // Rules in absolute time
	TIMEC_NONE                      // Rules with no temporal components
	TIMEC_REL                       // Relative time aka periodic
	TIMEC_ON_VALID                  // On valid signal
	TIMEC_ON_RECV                   // On receive signal
)

const (
	ACTION_SET    uint8 = 0 + iota // Set a value to an object
	ACTION_GET                     // Get a value from an object
	ACTION_SHOW                    // Show a value from an object
	ACTION_CONFIG                  // Configuration directive
)

type Rule struct {
	Timec  uint8  // Time constraint type: absolute, relative, none, on valid, on receive
	Tick   uint64 // Tick (if applicable)
	Action uint8  // Action: set, get, show, config
	Object string // Object: register, memory, io, config option
	Extra  string // Extra info: unsigned, signed, hex, etc.
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
			case "show_all":
				return "config:show_all:" + rule.Extra
			case "show_all_internal":
				return "config:show_all_internal:" + rule.Extra
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
	case TIMEC_ON_VALID:
		switch rule.Action {
		case ACTION_GET:
			return "onvalid:get:" + rule.Object + ":" + rule.Extra
		case ACTION_SHOW:
			return "onvalid:show:" + rule.Object + ":" + rule.Extra
		}
	case TIMEC_ON_RECV:
		switch rule.Action {
		case ACTION_GET:
			return "onrecv:get:" + rule.Object + ":" + rule.Extra
		case ACTION_SHOW:
			return "onrecv:show:" + rule.Object + ":" + rule.Extra
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
	return errors.New("index out of range")
}

func (r *Simbox) Add(adds string) error {
	words := strings.Split(adds, ":")
	if len(words) == 5 {
		if words[0] == "absolute" && words[2] == "set" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_ABS,
					Tick:   uint64(tick),
					Action: ACTION_SET,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "set" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_REL,
					Tick:   uint64(every),
					Action: ACTION_SET,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
		if words[0] == "absolute" && words[2] == "get" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_ABS,
					Tick:   uint64(tick),
					Action: ACTION_GET,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "get" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_REL,
					Tick:   uint64(every),
					Action: ACTION_GET,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
		if words[0] == "absolute" && words[2] == "show" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_ABS,
					Tick:   uint64(tick),
					Action: ACTION_SHOW,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "show" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_REL,
					Tick:   uint64(every),
					Action: ACTION_SHOW,
					Object: words[3],
					Extra:  words[4],
				})
				return nil
			}
		}
	} else if len(words) == 4 {
		if words[0] == "absolute" && words[2] == "get" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_ABS,
					Tick:   uint64(tick),
					Action: ACTION_GET,
					Object: words[3],
					Extra:  "unsigned",
				})
				return nil
			}
		} else if words[0] == "absolute" && words[2] == "show" {
			if tick, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_ABS,
					Tick:   uint64(tick),
					Action: ACTION_SHOW,
					Object: words[3],
					Extra:  "unsigned",
				})
				return nil
			}
		}
		if words[0] == "relative" && words[2] == "get" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_REL,
					Tick:   uint64(every),
					Action: ACTION_GET,
					Object: words[3],
					Extra:  "unsigned",
				})
				return nil
			}
		} else if words[0] == "relative" && words[2] == "show" {
			if every, err := strconv.Atoi(words[1]); err == nil {
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_REL,
					Tick:   uint64(every),
					Action: ACTION_SHOW,
					Object: words[3],
					Extra:  "unsigned",
				})
				return nil
			}
		} else if words[0] == "onvalid" && words[1] == "get" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_VALID,
				Tick:   0,
				Action: ACTION_GET,
				Object: words[2],
				Extra:  words[3],
			})
			return nil
		} else if words[0] == "onvalid" && words[1] == "show" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_VALID,
				Tick:   0,
				Action: ACTION_SHOW,
				Object: words[2],
				Extra:  words[3],
			})
			return nil
		} else if words[0] == "onrecv" && words[1] == "get" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_RECV,
				Tick:   0,
				Action: ACTION_GET,
				Object: words[2],
				Extra:  words[3],
			})
			return nil
		} else if words[0] == "onrecv" && words[1] == "show" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_RECV,
				Tick:   0,
				Action: ACTION_SHOW,
				Object: words[2],
				Extra:  words[3],
			})
			return nil
		}
	} else if len(words) == 3 {
		if words[0] == "config" {
			switch words[1] {
			case "get_all":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "get_all",
					Extra:  words[2],
				})
				return nil
			case "get_all_internal":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "get_all_internal",
					Extra:  words[2],
				})
				return nil
			case "show_all":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_all",
					Extra:  words[2],
				})
				return nil
			case "show_all_internal":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_all_internal",
					Extra:  words[2],
				})
				return nil
			}
		} else if words[0] == "onvalid" && words[1] == "get" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_VALID,
				Tick:   0,
				Action: ACTION_GET,
				Object: words[2],
				Extra:  "unsigned",
			})
			return nil
		} else if words[0] == "onvalid" && words[1] == "show" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_VALID,
				Tick:   0,
				Action: ACTION_SHOW,
				Object: words[2],
				Extra:  "unsigned",
			})
			return nil
		} else if words[0] == "onrecv" && words[1] == "get" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_RECV,
				Tick:   0,
				Action: ACTION_GET,
				Object: words[2],
				Extra:  "unsigned",
			})
			return nil
		} else if words[0] == "onrecv" && words[1] == "show" {
			r.Rules = append(r.Rules, Rule{
				Timec:  TIMEC_ON_RECV,
				Tick:   0,
				Action: ACTION_SHOW,
				Object: words[2],
				Extra:  "unsigned",
			})
			return nil
		}
	} else if len(words) == 2 {
		if words[0] == "config" {
			switch words[1] {
			case "show_pc":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_pc",
					Extra:  "",
				})
				return nil
			case "show_instruction":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_instruction",
					Extra:  "",
				})
				return nil
			case "show_disasm":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_disasm",
					Extra:  "",
				})
				return nil
			case "show_ticks":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_ticks",
					Extra:  "",
				})
				return nil
			case "get_ticks":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "get_ticks",
					Extra:  "",
				})
				return nil
			case "show_proc_regs_pre":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_proc_regs_pre",
					Extra:  "",
				})
				return nil
			case "show_proc_regs_post":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_proc_regs_post",
					Extra:  "",
				})
				return nil
			case "show_proc_io_pre":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_proc_io_pre",
					Extra:  "",
				})
				return nil
			case "show_proc_io_post":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_proc_io_post",
					Extra:  "",
				})
				return nil
			case "show_io_pre":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_io_pre",
					Extra:  "",
				})
				return nil
			case "show_io_post":
				r.Rules = append(r.Rules, Rule{
					Timec:  TIMEC_NONE,
					Tick:   0,
					Action: ACTION_CONFIG,
					Object: "show_io_post",
					Extra:  "",
				})
				return nil
			}
		}
	}
	return errors.New("rule cannot be decoded")
}
