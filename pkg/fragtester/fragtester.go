package fragtester

import (
	"fmt"
)

type Config struct {
	DataType      string
	TypePrefix    string
	Params        map[string]string
	Debug         bool
	Verbose       bool
	NeuronLibPath string
}

func NewConfig() *Config {
	return &Config{
		DataType:      "float32",
		TypePrefix:    "0f",
		Params:        make(map[string]string),
		Debug:         false,
		Verbose:       false,
		NeuronLibPath: "",
	}
}

func (c *Config) AnalyzeFragment(fragment string) error {
	fmt.Println("Analyzing fragment:", fragment)
	return nil
}

func (c *Config) WriteBasm() (string, error) {
	result := fmt.Sprintf("%%meta bmdef     global registersize:\n")
	return result, nil
}
