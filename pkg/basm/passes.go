package basm

import "errors"

const (
	passTemplateResolver      = uint64(1)
	passMetadataInfer1        = uint64(2)
	passMacroResolver         = uint64(4)
	passCallResolver          = uint64(8)
	passDynamicalInstructions = uint64(16)
	passSymbolTagger1         = uint64(32)
	passDataSections2Bytes    = uint64(64)
	passMetadataInfer2        = uint64(128)
	passFragmentAnalyzer      = uint64(256)
	passFragmentOptimizer1    = uint64(512)
	passFragmentPruner        = uint64(1024)
	passFragmentComposer      = uint64(2048)
	passMetadataInfer3        = uint64(4096)
	passEntryPoints           = uint64(8192)
	passMatcherResolver       = uint64(16384)
	passSymbolTagger2         = uint64(32768)
	passMemComposer           = uint64(65536)
	passSectionCleaner        = uint64(131072)
	passSymbolTagger3         = uint64(262144)
	passSymbolsResolver       = uint64(524288)
	LAST_PASS                 = uint64(524288)
)

func getPassFunction() map[uint64]func(*BasmInstance) error {
	return map[uint64]func(*BasmInstance) error{
		passTemplateResolver:      templateResolver,
		passDynamicalInstructions: dynamicalInstructions,
		passSymbolTagger1:         symbolTagger,
		passSymbolTagger2:         symbolTagger,
		passSymbolTagger3:         symbolTagger,
		passDataSections2Bytes:    dataSections2Bytes,
		passMetadataInfer1:        metadataInfer,
		passMetadataInfer2:        metadataInfer,
		passMetadataInfer3:        metadataInfer,
		passEntryPoints:           entryPoints,
		passSymbolsResolver:       symbolResolver,
		passMatcherResolver:       matcherResolver,
		passFragmentAnalyzer:      fragmentAnalyzer,
		passFragmentPruner:        fragmentPruner,
		passFragmentComposer:      fragmentComposer,
		passFragmentOptimizer1:    fragmentOptimizer,
		passMemComposer:           memComposer,
		passSectionCleaner:        sectionCleaner,
		passCallResolver:          callResolver,
		passMacroResolver:         macroResolver,
	}
}

func getPassFunctionName() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateResolver",
		passDynamicalInstructions: "dynamicalInstructions",
		passSymbolTagger1:         "symbolTagger (1)",
		passSymbolTagger2:         "symbolTagger (2)",
		passSymbolTagger3:         "symbolTagger (3)",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadataInfer (1)",
		passMetadataInfer2:        "metadataInfer (2)",
		passMetadataInfer3:        "metadataInfer (3)",
		passEntryPoints:           "entryPoints",
		passSymbolsResolver:       "symbolResolver",
		passMatcherResolver:       "matcherResolver",
		passFragmentAnalyzer:      "fragmentAnalyzer",
		passFragmentPruner:        "fragmentPruner",
		passFragmentComposer:      "fragmentComposer",
		passFragmentOptimizer1:    "fragmentOptimizer",
		passMemComposer:           "memComposer",
		passSectionCleaner:        "sectionCleaner",
		passCallResolver:          "callResolver",
		passMacroResolver:         "macroResolver",
	}
}

func IsOptionalPass() map[uint64]bool {
	return map[uint64]bool{
		passTemplateResolver:      false,
		passDynamicalInstructions: false,
		passSymbolTagger1:         false,
		passSymbolTagger2:         false,
		passSymbolTagger3:         false,
		passDataSections2Bytes:    false,
		passMetadataInfer1:        false,
		passMetadataInfer2:        false,
		passMetadataInfer3:        false,
		passEntryPoints:           false,
		passSymbolsResolver:       false,
		passMatcherResolver:       false,
		passFragmentAnalyzer:      false,
		passFragmentPruner:        false,
		passFragmentComposer:      false,
		passFragmentOptimizer1:    true,
		passMemComposer:           false,
		passSectionCleaner:        false,
		passCallResolver:          false,
		passMacroResolver:         false,
	}
}

func (bi *BasmInstance) ActivePass(active uint64) bool {
	return (bi.passes & active) != uint64(0)
}

func GetPassMnemonic() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateresolver",
		passDynamicalInstructions: "dynamicalinstructions",
		passSymbolTagger1:         "symboltagger1",
		passSymbolTagger2:         "symboltagger2",
		passSymbolTagger3:         "symboltagger3",
		passDataSections2Bytes:    "datasections2bytes",
		passMetadataInfer1:        "metadatainfer1",
		passMetadataInfer2:        "metadatainfer2",
		passMetadataInfer3:        "metadatainfer3",
		passEntryPoints:           "entrypoints",
		passSymbolsResolver:       "symbolresolver",
		passMatcherResolver:       "matcherresolver",
		passFragmentAnalyzer:      "fragmentanalyzer",
		passFragmentPruner:        "fragmentpruner",
		passFragmentComposer:      "fragmentcomposer",
		passFragmentOptimizer1:    "fragmentoptimizer",
		passMemComposer:           "memcomposer",
		passSectionCleaner:        "sectioncleaner",
		passCallResolver:          "callresolver",
		passMacroResolver:         "macroresolver",
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
