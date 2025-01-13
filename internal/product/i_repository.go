package product

type IProductRepository interface {
	// Returns the created product
	Create(ProductParams) (ProductModel, error)

	// Returns the found products
	GetAll() ([]ProductModel, error)

	// Returns the found product
	GetOne() (ProductModel, error)

	// Returns amount of deleted products
	DeleteOne() (uint, error)

	// Returns amount of deleted products
	DeleteAll() (uint, error)

	// Returns the updated product
	Update() (ProductModel, error)
}
