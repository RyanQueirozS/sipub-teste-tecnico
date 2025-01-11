//go:build exclude

// TODO For now

package product

import (
	models "sipub-test/internal/product/models"
)

// Implements the IProductRepository
type InMemoryProductRepository struct {
	// uuid -> string
	// models.ProductModel -> product
	products map[string]models.ProductModel
}

func (r *InMemoryProductRepository) Create(models.ProductParams) models.ProductModel {}

func (r *InMemoryProductRepository) GetAll() []models.ProductModel {}

func (r *InMemoryProductRepository) GetOne() models.ProductModel {}

func (r *InMemoryProductRepository) DeleteOne() uint {}

func (r *InMemoryProductRepository) DeleteAll() uint {}

func (r *InMemoryProductRepository) Update() product.ProductMode {}
