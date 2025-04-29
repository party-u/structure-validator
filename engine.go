package structure_validator

import (
	"sort"
	"time"
)

type Config func(*EngineConfig)

// RuleEngine defines a generic interface for adding validation rules that process an input of type T.
type RuleEngine[T any] interface {
	addRule(rule Rule[T])
	Rules() []Rule[T]
	Config() *EngineConfig
}

// ruleEngine is a generic type used for validating items of type T against a set of rules producing RuleError results.
// it holds a slice of rule functions and configuration for the analysis process.
type ruleEngine[T any] struct {
	r      []Rule[T]
	config *EngineConfig
}

// AddRule adds a validation rule to the ruleEngine, which processes items of type T and returns a RuleError if violated.
func (anz *ruleEngine[T]) addRule(rule Rule[T]) {
	if len(anz.r) >= anz.config.MaxRules {
		panic("exceeded maximum number of rules")
	}

	anz.r = append(anz.r, rule)
}

// EngineConfig defines the configuration settings for an ruleEngine, including maximum rules and timeout duration.
type EngineConfig struct {
	MaxRules int
	Timeout  time.Duration
}

func EngineConfiguration(config ...Config) *EngineConfig {
	configuration := &EngineConfig{
		MaxRules: 25,
		Timeout:  0,
	}
	for _, conf := range config {
		conf(configuration)
	}
	return configuration
}

func WithMaxRules(maxRules int) Config {
	return func(config *EngineConfig) {
		if maxRules <= 0 {
			maxRules = 1
		}
		config.MaxRules = maxRules
	}
}

// WithTimeout sets the timeout duration for EngineConfig. Defaults to 30 seconds if the provided duration is <= 0.
func WithTimeout(timeout time.Duration) Config {
	return func(config *EngineConfig) {
		if timeout < 0 {
			panic("invalid timeout")
		}
		config.Timeout = timeout
	}
}

// NewRuleEngine creates a new RuleEngine instance with a configurable set of rules for processing items of generic type T.
func NewRuleEngine[T any](config *EngineConfig, rules ...Rule[T]) RuleEngine[T] {
	anz := &ruleEngine[T]{
		r:      make([]Rule[T], 0, config.MaxRules),
		config: config,
	}

	for _, rule := range rules {
		anz.addRule(rule)

	}
	return anz
}

func (anz *ruleEngine[T]) Rules() []Rule[T] {
	return sortByPriority(anz.r)
}

func (anz *ruleEngine[T]) Config() *EngineConfig {
	return anz.config
}

func sortByPriority[T any](rules []Rule[T]) []Rule[T] {
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Priority <= rules[j].Priority
	})
	return rules
}
