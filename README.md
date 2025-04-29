# Structure Validator

A lightweight, generic validation library for Go that makes complex structure validation simple and type-safe.

[![Go Report Card](https://goreportcard.com/badge/github.com/party-u/structure-validator)](https://goreportcard.com/report/github.com/party-u/structure-validator)
[![GoDoc](https://godoc.org/github.com/party-u/structure-validator?status.svg)](https://godoc.org/github.com/party-u/structure-validator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Overview

Structure Validator leverages Go's generics to provide a flexible and type-safe approach to validating complex data structures. It allows you to define reusable validation rules for any type and compose them to validate even the most complex nested structures with detailed error reporting.

## Installation

```bash
go get github.com/party-u/structure-validator
```

## Key Features

- **Type-Safe** - Built with Go generics for compile-time type checking
- **Flexible** - Define custom validation rules for any data type
- **Nestable** - Validate complex nested structures at any depth
- **Configurable** - Customize engine behavior with timeouts and rule limits
- **Detailed Error Reporting** - Rich error information with support for nested errors
- **Critical Errors** - Mark validation failures as critical to stop further processing

## Usage Examples

### Basic Validation

Simple string validation with a maximum length rule:

```go
package main

import (
    "fmt"
    sv "github.com/party-u/structure-validator"
    "time"
)

func main() {
    // Define a rule for maximum string length
    maxLengthRule := sv.Rule[string]{
        Validate: func(v string) *sv.RuleError {
            if len(v) > 8 {
                return &sv.RuleError{
                    Message:    "too big max length",
                    IsCritical: true,
                }
            }
            return nil
        },
        Metadata: sv.RuleMetadata{
            Name: "max_length_string",
        },
    }

    // Configure the validation engine
    conf := sv.EngineConfiguration(
        sv.WithMaxRules(10),
        sv.WithTimeout(time.Second),
    )

    // Create the rule engine with our rule
    engine := sv.NewRuleEngine(conf, maxLengthRule)

    // Validate a string
    result := sv.NewRuleValidator(engine).Analyze("aaaaaaaaa")

    // Check and print results
    if len(result) > 0 {
        fmt.Println(result[0].String())
    } else {
        fmt.Println("Validation passed!")
    }
}
```

### Complex Structure Validation

Validating nested structures with composition:

```go
package main

import (
    "fmt"
    sv "github.com/party-u/structure-validator"
)

func main() {
    // Define our structures
    type nestedStruct struct {
        Value int
    }

    type mainStruct struct {
        Name   string
        Nested nestedStruct
    }

    // Create an instance to validate
    value := mainStruct{
        Name: "Main",
        Nested: nestedStruct{
            Value: 0, // This will fail validation
        },
    }

    // Rule for nested structure
    nestedRule := sv.Rule[nestedStruct]{
        Validate: func(ns nestedStruct) *sv.RuleError {
            if ns.Value != 1 {
                return &sv.RuleError{
                    Message: "Invalid value",
                }
            }
            return nil
        },
    }

    // Rule for main structure's name
    mainNameRule := sv.Rule[mainStruct]{
        Validate: func(ms mainStruct) *sv.RuleError {
            if len(ms.Name) == 0 {
                return &sv.RuleError{
                    Message: "name is empty",
                }
            }
            return nil
        },
    }

    // Rule that validates the nested structure within main
    mainNestedRule := sv.Rule[mainStruct]{
        Validate: func(ms mainStruct) *sv.RuleError {
            // Create a nested validator just for the nested structure
            conf := sv.EngineConfiguration(sv.WithMaxRules(1))
            engine := sv.NewRuleEngine(conf, nestedRule)
            result := sv.NewRuleValidator(engine).Analyze(ms.Nested)

            if len(result) > 0 {
                return &sv.RuleError{
                    Message: "nested struct validation failed",
                    Errors:  result, // Include nested validation errors
                }
            }
            return nil
        },
    }

    // Create main validator with both rules
    conf := sv.EngineConfiguration(sv.WithMaxRules(2))
    engine := sv.NewRuleEngine(conf, mainNameRule, mainNestedRule)
    
    // Run validation
    result := sv.NewRuleValidator(engine).Analyze(value)
    
    if len(result) > 0 {
        for _, err := range result {
            fmt.Println(err)
        }
    } else {
        fmt.Println("Validation passed!")
    }
}
```

## Creating Rules

Rules are defined using the generic `Rule` struct, which takes a type parameter:

```go
Rule[T]{
    Validate: func(value T) *RuleError {
        // Validation logic here
        if !isValid(value) {
            return &RuleError{
                Message:    "Error description",
                IsCritical: false, // Set to true to stop processing other rules
                Cause:      nil,   // Optional underlying error
                Errors:     nil,   // Optional nested validation errors
            }
        }
        return nil // Return nil for valid values
    },
    Metadata: RuleMetadata{
        Name:     "rule_name", // Optional but recommended
        Priority: 0,           // Optional, higher priority rules execute first
    },
}
```

## Engine Configuration

Customize the validation engine with these options:

```go
conf := EngineConfiguration(
    WithMaxRules(10),   // Maximum number of rules per engine
    WithTimeout(time.Second), // Timeout for validation operations
)
```

## Error Handling

Validation results are returned as a slice of `RuleError` objects:

- Check the length of the result to determine if validation passed
- Access detailed error information through the `RuleError` properties
- Nested errors are accessible via the `Errors` field
- Use `IsCritical` to identify serious validation failures

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
