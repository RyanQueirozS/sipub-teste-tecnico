package product

// This is what will be used to create/find/update the product model. The fields are used as pointers so they can be nullified
type ProductParams struct {
	IsActive    *bool
	IsDeleted   *bool // Soft deletion
	WeightGrams *float32
	CreatedAt   *string // TODO
	Price       *float32
	Name        *string
}

// This is what will be passed to user about the product model.
type ProductDTO struct {
	Id            string  `json:"Id"`
	CreatedAt     string  `json:"CreatedAt"`
	WeightInGrams float32 `json:"WeightInGrams"`
	Price         float32 `json:"Price"`
	Name          string  `json:"Name"`
}

type ProductModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id        string // ID will be a uuid
	isActive  bool
	isDeleted bool // Soft deletion
	createdAt string

	// Weight and price per product
	weightGrams float32
	price       float32
	name        string
}

func (p *ProductModel) ToDTO() ProductDTO {
	dtoProduct := ProductDTO{Id: p.id, CreatedAt: p.createdAt, WeightInGrams: p.weightGrams, Price: p.price, Name: p.name}
	return dtoProduct
}

func (p *ProductModel) GetWeight() float32 {
	return p.weightGrams
}

func (p *ProductModel) GetPrice() float32 {
	return p.price
}

func (p *ProductModel) GetName() string {
	return p.name
}

func (p *ProductModel) GetIsActive() bool {
	return p.isActive
}

func (p *ProductModel) SetWeight(newWeight float32) {
	p.price = newWeight
}

func (p *ProductModel) SetPrice(newPrice float32) {
	p.price = newPrice
}

func (p *ProductModel) SetName(newName string) {
	p.name = newName
}
