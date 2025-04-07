package bmbuilder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
	"github.com/BondMachineHQ/BondMachine/pkg/bondmachine"
)

// Find metadata by key, querying in a cascaded-fashion with order: BasmLine -> BasmElement -> BMBuilder
func findMeta(key string, l *bmline.BasmLine, e *bmline.BasmElement, b *BMBuilder) (string, bool) {
	var result string
	if result = l.BasmMeta.GetMeta(key); result != "" {
		return result, true
	}
	if result = e.BasmMeta.GetMeta(key); result != "" {
		return result, true
	}
	if result = b.global.BasmMeta.GetMeta(key); result != "" {
		return result, true
	}

	return "", false
}

// FIXME: Bad API, should separate the case where the metadata was not found from parsing/other errors
func findDatatypeMeta(key string, l *bmline.BasmLine, e *bmline.BasmElement, b *BMBuilder) (bmnumbers.BMNumberType, error) {
	value, found := findMeta(key, l, e, b)
	if !found {
		return nil, fmt.Errorf("Meta definition '%s' not found", key)
	}

	// get the instructions from the data type using bmnumbers
	dataType := bmnumbers.GetType(value)
	if dataType == nil {
		return nil, fmt.Errorf("Data type '%s' not found", value)
	}

	return dataType, nil
}

// TODO: This should be a BMNumberType method (?)
func metaFloatLiteral(value float64, numType bmnumbers.BMNumberType) string {
	return fmt.Sprintf("%s%f", numType.ShowPrefix(), value)
}

type ProperStringBuilder struct {
	strings.Builder
}

func (self ProperStringBuilder) WriteStringln(str string) (int, error) {
	return self.WriteString(fmt.Sprintln(str))
}

func (self ProperStringBuilder) WriteFmtStringln(format string, a ...any) (int, error) {
	return self.WriteStringln(fmt.Sprintf(format, a...))
}

// FIXME: Suboptimal data model of DenseLayer, weights len can vary across nodes when it obviously shouldn't
type DenseLayerNode struct {
	weights []float32
	bias    float32
}

type DenseLayer struct {
	activationFunc string
	nodes          []DenseLayerNode
}

func DenseGenerator(b *BMBuilder, e *bmline.BasmElement, l *bmline.BasmLine) (*bondmachine.Bondmachine, error) {
	// Overview:
	// - We take N inputs from the previous input layer, each input value being a NN incoming layer node
	// - We'll produce a series of CPs for each node that compute the weighted sums, add the bias
	//   and apply the activation function
	// - Output M result values from this layer
	// TODO:
	// - Load the fragments library (or have the user load it?)
	// - Actually parse the DenseLayer from BASM/BM directives

	// Pay attention to the following: This code works only if the Generator function is called on a sequential block
	// TODO Include here a check and error if the Generator function is called outside a sequential block

	if b.debug {
		fmt.Println(green("\t\t\tDense Generator - Start"))
		defer fmt.Println(green("\t\t\tDense Generator - End"))
	}

	cb := b.currentBlock
	if len(b.blocks[cb].blockBMs) == 0 {
		return nil, errors.New("Dense Generator: No previous BMs found, a Dense Generator cannot be the first block")
	}

	prevBM := b.blocks[cb].blockBMs[len(b.blocks[cb].blockBMs)-1]
	bmInputsCount := prevBM.Outputs

	// Get the data type, starting from the maxpool metadata and falling back to the global metadata
	dataType, err := findDatatypeMeta("datatype", l, e, b)
	if err != nil {
		return nil, fmt.Errorf("DenseGenerator: Could not retrieve datatype param: %w", err)
	}

	ops := dataType.ShowInstructions()
	multop, multopFound := ops["multop"]
	if !multopFound {
		return nil, fmt.Errorf("Dense Generator: Undefined multop for datatype '%s'", dataType.GetName())
	}
	addop, addopFound := ops["addop"]
	if !addopFound {
		return nil, fmt.Errorf("Dense Generator: Undefined addop for datatype '%s'", dataType.GetName())
	}

	// HARDCODED
	// TODO: Actually instantiate this struct by parsing incoming BASM directives
	denseLayer := new(DenseLayer)
	// --HARDCODED

	for _, node := range denseLayer.nodes {
		if bmInputsCount != len(node.weights) {
			return nil, errors.New("Dense Generator: Weights do not match input size")
		}
	}

	basmCode := new(ProperStringBuilder)
	basmCode.WriteFmtStringln("%%meta bmdef	global  registersize:%d", dataType.GetSize())
	basmCode.WriteStringln("%meta bmdef global iomode:async")

	// CP Name -> List of fragments to collapse
	cpList := make(map[string][]string)

	// Building the layer one node at a time
	for nodeIdx, node := range denseLayer.nodes {
		// Bookkeeping names of the weight fragments
		weightsFragments := make([]string, len(node.weights))

		// Inputs X Weights multiplications
		for weightIdx, weight := range node.weights {
			weightFragmentName := fmt.Sprintf("weight%d_%d", nodeIdx, weightIdx)
			weightsFragments[weightIdx] = weightFragmentName

			basmCode.WriteFmtStringln(
				"%%meta fidef %s fragment:weight, weight: %s, multop: %s",
				weightFragmentName, metaFloatLiteral(float64(weight), dataType), multop,
			)
			cpList[weightFragmentName] = []string{weightFragmentName}

			weightInputLinkName := weightFragmentName + "_input"
			basmCode.WriteFmtStringln("%%meta filinkdef %s type:fl", weightInputLinkName)
			basmCode.WriteFmtStringln(
				"%%meta filinkatt %s fi:ext, type:input, index:%d",
				weightInputLinkName, weightIdx,
			)
			basmCode.WriteFmtStringln(
				"%%meta filinkatt %s fi:%s, type:input, index:0",
				weightInputLinkName, weightFragmentName,
			)
		}

		// Summation and bias
		summationFragmentName := "sum"
		basmCode.WriteFmtStringln(
			"%%meta fidef %s fragment:summation, inputs:%d, bias:%s, addop:%s",
			summationFragmentName, len(node.weights), metaFloatLiteral(float64(node.bias), dataType), addop,
		)
		cpList[summationFragmentName] = []string{summationFragmentName}
		for weightIdx, weightFrag := range weightsFragments {
			// Linking each 'input X weight' result to an input of the summation fragment
			inputLinkName := fmt.Sprintf("%s_input%d_%d", summationFragmentName, nodeIdx, weightIdx)
			basmCode.WriteFmtStringln("%%meta filinkdef %s type:fl", inputLinkName)

			basmCode.WriteFmtStringln(
				"%%meta filinkatt %s fi:%s, type:output, index:0",
				inputLinkName, weightFrag,
			)
			basmCode.WriteFmtStringln(
				"%%meta filinkatt %s fi:%s, type:input, index:%d",
				inputLinkName, summationFragmentName, weightIdx,
			)
		}

		// Activation function
		// FIXME: Using a RELU is for prototyping purposes, obviously this should be generalized
		activationFragmentName := "activation"
		basmCode.WriteFmtStringln("%%meta fidef %s fragment:relu", activationFragmentName)
		cpList[activationFragmentName] = []string{activationFragmentName}

		activationInputLinkName := activationFragmentName + "_input"
		basmCode.WriteFmtStringln("%%meta filinkdef %s type:fl", activationInputLinkName)
		basmCode.WriteFmtStringln(
			"%%meta filinkatt %s fi:%s, type:output, index:0",
			activationInputLinkName,
			summationFragmentName,
		)
		basmCode.WriteFmtStringln(
			"%%meta filinkatt %s fi:%s, type:input, index:0",
			activationInputLinkName, activationFragmentName,
		)

		// Output node result
		outputLinkName := "output" + string(nodeIdx)
		basmCode.WriteFmtStringln("%%meta filinkdef %s type:fl", outputLinkName)
		basmCode.WriteFmtStringln(
			"%%meta filinkatt %s fi:%s, type:output, index:0",
			outputLinkName, activationFragmentName,
		)
		basmCode.WriteFmtStringln(
			"%%meta filinkatt %s fi:ext, type:output, index:%d",
			outputLinkName, nodeIdx,
		)
	}

	// Processors instantiation
	for cpName, frags := range cpList {
		// FIXME: With which character should the collapsing fragment names be separated by?
		basmCode.WriteFmtStringln("%%meta cpdef %s fragcollapse:%s", cpName, strings.Join(frags, ":"))
	}

	// Offloading the actual compilation to BasmGenerator
	l.BasmMeta = l.BasmMeta.AddMeta("basmcode", basmCode.String())
	return BasmGenerator(b, e, l)
}
