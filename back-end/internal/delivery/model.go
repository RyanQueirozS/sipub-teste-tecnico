package delivery

// This is what will be used to create/find/update the delivery model. The
// fields are used as pointers so they can be nullified
type DeliveryParams struct {
	IsActive  *bool
	IsDeleted *bool
	CreatedAt *string

	UserID    *string
	AddressID *string
}

type DeliveryDTO struct {
	Id        string `json:"Id"`
	IsActive  bool   `json:"IsActive"`
	IsDeleted bool   `json:"IsDeleted"`
	CreatedAt string `json:"CreatedAt"`

	UserID    string `json:"UserID"`
	AddressID string `json:"AddressID"`
}

type DeliveryModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id        string // ID will be a uuid
	isActive  bool
	isDeleted bool // Soft deletion
	createdAt string

	userID    string
	addressID string
}

func (d *DeliveryModel) ToDTO() DeliveryDTO {
	dtoDelivery := DeliveryDTO{
		Id:        d.id,
		CreatedAt: d.createdAt,
		IsActive:  d.isActive,
		IsDeleted: d.isDeleted,
		UserID:    d.userID,
		AddressID: d.addressID,
	}
	return dtoDelivery
}
