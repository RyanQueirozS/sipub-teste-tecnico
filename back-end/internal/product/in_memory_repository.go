//go:build ignore

// This was mainly used for testing the controller and developing the
// architecture due to it's simplicity.
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

func NewInMemoryRepository() *InMemoryProductRepository {
	return &InMemoryProductRepository{products: make(map[string]ProductModel)}
}

func (r *InMemoryProductRepository) Create(params ProductParams) (ProductModel, error) {
	product := ProductModel{
		id:       uuid.NewString(),
		isActive: true, isDeleted: false,
		createdAt:   time.Now(),
		weightGrams: *params.WeightGrams, price: *params.Price, name: *params.Name,
	}
	_, exists := r.products[product.id]
	if exists {
		return ProductModel{}, errors.New("Product id already exists")
	}
	r.products[product.id] = product

	return product, nil
}

func (r *InMemoryProductRepository) GetAll(filter ProductParams) ([]ProductModel, error) {
	var result []ProductModel

	for _, product := range r.products {
		// Check if the product matches the filter.
		if (filter.WeightGrams == nil || product.weightGrams == *filter.WeightGrams) &&
			(filter.Price == nil || product.price == *filter.Price) &&
			(filter.Name == nil || product.name == *filter.Name) {
			result = append(result, product)
		}
	}

	return result, nil
}

func (r *InMemoryProductRepository) GetOne(id string) (ProductModel, error) {
	product, exists := r.products[id]
	if !exists {
		return ProductModel{}, nil
	}
	return product, nil
}

func (r *InMemoryProductRepository) DeleteOne(id string) (uint, error) {
	return 0, nil
}

func (r *InMemoryProductRepository) DeleteAll(filter ProductParams) (uint, error) {
	return 0, nil
}

func (r *InMemoryProductRepository) Update(id string, newProduct ProductParams) (ProductModel, error) {
	return ProductModel{}, nil
}
