package delivery_product

import (
	"errors"
)

type DeliveryProductValidator struct{}

// Helper functions

func (v *DeliveryProductValidator) Validate(delivery DeliveryProductParams) error {
	if delivery.OrderID == nil {
		return errors.New("Empty OrderID")
	}
	return nil
}
