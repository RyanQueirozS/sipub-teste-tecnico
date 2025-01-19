package shopping_cart

// This is what will be used to create/find/update the ShoppingCart model. The
// fields are used as pointers so they can be nullified
type ShoppingCartParams struct {
	UserID        *string
	ProductID     *string
	ProductAmount *uint
}

type ShoppingCartDTO struct {
	Id            string `json:"Id"`
	UserID        string `json:"UserID"`
	ProductID     string `json:"ProductID"`
	ProductAmount uint   `json:"ProductAmount"`
}

type ShoppingCartModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id            string // ID will be a uuid
	userID        string
	productID     string
	productAmount uint
}

func (d *ShoppingCartModel) ToDTO() ShoppingCartDTO {
	dtoShoppingCart := ShoppingCartDTO{
		Id:            d.id,
		UserID:        d.userID,
		ProductID:     d.productID,
		ProductAmount: d.productAmount,
	}
	return dtoShoppingCart
}
