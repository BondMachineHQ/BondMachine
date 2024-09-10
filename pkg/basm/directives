
## __Directives__ ##

A .basm input file is made of directives starting with the % character. Some directive starts a block that ends with the respective %end directive. Others are single line directives.

| Directive     | End directive  | Description  |
| -------------:|:-------------| :-----|
| %section      | %endsection | Code section |
| %macro        | %endmacro |  Macro definition|
| %fragment     | %endfragment | Code fragment |
| %chunk        | %endchunk | Chunk of data |
| %meta bmdef   | none | BondMachine metadatas |
| %meta cpdef   | none | CP metadatas |
| %meta sodef   | none | SO metadatas |
| %meta iodef   | none | Input/Output metadatas |
| %meta fidef   | none | Fragment Instance metadata |
| %meta filinkdef   | none | Fragment Instance Link metadatas |
| %meta soatt   | none | SO attach directives |
| %meta ioatt   | none | IO attach directives |
| %meta filinkatt   | none | Fragment Instance link attach directives |

## __Sections__ ##

The main Basm structure is the section. A section is a block of code or data.
The directive %section is used start a section  and ends with the directive %endsection. The syntax is:

    %section section_name section_type [section metadata]
    ...
    %endsection

Within Basm four type of sections can be defined:

  - **romtext**: Read only code section
  - **romdata**: Read only data section
  - **ramtext**: Code section
  - **ramdata**: Data section

The section name is a valid identifier. It has to be unique in the current BasmInstance. The section type is one of the four listed above. The section metadata is a list of metadata that are associated to the section. The metadata are used to define the section properties. The list of metadata is given in the table below.

| Metadata     | Description  |
| -------------:|:-------------|
| **iomode** | The defualt I/O mode for the section. It can be **sync** or **async**. The default is **async**. |

### Code sections ###

### Data sections ###



TODO

## __Macros__ ##

TODO

## __Fragments__ ##

TODO

## __Chunks__ ##

TODO

## __%meta Directives__ ##
The %meta directives are used to define the BondMachine architecture. They are used to define the SOs, the CPs and the BM itself. The syntax is:

TODO
## __Instructions__ ##

Basm instructions are identified by an operation that requires operands to be processed. Operands can be registers, memory address or more generally any SOs I/O ports.
Some instructions do not require an operand, whereas some other instructions may require one, two, or three operands.

The first operand is generally the destination and the second operand is the source.

## Pseudo instructions ##

Pseudo instructions are not real instructions, they are directives that are translated into one or more real instructions. They are used in the instruction field anyway because that’s the most convenient place to put them. The most common pseudo instructions are the ones that define data and reserve storage space. They are used to define constants and variables and to reserve space for arrays and other data structures.

### Declaring data ###

The declaration of data is done in data sections by using the syntax:

    symbol directive expression [, expression]...

**symbol** is the name of the variable. It can be any valid identifier. It as to be unique in the current section.

**directive** is the name of the directive that declares the data. The list of directives is given in the table below.

 **expression** is a constant expression that evaluates to a byte (or other types according to the directives) value. The assembler assigns the first expression to the first byte of the variable, the second expression to the second byte, and so on. The assembler automatically advances the location counter by the number of bytes declared.

The following table lists the directives that can be used to declare data:

| Directive     | Type | Description  | Example |
| -------------:|:-------------|--|--|
| db            | Initiatialized data | Declare byte | examplevar db 1, 2, 3 |



## Expressing I/O ##

Apart from operations among internal elements, CP instructions can interact with external objects. These can be memories, like in standard processor, or more generally can be all kind of SOs. Basm has a special syntax to express such data movements. The most common instruction is _**mov**_ and  the general syntax for the I/O objects is:

    [IO object name][object index]:[I/O port]

If the object has no port the colon can be omitted. Similarly, if the object has no index, just the name can be used alone.

Following is a list of examples of I/O objects:

| Object name        | Description  |
| :-------------| :-----|
| *ram* | The random access memory |
| *rom* | The read only memory |
| *i* | The input |
| *o* | The output |
| *vtm* | The video textual memory |
| *q* | The shared queue |
| *s* | The shared stack |

The I/O port can be addressed in two ways:

- *Immediate addressing*: An immediate port has a constant value or an expression.

- *Register addressing*: In this addressing mode, a register contains the port. In this mode the register name is surrounded by brackets.

- *Symbol addressing*: In this addressing mode, a symbol refers to a port. 


Here is some examples along with theirs descriptions:

| Example        | Description  |
| :-------------| :-----|
| *mov r0 5* | The register _**r0**_ is filled with the decimal value of 5|
| *mov r1 i0* or *mov r1 i0:*| The register _**r0**_ is filled with the value read from the input _**i0**_
| *mov vtm0:48 r0* | Copy the value contained in the _**r0**_ register into the location 48 of the video textual memory using immediate addressing |
| *mov vtm0:[r1] r0* | Copy the value contained in the _**r0**_ register into the location of the video textual memory pointed by the _**r1**_ register using register addressing |

## Contexts ##

    %macro macro_name num_of_params
    %endmacro

Macros are not CP-bounded, they are assembler bounded

---

Every section takes a CP tag, if that is missing it means there is only 1 tag associated to CP 0

    sections .romtext cp_tag
         .romdata cp_tag
         .text cp_tag
         .data cp_tag
         .bss cp_tag
---

## Labels are CP-bounded ##

metas are directives specific of the BondMachine, they are used to create groups of processors or groups of groups (clusters)

    %meta sodef [so name] [so definition] [so tags]

    %meta cpdef [cp name] [cp tags]

    %meta bmdef 

---

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

---

## To Cut and Paste ##

Most assembly language instructions require operands to be processed. An operand address provides the location, where the data to be processed is stored. Some instructions do not require an operand, whereas some other instructions may require one, two, or three operands.

The three basic modes of addressing are −

    Register addressing
    Immediate addressing
    Memory addressing

MOV DX, TAX_RATE   ; Register in first operand
MOV COUNT, CX   ; Register in second operand
MOV EAX, EBX   ; Both the operands are in registers

As processing data between registers does not involve memory, it provides fastest processing of data.

An immediate operand has a constant value or an expression. When an instruction with two operands uses immediate addressing, the first operand may be a register or memory location, and the second operand is an immediate constant. The first operand defines the length of the data.

For example,

BYTE_VALUE  DB  150    ; A byte value is defined
WORD_VALUE  DW  300    ; A word value is defined
ADD  BYTE_VALUE, 65    ; An immediate operand 65 is added
MOV  AX, 45H           ; Immediate constant 45H is transferred to AX

Direct Memory Addressing

When operands are specified in memory addressing mode, direct access to main memory, usually to the data segment, is required. This way of addressing results in slower processing of data. To locate the exact location of data in memory, we need the segment start address, which is typically found in the DS register and an offset value. This offset value is also called effective address.

In direct addressing mode, the offset value is specified directly as part of the instruction, usually indicated by the variable name. The assembler calculates the offset value and maintains a symbol table, which stores the offset values of all the variables used in the program.

In direct memory addressing, one of the operands refers to a memory location and the other operand references a register.

For example,

ADD BYTE_VALUE, DL; Adds the register in the memory location
MOV BX, WORD_VALUE; Operation from the memory is added to register

Direct-Offset Addressing

This addressing mode uses the arithmetic operators to modify an address. For example, look at the following definitions that define tables of data −

BYTE_TABLE DB  14, 15, 22, 45      ; Tables of bytes
WORD_TABLE DW  134, 345, 564, 123  ; Tables of words

The following operations access data from the tables in the memory into registers −

MOV CL, BYTE_TABLE[2] ; Gets the 3rd element of the BYTE_TABLE
MOV CL, BYTE_TABLE + 2 ; Gets the 3rd element of the BYTE_TABLE
MOV CX, WORD_TABLE[3] ; Gets the 4th element of the WORD_TABLE
MOV CX, WORD_TABLE + 3 ; Gets the 4th element of the WORD_TABLE

Indirect Memory Addressing

This addressing mode utilizes the computer's ability of Segment:Offset addressing. Generally, the base registers EBX, EBP (or BX, BP) and the index registers (DI, SI), coded within square brackets for memory references, are used for this purpose.

Indirect addressing is generally used for variables containing several elements like, arrays. Starting address of the array is stored in, say, the EBX register.

The following code snippet shows how to access different elements of the variable.

MY_TABLE TIMES 10 DW 0  ; Allocates 10 words (2 bytes) each initialized to 0
MOV EBX, [MY_TABLE]     ; Effective Address of MY_TABLE in EBX
MOV [EBX], 110          ; MY_TABLE[0] = 110
ADD EBX, 2              ; EBX = EBX +2
MOV [EBX], 123          ; MY_TABLE[1] = 123

The MOV Instruction

We have already used the MOV instruction that is used for moving data from one storage space to another. The MOV instruction takes two operands.
Syntax

The syntax of the MOV instruction is −

MOV  destination, source

The MOV instruction may have one of the following five forms −

MOV  register, register
MOV  register, immediate
MOV  memory, immediate
MOV  register, memory
MOV  memory, register

Please note that −

    Both the operands in MOV operation should be of same size
    The value of source operand remains unchanged

The MOV instruction causes ambiguity at times. For example, look at the statements −

MOV  EBX, [MY_TABLE]  ; Effective Address of MY_TABLE in EBX
MOV  [EBX], 110      ; MY_TABLE[0] = 110

It is not clear whether you want to move a byte equivalent or word equivalent of the number 110. In such cases, it is wise to use a type specifier.

Following table shows some of the common type specifiers −
Type Specifier Bytes addressed
BYTE 1
WORD 2
DWORD 4
QWORD 8
TBYTE 10
