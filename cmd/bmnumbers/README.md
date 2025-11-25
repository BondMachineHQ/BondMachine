# BMNumbers

bmnumbers is part of BondMachine project. bmnumbers is both a command line tool to convert or cast numbers to and from different formats and a library to do the same. It is used within the BondMachine every time numbers are handled.

## Installation

```bash
go install github.com/BondMachineHQ/BondMachine/cmd/bmnumbers@latest
```

## Command Line Options

| Option | Description |
| ------ | ----------- |
| `-convert <type>` | Convert the input number to the specified type |
| `-cast <type>` | Cast the input number to the specified type (reinterpret bits) |
| `-show <format>` | Output format: `native` (default), `bin`, or `unsigned` |
| `-with-size` | Include size information in binary output |
| `-get-prefix <type>` | Get the prefix for a number type |
| `-get-size <type>` | Get the bit size of a number type |
| `-get-instructions <type>` | Get hardware instructions for a number type (JSON) |
| `-use-files` | Process CSV files instead of command line arguments |
| `-omit-prefix` | Omit type prefix in output |
| `-serve` | Start REST API server |
| `-linear-data-range <file>` | Load a linear data range file |
| `-v` | Enable verbose output |
| `-d` | Enable debug output |

## Usage Examples

### Basic Number Input

```bash
# Parse and display a number
bmnumbers 0u42        # unsigned integer: 42
bmnumbers 0x2a        # hexadecimal: 0x<8>2a
bmnumbers 0b101010    # binary: 0b<6>101010
bmnumbers 0f32.5      # float32: 0f<32>32.5
```

### Show Output in Different Formats

```bash
# Show as binary
bmnumbers -show bin 0u42           # Output: 101010
bmnumbers -show bin -with-size 0u42  # Output: 0b<64>101010

# Show as unsigned integer
bmnumbers -show unsigned 0f32.5    # Output: 1107427328
```

### Cast Between Types

```bash
# Cast binary to unsigned (reinterpret the bits)
bmnumbers -cast unsigned 0b101010  # Output: 42
```

### Get Type Information

```bash
# Get the prefix for a type
bmnumbers -get-prefix unsigned     # Output: 0u
bmnumbers -get-prefix float32      # Output: 0f<32>

# Get the size of a type
bmnumbers -get-size float32        # Output: 32

# Get hardware instructions for a type
bmnumbers -get-instructions float32
# Output: {"addop":"addf","divop":"divf","multop":"multf","powop":"multf"}
```

### Process CSV Files

```bash
# Process numbers in CSV files (creates output files with .out extension)
bmnumbers -use-files -cast unsigned input.csv
```

## Supported number types

The supported number types are listed in the following table.

| Type Name | Prefixes | Description | Static | Lenght |
| ---- | ------- | ----------- | ------ | ------ |
| unsigned | none <br> 0u <br> 0d | Unsigned integer | yes | any |
| signed | 0s <br> 0sd | Signed integer | yes | any |
| bin | 0b <br> 0b\<s\> | Binary number | yes | any <br> s bits|
| hex | 0x | Hexadecimal number | yes | any |
| float16 | 0f<16> | IEEE 754 half precision floating point number | yes | 16 bits |
| float32 | 0f <br> 0f<32> | IEEE 754 single precision floating point number | yes | 32 bits |
| lqs[s]t[t] | 0lq\<s.t\> | Linear quantized number with size s and type t | no | s bits |
| fps[s]f[f] | 0fp\<s.f\> | Fixed point number with size s and fraction f | no | s bits |
| flp[e]f[f] | 0flp\<e.f\> | FloPoCo floating point number with exponent e and mantissa f | no | e+f+3 bits |
