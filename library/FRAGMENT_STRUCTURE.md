# BASM Fragment Structure

## Overview
BASM fragments are reusable code blocks with configurable parameters and metadata, enclosed between `%fragment` and `%endfragment` directives.

## Basic Structure

```
%fragment <name> [attributes...]
    <assembly_code>
%endfragment
```

attributes are key-multi-value pairs that define fragment behavior and parameters. They are separated by spaces and formatted as `key:value1:value2:...`.


## Fragment Attributes

### Core Attributes

- **`template:true`** - Marks fragment as a template with parameter substitution
- **`resin:<reg_list>`** - Input registers (colon-separated: `r0:r1:r2`)
- **`resout:<reg_list>`** - Output registers (colon-separated: `r0:r1`)

- **`interface:<interface_style>`** - Defines fragment interface style (e.g., `regbased (default)`, `stackbased`)
- **`stack:<name>`** - Specifies stack name for stack-based fragments, the name also contains stack parameters like size and type.

### Default Parameters

Defined as `default_<param_name>:<value>`:

- **`default_mem:<value>`** - Default memory type (e.g., `rom`, `ram`)
- **`default_mem1:<value>`** - First memory parameter
- **`default_mem2:<value>`** - Second memory parameter
- **`default_out:<value>`** - Output device (e.g., `vtm0`)
- **`default_kbd:<value>`** - Input device (e.g., `i0`)
- **`default_feedout:<value>`** - Feedback output
- **`default_feedreg:<value>`** - Feedback register
- **`default_height:<value>`** - Height parameter
- **`default_width:<value>`** - Width parameter
- **`default_outer:<value>`** - Outer loop register
- **`default_inner:<value>`** - Inner loop register
- **`default_temp:<value>`** - Temporary register
- **`default_setop:<value>`** - Set operation (e.g., `rset`)
- **`default_prefix:<value>`** - Value prefix (e.g., `0f` for float)
- **`default_<precision>:<value>`** - Precision parameter (e.g., `default_sinprec:20`)
- **Custom parameters** - Any parameter specific to fragment logic

## Requirements Directives

some attributes specify requirements for the CP that will use the fragment:

- **`require_llmemname:memname1:memname2:...`** - Requires blocks of `local` linear memory with specified names (the number of names and sizes must match those in `require_llmemsize`)
- **`require_llmemsize:size1:size2:...`** - For every memory name in `require_llmemname`, specifies the required size in bytes
- **`require_clmemname:memname1:memname2:...`** - Requires blocks of linear memory in the common linear memory space of the target CP with specified names
- **`require_clmemsize:size1:size2:...`** - For every memory name in `require_clmemname`, specifies the required size in bytes.

When different fragments have the same memory name in their clmemname they need to be on the sale CP and have the same size. 
TODO: Think about shared memory across different CPs and concurrency.

## Meta Directives

### %meta literal
Defines literal metadata that is evaluated during template processing:

```
%meta literal resin {{template_expression}}
```

Used to dynamically generate register lists or other metadata.

## Parameter Substitution

Inside template fragments, parameters are accessed via Go template syntax:

```
{{ .Params.<param_name> }}
```

Examples:
- `{{ .Params.mem }}` - Access memory parameter
- `{{ .Params.addop }}` - Access operation parameter
- `{{ .Params.prefix }}` - Access prefix parameter

## Template Expressions

Common template patterns:

- **Range iteration**: `{{range $y := intRange "1" $last}}...{{end}}`
- **Conditional**: `{{with $last := incs .Params.inputs}}...{{end}}`
- **String formatting**: `{{printf "r%d:" $y}}`

## Comment Annotations

### Testing Annotations
```
;fragtester <directives>
;nofragtester <directives>
;sympy <python_code>
```

Used for automated testing and symbolic validation.

### Documentation
Standard assembly comments with `;` describing:
- Input/output register usage
- Parameter descriptions
- Algorithm explanations

## Section Fragments

Alternative form using `%section`:

```
%section <name> <section_type> [attributes...]
    <assembly_code>
%endsection
```

Where `<section_type>` can be `.romtext` or other section types.

## Common Patterns

### Simple Fragment
```
%fragment name default_param:value
    <code>
%endfragment
```

### Template Fragment
```
%fragment name template:true resin:r0 resout:r0 default_op:add
    {{ .Params.op }} r0, r1
%endfragment
```

### Multi-Register I/O
```
%fragment name template:true resin:r0:r1:r2:r3 resout:r0:r1
    <code>
%endfragment
```

### Dynamic Metadata
```
%fragment name template:true resout:r0
%meta literal resin {{range...}}
    <code>
%endfragment
```
