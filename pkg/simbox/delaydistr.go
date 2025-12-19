package simbox

import (
	"encoding/json"
	"math/rand/v2"
	"os"
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
	var distance float64 = 0.0
	// Get all unique keys
	keys := make(map[int32]struct{})
	for k := range d1 {
		keys[k] = struct{}{}
		if _, ok := d2[k]; !ok {
			d2[k] = 0.0
		}
	}
	for k := range d2 {
		keys[k] = struct{}{}
		if _, ok := d1[k]; !ok {
			d1[k] = 0.0
		}
	}
	// Compute distance
	for k := range keys {
		p1 := float64(d1[k])
		p2 := float64(d2[k])
		diff := p1 - p2
		distance += diff * diff
	}
	return distance
}
