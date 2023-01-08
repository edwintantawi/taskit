package validator

import "github.com/edwintantawi/taskit/internal/domain"

type validator struct{}

func New() *validator {
	return &validator{}
}

func (v *validator) Validate(validater domain.Validater) error {
	return validater.Validate()
}
