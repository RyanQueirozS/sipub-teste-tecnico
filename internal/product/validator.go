package product

import (
	"errors"
	"strings"
)

type ProductValidator struct{}

func (v *ProductValidator) Validate(params ProductParams) error {
	if strings.TrimSpace(params.Name) == "" {
		return errors.New("Name is empty")
	}
	if params.WeightGrams <= 0 {
		return errors.New("Weight is invalid")
	}
	if params.Price <= 0 {
		return errors.New("Price is invalid")
	}

	return nil
}
