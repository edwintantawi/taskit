package validator

import "github.com/edwintantawi/taskit/internal/domain"

type Validator struct{}

func New() Validator {
	return Validator{}
}

func (v *Validator) Validate(validater domain.Validater) error {
	return validater.Validate()
}
