Basm is the BondMachine assembler. It is a tool that translates a BondMachine assembly code into a BondMachine cluster object file (BCOF) or into a BondMachine architecture file. Basm is a tool that is part of the BondMachine toolchain.

Basm is very different from traditional assemblers. In the default case, the target architecture of a basm instance does not yet exists when it is running. The creation of the architecture is part of the assembling process. The architecture is created by the assembler itself by looking at the code and at the metadatas. The assembler creates the architecture and fills it with the data collected during the process.
 This is a novel concept of the BondMachine.

One or more .basm files contain the code to be compiled in an existing BM or to produce a suitable one. Basm files not only contain the code but also the metadatas that are used to define the architecture. The metadatas are used to define the SOs, the CPs and the BM itself. The metadatas are defined by using the %meta directives. 

## __Usage__ ##

    basm [options] [input files]

The options are:

[directives](directives)
