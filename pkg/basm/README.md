Basm is the BondMachine assembler. It is a tool that translates a BondMachine assembly code into a BondMachine cluster object file (BCOF) or into a BondMachine architecture file. Basm is a tool that is part of the BondMachine toolchain.

Basm is very different from traditional assemblers. In the default case, the target architecture of a basm instance does not yet exists when it is running. The creation of the architecture is part of the assembling process. The architecture is created by the assembler itself by looking at the code and at the metadatas. The assembler creates the architecture and fills it with the data collected during the process.
 This is a novel concept of the BondMachine.

One or more .basm files contain the code to be compiled in an existing BM or to produce a suitable one. Basm files not only contain the code but also the metadatas that are used to define the architecture. The metadatas are used to define the SOs, the CPs and the BM itself. The metadatas are defined by using the %meta directives. 

The page [Basm File Structure](docbasmfile.md) contains a detailed description of the structure of a .basm file. The page [Basm Instructions](docinstructions.md) contains a detailed description how to write instructions and I/O objects in a .basm file. The page [Basm Internals](docinternals.md) contains information about the internals of the assembler. The single instructions documentation is available in the [BASM Assembly Reference](reference/).

## __Usage__ ##

```bash
basm [options] [input files]
```

where options are:
```bash
Usage of basm:
  -activate-optimizations string
        List of comma separated optional optimizations to activate (default: none, everything: all)
  -activate-passes string
        List of comma separated optional passes to activate (default: none)
  -bminfo-file string
        Load additional information about the BondMachine
  -bo string
        BCOF Output file
  -bondmachine string
        Load a bondmachine JSON file
  -chooser-force-same-name
        Force the chooser to use the same name for the ROM and the RAM
  -chooser-min-word-size
        Choose the minimum word size for the chooser
  -d    Verbose
  -deactivate-passes string
        List of comma separated optional passes to deactivate (default: none)
  -disable-dynamical-matching
        Disable the dynamical matching
  -dump-requirements string
        Dump the requirements of the BondMachine in a JSON file
  -getmeta string
        Get the metadata of an internal parameter of the BondMachine
  -linear-data-range string
        Load a linear data range file (with the syntax index,filename)
  -list-optimizations
        List the available optimizations
  -list-passes
        List the available passes
  -o string
        BondMachine Output file
  -si string
        Load a symbols JSON file
  -so string
        Symbols Output file
  -v    Verbose
```
