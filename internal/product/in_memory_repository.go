//go:build exclude

// TODO For now

package product

// Implements the IProductRepository
type InMemoryProductRepository struct {
	// uuid -> string
	// ProductModel -> product
	products map[string]ProductModel
}

func (r *InMemoryProductRepository) Create(ProductParams) ProductModel {}

func (r *InMemoryProductRepository) GetAll() []ProductModel {}

func (r *InMemoryProductRepository) GetOne() ProductModel {}

func (r *InMemoryProductRepository) DeleteOne() uint {}

func (r *InMemoryProductRepository) DeleteAll() uint {}

func (r *InMemoryProductRepository) Update() product.ProductMode {}
