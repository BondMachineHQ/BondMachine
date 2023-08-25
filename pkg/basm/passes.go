package basm

import "errors"

const (
	passTemplateResolver      = uint64(1)
	passDynamicalInstructions = uint64(2)
	passSymbolTagger          = uint64(4)
	passDataSections2Bytes    = uint64(8)
	passMetadataInfer1        = uint64(16)
	passFragmentAnalyzer      = uint64(32)
	passFragmentOptimizer1    = uint64(64)
	passFragmentPruner        = uint64(128)
	passFragmentComposer      = uint64(256)
	passMetadataInfer2        = uint64(512)
	passEntryPoints           = uint64(1024)
	passRomComposer           = uint64(2048)
	passSymbolsResolver       = uint64(4096)
	passMatcherResolver       = uint64(8192)
	LAST_PASS                 = uint64(8192)
)

func getPassFunction() map[uint64]func(*BasmInstance) error {
	return map[uint64]func(*BasmInstance) error{
		passTemplateResolver:      templateResolver,
		passDynamicalInstructions: dynamicalInstructions,
		passSymbolTagger:          symbolTagger,
		passDataSections2Bytes:    dataSections2Bytes,
		passMetadataInfer1:        metadataInfer,
		passMetadataInfer2:        metadataInfer,
		passEntryPoints:           entryPoints,
		passSymbolsResolver:       symbolResolver,
		passMatcherResolver:       matcherResolver,
		passFragmentAnalyzer:      fragmentAnalyzer,
		passFragmentPruner:        fragmentPruner,
		passFragmentComposer:      fragmentComposer,
		passFragmentOptimizer1:    fragmentOptimizer,
		passRomComposer:           romComposer,
	}
}

func getPassFunctionName() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateResolver",
		passDynamicalInstructions: "dynamicalInstructions",
		passSymbolTagger:          "symbolTagger",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadataInfer (1)",
		passMetadataInfer2:        "metadataInfer (2)",
		passEntryPoints:           "entryPoints",
		passSymbolsResolver:       "symbolResolver",
		passMatcherResolver:       "matcherResolver",
		passFragmentAnalyzer:      "fragmentAnalyzer",
		passFragmentPruner:        "fragmentPruner",
		passFragmentComposer:      "fragmentComposer",
		passFragmentOptimizer1:    "fragmentOptimizer",
		passRomComposer:           "romComposer",
	}
}

func IsOptionalPass() map[uint64]bool {
	return map[uint64]bool{
		passTemplateResolver:      false,
		passDynamicalInstructions: false,
		passSymbolTagger:          false,
		passDataSections2Bytes:    false,
		passMetadataInfer1:        false,
		passMetadataInfer2:        false,
		passEntryPoints:           false,
		passSymbolsResolver:       false,
		passMatcherResolver:       false,
		passFragmentAnalyzer:      false,
		passFragmentPruner:        false,
		passFragmentComposer:      false,
		passFragmentOptimizer1:    true,
		passRomComposer:           false,
	}
}

func (bi *BasmInstance) ActivePass(active uint64) bool {
	return (bi.passes & active) != uint64(0)
}

func GetPassMnemonic() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateresolver",
		passDynamicalInstructions: "dynamicalinstructions",
		passSymbolTagger:          "symboltagger",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadatainfer1",
		passMetadataInfer2:        "metadatainfer2",
		passEntryPoints:           "entrypoints",
		passSymbolsResolver:       "symbolresolver",
		passMatcherResolver:       "matcherresolver",
		passFragmentAnalyzer:      "fragmentanalyzer",
		passFragmentPruner:        "fragmentpruner",
		passFragmentComposer:      "fragmentcomposer",
		passFragmentOptimizer1:    "fragmentoptimizer",
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
