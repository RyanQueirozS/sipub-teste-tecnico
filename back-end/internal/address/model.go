package address

// This is what will be used to create/find/update the address model. The
// fields are used as pointers so they can be nullified
type AddressParams struct {
	IsActive     *bool
	IsDeleted    *bool
	CreatedAt    *string
	Street       *string
	Number       *string
	Neighborhood *string
	Complement   *string
	City         *string
	State        *string
	Country      *string
	Latitude     *float32
	Longitude    *float32

	Name *string
}

type AddressDTO struct {
	Id           string  `json:"Id"`
	CreatedAt    string  `json:"CreatedAt"`
	Street       string  `json:"Street"`
	Number       string  `json:"Number"`
	Neighborhood string  `json:"Neighborhood"`
	Complement   string  `json:"Complement"`
	City         string  `json:"City"`
	State        string  `json:"State"`
	Country      string  `json:"Country"`
	Latitude     float32 `json:"Latitude"`
	Longitude    float32 `json:"Longitude"`

	Name string `json:"Name"`
}

type AddressModel struct {
	// Base of db models, included here because go doesn't allow for
	// inheritance. Explained in COMMENTS.md
	id        string // ID will be a uuid
	isActive  bool
	isDeleted bool // Soft deletion
	createdAt string

	street       string
	number       string
	neighborhood string
	complement   string
	city         string
	state        string
	country      string
	latitude     float32
	longitude    float32

	name string
}

func (a *AddressModel) ToDTO() AddressDTO {
	dtoAddress := AddressDTO{
		Id:           a.id,
		CreatedAt:    a.createdAt,
		Street:       a.street,
		Number:       a.number,
		Neighborhood: a.neighborhood,
		Complement:   a.complement,
		City:         a.city,
		State:        a.state,
		Country:      a.country,
		Latitude:     a.latitude,
		Longitude:    a.longitude,
		Name:         a.name,
	}
	return dtoAddress
}

func (a *AddressModel) GetID() string {
	return a.id
}

func (a *AddressModel) GetIsActive() bool {
	return a.isActive
}

func (a *AddressModel) GetStreet() string {
	return a.street
}

func (a *AddressModel) GetNumber() string {
	return a.number
}

func (a *AddressModel) GetNeighborhood() string {
	return a.neighborhood
}

func (a *AddressModel) GetComplement() string {
	return a.complement
}

func (a *AddressModel) GetCity() string {
	return a.city
}

func (a *AddressModel) GetState() string {
	return a.state
}

func (a *AddressModel) GetCountry() string {
	return a.country
}

func (a *AddressModel) GetLatitude() float32 {
	return a.latitude
}

func (a *AddressModel) GetLongitude() float32 {
	return a.longitude
}
