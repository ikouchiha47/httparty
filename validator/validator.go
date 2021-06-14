package validator

import (
	"httparty/engine"
	"httparty/errorrs"
)

type Validator interface {
	Validate(fleet engine.Fleet) error
}

type HttpValidator struct {
}

func NewHttpValidator() HttpValidator {
	return HttpValidator{}
}

func (hv HttpValidator) Validate(fleet engine.Fleet) error {
	if len(fleet.Steps) <= fleet.ExitPoint {
		return errorrs.IndexOutOfRange
	}

	stepper := map[string]map[string]bool{}

	for _, step := range fleet.Steps {
		if _, ok := stepper[step.Name]; ok {
			return errorrs.DuplicateEntry
		}

		stepper[step.Name] = toBoolMap(step.Requires)
	}

	for ruleKey, requireMap := range stepper {
		for key, _ := range requireMap {
			if _, ok := stepper[key][ruleKey]; ok {
				return errorrs.CyclicDependency
			}
		}
	}

	return nil
}

func toBoolMap(array []string) map[string]bool {
	arrMap := map[string]bool{}

	for _, ar := range array {
		arrMap[ar] = true
	}

	return arrMap
}
