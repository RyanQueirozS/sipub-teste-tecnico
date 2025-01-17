package address

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AddressController struct {
	// TODO
	validator  AddressValidator
	repository IAddressRepository
}

// Used for testing
func (c *AddressController) SetRepository(repo IAddressRepository) {
	c.repository = repo
}

func NewAddressController() *AddressController {
	return &AddressController{repository: NewMySQLAddressRepository()}
}

func (c *AddressController) Create(w http.ResponseWriter, r *http.Request) {
	var addressParam AddressParams
	err := json.NewDecoder(r.Body).Decode(&addressParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.validator.Validate(addressParam); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAddress, err := c.repository.Create(addressParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdAddress.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *AddressController) GetAll(w http.ResponseWriter, r *http.Request) {
	var addressParams AddressParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "isactive" &&
			strings.ToLower(key) != "isdeleted" &&
			strings.ToLower(key) != "createdAt" &&
			strings.ToLower(key) != "street" &&
			strings.ToLower(key) != "number" &&
			strings.ToLower(key) != "neighborhood" &&
			strings.ToLower(key) != "city" &&
			strings.ToLower(key) != "state" &&
			strings.ToLower(key) != "country" &&
			strings.ToLower(key) != "latitude" &&
			strings.ToLower(key) != "longitude" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if isActive := queryParams.Get("IsActive"); isActive != "" {
		value := isActive == "true"
		*addressParams.IsActive = value
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*addressParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*addressParams.CreatedAt = createdAt
	}

	if street := queryParams.Get("Street"); street != "" {
		*addressParams.Street = street
	}

	if number := queryParams.Get("Number"); number != "" {
		*addressParams.Number = number
	}

	if neighborhood := queryParams.Get("Neighborhood"); neighborhood != "" {
		*addressParams.Neighborhood = neighborhood
	}

	if complement := queryParams.Get("Complement"); complement != "" {
		*addressParams.Complement = complement
	}

	if city := queryParams.Get("City"); city != "" {
		*addressParams.City = city
	}

	if state := queryParams.Get("State"); state != "" {
		*addressParams.State = state
	}

	if country := queryParams.Get("Country"); country != "" {
		*addressParams.Country = country
	}

	if latitude := queryParams.Get("Latitude"); latitude != "" {
		lat, err := strconv.ParseFloat(latitude, 32)
		if err == nil {
			value := float32(lat)
			*addressParams.Latitude = value
		}
	}

	if longitude := queryParams.Get("Longitude"); longitude != "" {
		lon, err := strconv.ParseFloat(longitude, 32)
		if err == nil {
			value := float32(lon)
			*addressParams.Longitude = value
		}
	}

	if name := queryParams.Get("Name"); name != "" {
		*addressParams.Name = name
	}

	// It now passes the address param as a "filter" and gets the found addresss
	foundAddresses, err := c.repository.GetAll(addressParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundAddresses := []AddressDTO{}
	for i := 0; i < len(foundAddresses); i++ {
		dtoFoundAddresses = append(dtoFoundAddresses, foundAddresses[i].ToDTO())
	}
	// Returns the DTO addresss
	if err := json.NewEncoder(w).Encode(dtoFoundAddresses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *AddressController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	address, err := c.repository.GetOne(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(address.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *AddressController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var addressParams AddressParams
	queryParams := r.URL.Query()

	if isActive := queryParams.Get("IsActive"); isActive != "" {
		value := isActive == "true"
		*addressParams.IsActive = value
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*addressParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*addressParams.CreatedAt = createdAt
	}

	if street := queryParams.Get("Street"); street != "" {
		if addressParams.Street == nil {
			addressParams.Street = new(string)
		}
		*addressParams.Street = street
	}

	if number := queryParams.Get("Number"); number != "" {
		if addressParams.Number == nil {
			addressParams.Number = new(string)
		}
		*addressParams.Number = number
	}

	if neighborhood := queryParams.Get("Neighborhood"); neighborhood != "" {
		if addressParams.Neighborhood == nil {
			addressParams.Neighborhood = new(string)
		}
		*addressParams.Neighborhood = neighborhood
	}

	if complement := queryParams.Get("Complement"); complement != "" {
		if addressParams.Complement == nil {
			addressParams.Complement = new(string)
		}
		*addressParams.Complement = complement
	}

	if city := queryParams.Get("City"); city != "" {
		if addressParams.City == nil {
			addressParams.City = new(string)
		}
		*addressParams.City = city
	}

	if state := queryParams.Get("State"); state != "" {
		if addressParams.State == nil {
			addressParams.State = new(string)
		}
		*addressParams.State = state
	}

	if country := queryParams.Get("Country"); country != "" {
		if addressParams.Country == nil {
			addressParams.Country = new(string)
		}
		*addressParams.Country = country
	}

	if latitude := queryParams.Get("Latitude"); latitude != "" {
		lat, err := strconv.ParseFloat(latitude, 32)
		if err == nil {
			value := float32(lat)
			*addressParams.Latitude = value
		}
	}

	if longitude := queryParams.Get("Longitude"); longitude != "" {
		lon, err := strconv.ParseFloat(longitude, 32)
		if err == nil {
			value := float32(lon)
			*addressParams.Longitude = value
		}
	}

	if name := queryParams.Get("Name"); name != "" {
		if addressParams.Name == nil {
			addressParams.Name = new(string)
		}
		*addressParams.Name = name
	}

	count, err := c.repository.DeleteAll(addressParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AddressController) DeleteOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	count, err := c.repository.DeleteOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AddressController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var addressParams AddressParams
	err := json.NewDecoder(r.Body).Decode(&addressParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	address, err := c.repository.Update(id, addressParams)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(address.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
