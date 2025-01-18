package delivery_product

import (
	"errors"
)

type DeliveryProductValidator struct{}

// Helper functions

func (v *DeliveryProductValidator) Validate(delivery DeliveryProductParams) error {
	if delivery.DeliveryID == nil {
		return errors.New("Empty DeliveryID")
	}
	return nil
}
