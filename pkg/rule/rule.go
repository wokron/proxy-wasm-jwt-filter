package rule

import (
	"strings"
	"wasm-jwt-filter/pkg/types"
)

type Rule struct {
	Match    RouteMatch   `json:"match"`
	Requires *Requirement `json:"requires"`
}

type RuleList []Rule

func (rules *RuleList) Validate(path string, validateByName ValidateFunc) (bool, error) {
	require, err := rules.findFirstMatchedRule(path)
	if err != nil {
		return false, nil
	}
	return require.validate(validateByName)
}

func (rules *RuleList) findFirstMatchedRule(path string) (*Requirement, error) {
	for _, rule := range *rules {
		ok, err := rule.Match.isMatch(path)
		if err != nil {
			return nil, err
		}
		if ok {
			return rule.Requires, nil
		}
	}
	return nil, types.ErrorRulesNotMatch
}

type RouteMatch struct {
	Prefix *string `json:"prefix"`
	Path   *string `json:"path"`
}

func (match *RouteMatch) isMatch(path string) (bool, error) {
	if match.Path != nil {
		return *match.Path == path, nil
	} else if match.Prefix != nil {
		return strings.HasPrefix(path, *match.Prefix), nil
	} else {
		return false, types.ErrorInvalidConfig
	}
}

type Requirement struct {
	NameRequired *string        `json:"name"`
	RequiresAny  *[]Requirement `json:"requires_any"`
	RequiresAll  *[]Requirement `json:"requires_all"`
}

type ValidateFunc func(name string) bool

func (require *Requirement) validate(validateByName ValidateFunc) (bool, error) {
	if require == nil {
		return true, nil
	}

	if require.NameRequired != nil {
		ok := validateByName(*require.NameRequired)
		return ok, nil
	} else if require.RequiresAll != nil {
		return validateRequiresAll(require.RequiresAll, validateByName)
	} else if require.RequiresAny != nil {
		return validateRequiresAny(require.RequiresAny, validateByName)
	} else {
		return false, types.ErrorInvalidConfig
	}
}

func validateRequiresAny(requiresAnyList *[]Requirement, validateByName ValidateFunc) (bool, error) {
	for _, require := range *requiresAnyList {
		ok, err := require.validate(validateByName)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func validateRequiresAll(requiresAllList *[]Requirement, validateByName ValidateFunc) (bool, error) {
	for _, require := range *requiresAllList {
		ok, err := require.validate(validateByName)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
