package product

type IProductRepository interface {
	// Returns the created product
	Create(ProductParams) (ProductModel, error)

	// Returns the found products
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter ProductParams) ([]ProductModel, error)

	// Returns the found product
	GetOne(id string) (ProductModel, error)

	// Returns amount of deleted products
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted products
	DeleteAll(filter ProductParams) (uint, error)

	// Returns the updated product
	Update(id string, newProduct ProductParams) (ProductModel, error)
}
