package delivery_product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type DeliveryProductController struct {
	// TODO
	repository IDeliveryProductRepository
}

// Used for testing
func (c *DeliveryProductController) SetRepository(repo IDeliveryProductRepository) {
	c.repository = repo
}

func NewDeliveryProductController() *DeliveryProductController {
	return &DeliveryProductController{repository: NewMySQLDeliveryRepository()}
}

func (c *DeliveryProductController) Create(w http.ResponseWriter, r *http.Request) {
	var deliveryParam DeliveryProductParams
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

func (c *DeliveryProductController) GetAll(w http.ResponseWriter, r *http.Request) {
	var deliveryParams DeliveryProductParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid. IT ONLY CHECKS FOR ORDER_ID
	for key := range queryParams {
		if strings.ToLower(key) != "deliveryid" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if deliveryID := queryParams.Get("DeliveryID"); deliveryID != "" {
		*deliveryParams.DeliveryID = deliveryID
	}

	// It now passes the delivery param as a "filter" and gets the found deliveries
	foundDeliveryes, err := c.repository.GetAll(deliveryParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundDeliveries := []DeliveryProductDTO{}
	for i := 0; i < len(foundDeliveryes); i++ {
		dtoFoundDeliveries = append(dtoFoundDeliveries, foundDeliveryes[i].ToDTO())
	}
	// Returns the DTO deliveries
	if err := json.NewEncoder(w).Encode(dtoFoundDeliveries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *DeliveryProductController) GetOne(w http.ResponseWriter, r *http.Request) {
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

func (c *DeliveryProductController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var deliveryParams DeliveryProductParams
	queryParams := r.URL.Query()

	if deliveryID := queryParams.Get("AddressID"); deliveryID != "" {
		*deliveryParams.DeliveryID = deliveryID
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

func (c *DeliveryProductController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *DeliveryProductController) Update(w http.ResponseWriter, r *http.Request) {
	// This needs to be implemented because of the interface, altough it won't
	// be used since the delivery-product shouldn't be updated
}
