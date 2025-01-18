package delivery_product

// This is what will be used to create/find/update the deliveryProduct model. The
// fields are used as pointers so they can be nullified
type DeliveryProductParams struct {
	DeliveryID       *string
	ProductID     *string
	ProductAmount *uint
}

type DeliveryProductDTO struct {
	Id            string `json:"Id"`
	DeliveryID       string `json:"DeliveryID"`
	ProductID     string `json:"ProductID"`
	ProductAmount uint   `json:"ProductAmount"`
}

type DeliveryProductModel struct {
	id            string // ID will be a uuid
	deliveryID       string
	productID     string
	productAmount uint
}

func (d *DeliveryProductModel) ToDTO() DeliveryProductDTO {
	dtoDelivery := DeliveryProductDTO{
		Id:            d.id,
		DeliveryID:       d.deliveryID,
		ProductID:     d.productID,
		ProductAmount: d.productAmount,
	}
	return dtoDelivery
}
