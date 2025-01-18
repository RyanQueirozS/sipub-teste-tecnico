package payment

// This is what will be used to create/find/update the payment model. The
// fields are used as pointers so they can be nullified
type PaymentParams struct {
	IsDeleted  *bool
	CreatedAt  *string
	DeliveryID *string
	Value      *float32
	UserID     *string // This is what will be used for getAll/deleteAll
}

type PaymentDTO struct {
	Id         string  `json:"Id"`
	IsDeleted  bool    `json:"IsDeleted"`
	CreatedAt  string  `json:"CreatedAt"`
	DeliveryID string  `json:"OrderID"`
	Value      float32 `json:"Value"`
}

type PaymentModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id         string // ID will be a uuid
	isDeleted  bool   // Soft deletion
	createdAt  string
	deliveryID string
	value      float32
}

func (a *PaymentModel) ToDTO() PaymentDTO {
	dtoPayment := PaymentDTO{
		Id:         a.id,
		CreatedAt:  a.createdAt,
		DeliveryID: a.deliveryID,
		Value:      a.value,
	}
	return dtoPayment
}
