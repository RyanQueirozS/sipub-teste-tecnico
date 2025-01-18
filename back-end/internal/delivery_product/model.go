package delivery_product

// This is what will be used to create/find/update the deliveryProduct model. The
// fields are used as pointers so they can be nullified
type DeliveryProductParams struct {
	OrderID       *string
	ProductID     *string
	ProductAmount *uint
}

type DeliveryProductDTO struct {
	Id            string `json:"Id"`
	OrderID       string `json:"OrderID"`
	ProductID     string `json:"ProductID"`
	ProductAmount uint   `json:"ProductAmount"`
}

type DeliveryProductModel struct {
	id            string // ID will be a uuid
	orderID       string
	productID     string
	productAmount uint
}

func (d *DeliveryProductModel) ToDTO() DeliveryProductDTO {
	dtoDelivery := DeliveryProductDTO{
		Id:            d.id,
		OrderID:       d.orderID,
		ProductID:     d.productID,
		ProductAmount: d.productAmount,
	}
	return dtoDelivery
}
