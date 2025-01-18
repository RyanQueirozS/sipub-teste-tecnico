package user_delivery

import (
	"errors"
)

type UserDeliveryValidator struct{}

// Helper functions

func (v *UserDeliveryValidator) Validate(delivery UserDeliveryParams) error {
	if delivery.DeliveryID == nil {
		return errors.New("Empty DeliveryID")
	}
	return nil
}
