package structure_validator

type Rule[T any] struct {
	Validate func(T) *RuleError
	Priority int
	Metadata RuleMetadata
}

type RuleMetadata struct {
	Name        string
	Description string
}
