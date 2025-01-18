package delivery

import (
	"errors"
)

type DeliveryValidator struct{}

// Helper functions

func (v *DeliveryValidator) Validate(delivery DeliveryParams) error {
	if delivery.UserID == nil {
		return errors.New("Empty UserID")
	}
	return nil
}
