package simbox

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
)

// GeneticConfig holds configuration for the genetic algorithm
type GeneticConfig struct {
	Debug            bool    // Enable debug output
	PopulationSize   int     // Number of individuals in the population
	Generations      int     // Number of generations to evolve
	MutationRate     float64 // Probability of mutation (0.0 to 1.0)
	CrossoverRate    float64 // Probability of crossover (0.0 to 1.0)
	ElitismCount     int     // Number of best individuals to preserve
	MinDelay         int32   // Minimum delay value for distributions
	MaxDelay         int32   // Maximum delay value for distributions
	DistributionSize int     // Number of delay entries per distribution
}

// Individual represents a single solution in the population
type Individual struct {
	SimDelays *SimDelays
	Fitness   float64
}

// Population represents a collection of individuals
type Population struct {
	Individuals []*Individual
	Config      GeneticConfig
	Opcodes     []string
}

// NewPopulation creates a new population with random individuals
func NewPopulation(opcodes []string, config GeneticConfig) *Population {
	pop := &Population{
		Individuals: make([]*Individual, config.PopulationSize),
		Config:      config,
		Opcodes:     opcodes,
	}

	for i := 0; i < config.PopulationSize; i++ {
		pop.Individuals[i] = pop.generateRandomIndividual()
	}

	return pop
}

// generateRandomIndividual creates a random SimDelays individual
func (p *Population) generateRandomIndividual() *Individual {
	simDelays := NewSimDelays()

	for _, opcode := range p.Opcodes {
		distr := make(DelayDistribution)

		// Generate random delay distribution
		for j := 0; j < p.Config.DistributionSize; j++ {
			delay := p.Config.MinDelay + rand.Int32N(p.Config.MaxDelay-p.Config.MinDelay+1)
			prob := rand.Float32()
			distr[delay] = prob
		}

		// Normalize the distribution
		distr.Normalize()
		simDelays.OpcodeDelays[opcode] = distr
	}

	return &Individual{
		SimDelays: simDelays,
		Fitness:   0.0,
	}
}

// Crossover performs crossover between two SimDelays and returns a new offspring
func (sd *SimDelays) Crossover(other *SimDelays, config GeneticConfig) *SimDelays {
	offspring := NewSimDelays()

	for opcode := range sd.OpcodeDelays {
		if rand.Float64() < config.CrossoverRate {
			// Inherit from this parent
			if distr, exists := sd.OpcodeDelays[opcode]; exists {
				offspring.OpcodeDelays[opcode] = copyDistribution(distr)
			}
		} else {
			// Inherit from other parent
			if distr, exists := other.OpcodeDelays[opcode]; exists {
				offspring.OpcodeDelays[opcode] = copyDistribution(distr)
			}
		}
	}

	// Handle opcodes that might be in other but not in this
	for opcode := range other.OpcodeDelays {
		if _, exists := offspring.OpcodeDelays[opcode]; !exists {
			if rand.Float64() < 0.5 {
				if distr, exists := other.OpcodeDelays[opcode]; exists {
					offspring.OpcodeDelays[opcode] = copyDistribution(distr)
				}
			}
		}
	}

	return offspring
}

// Mutate performs mutation on a SimDelays
func (sd *SimDelays) Mutate(config GeneticConfig) {
	for opcode, distr := range sd.OpcodeDelays {
		if rand.Float64() < config.MutationRate {
			// Mutate this distribution
			mutatedDistr := make(DelayDistribution)

			for delay, prob := range distr {
				// Randomly modify delay or probability
				if rand.Float64() < 0.5 {
					// Modify delay value
					newDelay := delay + (rand.Int32N(3) - 1) // -1, 0, or +1
					if newDelay < config.MinDelay {
						newDelay = config.MinDelay
					}
					if newDelay > config.MaxDelay {
						newDelay = config.MaxDelay
					}
					mutatedDistr[newDelay] = prob
				} else {
					// Modify probability
					newProb := prob + (rand.Float32()*0.2 - 0.1) // ±10%
					if newProb < 0 {
						newProb = 0
					}
					mutatedDistr[delay] = newProb
				}
			}

			// Maybe add a new random entry
			if rand.Float64() < 0.3 && len(mutatedDistr) < config.DistributionSize*2 {
				newDelay := config.MinDelay + rand.Int32N(config.MaxDelay-config.MinDelay+1)
				mutatedDistr[newDelay] = rand.Float32()
			}

			// Maybe remove an entry
			if rand.Float64() < 0.3 && len(mutatedDistr) > 1 {
				// Remove a random entry
				keys := make([]int32, 0, len(mutatedDistr))
				for k := range mutatedDistr {
					keys = append(keys, k)
				}
				if len(keys) > 0 {
					removeKey := keys[rand.IntN(len(keys))]
					delete(mutatedDistr, removeKey)
				}
			}

			// Normalize and update
			mutatedDistr.Normalize()
			sd.OpcodeDelays[opcode] = mutatedDistr
		}
	}
}

// copyDistribution creates a deep copy of a DelayDistribution
func copyDistribution(distr DelayDistribution) DelayDistribution {
	newDistr := make(DelayDistribution)
	for k, v := range distr {
		newDistr[k] = v
	}
	return newDistr
}

// FitnessFunc is the function signature for fitness evaluation
type FitnessFunc func(*SimDelays) float64

// BlankFitness is a placeholder fitness function that returns 0
// This should be replaced with actual fitness evaluation logic
func BlankFitness(sd *SimDelays) float64 {
	// TODO: Implement actual fitness evaluation
	// This function should evaluate how good the SimDelays configuration is
	// Higher values indicate better fitness
	return 0.0
}

// EvaluatePopulation evaluates the fitness of all individuals in the population
func (p *Population) EvaluatePopulation(fitnessFunc FitnessFunc) {
	for _, individual := range p.Individuals {
		individual.Fitness = fitnessFunc(individual.SimDelays)
		if p.Config.Debug {
			fmt.Printf("Individual Fitness: %.16f\n", individual.Fitness)
		}
	}
}

// SortByFitness sorts the population by fitness in descending order (best first)
func (p *Population) SortByFitness() {
	sort.Slice(p.Individuals, func(i, j int) bool {
		return p.Individuals[i].Fitness > p.Individuals[j].Fitness
	})
}

// SelectParent performs tournament selection to choose a parent
func (p *Population) SelectParent(tournamentSize int) *Individual {
	best := p.Individuals[rand.IntN(len(p.Individuals))]
	for i := 1; i < tournamentSize; i++ {
		competitor := p.Individuals[rand.IntN(len(p.Individuals))]
		if competitor.Fitness > best.Fitness {
			best = competitor
		}
	}
	return best
}

// Evolve runs the genetic algorithm for the configured number of generations
func (p *Population) Evolve(fitnessFunc FitnessFunc) *SimDelays {
	for generation := 0; generation < p.Config.Generations; generation++ {

		fmt.Printf("Generation %d: Best Fitness = %.16f, Average Fitness = %.16f\n",
			generation,
			p.GetBest().Fitness,
			p.GetAverageFitness())

		// Evaluate fitness
		p.EvaluatePopulation(fitnessFunc)

		// Sort by fitness
		p.SortByFitness()

		// Create new generation
		newIndividuals := make([]*Individual, p.Config.PopulationSize)

		// Elitism: preserve best individuals
		for i := 0; i < p.Config.ElitismCount && i < len(p.Individuals); i++ {
			newIndividuals[i] = &Individual{
				SimDelays: copySimDelays(p.Individuals[i].SimDelays),
				Fitness:   p.Individuals[i].Fitness,
			}
		}

		// Generate offspring
		for i := p.Config.ElitismCount; i < p.Config.PopulationSize; i++ {
			// Select parents
			parent1 := p.SelectParent(3) // Tournament size of 3
			parent2 := p.SelectParent(3)

			// Crossover
			offspring := parent1.SimDelays.Crossover(parent2.SimDelays, p.Config)

			// Mutation
			offspring.Mutate(p.Config)

			newIndividuals[i] = &Individual{
				SimDelays: offspring,
				Fitness:   0.0,
			}
		}

		p.Individuals = newIndividuals
	}

	// Final evaluation
	p.EvaluatePopulation(fitnessFunc)
	p.SortByFitness()

	// Return the best individual
	return copySimDelays(p.Individuals[0].SimDelays)
}

// copySimDelays creates a deep copy of a SimDelays
func copySimDelays(sd *SimDelays) *SimDelays {
	newSD := NewSimDelays()
	for opcode, distr := range sd.OpcodeDelays {
		newSD.OpcodeDelays[opcode] = copyDistribution(distr)
	}
	return newSD
}

// GetBest returns the best individual in the population
func (p *Population) GetBest() *Individual {
	p.SortByFitness()
	return p.Individuals[0]
}

// GetAverageFitness returns the average fitness of the population
func (p *Population) GetAverageFitness() float64 {
	total := 0.0
	for _, individual := range p.Individuals {
		total += individual.Fitness
	}
	return total / float64(len(p.Individuals))
}

// RunGeneticAlgorithm is a convenience function to run the genetic algorithm
// with the given opcodes, configuration, and fitness function
func RunGeneticAlgorithm(opcodes []string, config GeneticConfig, fitnessFunc FitnessFunc) (*SimDelays, *Population) {
	population := NewPopulation(opcodes, config)
	bestSimDelays := population.Evolve(fitnessFunc)
	return bestSimDelays, population
}

func GetDefaultGeneticConfig() GeneticConfig {
	return GeneticConfig{
		Debug:            false,
		PopulationSize:   50,
		Generations:      100,
		MutationRate:     0.1,
		CrossoverRate:    0.7,
		ElitismCount:     5,
		MinDelay:         1,
		MaxDelay:         10,
		DistributionSize: 5,
	}
}

func GetGeneticConfigFromJSON(fileName string) (GeneticConfig, error) {
	config := GetDefaultGeneticConfig()

	data, err := os.ReadFile(fileName)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (gc *GeneticConfig) ToJSONFile(fileName string) error {
	data, err := json.MarshalIndent(gc, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
