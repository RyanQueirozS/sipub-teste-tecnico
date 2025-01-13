package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Implements the IProductRepository
type InMemoryProductRepository struct {
	// uuid : string
	// ProductModel : product
	products map[string]ProductModel
}

func (r *InMemoryProductRepository) Create(params ProductParams) (ProductModel, error) {
	product := ProductModel{
		id:       uuid.NewString(),
		isActive: true, isDeleted: false,
		createdAt:   time.Now(),
		weightGrams: params.WeightGrams, price: params.Price, name: params.Name,
	}
	_, exists := r.products[product.id]
	if exists {
		return ProductModel{}, errors.New("Product id already exists")
	}
	r.products[product.id] = product

	return product, nil
}

func (r *InMemoryProductRepository) GetAll() ([]ProductModel, error) {
	// Todo have a querystring
	return nil, nil
}

func (r *InMemoryProductRepository) GetOne() (ProductModel, error) {
	return ProductModel{}, nil
}

func (r *InMemoryProductRepository) DeleteOne() (uint, error) {
	return 0, nil
}

func (r *InMemoryProductRepository) DeleteAll() (uint, error) {
	return 0, nil
}

func (r *InMemoryProductRepository) Update() (ProductModel, error) {
	return ProductModel{}, nil
}
