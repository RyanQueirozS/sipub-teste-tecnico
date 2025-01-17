package user

type IUserRepository interface {
	// Returns the created user
	Create(UserParams) (UserModel, error)

	// Returns the found users
	// NOTE: reusing the same type for a filter and a "constructor" is not
	// ideal at all, but it will save on code repetition
	GetAll(filter UserParams) ([]UserModel, error)

	// Returns the found user
	GetOne(id string) (UserModel, error)

	// Returns amount of deleted users
	DeleteOne(id string) (uint, error)

	// Returns amount of deleted users
	DeleteAll(filter UserParams) (uint, error)

	// Returns the updated user
	Update(id string, newUser UserParams) (UserModel, error)
}
