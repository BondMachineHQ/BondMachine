# Simbox Rules Documentation

## Overview

Simbox is a simulation control system that allows you to define rules for interacting with a BondMachine simulation. Rules can be time-based (absolute or relative), event-driven (on valid/receive signals), or configuration-based.

## Rule Format

Rules follow a colon-separated format: `type:parameter:action:object:extra`

## Time Constraint Types

- **absolute** - Execute at a specific tick number
- **relative** - Execute periodically every N ticks
- **onvalid** - Execute when a valid signal is received
- **onrecv** - Execute when a receive signal is triggered
- **onexit** - Execute when the simulation ends
- **config** - Configuration rules (no time constraint)

## Actions

- **set** - Set a value to an object
- **get** - Get a value from an object (machine-readable output)
- **show** - Display a value from an object (human-readable output)
- **config** - Configure simulation behavior

## Rule Types

### 1. Absolute Time Rules

Execute actions at a specific simulation tick.

**Format:** `absolute:<tick>:<action>:<object>:<extra>`

**Examples:**
```
absolute:100:set:r0:42
absolute:200:get:r1:unsigned
absolute:300:show:r2:hex
absolute:500:get:io_input:signed
```

### 2. Relative Time Rules

Execute actions periodically every N ticks.

**Format:** `relative:<period>:<action>:<object>:<extra>`

**Examples:**
```
relative:10:set:r0:100
relative:50:get:r1:unsigned
relative:100:show:r2:hex
relative:25:get:memory_0:signed
```

### 3. On Valid Signal Rules

Execute actions when a valid signal is received.

**Format:** `onvalid:<action>:<object>:<extra>` or `onvalid:<action>:<object>`

**Examples:**
```
onvalid:get:r0:unsigned
onvalid:show:r1:hex
onvalid:get:io_input:signed
onvalid:show:r2              # Uses default 'unsigned' format
```

### 4. On Receive Signal Rules

Execute actions when a receive signal is triggered.

**Format:** `onrecv:<action>:<object>:<extra>` or `onrecv:<action>:<object>`

**Examples:**
```
onrecv:get:r0:unsigned
onrecv:show:r1:hex
onrecv:get:io_input:signed
onrecv:show:r2               # Uses default 'unsigned' format
```

### 5. On Exit Rules

Execute actions when the simulation ends. This is useful for capturing final state or generating summary reports at the end of a simulation run.

**Format:** `onexit:<action>:<object>:<extra>` or `onexit:<action>:<object>`

**Examples:**
```
onexit:get:r0:unsigned
onexit:show:r1:hex
onexit:get:io_output:signed
onexit:show:r2               # Uses default 'unsigned' format
```

### 6. Configuration Rules

Control simulation display and behavior.

**Format:** `config:<option>` or `config:<option>:<parameter>`

#### Display Options (no parameters)

```
config:show_pc                    # Show program counter
config:show_instruction           # Show current instruction
config:show_disasm               # Show disassembly
config:show_ticks                # Show tick counter
config:get_ticks                 # Get tick counter (machine-readable)
config:show_proc_regs_pre        # Show processor registers before execution
config:show_proc_regs_post       # Show processor registers after execution
config:show_proc_io_pre          # Show processor I/O before execution
config:show_proc_io_post         # Show processor I/O after execution
config:show_io_pre               # Show I/O before execution
config:show_io_post              # Show I/O after execution
```

#### Bulk Display Options (with parameters)

```
config:get_all:<format>          # Get all registers/objects
config:get_all_internal:<format> # Get all internal state
config:show_all:<format>         # Show all registers/objects
config:show_all_internal:<format> # Show all internal state
```

**Format options:** `unsigned`, `signed`, `hex`, `binary`

**Examples:**
```
config:get_all:hex
config:show_all_internal:unsigned
```

## Objects

Objects that can be manipulated with set/get/show actions:

- **rN** - Register N (e.g., `r0`, `r1`, `r2`)
- **memory_N** - Memory location N
- **io_N** - I/O port N
- **io_input** - Input I/O
- **io_output** - Output I/O

## Extra Parameters

Optional formatting for get/show/set actions:

- **unsigned** - Display as unsigned integer (default)
- **signed** - Display as signed integer
- **hex** - Display in hexadecimal format
- **binary** - Display in binary format

## Complete Examples

### Basic Monitoring
```
config:show_pc
config:show_ticks
config:show_instruction
relative:1:show:r0:hex
relative:1:show:r1:hex
```

### Periodic Testing
```
relative:10:set:r0:100
relative:10:get:r1:unsigned
relative:50:show:r2:hex
config:show_all:unsigned
```

### Absolute Time Testing
```
absolute:0:set:r0:0
absolute:100:set:r0:50
absolute:200:get:r0:unsigned
absolute:300:show:r1:hex
absolute:500:get:io_input:signed
```

### Event-Driven Testing
```
onvalid:get:r0:unsigned
onvalid:show:r1:hex
onrecv:get:io_input:signed
onrecv:show:r2
onexit:get:r0:unsigned
onexit:show:r1:hex
config:show_ticks
```

### Mixed Rules Example
```
config:show_pc
config:show_instruction
config:show_ticks
absolute:0:set:r0:0
relative:10:show:r0:unsigned
onvalid:get:r1:hex
onrecv:show:r2:signed
onexit:get:r0:unsigned
onexit:show:r1:hex
config:get_all:hex
```

## Usage in Code

### Adding Rules
```go
simbox := &Simbox{}
simbox.Add("absolute:100:set:r0:42")
simbox.Add("relative:10:get:r1:unsigned")
simbox.Add("onvalid:show:r2:hex")
simbox.Add("onrecv:get:r3:signed")
simbox.Add("onexit:get:r0:unsigned")
simbox.Add("onexit:show:r1:hex")
simbox.Add("config:show_pc")
```

### Deleting Rules
```go
simbox.Del(0)  // Delete first rule
```

### Printing Rules
```go
fmt.Println(simbox.Print())
```

## Error Handling

Rules that cannot be decoded will return an error: `"rule cannot be decoded"`

Common errors:
- Invalid time constraint type
- Invalid action type
- Missing required parameters
- Invalid tick/period values (must be integers for absolute/relative rules)
- Unsupported object type
- Invalid format specifier

## Rule Format Summary

| Rule Type | Format | Example |
|-----------|--------|---------|
| Absolute Time | `absolute:<tick>:<action>:<object>:<extra>` | `absolute:100:set:r0:42` |
| Relative Time | `relative:<period>:<action>:<object>:<extra>` | `relative:10:get:r1:unsigned` |
| On Valid | `onvalid:<action>:<object>:<extra>` | `onvalid:show:r2:hex` |
| On Receive | `onrecv:<action>:<object>:<extra>` | `onrecv:get:r3:signed` |
| On Exit | `onexit:<action>:<object>:<extra>` | `onexit:get:r0:unsigned` |
| Config (simple) | `config:<option>` | `config:show_pc` |
| Config (with param) | `config:<option>:<parameter>` | `config:get_all:hex` |

**Note:** For `onvalid`, `onrecv`, and `onexit` rules, if the `<extra>` parameter is omitted, it defaults to `unsigned`.