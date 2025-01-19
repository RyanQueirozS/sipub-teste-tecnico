package product

import (
	"errors"
	"strings"
)

type ProductValidator struct{}

func (v *ProductValidator) Validate(params ProductParams) error {
	/*
		IsActive    *bool
		IsDeleted   *bool // Soft deletion
		WeightGrams *float32
		CreatedAt   *string // TODO
		Price       *float32
		Name        *string
	*/
	if params.IsActive == nil {
		return errors.New("IsActive is empty")
	}
	if params.IsDeleted == nil {
		return errors.New("IsDeleted is empty")
	}
	if params.Name == nil {
		return errors.New("Name is empty")
	}
	if strings.TrimSpace(*params.Name) == "" {
		return errors.New("Name is empty")
	}
	if params.WeightGrams == nil {
		return errors.New("Weight is invalid")
	}
	if *params.WeightGrams <= 0 {
		return errors.New("Weight is invalid")
	}
	if params.Price == nil {
		return errors.New("Price is invalid")
	}
	if *params.Price <= 0 {
		return errors.New("Price is invalid")
	}

	return nil
}
