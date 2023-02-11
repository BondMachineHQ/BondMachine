package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmstack"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var typeR = flag.String("memory-type", "queue", "Memory type, either stack or queue")
var dataWidth = flag.Int("data-width", 32, "Width of the data bus")
var depth = flag.Int("depth", 8, "Depth of the stack/queue")

var senders = flag.String("senders", "", "Comma separated list of names of signal tags that will send data to the stack/queue")
var receivers = flag.String("receivers", "", "Comma separated list of names of signal tags that will receive data from the stack/queue")

var hdlFile = flag.String("hdl-file", "stack.v", "Name of the file to write the HDL to (default: stack.v, empty string to disable)")
var tbFile = flag.String("tb-file", "", "Name of the file to write the testbench to (default: empty, empty string to disable)")

var stimulusFile = flag.String("stimulus-file", "", "Name of the JSON file to load the stimulus from (default: empty, empty string to disable)")
var randomStimulus = flag.Int("random-stimulus", 0, "Generate random stimulus including N pushes and pops for every agent (default: 0, 0 to disable)")
var simLength = flag.Int("sim-length", 1000, "Length of the simulation in clock cycles (default: 1000)")

func init() {
	flag.Parse()
}

func main() {
	bmStack := bmstack.CreateBasicStack()

	switch *typeR {
	case "stack":
		bmStack.MemType = "LIFO"
	case "queue":
		bmStack.MemType = "FIFO"
	default:
		log.Fatal("Invalid memory type")
	}

	if *dataWidth < 1 && *dataWidth > 128 {
		log.Fatal("Invalid data width")
	}
	bmStack.DataSize = *dataWidth

	if *depth < 1 {
		log.Fatal("Invalid depth")
	}
	bmStack.Depth = *depth

	allAgents := make(map[string]struct{})

	if *senders != "" {
		for _, sender := range strings.Split(*senders, ",") {
			if _, ok := allAgents[sender]; !ok {
				bmStack.Senders = append(bmStack.Senders, sender)
				allAgents[sender] = struct{}{}
			} else {
				log.Fatal("Duplicate agent name " + sender)
			}
		}
	} else {
		log.Fatal("No senders specified")
	}

	if *receivers != "" {
		for _, receiver := range strings.Split(*receivers, ",") {
			if _, ok := allAgents[receiver]; !ok {
				bmStack.Receivers = append(bmStack.Receivers, receiver)
				allAgents[receiver] = struct{}{}
			} else {
				log.Fatal("Duplicate agent name " + receiver)
			}
		}
	} else {
		log.Fatal("No receivers specified")
	}

	if *hdlFile != "" {
		hdl, err := bmStack.WriteHDL()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(*hdlFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(hdl)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *stimulusFile != "" {
		err := bmStack.LoadJSON(*stimulusFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *randomStimulus > 0 {
		for i := 0; i < *randomStimulus; i++ {
			for _, sender := range bmStack.Senders {
				tick := rand.Intn(*simLength)
				value := rand.Intn(10)
				valueS := strconv.Itoa(value)
				bmStack.Pushes = append(bmStack.Pushes, bmstack.Push{Agent: sender, Tick: uint64(tick), Value: valueS})
			}
			for _, receiver := range bmStack.Receivers {
				tick := rand.Intn(*simLength)
				bmStack.Pops = append(bmStack.Pops, bmstack.Pop{Agent: receiver, Tick: uint64(tick)})
			}
		}
	}

	if *tbFile != "" {
		tb, err := bmStack.WriteTestBench()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(*tbFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(tb)
		if err != nil {
			log.Fatal(err)
		}
	}

}
