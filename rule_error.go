package structure_validator

import (
	"strconv"
	"strings"
)

// RuleError represents an error resulting from a rule violation, optionally containing nested errors.
type RuleError struct {
	Message    string
	IsCritical bool
	Cause      error
	Errors     []*RuleError
}

func (e *RuleError) String() string {
	if e == nil {
		return "<nil>"
	}

	var sb strings.Builder

	// Add critical flag if true
	if e.IsCritical {
		sb.WriteString("[CRITICAL] ")
	} else {
		sb.WriteString("[WARNING] ")
	}

	// Add main message
	sb.WriteString(e.Message)

	// Add cause if present
	if e.Cause != nil {
		sb.WriteString("\nCaused by: ")
		sb.WriteString(e.Cause.Error())
	}

	// Add nested errors if present
	if len(e.Errors) > 0 {
		sb.WriteString("\nNested errors:")
		for i, err := range e.Errors {
			if err != nil {
				sb.WriteString("\n  ")
				sb.WriteString(strconv.Itoa(i + 1))
				sb.WriteString(". ")
				sb.WriteString(err.String())
			}
		}
	}

	return sb.String()
}
