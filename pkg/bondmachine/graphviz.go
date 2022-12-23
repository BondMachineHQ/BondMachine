package bondmachine

const (
	GVNODE = uint8(0) + iota
	GVCLUS
	GVCLUSPROC
	GVCLUSIN
	GVCLUSOUT
	GVNODEIN
	GVNODEOUT

	GVNODEINPROC
	GVCLUSINPROC
	GVNODEININPROC
	GVNODEOUTINPROC
	GVCLUSININPROC
	GVCLUSOUTINPROC

	GVEDGE

	GVPEER

	GVNODEININPEER
	GVNODEOUTINPEER
	GVNODECHINPEER

	GVCLUSININPEER
	GVCLUSOUTINPEER
	GVCLUSCHINPEER

	GVINFOPROCREGS
	GVINFOPROCOPCODES
	GVINFOPROCPROG
	GVINFOPROCPC // Processor Program Counter style

	GVINFOPROCPROGLINE
	GVINFOPROCPROGLINESEL
)

func GV_config(element uint8) string {
	result := ""
	switch element {
	case GVNODEIN:
		result += "style=filled fillcolor=greenyellow color=black"
	case GVNODEOUT:
		result += "style=filled fillcolor=lightcoral color=black"
	case GVNODEININPROC:
		result += "style=filled fillcolor=greenyellow color=black"
	case GVNODEOUTINPROC:
		result += "style=filled fillcolor=lightcoral color=black"
	case GVCLUSININPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey80"
	case GVCLUSOUTINPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey50"
	case GVCLUSIN:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey80"
	case GVCLUSOUT:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=grey50"
	case GVCLUSPROC:
		result += "style=filled;\n\t\tcolor=black;\n\t\tfillcolor=aquamarine3"

	case GVPEER:
		result += "style=\"filled, rounded\" fillcolor=coral color=grey30"
	case GVNODEININPEER:
		result += "style=\"filled\" shape=box fillcolor=lightskyblue color=black"
	case GVNODEOUTINPEER:
		result += "style=\"filled\" shape=box fillcolor=indianred3 color=black"
	case GVNODECHINPEER:
		result += "style=\"filled\" shape=box fillcolor=red color=black"
	case GVCLUSININPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey70"
	case GVCLUSOUTINPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey50"
	case GVCLUSCHINPEER:
		result += "style=\"filled, rounded\";\n\t\tcolor=black;\n\t\tfillcolor=grey30"

	case GVINFOPROCREGS:
		result += "style=\"filled\" shape=box color=black fillcolor=white"
	case GVINFOPROCOPCODES:
		result += "style=\"filled\" shape=box color=black fillcolor=white"
	case GVINFOPROCPROG:
		result += "style=\"filled\" shape=box color=black fillcolor=white"
	case GVINFOPROCPC:
		result += "style=\"filled\" shape=box color=black fillcolor=red"
	case GVINFOPROCPROGLINE:
		result += "style=\"filled\" shape=none color=black fillcolor=white"
	case GVINFOPROCPROGLINESEL:
		result += "style=\"filled\" shape=box color=black fillcolor=white"
	}
	return result
}
