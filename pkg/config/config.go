package config

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

func ParseConfig(file string) (*GeneratorConfig, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	config := &GeneratorConfig{}
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	for _, c := range config.Counters {
		if c.Number == 0 {
			c.Number = 1
		}
	}
	for _, g := range config.Gauges {
		if g.Number == 0 {
			g.Number = 1
		}
		if g.Range.Upper == 0 {
			g.Range.Upper = 100
		}
	}
	return config, nil
}

type GeneratorConfig struct {
	Counters []*CounterConfig `yaml:"counters,omitempty"`
	Gauges   []*GaugeConfig   `yaml:"gauges,omitempty"`
}

type CounterConfig struct {
	Prefix string         `yaml:"prefix"`
	Number  int            `yaml:"number,omitempty"`
	Labels []LabelSetting `yaml:"labels,omitempty"`
}

type GaugeConfig struct {
	Prefix string         `yaml:"prefix"`
	Number  int            `yaml:"number,omitempty"`
	Labels []LabelSetting `yaml:"labels,omitempty"`
	Range  GaugeRange     `yaml:"range,omitempty"`
}
type LabelSetting struct {
	Name     string   `yaml:"name"`
	ValueSet []string `yaml:"valueset"`
}
type GaugeRange struct {
	Upper int64 `yaml:"upper"`
	Lower int64 `yaml:"lower"`
}
