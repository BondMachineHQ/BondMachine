package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
	"github.com/BondMachineHQ/BondMachine/pkg/simbox"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var bondmachineFile = flag.String("bondmachine-file", "", "Bondmachine in JSON format")
var inputsFile = flag.String("inputs-file", "", "Inputs in CSV format")
var outputsFile = flag.String("outputs-file", "", "Outputs in CSV format, with expected results and latencies")
var delaysInFile = flag.String("delays-input-file", "", "Input delays parameters in JSON format")
var delaysOutFile = flag.String("delays-output-file", "delaysout.json", "Output delays parameters in JSON format")
var geneticConfigFile = flag.String("genetic-config-file", "", "Genetic algorithm configuration in JSON format")
var workers = flag.Int("workers", 4, "Number of concurrent workers for simulation")
var includeOpcodes = flag.String("include-opcodes", "", "Comma-separated list of opcodes to include in the optimization")
var excludeOpcodes = flag.String("exclude-opcodes", "", "Comma-separated list of opcodes to exclude from the optimization")

type record []string

type FitnessEnv struct {
	Inputs              []record
	Outputs             []record
	RealLatencies       []uint32
	LatencyDistribution *simbox.DelayDistribution
	Workers             int
	BM                  *bondmachine.Bondmachine
	SimDelays           *simbox.SimDelays
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func init() {
	flag.Parse()
	if *bondmachineFile == "" || *inputsFile == "" || *outputsFile == "" {
		flag.Usage()
		panic("Missing required arguments")
	}
	if *delaysOutFile == "" {
		flag.Usage()
		panic("Missing required argument for delays file")
	}
	if *includeOpcodes != "" && *excludeOpcodes != "" {
		flag.Usage()
		panic("Cannot use both include-opcodes and exclude-opcodes at the same time")
	}
}

func (fe *FitnessEnv) FitnessFunction(simDelays *simbox.SimDelays) float64 {
	bm := fe.BM
	computedLatencies := make([]uint32, len(fe.Outputs))

	chanExits := make(chan struct{})
	chanIn := make(chan int)
	chanOut := make(chan struct {
		Index   int
		Latency uint32
	})

	// Merge the base simDelays with the candidate simDelays updates
	mergedSimDelays := simbox.MergeSimDelays(fe.SimDelays, simDelays)

	workers := fe.Workers
	if workers <= 0 {
		workers = 1
	}
	for w := 0; w < workers; w++ {
		go func() {
			for {
				select {
				case <-chanExits:
					return
				case i := <-chanIn:
					rec := fe.Inputs[i]
					if out, err := bm.SinglePipelineSimulate("float32", rec, mergedSimDelays); err == nil {
						latency := out[bm.Outputs-1]
						var latencyVal uint32
						fmt.Sscanf(latency, "%d", &latencyVal)
						chanOut <- struct {
							Index   int
							Latency uint32
						}{Index: i, Latency: latencyVal}
					}
				}
			}
		}()
	}

	for i, j := 0, 0; i < len(fe.Inputs)*2; i++ {
		if j >= len(fe.Inputs) {
			result := <-chanOut
			computedLatencies[result.Index] = result.Latency
			continue
		}

		select {
		case chanIn <- j:
			j++
		case result := <-chanOut:
			computedLatencies[result.Index] = result.Latency
		}
	}

	for w := 0; w < workers; w++ {
		chanExits <- struct{}{}
	}

	latencyDistribution := make(simbox.DelayDistribution)
	for _, cl := range computedLatencies {
		if val, ok := (latencyDistribution)[int32(cl)]; ok {
			latencyDistribution[int32(cl)] = val + 1
		} else {
			latencyDistribution[int32(cl)] = 1
		}
	}

	// Compute fitness based on the difference between computed and real latencies
	var totalError float64
	// for i, realLatency := range fe.RealLatencies {
	// 	computedLatency := computedLatencies[i]
	// 	error := float64(realLatency) - float64(computedLatency)
	// 	totalError += error * error // Squared error
	// }
	totalError = simbox.DistributionDistance(latencyDistribution, *fe.LatencyDistribution)

	// Lower error means better fitness; we can invert it
	if totalError > 0 {
		// fmt.Println(1.0 / totalError)
		return 1.0 / totalError
	}
	return 1.0
}

func main() {

	// Load the Bondmachine
	var bm *bondmachine.Bondmachine

	if _, err := os.Stat(*bondmachineFile); err == nil {
		// Open the bondmachine file is exists
		if bondmachineJSON, err := os.ReadFile(*bondmachineFile); err == nil {
			var bmj bondmachine.Bondmachine_json
			if err := json.Unmarshal([]byte(bondmachineJSON), &bmj); err == nil {
				bm = (&bmj).Dejsoner()
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
		bm.Init()
	} else {
		panic(err)
	}

	// Load the inputs file into a suitable structure
	inputs := make([]record, 0)

	file, err := os.Open(*inputsFile)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ",")

		var rec record
		for _, field := range fields {
			rec = append(rec, strings.TrimSpace(field))
		}
		inputs = append(inputs, rec)
	}
	check(scanner.Err())

	// Load the outputs file into a suitable structure
	outputs := make([]record, 0)

	fileOut, err := os.Open(*outputsFile)
	check(err)
	defer fileOut.Close()

	scannerOut := bufio.NewScanner(fileOut)
	for scannerOut.Scan() {
		line := scannerOut.Text()
		fields := strings.Split(line, ",")

		var rec record
		for _, field := range fields {
			rec = append(rec, strings.TrimSpace(field))
		}
		outputs = append(outputs, rec)
	}
	check(scannerOut.Err())

	// Load optional genetic configuration
	var geneticConfig simbox.GeneticConfig
	if *geneticConfigFile != "" {
		if gc, err := simbox.GetGeneticConfigFromJSON(*geneticConfigFile); err != nil {
			panic(err)
		} else {
			geneticConfig = gc
		}
	} else {
		geneticConfig = simbox.GetDefaultGeneticConfig()
	}

	// Whenever available, load existing delays from file, to be used as base for the optimization
	// the opcodes under optimization will be modified, while the others will remain as in the loaded file
	var simDelays *simbox.SimDelays
	if *delaysInFile != "" {
		if sd, err := simbox.LoadSimDelaysFromFile(*delaysInFile); err != nil {
			panic(err)
		} else {
			simDelays = sd
		}
	} else {
		simDelays = simbox.NewSimDelays()
	}

	// Determine the used opcodes in the Bondmachine
	usedOpcodes := bm.GetUsedOpcodes()
	if *includeOpcodes != "" {
		included := strings.Split(*includeOpcodes, ",")
		opcodeSet := make(map[string]struct{})
		for _, op := range included {
			opcodeSet[strings.TrimSpace(op)] = struct{}{}
		}
		var filtered []string
		for _, op := range usedOpcodes {
			if _, ok := opcodeSet[op]; ok {
				filtered = append(filtered, op)
			}
		}
		usedOpcodes = filtered
	} else if *excludeOpcodes != "" {
		excluded := strings.Split(*excludeOpcodes, ",")
		opcodeSet := make(map[string]struct{})
		for _, op := range excluded {
			opcodeSet[strings.TrimSpace(op)] = struct{}{}
		}
		var filtered []string
		for _, op := range usedOpcodes {
			if _, ok := opcodeSet[op]; !ok {
				filtered = append(filtered, op)
			}
		}
		usedOpcodes = filtered
	}

	fe := &FitnessEnv{
		Inputs:    inputs,
		Outputs:   outputs,
		BM:        bm,
		Workers:   *workers,
		SimDelays: simDelays,
	}

	// Convert the real latencies from outputs
	realLatencies := make([]uint32, 0)
	for _, outRec := range outputs {
		if len(outRec) < 2 {
			panic("Outputs file must have at least two columns: expected output and real latency")
		}
		var latency uint32
		_, err := fmt.Sscanf(outRec[len(outRec)-1], "%d", &latency)
		if err != nil {
			panic(err)
		}
		realLatencies = append(realLatencies, latency)
	}
	fe.RealLatencies = realLatencies

	latencyDistribution := make(simbox.DelayDistribution)
	for _, rl := range realLatencies {
		if val, ok := (latencyDistribution)[int32(rl)]; ok {
			latencyDistribution[int32(rl)] = val + 1
		} else {
			latencyDistribution[int32(rl)] = 1
		}
	}

	fe.LatencyDistribution = &latencyDistribution
	// At this point, we have the Bondmachine, inputs, and outputs loaded
	// Further processing would go here

	bestSimDelays, _ := simbox.RunGeneticAlgorithm(usedOpcodes, geneticConfig, fe.FitnessFunction)

	// Save the best delays to the delays file
	simDelaysJSON, err := json.MarshalIndent(simbox.MergeSimDelays(simDelays, bestSimDelays), "", "  ")
	check(err)
	err = os.WriteFile(*delaysOutFile, simDelaysJSON, 0644)
	check(err)

}
