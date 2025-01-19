package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DeliveryController struct {
	// TODO
	repository IDeliveryRepository
}

// Used for testing
func (c *DeliveryController) SetRepository(repo IDeliveryRepository) {
	c.repository = repo
}

func NewDeliveryController() *DeliveryController {
	return &DeliveryController{repository: NewMySQLDeliveryRepository()}
}

func (c *DeliveryController) Create(w http.ResponseWriter, r *http.Request) {
	var deliveryParam DeliveryParams
	err := json.NewDecoder(r.Body).Decode(&deliveryParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validation

	createdDelivery, err := c.repository.Create(deliveryParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdDelivery.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DeliveryController) GetAll(w http.ResponseWriter, r *http.Request) {
	var deliveryParams DeliveryParams
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
	if isActive := queryParams.Get("IsActive"); isActive != "" {
		value := isActive == "true"
		*deliveryParams.IsActive = value
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*deliveryParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*deliveryParams.CreatedAt = createdAt
	}

	if userID := queryParams.Get("UserID"); userID != "" {
		*deliveryParams.UserID = userID
	}

	if addressID := queryParams.Get("AddressID"); addressID != "" {
		*deliveryParams.AddressID = addressID
	}

	// It now passes the delivery param as a "filter" and gets the found deliveries
	foundDeliveryes, err := c.repository.GetAll(deliveryParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundDeliveries := []DeliveryDTO{}
	for i := 0; i < len(foundDeliveryes); i++ {
		dtoFoundDeliveries = append(dtoFoundDeliveries, foundDeliveryes[i].ToDTO())
	}
	// Returns the DTO deliveries
	if err := json.NewEncoder(w).Encode(dtoFoundDeliveries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DeliveryController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	delivery, err := c.repository.GetOne(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(delivery.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DeliveryController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var deliveryParams DeliveryParams
	queryParams := r.URL.Query()

	if isActive := queryParams.Get("IsActive"); isActive != "" {
		value := isActive == "true"
		*deliveryParams.IsActive = value
	}

	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*deliveryParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*deliveryParams.CreatedAt = createdAt
	}

	if userID := queryParams.Get("UserID"); userID != "" {
		*deliveryParams.UserID = userID
	}

	if addressID := queryParams.Get("AddressID"); addressID != "" {
		*deliveryParams.AddressID = addressID
	}
	count, err := c.repository.DeleteAll(deliveryParams)
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

func (c *DeliveryController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *DeliveryController) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var deliveryParams DeliveryParams
	err := json.NewDecoder(r.Body).Decode(&deliveryParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	delivery, err := c.repository.Update(id, deliveryParams)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // Did go through, none found
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(delivery.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
