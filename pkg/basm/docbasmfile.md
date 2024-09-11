
## __Directives__ ##

A .basm input file is made of directives starting with the % character. Some directive starts a block that ends with the respective %end directive. Others are  single line directives. The table below lists the directives that can be used in a .basm file.

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

- __romtext__: Read only code section
- __romdata__: Read only data section
- __ramtext__: Code section
- __ramdata__: Data section

The section name is a valid identifier. It has to be unique in the current BasmInstance. The section type is one of the four listed above. The section metadata is a list of metadata that are associated to the section. The metadata are used to define the section properties. The list of metadata is given in the table below.

| Metadata     | Description  |
| -------------:|:-------------|
| __iomode__ | The defualt I/O mode for the section. It can be __sync__ or __async__. The default is __async__. |

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
