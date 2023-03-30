package procbuilder

const (
	HwOptimize = uint64(1)
)

type HwOptimizations uint64

func HwOptimizationId(name string) HwOptimizations {
	switch name {
	case "hwoptimize":
		return HwOptimizations(HwOptimize)
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
