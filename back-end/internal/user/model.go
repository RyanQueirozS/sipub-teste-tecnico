package user

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
