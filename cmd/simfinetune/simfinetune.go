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
var delaysFile = flag.String("delays-file", "delaysout.json", "Output delays parameters in JSON format")
var geneticConfigFile = flag.String("genetic-config-file", "", "Genetic algorithm configuration in JSON format")

type record []string

type FitnessEnv struct {
	Inputs        []record
	Outputs       []record
	RealLatencies []uint32
	BM            *bondmachine.Bondmachine
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
	if *delaysFile == "" {
		flag.Usage()
		panic("Missing required argument for delays file")
	}
}

func (fe *FitnessEnv) FitnessFunction(simDelays *simbox.SimDelays) float64 {
	bm := fe.BM
	computedLatencies := make([]uint32, len(fe.Outputs))
	for i, rec := range fe.Inputs {
		if out, err := bm.SinglePipelineSimulate("float32", rec, simDelays); err == nil {
			latency := out[bm.Outputs-1]
			fmt.Sscanf(latency, "%d", &computedLatencies[i])
			// fmt.Printf("Input: %v, Output: %v, Expected: %v\n", rec, out, fe.Outputs[i])
		} else {
			return 0.0
		}
	}

	// Compute fitness based on the difference between computed and real latencies
	var totalError float64
	for i, realLatency := range fe.RealLatencies {
		computedLatency := computedLatencies[i]
		error := float64(realLatency) - float64(computedLatency)
		totalError += error * error // Squared error
	}
	// Lower error means better fitness; we can invert it
	if totalError > 0 {
		fmt.Println(1.0 / totalError)
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

	usedOpcodes := bm.GetUsedOpcodes()

	fe := &FitnessEnv{
		Inputs:  inputs,
		Outputs: outputs,
		BM:      bm,
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

	// At this point, we have the Bondmachine, inputs, and outputs loaded
	// Further processing would go here

	bestSimDelays, _ := simbox.RunGeneticAlgorithm(usedOpcodes, geneticConfig, fe.FitnessFunction)

	// Save the best delays to the delays file
	fmt.Println("Best SimDelays:", bestSimDelays)

}
