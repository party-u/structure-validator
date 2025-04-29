package structure_validator

// Rule represents a validation rule that processes an input of type T and returns a RuleError if the rule is violated.
type Rule[T any] struct {
	Validate func(T) *RuleError
	Priority int
	Metadata RuleMetadata
}

type RuleMetadata struct {
	Name        string
	Description string
}
