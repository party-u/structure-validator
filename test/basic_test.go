package test

import (
	. "github.com/party-u/structure-validator"
	"testing"
	"time"
)

func TestBasicValidator(t *testing.T) {
	maxLengthRule := Rule[string]{
		Validate: func(v string) *RuleError {
			if len(v) > 8 {
				return &RuleError{
					Message:    "too big max length",
					IsCritical: true,
					Cause:      nil,
					Errors:     nil,
				}
			}
			return nil
		},
		Metadata: RuleMetadata{
			Name: "max_length_string",
		},
	}

	conf := EngineConfiguration(
		WithMaxRules(10),
		WithTimeout(time.Second),
	)

	anz := NewRuleEngine(conf, maxLengthRule)

	result := NewRuleValidator(anz).Analyze("aaaaaaaaa")

	if len(result) != 1 {
		t.Errorf("expect result length 1, got %d", len(result))
	}

}
