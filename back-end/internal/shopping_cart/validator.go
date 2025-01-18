package shopping_cart

import (
	"errors"
)

type ShoppingCartValidator struct{}

// Helper functions

func (v *ShoppingCartValidator) Validate(shoppingCart ShoppingCartParams) error {
	if shoppingCart.UserID == nil {
		return errors.New("Empty UserID")
	}
	return nil
}
