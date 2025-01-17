package user

import (
	"errors"
)

type UserValidator struct{}

// Helper functions

func (v *UserValidator) Validate(user UserParams) error {
	if user.Name == nil {
		return errors.New("Invalid name")
	}
	if *user.Name == "" {
		return errors.New("invalid Name")
	}
	if user.Email == nil {
		return errors.New("Invalid Price")
	}
	if *user.Email == "" {
		return errors.New("Invalid Price")
	}
	if user.Cpf == nil {
		return errors.New("Invalid WeightGrams")
	}
	if *user.Cpf <= "" {
		return errors.New("Invalid WeightGrams")
	}
	if user.IsActive == nil {
		return errors.New("Invalid IsActive")
	}
	if user.IsDeleted == nil {
		return errors.New("Invalid IsDeleted")
	}
	return nil
}
