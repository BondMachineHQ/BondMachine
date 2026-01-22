package basm

import "errors"

const (
	passClusterChecker = uint64(1) << iota
	passTemplateResolver
	passDependencyResolver
	passMetadataInfer1
	passMacroResolver
	passCallResolver
	passDynamicalInstructions
	passSymbolTagger1
	passDataSections2Bytes
	passMetadataInfer2
	passFragmentAnalyzer
	passFragmentOptimizer1
	passFragmentPruner
	passFragmentComposer
	passMetadataInfer3
	passEntryPoints
	passTemplateFinalizer
	passMatcherResolver
	passSymbolTagger2
	passMemComposer
	passSectionCleaner
	passSymbolTagger3
	passSymbolsResolver
	LAST_PASS = passSymbolsResolver
)

func getPassFunction() map[uint64]func(*BasmInstance) error {
	return map[uint64]func(*BasmInstance) error{
		passTemplateResolver:      templateResolver,
		passDependencyResolver:    dependencyResolver,
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
		passTemplateFinalizer:     templateFinalizer,
		passClusterChecker:        clusterChecker,
	}
}

func getPassFunctionName() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateResolver",
		passDependencyResolver:    "dependencyResolver",
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
		passTemplateFinalizer:     "templateFinalizer",
		passClusterChecker:        "clusterChecker",
	}
}

func IsOptionalPass(frontEnd string) map[uint64]bool {
	if frontEnd == "bondbits" {
		return map[uint64]bool{
			passTemplateResolver:      true,
			passDependencyResolver:    true,
			passDynamicalInstructions: true,
			passSymbolTagger1:         true,
			passSymbolTagger2:         true,
			passSymbolTagger3:         true,
			passDataSections2Bytes:    true,
			passMetadataInfer1:        true,
			passMetadataInfer2:        true,
			passMetadataInfer3:        true,
			passEntryPoints:           true,
			passSymbolsResolver:       true,
			passMatcherResolver:       true,
			passFragmentAnalyzer:      true,
			passFragmentPruner:        true,
			passFragmentComposer:      true,
			passFragmentOptimizer1:    true,
			passMemComposer:           true,
			passSectionCleaner:        true,
			passCallResolver:          true,
			passMacroResolver:         true,
			passTemplateFinalizer:     true,
			passClusterChecker:        true,
		}
	} else {
		return map[uint64]bool{
			passTemplateResolver:      false,
			passDependencyResolver:    false,
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
			passTemplateFinalizer:     true,
			passClusterChecker:        false,
		}
	}
}

func (bi *BasmInstance) ActivePass(active uint64) bool {
	return (bi.passes & active) != uint64(0)
}

func GetPassMnemonic() map[uint64]string {
	return map[uint64]string{
		passTemplateResolver:      "templateresolver",
		passDependencyResolver:    "dependencyresolver",
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
		passTemplateFinalizer:     "templatefinalizer",
		passClusterChecker:        "clusterchecker",
	}

}

func (bi *BasmInstance) SetActive(pass string, frontEnd string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass(frontEnd)[passN]; ok {
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

func (bi *BasmInstance) UnsetActive(pass string, frontend string) error {
	for passN, v := range GetPassMnemonic() {
		if v == pass {
			if ch, ok := IsOptionalPass(frontend)[passN]; ok {
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
