# BondMachine

The `bondmachine` command is the primary command-line tool for creating, manipulating, simulating, and generating HDL code for BondMachine architectures. A BondMachine is a heterogeneous multi-processor architecture where each processor can be configured and specialized for specific tasks.

## Installation

```bash
go install github.com/BondMachineHQ/BondMachine/cmd/bondmachine@latest
```

## Overview

The `bondmachine` command provides a comprehensive interface for working with BondMachine architectures. It supports operations ranging from basic machine construction (adding/removing processors, inputs, outputs) to advanced features like simulation, emulation, and Verilog generation for FPGA deployment.

## Required Flags

- `-bondmachine-file <file>`: Path to the BondMachine JSON file (required for all operations). If the file doesn't exist, a new BondMachine will be created.

## Core Features

### Machine Construction

#### Domains

Domains define processor architectures that can be instantiated within a BondMachine.

- `-list-domains`: List all domains in the BondMachine
- `-add-domains <files>`: Add processor domains from comma-separated JSON machine files
- `-del-domains <ids>`: Delete domains by comma-separated IDs

#### Processors

- `-list-processors`: List all processors in the BondMachine
- `-add-processor <domain_id>`: Add a processor from the specified domain
- `-enum-processors`: Enumerate all processors with detailed information

#### Inputs and Outputs

- `-list-inputs`: List all inputs
- `-add-inputs <count>`: Add the specified number of inputs
- `-del-inputs <ids>`: Delete inputs by comma-separated IDs
- `-list-outputs`: List all outputs
- `-add-outputs <count>`: Add the specified number of outputs
- `-del-outputs <ids>`: Delete outputs by comma-separated IDs

#### Internal I/O

- `-list-internal-inputs`: List internal inputs
- `-list-internal-outputs`: List internal outputs

#### Bonds (Inter-processor Connections)

Bonds connect processors together for communication.

- `-list-bonds`: List all bonds (connections between processors)
- `-add-bond <endpoint1>,<endpoint2>`: Add a bond between two endpoints
- `-del-bonds <ids>`: Delete bonds by comma-separated IDs
- `-enum-bonds`: Enumerate all bonds with detailed information

#### Shared Objects

- `-list-shared-objects`: List all shared objects
- `-add-shared-objects <objects>`: Add shared objects
- `-del-shared-objects <ids>`: Delete shared objects by IDs
- `-list-processor-shared-object-links`: List links between processors and shared objects
- `-connect-processor-shared-object <processor>,<object>`: Connect a processor to a shared object
- `-disconnect-processor-shared-object <processor>,<object>`: Disconnect a processor from a shared object

### Simulation and Emulation

#### Simulation

The simulation mode allows you to test your BondMachine without hardware.

- `-sim`: Enable simulation mode
- `-sim-interactions <count>`: Number of simulation cycles (default: 10)
- `-sim-stop-on-valid-of <output_id>`: Stop simulation when the specified output becomes valid
- `-sim-report <file>`: Generate a CSV simulation report
- `-simbox-file <file>`: Simulation data file with rules and configuration

#### Emulation

- `-emu`: Enable emulation mode (with hardware-like drivers)
- `-emu-interactions <count>`: Number of emulation cycles (0 = infinite, default: 10)

### Verilog/HDL Generation

Generate hardware description language files for FPGA deployment.

#### Basic Verilog Options

- `-create-verilog`: Generate Verilog files
- `-verilog-flavor <flavor>`: Verilog device type (default: "iverilog", also supports: "de10nano")
- `-verilog-mapfile <file>`: JSON file mapping device I/O to BondMachine I/O
- `-verilog-simulation`: Generate simulation-oriented Verilog
- `-comment-verilog`: Add comments to generated Verilog code

#### Board-Specific Support

**Basys3 FPGA Board:**
- `-basys3-7segment`: Enable Basys3 7-segment display support
- `-basys3-7segment-map <mapping>`: 7-segment display mappings
- `-basys3-leds`: Enable Basys3 LED support
- `-basys3-leds-map <mapping>`: LED mappings
- `-basys3-leds-name <name>`: LED signal name (default: "led")

**iCE40 FPGA Boards:**
- `-icebreaker-leds`: Enable iCEBreaker LED support
- `-icebreaker-leds-map <mapping>`: iCEBreaker LED mappings
- `-icefun-leds`: Enable iCEFun LED support
- `-icefun-leds-map <mapping>`: iCEFun LED mappings
- `-ice40lp1k-leds`: Enable iCE40-LP1K LED support
- `-ice40lp1k-leds-map <mapping>`: iCE40-LP1K LED mappings

#### Peripheral Support

**VGA Text Display:**
- `-vgatext`: Enable multi-CPU VGA textual support
- `-vgatext-flavor <flavor>`: VGA flavor (currently supported: "800x600")
- `-vgatext-fonts <file>`: VGA fonts file
- `-vgatext-header <file>`: VGA header ROM file

**PS2 Keyboard:**
- `-ps2-keyboard`: Enable PS2 keyboard support via shared object
- `-ps2-keyboard-io`: Enable PS2 keyboard support via I/O
- `-ps2-keyboard-io-map <mapping>`: PS2 keyboard I/O mappings

**UART:**
- `-uart`: Enable UART support
- `-uart-mapfile <file>`: UART mappings file

**Counter:**
- `-counter`: Enable counter support
- `-counter-map <mapping>`: Counter mappings
- `-counter-slow-factor <factor>`: Counter slow factor (default: "23")

**Board Clock Control:**
- `-board-slow`: Enable board slow mode
- `-board-slow-factor <factor>`: Board slow factor (default: 1)

### Clustering and Networking

BondMachines can be connected together in distributed clusters.

#### Cluster Configuration

- `-cluster-spec <file>`: Cluster specification file
- `-peer-id <id>`: Peer ID of this BondMachine within the cluster

#### Etherbond (Ethernet-based Clustering)

- `-use-etherbond`: Enable etherbond support
- `-etherbond-flavor <flavor>`: Ethernet device type (currently supported: "enc60j28")
- `-etherbond-mapfile <file>`: Etherbond I/O mapping file
- `-etherbond-macfile <file>`: File mapping peers to MAC addresses

#### UDPBond (UDP-based Clustering)

- `-use-udpbond`: Enable UDP bond support
- `-udpbond-flavor <flavor>`: Network device type (currently supported: "esp8266")
- `-udpbond-mapfile <file>`: UDP bond I/O mapping file
- `-udpbond-ipfile <file>`: File mapping peers to IP addresses
- `-udpbond-netconfig <file>`: JSON file with network configuration

#### Bondirect (Direct Connection)

- `-use-bondirect`: Enable bondirect support
- `-bondirect-flavor <flavor>`: Bondirect device type (currently supported: "basic")
- `-bondirect-mapfile <file>`: Bondirect I/O mapping file
- `-bondirect-mesh <file>`: Bondirect mesh specification file

### BMAPI (BondMachine API)

Generate API libraries for interfacing with BondMachine from software.

- `-use-bmapi`: Build a BMAPI interface
- `-bmapi-language <language>`: API language (go, c, python)
- `-bmapi-framework <framework>`: API framework (e.g., "pynq")
- `-bmapi-flavor <flavor>`: BMAPI interconnect type
- `-bmapi-flavor-version <version>`: BMAPI interconnect version
- `-bmapi-mapfile <file>`: BMAPI I/O mapping file
- `-bmapi-liboutdir <dir>`: Output directory for BMAPI library
- `-bmapi-modoutdir <dir>`: Output directory for BMAPI kernel module
- `-bmapi-auxoutdir <dir>`: Output directory for BMAPI auxiliary material
- `-bmapi-packagename <name>`: Go package name
- `-bmapi-modulename <name>`: Go module name
- `-bmapi-generate-example <type>`: Generate example program using BMAPI
- `-bmapi-data-type <type>`: Data type for BMAPI (default: "float32")

### Information and Analysis

- `-specs`: Display BondMachine specifications
- `-emit-dot`: Generate GraphViz dot file (outputs to stdout)
- `-dot-detail <level>`: Detail level for dot file (1-5, default: 1)
- `-show-program-disassembled`: Show disassembled programs
- `-show-program-alias`: Show program aliases for processors (saves to p*.alias files)
- `-multi-abstract-assembly-file <file>`: Save the BondMachine as multi-abstract assembly file

### Advanced Options

#### Configuration Files

- `-bcof-file <file>`: Use a BCOF (BondMachine Common Object Format) file for RAM initialization
- `-bminfo-file <file>`: File containing BondMachine extra info (JSON)
- `-bmrequirements-file <file>`: File containing BondMachine requirements (JSON)
- `-linear-data-range <file>`: Load linear data range file (syntax: index,filename)

#### Hardware Optimizations

- `-hw-optimizations <opts>`: Comma-separated list of hardware optimizations

#### Benchmark Core

- `-attach-benchmark-core <args>`: Attach a benchmark core (requires 2 arguments)
- `-attach-benchmark-core-v2 <args>`: Attach a benchmark core v2 (requires 2 arguments)

#### General Options

- `-register-size <bits>`: Number of bits per register (default: 8)
- `-d`: Enable debug mode
- `-v`: Enable verbose output

## Usage Examples

### Create a New BondMachine

```bash
# Create a new BondMachine with 8-bit registers
bondmachine -bondmachine-file my_machine.json -register-size 8
```

### Build a Simple Machine

```bash
# Create a machine file
bondmachine -bondmachine-file simple.json -register-size 8

# Add a domain (processor architecture)
bondmachine -bondmachine-file simple.json -add-domains processor1.json

# Add processors from domain 0
bondmachine -bondmachine-file simple.json -add-processor 0
bondmachine -bondmachine-file simple.json -add-processor 0

# Add inputs and outputs
bondmachine -bondmachine-file simple.json -add-inputs 2
bondmachine -bondmachine-file simple.json -add-outputs 1

# List processors
bondmachine -bondmachine-file simple.json -list-processors

# Add a bond between processors
bondmachine -bondmachine-file simple.json -add-bond "p0o0,p1i0"

# View the machine structure
bondmachine -bondmachine-file simple.json -specs
```

### Generate Verilog for FPGA

```bash
# Generate Verilog files for iVerilog simulation
bondmachine -bondmachine-file my_machine.json \
  -create-verilog \
  -verilog-flavor iverilog \
  -verilog-simulation \
  -comment-verilog

# Generate Verilog for DE10-Nano FPGA
bondmachine -bondmachine-file my_machine.json \
  -create-verilog \
  -verilog-flavor de10nano \
  -verilog-mapfile io_mapping.json
```

### Simulate a BondMachine

```bash
# Run a simulation with 100 cycles
bondmachine -bondmachine-file my_machine.json \
  -sim \
  -sim-interactions 100 \
  -simbox-file simulation_rules.json

# Run simulation with CSV report
bondmachine -bondmachine-file my_machine.json \
  -sim \
  -sim-interactions 1000 \
  -simbox-file simulation_rules.json \
  -sim-report results.csv

# Run simulation with GraphViz output
bondmachine -bondmachine-file my_machine.json \
  -sim \
  -sim-interactions 50 \
  -emit-dot \
  -dot-detail 3 > machine_state.dot
```

### Create a Cluster

```bash
# Create a clustered BondMachine with etherbond (for peer 0)
bondmachine -bondmachine-file peer0.json \
  -create-verilog \
  -use-etherbond \
  -etherbond-flavor enc60j28 \
  -cluster-spec cluster.json \
  -peer-id 0 \
  -etherbond-mapfile eth_mapping.json \
  -etherbond-macfile mac_addresses.json

# Create peer 1
bondmachine -bondmachine-file peer1.json \
  -create-verilog \
  -use-etherbond \
  -etherbond-flavor enc60j28 \
  -cluster-spec cluster.json \
  -peer-id 1 \
  -etherbond-mapfile eth_mapping.json \
  -etherbond-macfile mac_addresses.json
```

### Generate BMAPI

```bash
# Generate a Go API library
bondmachine -bondmachine-file my_machine.json \
  -use-bmapi \
  -bmapi-language go \
  -bmapi-mapfile api_mapping.json \
  -bmapi-liboutdir ./bmapi_lib \
  -bmapi-packagename mybmapi \
  -bmapi-generate-example test

# Generate Python API with PYNQ framework
bondmachine -bondmachine-file my_machine.json \
  -use-bmapi \
  -bmapi-language python \
  -bmapi-framework pynq \
  -bmapi-mapfile api_mapping.json \
  -bmapi-liboutdir ./bmapi_python
```

### Board-Specific Examples

```bash
# Generate Verilog for Basys3 with LED support
bondmachine -bondmachine-file basys3_machine.json \
  -create-verilog \
  -basys3-leds \
  -basys3-leds-map "o0" \
  -basys3-leds-name "led"

# Generate for iCEBreaker FPGA
bondmachine -bondmachine-file icebreaker_machine.json \
  -create-verilog \
  -icebreaker-leds \
  -icebreaker-leds-map "o0,o1,o2"
```

### Emulation Mode

```bash
# Run in emulation mode with VGA text output
bondmachine -bondmachine-file my_machine.json \
  -emu \
  -emu-interactions 0 \
  -vgatext \
  -vgatext-flavor 800x600 \
  -vgatext-fonts font.bin \
  -vgatext-header header.bin
```

## File Formats

### BondMachine JSON File

The BondMachine file stores the complete machine state including processors, bonds, inputs, outputs, and shared objects. It is automatically created/updated when you run commands.

### I/O Mapping Files

JSON files that map BondMachine I/O to physical device pins or protocol endpoints. Example structure:

```json
{
  "inputs": {
    "i0": "physical_pin_name"
  },
  "outputs": {
    "o0": "physical_pin_name"
  }
}
```

### Cluster Specification File

Defines the cluster topology including peer information and connections.

### Simbox File

Contains simulation rules for testing. See the `simbox` command documentation for details on creating simulation files.

## Architecture Flow

1. **Design**: Create processor domains using `procbuilder` command
2. **Build**: Construct BondMachine with `bondmachine` command (add processors, I/O, bonds)
3. **Test**: Simulate with `-sim` flag
4. **Deploy**: Generate Verilog with `-create-verilog` for FPGA deployment
5. **Interface**: Create API with `-use-bmapi` for software integration

## Tips

1. **Always specify** `-bondmachine-file` - it's required for all operations
2. **Test with simulation** before generating expensive Verilog synthesis
3. **Use verbose mode** (`-v`) to see detailed operation results
4. **List before delete** - use list commands to check IDs before deletion
5. **Start simple** - build incrementally and test at each stage
6. **Backup your files** - BondMachine files are overwritten on each operation

## See Also

- [BondMachine Website](https://www.bondmachine.it)
- Related commands:
  - `procbuilder`: Build processor architectures
  - `simbox`: Manage simulation rules
  - `basm`: BondMachine assembler
  - `bmbuilder`: Higher-level BondMachine builder
- Main repository documentation
- Package documentation: `pkg/bondmachine/`

## Error Handling

The command will panic with descriptive error messages if:

- Required files are missing or invalid
- Invalid JSON in configuration files
- Invalid domain/processor/bond IDs
- Incompatible flag combinations
- Hardware optimization names are unknown

Always check error messages for guidance on fixing issues.
