package test

import (
	"fmt"
	structure_validator "github.com/party-u/structure-validator"
	"testing"
)

func TestComplexStructureTest(t *testing.T) {

	type nestedStruct struct {
		Value int
	}

	type mainStruct struct {
		Name   string
		Nested nestedStruct
	}

	value := mainStruct{
		Name: "Main",
		Nested: nestedStruct{
			Value: 0,
		},
	}

	nestedRule := structure_validator.Rule[nestedStruct]{
		Validate: func(nestedStruct nestedStruct) *structure_validator.RuleError {
			if nestedStruct.Value != 1 {
				return &structure_validator.RuleError{
					Message: "Invalid value",
				}
			}
			return nil

		},
	}

	mainRuleName := structure_validator.Rule[mainStruct]{
		Validate: func(mainStruct mainStruct) *structure_validator.RuleError {
			if len(mainStruct.Name) == 0 {
				return &structure_validator.RuleError{
					Message: "name is empty",
				}
			}
			return nil
		},
	}

	mainNestedStructRule := structure_validator.Rule[mainStruct]{
		Validate: func(mainStruct mainStruct) *structure_validator.RuleError {

			conf := structure_validator.EngineConfiguration(
				structure_validator.WithMaxRules(1))

			engine := structure_validator.NewRuleEngine(conf, nestedRule)

			result := structure_validator.NewRuleValidator(engine).Analyze(mainStruct.Nested)

			if len(result) > 0 {
				return &structure_validator.RuleError{
					Message:    "nested struct validation failed",
					IsCritical: false,
					Cause:      nil,
					Errors:     result,
				}

			}
			return nil

		},
	}

	conf := structure_validator.EngineConfiguration(
		structure_validator.WithMaxRules(2))

	engine := structure_validator.NewRuleEngine(conf, mainRuleName, mainNestedStructRule)

	result := structure_validator.NewRuleValidator(engine).Analyze(value)

	fmt.Println(result[0])

}
