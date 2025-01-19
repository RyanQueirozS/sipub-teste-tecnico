package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type UserController struct {
	// TODO
	validator  UserValidator
	repository IUserRepository
}

// Used for testing
func (c *UserController) SetRepository(repo IUserRepository) {
	c.repository = repo
}

func NewUserController() *UserController {
	return &UserController{repository: NewMySQLUserRepository()}
}

func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var userParam UserParams
	err := json.NewDecoder(r.Body).Decode(&userParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := c.validator.Validate(userParam); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := c.repository.Create(userParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdUser.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	var userParams UserParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "isactive" &&
			strings.ToLower(key) != "isdeleted" &&
			strings.ToLower(key) != "createdAt" &&
			strings.ToLower(key) != "email" &&
			strings.ToLower(key) != "cpf" &&
			strings.ToLower(key) != "name" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if isActive := queryParams.Get("IsActive"); isActive != "" {
		isActiveValue, err := strconv.ParseBool(isActive)
		if err == nil {
			*userParams.IsActive = isActiveValue
		}
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		isDeletedValue, err := strconv.ParseBool(isDeleted)
		if err == nil {
			*userParams.IsDeleted = isDeletedValue
		}
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*userParams.CreatedAt = createdAt
	}

	if email := queryParams.Get("WeightGrams"); email != "" {
		*userParams.Email = email
	}

	if cpf := queryParams.Get("Price"); cpf != "" {
		*userParams.Name = cpf
	}

	if name := queryParams.Get("Name"); name != "" {
		*userParams.Name = name
	}

	// NOTE: Why don't I just parse it from the json? Well, if the json is nil,
	// then it will generate an error and if the json isn't but a field is,
	// that is another possible error. This way, although repetitive, will make
	// it simple to understand. Where as having multiple nested `if`s might not

	// It now passes the user param as a "filter" and gets the found users
	foundUsers, err := c.repository.GetAll(userParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundUsers := []UserDTO{}
	for i := 0; i < len(foundUsers); i++ {
		dtoFoundUsers = append(dtoFoundUsers, foundUsers[i].ToDTO())
	}
	// Returns the DTO users
	if err := json.NewEncoder(w).Encode(dtoFoundUsers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := c.repository.GetOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(user.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *UserController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var userParams UserParams
	queryParams := r.URL.Query()

	if isActive := queryParams.Get("IsActive"); isActive != "" {
		isActiveValue, err := strconv.ParseBool(isActive)
		if err == nil {
			*userParams.IsActive = isActiveValue
		}
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		isDeletedValue, err := strconv.ParseBool(isDeleted)
		if err == nil {
			*userParams.IsDeleted = isDeletedValue
		}
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*userParams.CreatedAt = createdAt
	}

	if email := queryParams.Get("WeightGrams"); email != "" {
		*userParams.Email = email
	}

	if cpf := queryParams.Get("Price"); cpf != "" {
		*userParams.Name = cpf
	}

	if name := queryParams.Get("Name"); name != "" {
		*userParams.Name = name
	}

	count, err := c.repository.DeleteAll(userParams)
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

func (c *UserController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var userParams UserParams
	err := json.NewDecoder(r.Body).Decode(&userParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := c.repository.Update(id, userParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
