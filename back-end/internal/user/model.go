package user

// I know that there is a lot of code repetition, and there is a possibility of
// just letting the main model to have all of it's fields public. This code
// repeats itself often because of the no inheritance that golang provides, not
// that it is a bad thing - I personally don't like it - it's just that some
// strategies wont work.

// This is what will be used to create/find/update the user model. The
// fields are used as pointers so they can be nullified
type UserParams struct {
	IsActive  *bool
	IsDeleted *bool
	CreatedAt *string
	Email     *string
	Cpf       *string
	Name      *string
}

type UserDTO struct {
	Id        string `json:"Id"`
	CreatedAt string `json:"CreatedAt"`
	Email     string `json:"Email"`
	Cpf       string `json:"Cpf"`
	Name      string `json:"Name"`
}

type UserModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id        string // ID will be a uuid
	isActive  bool
	isDeleted bool // Soft deletion
	createdAt string

	// Weight and price per user
	email string
	cpf   string
	name  string
}

func (u *UserModel) ToDTO() UserDTO {
	dtoUser := UserDTO{Id: u.id, CreatedAt: u.createdAt, Email: u.email, Cpf: u.cpf, Name: u.name}
	return dtoUser
}

func (u *UserModel) GetIsActive() bool {
	return u.isActive
}
