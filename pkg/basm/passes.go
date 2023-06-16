package basm

const (
	passTemplateResolver      = uint64(1)
	passDynamicalInstructions = uint64(2)
	passLabelTagger           = uint64(4)
	passMetadataInfer1        = uint64(8)
	passFragmentAnalyzer      = uint64(16)
	passFragmentPruner        = uint64(32)
	passFragmentComposer      = uint64(64)
	passMetadataInfer2        = uint64(128)
	passDataSections2Bytes    = uint64(256)
	passEntryPoints           = uint64(512)
	passLabelsResolver        = uint64(1024)
	passRomComposer           = uint64(2048)
	passMatcherResolver       = uint64(4096)
	LAST_PASS                 = uint64(4096)
)

func getPassFunction() map[uint64]func(*BasmInstance) error {
	return map[uint64]func(*BasmInstance) error{
		passTemplateResolver:      templateResolver,
		passDynamicalInstructions: dynamicalInstructions,
		passLabelTagger:           labelTagger,
		passDataSections2Bytes:    dataSections2Bytes,
		passMetadataInfer1:        metadataInfer,
		passMetadataInfer2:        metadataInfer,
		passEntryPoints:           entryPoints,
		passLabelsResolver:        labelResolver,
		passMatcherResolver:       matcherResolver,
		passFragmentAnalyzer:      fragmentAnalyzer,
		passFragmentPruner:        fragmentPruner,
		passFragmentComposer:      fragmentComposer,
		passRomComposer:           romComposer,
	}
}

func getPassFunctionName() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateResolver",
		passDynamicalInstructions: "dynamicalInstructions",
		passLabelTagger:           "labelTagger",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadataInfer (1)",
		passMetadataInfer2:        "metadataInfer (2)",
		passEntryPoints:           "entryPoints",
		passLabelsResolver:        "labelResolver",
		passMatcherResolver:       "matcherResolver",
		passFragmentAnalyzer:      "fragmentAnalyzer",
		passFragmentPruner:        "fragmentPruner",
		passFragmentComposer:      "fragmentComposer",
		passRomComposer:           "romComposer",
	}
}

func activePass(passes uint64, active uint64) bool {
	if (passes & active) == uint64(0) {
		return false
	}
	return true
}

func setActive(passes uint64, pass uint64) uint64 {
	return passes | pass
}

func unsetActive(passes uint64, pass uint64) uint64 {
	return passes & ^pass
}
