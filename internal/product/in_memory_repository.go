package product

import (
	"time"

	"github.com/google/uuid"
)

// Implements the IProductRepository
type InMemoryProductRepository struct {
	// uuid : string
	// ProductModel : product
	products map[string]ProductModel
}

func (r *InMemoryProductRepository) Create(params ProductParams) ProductModel {
	product := ProductModel{
		id:       uuid.NewString(),
		isActive: true, isDeleted: false,
		createdAt:   time.Now(),
		weightGrams: params.WeightGrams, price: params.Price,
	}
	r.products[product.id] = product

	return product
}

func (r *InMemoryProductRepository) GetAll() []ProductModel {
	// Todo have a querystring
	return nil
}

func (r *InMemoryProductRepository) GetOne() ProductModel {
	return ProductModel{}
}

func (r *InMemoryProductRepository) DeleteOne() uint {
	return 0
}

func (r *InMemoryProductRepository) DeleteAll() uint {
	return 0
}

func (r *InMemoryProductRepository) Update() ProductModel {
	return ProductModel{}
}
