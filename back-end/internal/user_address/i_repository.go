package user_address

type IUserAddressRepository interface {
	// Returns the created user
	Create(UserAddressParams) (UserAddressModel, error)

	// Returns the found users
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter UserAddressParams) ([]UserAddressModel, error)

	// Returns the found user
	GetOne(id string) (UserAddressModel, error)

	// Returns amount of deleted users
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted users
	DeleteAll(filter UserAddressParams) (uint, error)

	// Won't be used
	// Update(id string, newUserAddress UserAddressParams) (UserAddressModel, error)
}
