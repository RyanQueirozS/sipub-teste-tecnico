package address

type IAddressRepository interface {
	// Returns the created address
	Create(AddressParams) (AddressModel, error)

	// Returns the found addresses
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter AddressParams) ([]AddressModel, error)

	// Returns the found address
	GetOne(id string) (AddressModel, error)

	// Returns amount of deleted addresses
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted addresses
	DeleteAll(filter AddressParams) (uint, error)

	// Returns the updated address
	Update(id string, newAddress AddressParams) (AddressModel, error)
}
