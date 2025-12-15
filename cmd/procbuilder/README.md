# Procbuilder

The `procbuilder` command is a command-line tool for building, configuring, and generating custom processor (CP) architectures for the BondMachine ecosystem. It provides a complete interface to create processors with specific instruction sets, assemble programs, generate hardware description language (HDL) code, and simulate execution.

## Installation

Build the procbuilder CLI:

```bash
go build github.com/BondMachineHQ/BondMachine/cmd/procbuilder
```

## Overview

Procbuilder supports two main workflows:

1. **Creating a new processor**: Define processor architecture, specify instruction sets, assemble programs, and save the configuration
2. **Loading an existing processor**: Load a saved processor from JSON, inspect it, modify it, generate HDL code, and simulate execution

The tool operates on a "machine" concept that includes:
- **Architecture (Arch)**: The processor architecture definition including registers, memory, I/O, and supported opcodes
- **Program**: The compiled program to run on the processor

## Usage

```bash
procbuilder [flags]
```

## Architecture Configuration Flags

These flags define the processor architecture when creating a new processor:

### Execution Model

- `-execution-model <model>`: Execution model for the processor (default: `ha`)
  - `vn`: Von Neumann architecture (shared memory for code and data)
  - `ha`: Harvard architecture (separate memory for code and data)
  - `hy`: Hybrid architecture

### Register and Memory Configuration

- `-register-size <n>`: Number of bits per register (default: `8`)
- `-registers <n>`: Number of registers as 2^n (default: `3`, meaning 2^3 = 8 registers)
- `-ram <n>`: Number of RAM memory cells as 2^n (default: `8`, meaning 2^8 = 256 cells)
- `-rom <n>`: Number of ROM memory cells as 2^n (default: `8`, meaning 2^8 = 256 cells)

### I/O Configuration

- `-inputs <n>`: Number of n-bit inputs (default: `1`)
- `-outputs <n>`: Number of n-bit outputs (default: `1`)

### Instruction Set Configuration

- `-opcodes <list>`: Comma-separated list of enabled opcodes (default: `nop`)
- `-list-opcodes`: List all available opcodes and exit
- `-opcode-optimizer`: Automatically determine required opcodes from assembly input

### Shared Resources

- `-shared-constraints <list>`: List of shared objects connected to the processor

## Program Input Options

Procbuilder can load programs from various sources:

- `-input-assembly <file>`: Load assembly program from file
- `-input-binary <file>`: Load binary program from file (limited support - not fully implemented)
- `-input-random`: Generate a random program

## Machine State Management

### Loading Machines

- `-load-machine <file>`: Load a processor and program from a JSON file

When loading a machine, the processor architecture and program are restored from the saved state. This is useful for:
- Continuing work on an existing processor
- Generating HDL code from a saved configuration
- Running simulations on pre-configured processors

### Saving Machines

- `-save-machine <file>`: Save the processor and program to a JSON file

The saved JSON file contains the complete machine state including architecture configuration and compiled program. The file will only be created if it doesn't already exist.

## Display and Inspection Flags

### Instruction Information

- `-show-instructions-alias`: Display instruction aliases for the processor
- `-show-opcodes`: List loaded opcodes
- `-show-opcodes-details`: Show detailed information for each loaded opcode

### Program Information

- `-show-program-alias`: Display program with instruction aliases
- `-show-program-binary`: Show the binary representation of the program
- `-show-program-disassembled`: Disassemble and display the program
- `-numlines`: Add line numbers to output
- `-hex`: Display values in hexadecimal format

## Verilog/HDL Generation

Procbuilder can generate Verilog HDL code for hardware implementation:

### Generate All Verilog Files

- `-create-verilog`: Generate all default Verilog files (`processor.v`, `ram.v`, `rom.v`, `arch.v`, `testbench.v`, `main.v`)

### Generate Individual Verilog Files

- `-create-verilog-processor <file>`: Generate processor Verilog file
- `-create-verilog-ram <file>`: Generate RAM Verilog file
- `-create-verilog-rom <file>`: Generate ROM Verilog file
- `-create-verilog-arch <file>`: Generate architecture Verilog file
- `-create-verilog-testbench <file>`: Generate testbench Verilog file
- `-create-verilog-main <file>`: Generate main Verilog file for FPGA

### Verilog Flavor

- `-verilog-flavor <flavor>`: Target Verilog device (default: `iverilog`)
  - `iverilog`: Icarus Verilog simulator
  - `kintex7`: Xilinx Kintex-7 FPGA

## Simulation and Execution

### Simulation Mode

- `-sim`: Enable simulation mode
- `-sim-interactions <n>`: Number of simulation interactions (default: `10`)
- `-simbox-file <file>`: Load simulation configuration from simbox JSON file

Simulation mode provides detailed step-by-step execution with register and I/O state before and after each instruction.

### Run Mode

- `-run`: Enable run mode (faster execution without detailed output)
- `-run-interactions <n>`: Number of run interactions (default: `1000`)

Run mode executes the program quickly without detailed output, suitable for performance testing.

## Other Flags

- `-d`: Enable debug mode
- `-v`: Enable verbose output

## Examples

### Example 1: List Available Opcodes

```bash
procbuilder -list-opcodes
```

### Example 2: Create a Simple Processor

Create an 8-bit processor with Harvard architecture, 8 registers, and basic opcodes:

```bash
procbuilder \
  -execution-model ha \
  -register-size 8 \
  -registers 3 \
  -ram 8 \
  -rom 8 \
  -opcodes nop,mov,add,sub,j \
  -save-machine myprocessor.json
```

### Example 3: Create a Processor and Assemble a Program

Create a processor with opcodes optimized for the assembly program:

```bash
procbuilder \
  -execution-model ha \
  -register-size 8 \
  -registers 3 \
  -input-assembly program.asm \
  -opcode-optimizer \
  -save-machine myprocessor.json
```

### Example 4: Load and Inspect a Processor

Load a saved processor and display its opcodes and disassembled program:

```bash
procbuilder \
  -load-machine myprocessor.json \
  -show-opcodes \
  -show-program-disassembled
```

### Example 5: Generate Verilog from a Saved Processor

```bash
procbuilder \
  -load-machine myprocessor.json \
  -create-verilog
```

This generates: `processor.v`, `ram.v`, `rom.v`, `arch.v`, `testbench.v`, and `main.v`

### Example 6: Generate Verilog for Specific FPGA

```bash
procbuilder \
  -load-machine myprocessor.json \
  -create-verilog \
  -verilog-flavor kintex7
```

### Example 7: Simulate a Processor

```bash
procbuilder \
  -load-machine myprocessor.json \
  -sim \
  -sim-interactions 100
```

### Example 8: Run a Processor with Simbox Configuration

```bash
procbuilder \
  -load-machine myprocessor.json \
  -sim \
  -simbox-file simulation.json \
  -sim-interactions 50
```

### Example 9: Complete Workflow

Create a processor, assemble a program, generate Verilog, and simulate:

```bash
# Step 1: Create processor with program
procbuilder \
  -execution-model ha \
  -register-size 8 \
  -registers 3 \
  -ram 8 \
  -rom 8 \
  -input-assembly program.asm \
  -opcode-optimizer \
  -save-machine myprocessor.json

# Step 2: Inspect the processor
procbuilder \
  -load-machine myprocessor.json \
  -show-opcodes-details \
  -show-program-disassembled

# Step 3: Generate Verilog
procbuilder \
  -load-machine myprocessor.json \
  -create-verilog

# Step 4: Simulate
procbuilder \
  -load-machine myprocessor.json \
  -sim \
  -sim-interactions 50
```

## Typical Workflows

### Workflow 1: New Processor Development

1. Create a processor with specific architecture and instruction set
2. Load an assembly program
3. Save the machine configuration
4. Inspect the generated program
5. Simulate execution
6. Generate HDL code for hardware implementation

### Workflow 2: Hardware Implementation

1. Load an existing processor configuration
2. Generate Verilog files for target platform
3. Use the generated files in your FPGA or ASIC workflow

### Workflow 3: Program Development and Testing

1. Load a processor architecture
2. Assemble and test different programs
3. Use simulation mode to debug programs
4. Save successful configurations

## JSON Machine File Format

The machine JSON file contains the complete processor state:

```json
{
  "Arch": {
    "Modes": ["ha"],
    "Rsize": 8,
    "R": 3,
    "L": 8,
    "N": 1,
    "M": 1,
    "O": 8,
    "Op": [...],
    "Shared_constraints": ""
  },
  "Program": {
    "Slocs": [...]
  }
}
```

## Notes

- Files are only created if they don't already exist (no overwriting)
- The tool panics on errors (file not found, parse errors, constraint violations, etc.)
- Constraint checking is performed before any operations
- The opcode optimizer (`-opcode-optimizer`) analyzes the assembly file and includes only the opcodes actually used

## See Also

- [procbuilder package documentation](../../pkg/procbuilder/README.md)
- [Simbox CLI](../simbox/README.md)
- [BondMachine Documentation](https://www.bondmachine.it)