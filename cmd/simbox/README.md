# Simbox CLI

The `simbox` command-line tool allows you to manage simulation rules for BondMachine simulations. You can add, remove, suspend, and reactivate rules stored in a simbox JSON file.

## Installation

Build the simbox CLI:

```bash
go build github.com/BondMachineHQ/BondMachine/cmd/simbox
```

## Usage

```bash
simbox -simbox-file <file> [options]
```

### Required Flag

- `-simbox-file <file>`: Path to the simulation data file (JSON format). This flag is mandatory for all operations.

### Operations

#### List Rules

Display all rules in the simbox file:

```bash
simbox -simbox-file rules.json -list
```

**Example output:**

```
000 - absolute:10:set:i0:5
001 - relative:5:get:o0:unsigned
002 - config:show_ticks [SUSPENDED]
```

Rules marked with `[SUSPENDED]` are inactive and will be skipped during simulation.

#### Add a Rule

Add a new rule to the simbox:

```bash
simbox -simbox-file rules.json -add "absolute:10:set:i0:5"
```

**Rule syntax:**

- Absolute time rule: `absolute:<tick>:<action>:<object>:<extra>`
- Relative/periodic rule: `relative:<interval>:<action>:<object>:<extra>`
- On valid signal: `onvalid:<action>:<object>:<extra>`
- On receive signal: `onrecv:<action>:<object>:<extra>`
- On exit: `onexit:<action>:<object>:<extra>`
- Configuration rule: `config:<option>`

**Actions:**

- `set`: Set a value to an object
- `get`: Get a value from an object
- `show`: Show a value from an object

**Examples:**

```bash
# Set input i0 to 5 at absolute tick 10
simbox -simbox-file rules.json -add "absolute:10:set:i0:5"

# Get output o0 every 5 ticks
simbox -simbox-file rules.json -add "relative:5:get:o0:unsigned"

# Show processor 0 register 1 at tick 20
simbox -simbox-file rules.json -add "absolute:20:show:p0r1:hex"

# Get register r0 when simulation exits
simbox -simbox-file rules.json -add "onexit:get:r0:unsigned"

# Show register r1 in hex format when simulation exits
simbox -simbox-file rules.json -add "onexit:show:r1:hex"

# Configuration to show ticks
simbox -simbox-file rules.json -add "config:show_ticks"
```

#### Remove a Rule

Delete a rule by its index:

```bash
simbox -simbox-file rules.json -del 2
```

**Note:** This permanently removes the rule from the simbox file.

#### Suspend a Rule

Temporarily disable a rule without deleting it:

```bash
simbox -simbox-file rules.json -suspend 1
```

Suspended rules remain in the file but are marked as inactive. They will be skipped during simulation setup and appear with a `[SUSPENDED]` marker when listed.

**Use case:** Suspend rules during debugging or testing without losing them.

#### Unsuspend (Reactivate) a Rule

Reactivate a previously suspended rule:

```bash
simbox -simbox-file rules.json -unsuspend 1
```

The rule becomes active again and will be applied during simulation.

### Other Options

- `-verify`: Verify the simbox against a machine or bondmachine file
- `-machine-file <file>`: Machine file in JSON format (used with `-verify`)
- `-bondmachine-file <file>`: BondMachine file in JSON format (used with `-verify`)
- `-d`: Enable debug mode
- `-v`: Enable verbose output

## Complete Workflow Example

```bash
# Create a new simbox with rules
simbox -simbox-file simulation.json -add "absolute:10:set:i0:5"
simbox -simbox-file simulation.json -add "absolute:20:set:i0:10"
simbox -simbox-file simulation.json -add "relative:5:get:o0:unsigned"
simbox -simbox-file simulation.json -add "config:show_ticks"

# List all rules
simbox -simbox-file simulation.json -list
# Output:
# 000 - absolute:10:set:i0:5
# 001 - absolute:20:set:i0:10
# 002 - relative:5:get:o0:unsigned
# 003 - config:show_ticks

# Suspend a rule temporarily
simbox -simbox-file simulation.json -suspend 1

# Verify the rule is suspended
simbox -simbox-file simulation.json -list
# Output:
# 000 - absolute:10:set:i0:5
# 001 - absolute:20:set:i0:10 [SUSPENDED]
# 002 - relative:5:get:o0:unsigned
# 003 - config:show_ticks

# Reactivate the rule
simbox -simbox-file simulation.json -unsuspend 1

# Remove a rule permanently
simbox -simbox-file simulation.json -del 3
```

## Rule Indexing

Rules are indexed starting from 0. Use the `-list` option to see the current indices before performing operations on specific rules.

**Important:** After deleting a rule, subsequent rules shift down in index. Always list rules before removing or modifying by index.

## Error Handling

The tool will panic with an error message if:

- The simbox file is not specified
- An invalid rule syntax is provided
- A rule index is out of range
- The simbox file cannot be read or written

## File Format

The simbox file is stored in JSON format. Example:

```json
{
  "Rules": [
    {
      "Timec": 0,
      "Tick": 10,
      "Action": 0,
      "Object": "i0",
      "Extra": "5",
      "Suspended": false
    },
    {
      "Timec": 2,
      "Tick": 5,
      "Action": 1,
      "Object": "o0",
      "Extra": "unsigned",
      "Suspended": true
    }
  ]
}
```

You can edit this file manually, but it's recommended to use the CLI tool to maintain consistency.

## Tips

1. **Use suspend instead of delete** when you want to temporarily disable a rule during testing or debugging.
2. **Always list rules** before suspending/unsuspending by index to ensure you're targeting the correct rule.
3. **Back up your simbox file** before making significant changes, especially when deleting rules.
4. **Use descriptive extra fields** (like "unsigned", "signed", "hex") to make the output more readable.

## See Also

- BondMachine documentation
- Simbox package: `pkg/simbox/simbox.go`
- Rule suspension feature details in the main repository
