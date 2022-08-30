package neuralbond

type TrainedNet struct {
	Nodes   []Node
	Weights []Weight
}

type Weight struct {
	Layer        int
	PosCurrLayer int
	PosPrevLayer int
	Value        float32
}

type Node struct {
	Layer int
	Pos   int
	Type  string
	Bias  float32
}
