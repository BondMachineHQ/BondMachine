# BondMachine Assembly Fragment Specifications

This document provides a comprehensive guide to the BondMachine Assembly (BASM) fragment composition format and a catalog of all available fragments in the library.

## Table of Contents

### Part I: Fragment Composition Format
1. [Fragment Structure](#fragment-structure)
2. [Fragment Directives](#fragment-directives)
3. [Template Parameters](#template-parameters)
4. [Metadata Attributes](#metadata-attributes)
5. [Comments and Documentation](#comments-and-documentation)
6. [Instructions and Addressing Modes](#instructions-and-addressing-modes)
7. [Labels and Flow Control](#labels-and-flow-control)
8. [Register Conventions](#register-conventions)

### Part II: Fragment Catalog
1. [comparestring](#comparestring---string-comparison)
2. [newline](#newline---new-line-with-scrolling)
3. [putline](#putline---write-string-with-counter)
4. [putstring](#putstring---write-string-to-output)
5. [read](#read---read-string-from-keyboard)
6. [readline](#readline---read-line-from-keyboard)
7. [readstring](#readstring---read-string-with-clear)
8. [readword](#readword---read-word-from-keyboard)
9. [sleep](#sleep---delay-function)
10. [sync](#sync---synchronize-memory-to-output)

---

## Part I: Fragment Composition Format

### Fragment Structure

A BASM fragment is a reusable block of parameterized assembly code defined using the following structure:

```assembly
; Comments describing the fragment
; Input/output specifications
%fragment <name> [default_param1:value1] [default_param2:value2] ... [metadata]
label1:
    instruction arg1, arg2
    instruction arg3, arg4
label2:
    instruction arg5
    ...
%endfragment
```

**Components:**
- **Fragment Declaration** - `%fragment` directive with name and parameters
- **Fragment Body** - Assembly instructions and labels
- **Fragment Termination** - `%endfragment` directive

---

### Fragment Directives

#### %fragment Directive

**Syntax:**
```assembly
%fragment <fragment_name> [parameter:value ...]
```

**Components:**
1. **Fragment Name** (required) - Identifier for the fragment
   - Must be a valid identifier (alphanumeric + underscore)
   - Used when calling the fragment: `%call fragment_name`

2. **Parameters** (optional) - Space-separated `key:value` pairs
   - **Default Parameters** - `default_<name>:<value>`
   - **Metadata** - `resin:<inputs>`, `resout:<outputs>`, etc.

**Examples:**
```assembly
%fragment putstring default_mem:rom default_out:vtm0

%fragment comparestring default_mem1:ram default_mem2:rom

%fragment newline default_out:vtm0 default_height:16 default_width:16

%fragment sleep default_outer:r0 default_inner:r1 default_temp:r2

%fragment readline default_out:vtm0 default_kbd:i0 default_feedout:vtm0 default_feedreg:r4
```

#### %endfragment Directive

**Syntax:**
```assembly
%endfragment
```

Marks the end of a fragment definition. Every `%fragment` must have a corresponding `%endfragment`.

---

### Template Parameters

Fragments use Go template syntax for parameterization, allowing customization at instantiation time.

#### Template Syntax

**Basic Parameter Reference:**
```assembly
{{ .Params.parameter_name }}
```

**Usage in Instructions:**

1. **Memory Type Parameters:**
```assembly
mov r3, {{ .Params.mem1 }}:[r0]      ; Read from parameterized memory
mov {{ .Params.out }}:[r1], r2       ; Write to parameterized output
```

2. **Register Parameters:**
```assembly
jz  {{ .Params.outer }}, endouter    ; Jump if parameterized register is zero
dec {{ .Params.temp }}               ; Decrement parameterized register
```

3. **Constant Parameters:**
```assembly
mov r1, {{ .Params.height }}         ; Load parameterized constant
mov r3, {{ .Params.width }}          ; Load another constant
```

#### Default Parameters

Default parameters provide fallback values when the fragment is called without explicit parameter values.

**Syntax:**
```assembly
%fragment <name> default_<param>:<value>
```

**Parameter Resolution Order:**
1. Explicit parameters passed at call site
2. Default parameters from fragment definition
3. Error if required parameter is missing

**Example:**
```assembly
%fragment putstring default_mem:rom default_out:vtm0
    mov r2, {{ .Params.mem }}:[r0]    ; Uses 'rom' if not overridden
    mov {{ .Params.out }}:[r1], r2    ; Uses 'vtm0' if not overridden
%endfragment
```

**Calling with Custom Parameters:**
```assembly
%call putstring mem:ram out:vtm1      ; Override defaults
%call putstring                       ; Use defaults (rom, vtm0)
```

#### Parameter Types

1. **Memory Types:**
   - `rom` - Read-only memory
   - `ram` - Random access memory
   - `vtm0`, `vtm1`, ... - Video textual memory
   - Custom memory types

2. **I/O Types:**
   - `i0`, `i1`, ... - Input ports
   - `o0`, `o1`, ... - Output ports
   - `kbd` - Keyboard input
   - Custom I/O devices

3. **Register Names:**
   - `r0`-`rN` - General purpose registers
   - `t0`-`tN` - Temporary registers (used internally)

4. **Numeric Constants:**
   - Integer values (decimal): `16`, `256`
   - Floating point values: `3.14`

---

### Metadata Attributes

Metadata provides information about fragment behavior and resource usage.

#### Standard Metadata Attributes

1. **resin:<res1>:<res2>:...** - Input resources
   - Lists resources consumed by the fragment
   - Example: `resin:r0:r1:r2`

2. **resout:<res1>:<res2>:...** - Output resources
   - Lists resources produced by the fragment
   - Example: `resout:r2:r4`

3. **resused:<res1>:<res2>:...** - Used resources
   - Lists all resources modified (both input and internal)
   - Example: `resused:r0:r1:r2:r3:r4`

**Example with Metadata:**
```assembly
%fragment example default_mem:ram resin:r0:r1 resout:r2 resused:r0:r1:r2:r3
```

#### fragtester Metadata

Special comments for fragment testing and validation:

```assembly
;fragtester range <param> arange(<from>,<to>,<step>)
;fragtester range <param> <val1>,<val2>,<val3>,...
;fragtester instance <param> <val1>,<val2>,...
;sympy <mathematical_expression>
```

**Examples:**
```assembly
;fragtester range width 1,2,4,8,16
;fragtester range height arange(8,32,4)
;fragtester instance mem ram,rom
;sympy result = input1 + input2
```

---

### Comments and Documentation

#### Comment Syntax

```assembly
; This is a single-line comment
```

- Comments start with `;`
- Everything after `;` to end of line is ignored
- Used for documentation and metadata

#### Documentation Conventions

**Header Comments:**
```assembly
; Description of what the fragment does
; Input:  r0 - description
; Input:  r1 - description
; Output: r2 - description
%fragment name ...
```

**Inline Comments:**
```assembly
mov r1, {{ .Params.height }}    ; Load screen height
dec r1                           ; Adjust for 0-indexing
```

**Section Comments:**
```assembly
; Docs (GitHub copilot generated)
;
; Extended documentation explaining the fragment's
; purpose, algorithm, and usage patterns.
```

---

### Instructions and Addressing Modes

#### Common BASM Instructions

**Data Movement:**
- `mov dest, src` - Move data
- `rsets8 reg, value` - Set register to 8-bit value

**Arithmetic:**
- `inc reg` - Increment register
- `dec reg` - Decrement register
- `add dest, src` - Addition
- `sub dest, src` - Subtraction
- `mult dest, src` - Multiplication

**Comparison:**
- `cmpr reg1, reg2` - Compare registers (sets flags)
- `cmprlt reg1, reg2` - Compare less than

**Flow Control:**
- `j label` - Unconditional jump
- `jz reg, label` - Jump if zero
- `jcmp label` - Jump if comparison true
- `jmp label` - Jump (alternate syntax)
- `call label` - Call subroutine
- `nop` - No operation

**I/O:**
- `i2rw dest, port` - Input to register (word)
- `r2ow port, src` - Register to output (word)

#### Addressing Modes

1. **Direct Register:**
```assembly
mov r1, r2          ; Register to register
```

2. **Immediate Value:**
```assembly
mov r1, 42          ; Constant to register
rsets8 r0, 0        ; Set to immediate 8-bit value
```

3. **Memory Indirect (Register as Address):**
```assembly
mov r1, ram:[r0]    ; Read from RAM at address in r0
mov rom:[r2], r3    ; Write to ROM at address in r2
```

4. **Memory Direct with Template:**
```assembly
mov r1, {{ .Params.mem }}:[r0]         ; Parameterized memory type
mov {{ .Params.out }}:[r1], r2         ; Parameterized output
```

5. **I/O Port:**
```assembly
i2rw r1, i0         ; Read from input port 0
r2ow o0, r1         ; Write to output port 0
```

#### Memory Addressing Syntax

**Format:**
```assembly
<memory_type>:[<address_register>]
```

**Examples:**
```assembly
ram:[r0]            ; RAM at address in r0
rom:[r1]            ; ROM at address in r1
vtm0:[r2]           ; Video memory at address in r2
{{ .Params.mem }}:[r0]  ; Template parameter memory
```

---

### Labels and Flow Control

#### Label Definition

**Syntax:**
```assembly
label_name:
    instruction ...
```

**Rules:**
- Label names must be unique within a fragment
- Labels end with colon `:`
- Can be referenced by jump/branch instructions

**Examples:**
```assembly
start:
    mov r0, 0
loop:
    inc r0
    cmpr r0, r1
    jcmp loop
end:
    nop
```

#### Label References

**Jump to Label:**
```assembly
j label_name        ; Unconditional jump
jz r0, label_name   ; Conditional jump (if r0 == 0)
jcmp label_name     ; Jump if comparison flag set
call label_name     ; Call subroutine at label
```

#### Common Control Flow Patterns

**1. Simple Loop:**
```assembly
loop:
    ; loop body
    dec r0
    jz r0, end
    j loop
end:
    nop
```

**2. Conditional Execution:**
```assembly
    cmpr r0, r1
    jcmp equal
    ; not equal code
    j continue
equal:
    ; equal code
continue:
    nop
```

**3. Nested Loops:**
```assembly
outerloop:
    mov r2, {{ .Params.inner }}
innerloop:
    dec r2
    jz r2, endinner
    j innerloop
endinner:
    dec r1
    jz r1, endouter
    j outerloop
endouter:
    nop
```

---

### Register Conventions

#### Register Types

**General Purpose Registers:** `r0` through `rN`
- Available to all code
- Must be declared/allocated by the system
- Used for data and computation

**Temporary Registers:** `t0` through `tN`
- Used internally within fragments
- Resolved to actual registers during fragment composition
- Not visible outside the fragment

#### Register Usage Patterns

**Input Registers:**
- Typically `r0`, `r1`, `r2`
- Contain input parameters before fragment execution

**Output Registers:**
- Can be any register
- Contain results after fragment execution
- Should be documented in fragment header

**Working Registers:**
- Higher-numbered registers (`r3`-`r7`)
- Used for temporary storage and computation
- May be modified during fragment execution

#### Example Fragment with Register Convention

```assembly
; Input:  r0 - buffer pointer (ram)
; Input:  r1 - source pointer (rom)
; Input:  r2 - length
; Output: Flags set based on comparison
; Modified: r3, r4 (working registers)
%fragment comparestring default_mem1:ram default_mem2:rom
cmpchar:
    mov r3, {{ .Params.mem1 }}:[r0]    ; r3 = working register
    mov r4, {{ .Params.mem2 }}:[r1]    ; r4 = working register
    cmpr r3, r4
    dec r2
    jz r2, end
    inc r0
    inc r1
    jcmp cmpchar
end:
    nop
%endfragment
```

---

## Part II: Fragment Catalog

---

## comparestring - String Comparison

**Purpose:** Compare two strings, one stored in ROM and one in RAM, to determine if they are equal.

### Input Parameters
- `r0` - Pointer to the buffer beginning in RAM (e.g., `mov r0, ram:buff`)
- `r1` - Pointer to the ROM string
- `r2` - Length of the string to compare

### Output
- Comparison result via flags (equal/not equal)

### Default Parameters
- `default_mem1: ram` - First memory type (typically RAM)
- `default_mem2: rom` - Second memory type (typically ROM)

### Operation
The fragment compares each character of two strings by:
1. Loading current character from RAM into `r3`
2. Loading current character from ROM into `r4`
3. Comparing the characters using `cmpr`
4. If characters don't match, exits the loop
5. Continues until all characters are compared or a mismatch is found

### Registers Modified
- `r3` - Temporary for RAM character
- `r4` - Temporary for ROM character

---

## newline - New Line with Scrolling

**Purpose:** Move to a new line in a video textual memory, with automatic scrolling when at the bottom of the screen.

### Input Parameters
- `r0` - Pointer to the memory start (0 for vtm0)
- `r2` - Current x position within the memory

### Output
- `r2` - Updated to new line position

### Default Parameters
- `default_out: vtm0` - Output device (video textual memory)
- `default_height: 16` - Screen height in characters
- `default_width: 16` - Screen width in characters

### Operation
The fragment performs:
1. Checks if scrolling is needed (current position >= last line)
2. If scrolling needed:
   - Copies each line up one position
   - Clears the last line
3. Finds the next line position
4. Fills any blank spaces with 0x00
5. Updates `r2` to point to the start of the new line

### Registers Modified
- `r1`, `r3`, `r4`, `r5`, `r6`, `r7` - Temporary registers for calculations and data movement

---

## putline - Write String with Counter

**Purpose:** Write a null-terminated string to the screen while maintaining a character counter.

### Input Parameters
- `r0` - Pointer to the null-terminated string within memory (e.g., `mov r0, rom:message`)
- `r2` - Starting x position within the output memory

### Output
- `r2` - Updated to position after the string
- `r4` - Number of characters written

### Default Parameters
- `default_mem: rom` - Source memory type
- `default_out: vtm0` - Output device (video textual memory)

### Operation
The fragment:
1. Reads character from source memory
2. Writes character to output memory
3. Increments pointers and counter
4. Continues until null terminator (0) is found

### Registers Modified
- `r1` - Temporary for character data
- `r4` - Character counter

---

## putstring - Write String to Output

**Purpose:** Write a null-terminated string to the screen (simplified version without counter).

### Input Parameters
- `r0` - Pointer to the null-terminated string (e.g., `mov r0, rom:message`)
- `r1` - Starting x position within the screen

### Output
- `r1` - Updated to position after the string

### Default Parameters
- `default_mem: rom` - Source memory type
- `default_out: vtm0` - Output device (video textual memory)

### Operation
The fragment:
1. Reads character from source memory
2. Checks for null terminator
3. Writes character to output memory
4. Increments pointers
5. Continues until null terminator is found

### Registers Modified
- `r2` - Temporary for character data

---

## read - Read String from Keyboard

**Purpose:** Read a string from the keyboard until Enter key is pressed, storing it in memory and displaying it on screen.

### Input Parameters
- `r0` - Pointer to the buffer beginning (e.g., `mov r0, ram:buff`)
- `r1` - Pointer to the screen position (e.g., `mov r1, vtm0:pos`)

### Output
- `r2` - Length of the string read

### Default Parameters
- `default_mem: ram` - Memory type for storage
- `default_out: vtm0` - Output device for display

### Operation
The fragment:
1. Initializes length counter to 0
2. Reads character from input (i0)
3. Checks for Enter key (ASCII 13)
4. Writes character to both screen and memory buffer
5. Increments pointers and length counter
6. Continues until Enter is pressed

### Registers Modified
- `r2` - Length counter
- `r3` - Enter key constant (13)
- `r4` - Temporary for character data

---

## readline - Read Line from Keyboard

**Purpose:** Read a line from the keyboard with configurable keyboard input and optional echo output.

### Input Parameters
- `r0` - Pointer to the buffer beginning (e.g., `mov r0, ram:buff`)
- `r2` - Pointer to the start of the memory where string will be stored

### Output
- `r0` - Updated to point after the stored string
- `r2` - Updated to point after the output position

### Default Parameters
- `default_out: vtm0` - Primary output device
- `default_kbd: i0` - Keyboard input device
- `default_feedout: vtm0` - Echo output device
- `default_feedreg: r4` - Register for feed tracking

### Operation
The fragment:
1. Reads character from keyboard
2. Writes to output memory
3. Writes to RAM buffer
4. Echoes to feed output
5. Increments all pointers
6. Continues until Enter key (ASCII 13)
7. Null-terminates the string in RAM

### Registers Modified
- `r1` - Temporary for character data
- `r3` - Enter key constant (13)
- `r4` (or feedreg) - Feed position tracker

---

## readstring - Read String with Clear

**Purpose:** Read a string from keyboard and clear the display after reading.

### Input Parameters
- `r0` - Pointer to the buffer beginning (e.g., `mov r0, ram:buff`)

### Output
- `r2` - Length of the string read

### Default Parameters
- None explicitly defined (uses hardcoded `vtm0` and `ram`)

### Operation
The fragment:
1. Initializes length counter to 0
2. Reads characters until Enter (ASCII 13)
3. Displays each character on vtm0
4. Stores characters in RAM buffer
5. After reading, clears the display by writing 0x00 to all positions

### Registers Modified
- `r1` - Temporary for character data
- `r2` - Length counter / position tracker
- `r3` - Enter key constant (13) / clear character (0)

---

## readword - Read Word from Keyboard

**Purpose:** Read a word from keyboard (similar to readstring, reads until Enter and clears display).

### Input Parameters
- `r0` - Pointer to the buffer beginning (e.g., `mov r0, ram:buff`)

### Output
- `r2` - Length of the word read

### Default Parameters
- None explicitly defined (uses hardcoded `vtm0` and `ram`)

### Operation
The fragment:
1. Initializes length counter to 0
2. Reads characters until Enter (ASCII 13)
3. Displays each character on vtm0
4. Stores characters in RAM buffer
5. After reading, clears the display by writing 0x00 to all positions

### Registers Modified
- `r1` - Temporary for character data
- `r2` - Length counter / position tracker
- `r3` - Enter key constant (13) / clear character (0)

**Note:** This fragment appears functionally identical to `readstring`.

---

## sleep - Delay Function

**Purpose:** Implement a simple delay using nested loops.

### Input Parameters
- `r0` (or outer) - Outer loop counter
- `r1` (or inner) - Inner loop counter value (copied to temp each iteration)
- `r2` (or temp) - Temporary loop counter

### Output
- Counters are decremented to zero

### Default Parameters
- `default_outer: r0` - Outer loop register
- `default_inner: r1` - Inner loop register
- `default_temp: r2` - Temporary register

### Operation
The fragment:
1. Checks if outer loop counter is zero (done)
2. Copies inner loop value to temp register
3. Decrements outer loop counter
4. Inner loop: decrements temp register until zero
5. Repeats until outer loop counter is zero

### Delay Calculation
Total iterations = `outer × inner`

### Registers Modified
- Outer counter (default r0) - Decremented to 0
- Temp counter (default r2) - Used for inner loop

---

## sync - Synchronize Memory to Output

**Purpose:** Copy a block of memory to video output (useful for refreshing screen from memory buffer).

### Input Parameters
- `r0` - Starting position within the source memory

### Output
- Screen is updated with memory contents

### Default Parameters
- `default_out: vtm0` - Output device (video textual memory)
- `default_height: 16` - Screen height in characters
- `default_width: 16` - Screen width in characters
- `default_mem: ram` - Source memory type

### Operation
The fragment:
1. Calculates total characters to copy (height × width)
2. Loops through each position:
   - Reads character from source memory
   - Writes character to output memory
   - Increments both pointers
3. Continues until all characters are copied

### Registers Modified
- `r1` - Total character counter
- `r2` - Output position
- `r3` - Temporary for character data

---

## Fragment Usage Notes

### General Guidelines

1. **Memory Types**: Fragments use template parameters (e.g., `{{ .Params.mem }}`) that can be configured with different memory types:
   - `rom` - Read-only memory
   - `ram` - Random access memory
   - `vtm0` - Video textual memory
   - Custom memory types as defined in your BondMachine configuration

2. **Register Conventions**: 
   - `r0`, `r1`, `r2` are commonly used for input parameters
   - Higher registers (`r3`-`r7`) are typically used as temporaries
   - Always check which registers are modified by each fragment

3. **Null Termination**: String fragments expect or produce null-terminated strings (ending with 0x00)

4. **Enter Key**: Keyboard reading fragments use ASCII 13 (carriage return) as the line terminator

5. **Parameterization**: Use the `default_*` parameters to customize fragment behavior for your specific hardware configuration

### Example Usage

```assembly
; Example: Write a message to screen
mov r0, rom:message
mov r1, 0
%call putstring

; Example: Read input from keyboard
mov r0, ram:buffer
mov r1, 0
%call read
; r2 now contains the length of input

; Example: Sleep for a delay
mov r0, 100    ; outer loop
mov r1, 1000   ; inner loop
%call sleep    ; total: 100,000 iterations
```

---

## Fragment Location

All fragments are located in: `/home/mirko/Projects/BondMachine/library/fragments/`

Each fragment is stored in its own `.basm` file with the fragment name as the filename.
