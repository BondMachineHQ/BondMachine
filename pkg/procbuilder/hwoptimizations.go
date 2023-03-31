package procbuilder

const (
	OnlyDestRegs = uint64(1)
	OnlySrcRegs  = uint64(2)
)

type HwOptimizations uint64

func HwOptimizationId(name string) HwOptimizations {
	switch name {
	case "onlydestregs":
		return HwOptimizations(OnlyDestRegs)
	case "onlysrcregs":
		return HwOptimizations(OnlySrcRegs)
	}
	return 0
}

func SetHwOptimization(current HwOptimizations, optimization HwOptimizations) HwOptimizations {
	return current | optimization
}

func UnsetHwOptimization(current HwOptimizations, optimization HwOptimizations) HwOptimizations {
	return current &^ optimization
}

func IsHwOptimizationSet(current HwOptimizations, optimization HwOptimizations) bool {
	return (current & optimization) != 0
}
