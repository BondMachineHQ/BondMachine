package basm

import "errors"

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

func IsOptionalPass() map[uint64]bool {
	return map[uint64]bool{
		passTemplateResolver:      false,
		passDynamicalInstructions: false,
		passLabelTagger:           false,
		passDataSections2Bytes:    false,
		passMetadataInfer1:        false,
		passMetadataInfer2:        false,
		passEntryPoints:           false,
		passLabelsResolver:        false,
		passMatcherResolver:       false,
		passFragmentAnalyzer:      false,
		passFragmentPruner:        false,
		passFragmentComposer:      false,
		passRomComposer:           false,
	}
}

func (bi *BasmInstance) activePass(active uint64) bool {
	return (bi.passes & active) != uint64(0)
}

func GetPassMnemonic() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateresolver",
		passDynamicalInstructions: "dynamicalinstructions",
		passLabelTagger:           "labeltagger",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadatainfer1",
		passMetadataInfer2:        "metadatainfer2",
		passEntryPoints:           "entrypoints",
		passLabelsResolver:        "labelresolver",
		passMatcherResolver:       "matcherresolver",
		passFragmentAnalyzer:      "fragmentanalyzer",
		passFragmentPruner:        "fragmentpruner",
		passFragmentComposer:      "fragmentcomposer",
		passRomComposer:           "romcomposer",
	}

}

func (bi *BasmInstance) SetActive(pass string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass()[passN]; ok {
				if ch {
					bi.passes = bi.passes | passN
					return nil
				} else {
					return errors.New("pass is not optional")
				}
			} else {
				return errors.New("pass is not defined")
			}
		}
	}
	return errors.New("pass not found")
}

func (bi *BasmInstance) UnsetActive(pass string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass()[passN]; ok {
				if ch {
					bi.passes = bi.passes & ^passN
					return nil
				} else {
					return errors.New("pass is not optional")
				}
			} else {
				return errors.New("pass is not defined")
			}
		}
	}
	return errors.New("pass not found")
}
