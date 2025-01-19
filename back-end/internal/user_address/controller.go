package user_address

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type UserAddressController struct {
	// TODO
	repository IUserAddressRepository
}

// Used for testing
func (c *UserAddressController) SetRepository(repo IUserAddressRepository) {
	c.repository = repo
}

func NewUserAddressController() *UserAddressController {
	return &UserAddressController{repository: NewMySQLUserAddressRepository()}
}

func (c *UserAddressController) Create(w http.ResponseWriter, r *http.Request) {
	var userAddressParam UserAddressParams
	err := json.NewDecoder(r.Body).Decode(&userAddressParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	{ // Validation is simpler because there are only two fields
		if userAddressParam.UserID == "" || userAddressParam.AddressID == "" {
			http.Error(w, "Invalid userID or addressID", http.StatusInternalServerError)
		}
	}

	createdUserAddress, err := c.repository.Create(userAddressParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdUserAddress); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserAddressController) GetAll(w http.ResponseWriter, r *http.Request) {
	var userAddressParams UserAddressParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "userid" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	if userID := queryParams.Get("userID"); userID != "" {
		userAddressParams.UserID = userID
	}

	// NOTE: Why don't I just parse it from the json? Well, if the json is nil,
	// then it will generate an error and if the json isn't but a field is,
	// that is another possible error. This way, although repetitive, will make
	// it simple to understand. Where as having multiple nested `if`s might not

	// It now passes the userAddress param as a "filter" and gets the found userAddresses
	foundUserAddresses, err := c.repository.GetAll(userAddressParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(foundUserAddresses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserAddressController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userAddress, err := c.repository.GetOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(userAddress); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserAddressController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var userAddressParams UserAddressParams
	queryParams := r.URL.Query()

	if userID := queryParams.Get("userID"); userID != "" {
		userAddressParams.UserID = userID
	}

	if addressID := queryParams.Get("addressID"); addressID != "" {
		userAddressParams.AddressID = addressID
	}

	count, err := c.repository.DeleteAll(userAddressParams)
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

func (c *UserAddressController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *UserAddressController) Update(w http.ResponseWriter, r *http.Request) {
	// Does nothing, but needs to be implemented since it is a IController interface
}
