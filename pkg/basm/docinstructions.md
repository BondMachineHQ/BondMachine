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

-