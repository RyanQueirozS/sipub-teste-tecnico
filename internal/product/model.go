package product

import "time"

type ProductParams struct {
	WeightGrams float32
	Price       float32
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
}

func (p *ProductModel) GetWeight() float32 {
	return p.weightGrams
}

func (p *ProductModel) GetPrice() float32 {
	return p.price
}

func (p *ProductModel) SetWeight(newWeight float32) {
	p.price = newWeight
}

func (p *ProductModel) SetPrice(newPrice float32) {
	p.price = newPrice
}
