
## Phases ##

The assembling process in divided into phases. The first phase (phase 0) read assembly file/s and create a raw BasmInstance. Subsequent passes operate transformations on BasmInstances. Specific steps can be removed using the dedicated arguments from the command line. The table following shows all the phases along with their description.

| Phase short name        | Description  |
| :-------------| :-----|
| TemplateResolver | The templated code is found and expanded. New (untemplated) elements (section or fragment) are created |
| DynamicalInstructions | The dynamical instructions are found according the name convention. They are created and inserted into the instruction database |
| SymbolsTagger1 | Map the sections and fragments symbols, creates the relative metadata within the instruction arguments |
| DataSections2Bytes | Compute the offsets of the data sections and convert the data into bytes |
| MetadataInfer1 | Infer the metadata by looking at the code and matching the instructions with the instruction database |
| FragmentAnalyzer | Analyze the fragments resources and create the relative metadata |
| FragmentOptimizer | Apply several customizable optimizations to the fragments |
| FragmentPruner | Prune the fragments that are specified in the command line |
| FragmentComposer | Compose the fragments into the sections as specified in the command line |
| MetadataInfer2 | Infer the metadata for the second time since news sections and fragments may have been created |
| EntryPoints | The programs entry points is detected for the sections where it is relevant |
| MatcherResolver | Resolv the pseudo-insructions and traslate the instructions into the real ones. If more than one instruction is matched, alternative sections are created to be evaluated in the next phases |
| SymbolsTagger2 | Map the sections and fragments symbols, creates the relative metadata within the instruction arguments |
| MemComposer | Associate the memory to the sections according to the final disposition of the sections within the BondMachine. Only section relevant for the cps metadata are considered and the others are discarded |
| SectionCleaner | Remove the sections that are not relevant for the cps metadata |
| SymbolsTagger3 | Map the sections and fragments symbols, creates the relative metadata within the instruction arguments |
| SymbolsResolver | Symbols are detected, removed from the actual code and written as locations |

After the last phase the BasmInstance is ready to be translated into a BondMachine or to a BCOF file. The structure of the BondMachine is defined by the SOs and CPs metas. While the code is processed, the assembler keeps track of the SOs and CPs that are used. At the end of the process, the assembler creates the BondMachine structure and fills it with the data collected during the process.
