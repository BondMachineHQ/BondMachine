# bmstack

bmstack is part of BondMachine project. Within the project It is used to create stacks and queues distributed among several BondMachine cores. due to its generality it can also be used as standalone tool to produce HDL shared stacks and queues to be used in other projects.
The HDL code is build using golang templates that creates the code starting from the following data structure. The comments describe the meaning of each field.

```go
type Push struct {
	Agent string // The name of the agent that is pushing
	Tick  uint64 // The tick at which the push occurs
	Value string // The value that is pushed
}

type Pop struct {
	Agent string // The name of the agent that is popping
	Tick  uint64 // The tick at which the pop occurs
}

type TestBenchData struct {
	Pops         []Pop    // List of pops
	Pushes       []Push   // List of pushes
	TestSequence []string // Pushes and pops in order
}

type BmStack struct {
	ModuleName string   // The name of the module
	DataSize   int      // The size of the data bus
	Depth      int      // The depth of the stack
	Senders    []string // The names of the agents that can send data to the stack
	Receivers  []string // The names of the agents that can receive data from the stack
	MemType    string   // "LIFO" for a stack or "FIFO" for a queue
	funcMap    template.FuncMap

	// TestBench data
	TestBenchData
}
```

## API 

The library can be used in two ways.
The first one is from a go program. After creating a BmStack structure, the user can call the WriteHDL function to produce the HDL code.
An example on how the library can be used this way can be seen in the go test file (bmstack_test.go).

```go
func TestWriteHDL(t *testing.T) {
	// Create a new stack with 4 agents
	stack := BmStack{
		ModuleName: "test",
		DataSize:   32,
		Depth:      8,
		Senders:    []string{"sender1", "sender2"},
		Receivers:  []string{"receiver1", "receiver2"},
		MemType:    "LIFO"
	}

	// write the HDL code
	stack.WriteHDL()
```

The second part of the struct can be filled with test data. These data are used to produce a test bench for the stack. Using the test bench the user can verify the correctness of the stack. The test bench is produced by calling the WriteTestBench function.

```go
func TestWriteTestBench(t *testing.T) {
	stack := BmStack{
	// ...
		s.Pushes = []Push{
			Push{"sender1", 200, "32'd1"},
		}
		s.Pops = []Pop{
			Pop{"receiver1", 60},
		}
	}
	// write the test bench
	stack.WriteTestBench()
}
```

The two functions (WriteHDL and WriteTestBench) return a string that contains the HDL code or an error if something went wrong (string, error).
The provided go test file (bmstack_test.go) shows how the library can be used in this way. When invoked (go test), it produces the two files (bmstack.v and bmstack_tb.v).

## CLI

the library also came with a companion CLI executable called bmstack that provides the basic interface to the library. This is the second Possible way to use it.

```bash
$ bmstack -h
Usage of bmstack:
  -d    Verbose
  -data-width int
        Width of the data bus (default 32)
  -depth int
        Depth of the stack/queue (default 8)
  -hdl-file string
        Name of the file to write the HDL to (empty string to disable) (default "stack.v")
  -memory-type string
        Memory type, either stack or queue (default "queue")
  -random-stimulus int
        Generate random stimulus including N pushes and pops for every agent (0 to disable)
  -receivers string
        Comma separated list of names of signal tags that will receive data from the stack/queue
  -senders string
        Comma separated list of names of signal tags that will send data to the stack/queue
  -sim-length int
        Length of the simulation in clock cycles (default 1000)
  -stimulus-file string
        Name of the JSON file to load the stimulus from (empty string to disable)
  -tb-file string
        Name of the file to write the testbench to (empty string to disable)
  -v    Verbose
```

Examples of usage:


## Simulation

WIP

To see how the library is used in the BondMachine project, please refer to the BondMachine examples directory that contains many examples of BondMachines using stacks and queues.
