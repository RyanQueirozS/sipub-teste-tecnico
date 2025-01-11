package product

type IProductRespository interface {
	// Returns the created product
	Create(ProductParams) ProductModel

	// Returns the found products
	GetAll() []ProductModel

	// Returns the found product
	GetOne() ProductModel

	// Returns amount of deleted products
	DeleteOne() uint

	// Returns amount of deleted products
	DeleteAll() uint

	// Returns the updated product
	Update() ProductModel
}
