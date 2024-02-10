package types

import "errors"

var (
	ErrorInvalidConfig = errors.New("invalid jwt-filter configuration")
	ErrorIllegalArgument = errors.New("illegal argument")
	ErrorRulesNotMatch = errors.New("no rules were matched")
)
