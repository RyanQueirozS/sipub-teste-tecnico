package address

import (
	"errors"
)

type AddressValidator struct{}

func (v *AddressValidator) Validate(address AddressParams) error {
	if address.IsActive == nil {
		return errors.New("Invalid IsActive")
	}
	if address.IsDeleted == nil {
		return errors.New("Invalid IsDeleted")
	}

	if address.Street == nil {
		return errors.New("Invalid Street")
	}
	if *address.Street == "" {
		return errors.New("Street cannot be empty")
	}

	if address.Number == nil {
		return errors.New("Invalid Number")
	}
	if *address.Number == "" {
		return errors.New("Number cannot be empty")
	}

	if address.Neighborhood == nil {
		return errors.New("Invalid Neighborhood")
	}
	if *address.Neighborhood == "" {
		return errors.New("Neighborhood cannot be empty")
	}

	if address.City == nil {
		return errors.New("Invalid City")
	}
	if *address.City == "" {
		return errors.New("City cannot be empty")
	}

	if address.State == nil {
		return errors.New("Invalid State")
	}
	if *address.State == "" {
		return errors.New("State cannot be empty")
	}

	if address.Country == nil {
		return errors.New("Invalid Country")
	}
	if *address.Country == "" {
		return errors.New("Country cannot be empty")
	}

	if address.Latitude == nil {
		return errors.New("Invalid Latitude")
	}

	if address.Longitude == nil {
		return errors.New("Invalid Longitude")
	}

	return nil
}
