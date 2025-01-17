package user_address

// This is what will be used to create/find/update the userAddress model. The
// fields are used as pointers so they can be nullified
type UserAddressParams struct {
	UserID    string
	AddressID string
}

type UserAddressModel struct {
	id        string
	UserID    string `json:"UserID"`
	AddressID string `json:"AddressID"`
}

func (u *UserAddressModel) GetID() string {
	return u.id
}
