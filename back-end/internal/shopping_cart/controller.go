package shopping_cart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type ShoppingCartController struct {
	// TODO
	repository IShoppingCartRepository
}

// Used for testing
func (c *ShoppingCartController) SetRepository(repo IShoppingCartRepository) {
	c.repository = repo
}

func NewShoppingCartController() *ShoppingCartController {
	return &ShoppingCartController{repository: NewMySQLShoppingCartRepository()}
}

func (c *ShoppingCartController) Create(w http.ResponseWriter, r *http.Request) {
	var shoppingCartParam ShoppingCartParams
	err := json.NewDecoder(r.Body).Decode(&shoppingCartParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation

	createdShoppingCart, err := c.repository.Create(shoppingCartParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdShoppingCart.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ShoppingCartController) GetAll(w http.ResponseWriter, r *http.Request) {
	var shoppingCartParams ShoppingCartParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "isactive" &&
			strings.ToLower(key) != "isdeleted" &&
			strings.ToLower(key) != "createdat" &&
			strings.ToLower(key) != "userid" &&
			strings.ToLower(key) != "addressid" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if userID := queryParams.Get("AddressID"); userID != "" {
		*shoppingCartParams.UserID = userID
	}

	// It now passes the ShoppingCart param as a "filter" and gets the found deliveries
	foundShoppingCartes, err := c.repository.GetAll(shoppingCartParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundDeliveries := []ShoppingCartDTO{}
	for i := 0; i < len(foundShoppingCartes); i++ {
		dtoFoundDeliveries = append(dtoFoundDeliveries, foundShoppingCartes[i].ToDTO())
	}
	// Returns the DTO deliveries
	if err := json.NewEncoder(w).Encode(dtoFoundDeliveries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ShoppingCartController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	shoppingCart, err := c.repository.GetOne(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(shoppingCart.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *ShoppingCartController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var shoppingCartParams ShoppingCartParams
	queryParams := r.URL.Query()

	if userID := queryParams.Get("UserID"); userID != "" {
		*shoppingCartParams.UserID = userID
	}
	count, err := c.repository.DeleteAll(shoppingCartParams)
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

func (c *ShoppingCartController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *ShoppingCartController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var shoppingCartParams ShoppingCartParams
	err := json.NewDecoder(r.Body).Decode(&shoppingCartParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shoppingCart, err := c.repository.Update(id, shoppingCartParams)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(shoppingCart.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
