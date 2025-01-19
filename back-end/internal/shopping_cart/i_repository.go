package shopping_cart

type IShoppingCartRepository interface {
	// Returns the created ShoppingCart
	Create(ShoppingCartParams) (ShoppingCartModel, error)

	// Returns the found deliveriees
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter ShoppingCartParams) ([]ShoppingCartModel, error)

	// Returns the found ShoppingCart
	GetOne(id string) (ShoppingCartModel, error)

	// Returns amount of deleted deliveries
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted deliveries
	DeleteAll(filter ShoppingCartParams) (uint, error)

	// Returns the updated ShoppingCart
	Update(id string, newShoppingCart ShoppingCartParams) (ShoppingCartModel, error)
}
