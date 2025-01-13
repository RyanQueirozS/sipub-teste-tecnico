package product

import "time"

type ProductParams struct {
	WeightGrams float32
	Price       float32
	Name        string
}

type ProductDTO struct {
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
	deletedAt time.Time
	createdAt time.Time

	// Weight and price per product
	weightGrams float32
	price       float32
	name        string
}

func (p *ProductModel) ToDTO() ProductDTO {
	dtoProduct := ProductDTO{WeightInGrams: p.weightGrams, Price: p.price, Name: p.name}
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

func (p *ProductModel) SetWeight(newWeight float32) {
	p.price = newWeight
}

func (p *ProductModel) SetPrice(newPrice float32) {
	p.price = newPrice
}

func (p *ProductModel) SetName(newName string) {
	p.name = newName
}
