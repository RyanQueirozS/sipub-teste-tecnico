package user_delivery

// This is what will be used to create/find/update the userDelivery model. The
// fields are used as pointers so they can be nullified
type UserDeliveryParams struct {
	DeliveryID *string
	UserID     *string
}

type UserDeliveryDTO struct {
	Id         string `json:"Id"`
	DeliveryID string `json:"DeliveryID"`
	UserID     string `json:"UserID"`
}

type UserDeliveryModel struct {
	id         string // ID will be a uuid
	deliveryID string
	userID     string
}

func (d *UserDeliveryModel) ToDTO() UserDeliveryDTO {
	dtoDelivery := UserDeliveryDTO{
		Id:         d.id,
		DeliveryID: d.deliveryID,
		UserID:     d.userID,
	}
	return dtoDelivery
}
