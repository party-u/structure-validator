package example

import (
	"fmt"
	. "github.com/party-u/structure-validator"
	"time"
)

func simpleExample() {

	rule1 := Rule[string]{
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

	anz := NewAnalyzer(conf, rule1)

	result := NewRuleValidator(anz).Analyze("aaaaaaaaa")

	fmt.Println(result[0].String())

}
