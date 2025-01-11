package product

import models "sipub-test/internal/product/models"

type IProductRespository interface {
	// Returns the created product
	Create(models.ProductParams) models.ProductModel

	// Returns the found products
	GetAll() []models.ProductModel

	// Returns the found product
	GetOne() models.ProductModel

	// Returns amount of deleted products
	DeleteOne() uint

	// Returns amount of deleted products
	DeleteAll() uint

	// Returns the updated product
	Update() models.ProductModel
}
