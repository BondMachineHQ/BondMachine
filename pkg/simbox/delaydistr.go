package simbox

import (
	"encoding/json"
	"math"
	"math/rand/v2"
	"os"
	"sort"
)

type DelayDistribution map[int32]float32 // map of delay in clocks to probability

type SimDelays struct {
	OpcodeDelays map[string]DelayDistribution // json:"opcode_delays"`
}

func NewSimDelays() *SimDelays {
	return &SimDelays{
		OpcodeDelays: make(map[string]DelayDistribution),
	}
}

func LoadSimDelaysFromFile(filename string) (*SimDelays, error) {
	simDelays := NewSimDelays()
	if filename != "" {
		if _, err := os.Stat(filename); err == nil {
			if sbJSON, err := os.ReadFile(filename); err == nil {
				if err := json.Unmarshal([]byte(sbJSON), simDelays); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}
	// Normalize all delay distributions
	for _, delayDistr := range simDelays.OpcodeDelays {
		delayDistr.Normalize()
	}
	return simDelays, nil
}

func (d *DelayDistribution) Normalize() {
	var total float32 = 0.0
	for _, prob := range *d {
		total += prob
	}
	if total > 0.0 {
		for delay, prob := range *d {
			(*d)[delay] = prob / total
		}
	}
}

func (d *DelayDistribution) GetValue() int32 {
	var cumulative float32 = 0.0
	var lastDelay int32 = 0
	var r float32 = rand.Float32()
	for delay, prob := range *d {
		cumulative += prob
		if r <= cumulative {
			return delay
		}
		lastDelay = delay
	}
	return lastDelay
}

func DistributionDistance(d1, d2 DelayDistribution) float64 {
	// Compute 1-Wasserstein distance (Earth Mover's Distance) between two probability distributions
	// This metric is more appropriate for comparing distributions as it considers the
	// "distance" between delay values, not just probability differences

	// Get all unique delay values
	delaySet := make(map[int32]struct{})
	for k := range d1 {
		delaySet[k] = struct{}{}
	}
	for k := range d2 {
		delaySet[k] = struct{}{}
	}

	// Convert to sorted slice
	delays := make([]int32, 0, len(delaySet))
	for k := range delaySet {
		delays = append(delays, k)
	}
	sort.Slice(delays, func(i, j int) bool { return delays[i] < delays[j] })

	// Compute Wasserstein distance using cumulative distributions
	var distance float64 = 0.0
	var cdf1, cdf2 float64 = 0.0, 0.0

	for i := 0; i < len(delays); i++ {
		delay := delays[i]

		// Add probabilities to cumulative distribution functions
		cdf1 += float64(d1[delay]) // Will be 0.0 if key doesn't exist
		cdf2 += float64(d2[delay]) // Will be 0.0 if key doesn't exist

		// Compute the step size for integration
		var stepSize float64
		if i < len(delays)-1 {
			stepSize = float64(delays[i+1] - delays[i])
		} else {
			// For the last element, we use step size of 1
			stepSize = 1.0
		}

		// Add the area between CDFs
		distance += math.Abs(cdf1-cdf2) * stepSize
	}

	return distance
}
