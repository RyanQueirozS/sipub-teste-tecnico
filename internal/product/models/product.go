package product

type ProductParams struct {
	WeightGrams float32
	Price       float32
}

type ProductModel struct {
	weightGrams float32
	price       float32
}

func (p *ProductModel) GetWeight() float32 {
	return p.weightGrams
}

func (p *ProductModel) GetPrice() float32 {
	return p.price
}
