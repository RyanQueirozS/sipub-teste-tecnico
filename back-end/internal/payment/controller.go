package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PaymentController struct {
	// TODO
	repository IPaymentRepository
}

// Used for testing
func (c *PaymentController) SetRepository(repo IPaymentRepository) {
	c.repository = repo
}

func NewPaymentController() *PaymentController {
	return &PaymentController{repository: NewMySQLPaymentRepository()}
}

func (c *PaymentController) Create(w http.ResponseWriter, r *http.Request) {
	var paymentParam PaymentParams
	err := json.NewDecoder(r.Body).Decode(&paymentParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdPayment, err := c.repository.Create(paymentParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdPayment.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PaymentController) GetAll(w http.ResponseWriter, r *http.Request) {
	var paymentParams PaymentParams
	queryParams := r.URL.Query()

	// First it checks to see if the values in the querystring are valid
	for key := range queryParams {
		if strings.ToLower(key) != "userID" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*paymentParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*paymentParams.CreatedAt = createdAt
	}

	if valQuery := queryParams.Get("Value"); valQuery != "" {
		value, err := strconv.ParseFloat(valQuery, 32)
		if err == nil {
			val := float32(value)
			*paymentParams.Value = val
		}
	}

	if deliveryID := queryParams.Get("DeliveryID"); deliveryID != "" {
		*paymentParams.DeliveryID = deliveryID
	}

	if userID := queryParams.Get("UserID"); userID != "" {
		*paymentParams.UserID = userID
	}

	// It now passes the payment param as a "filter" and gets the found payment
	foundPaymentes, err := c.repository.GetAll(paymentParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	dtoFoundPaymentes := []PaymentDTO{}
	for i := 0; i < len(foundPaymentes); i++ {
		dtoFoundPaymentes = append(dtoFoundPaymentes, foundPaymentes[i].ToDTO())
	}
	// Returns the DTO payments
	if err := json.NewEncoder(w).Encode(dtoFoundPaymentes); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PaymentController) GetOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	payment, err := c.repository.GetOne(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound) // The Id was not found but the request did go though
		return
	}
	if err := json.NewEncoder(w).Encode(payment.ToDTO()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *PaymentController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	var paymentParams PaymentParams
	queryParams := r.URL.Query()

	for key := range queryParams {
		if strings.ToLower(key) != "userID" {
			http.Error(w, fmt.Sprintf("Invalid query parameter: %s", key), http.StatusBadRequest)
			return
		}
	}

	// If the values are valid it will check what each value is
	if isDeleted := queryParams.Get("IsDeleted"); isDeleted != "" {
		value := isDeleted == "true"
		*paymentParams.IsDeleted = value
	}

	if createdAt := queryParams.Get("CreatedAt"); createdAt != "" {
		*paymentParams.CreatedAt = createdAt
	}

	if valQuery := queryParams.Get("Value"); valQuery != "" {
		value, err := strconv.ParseFloat(valQuery, 32)
		if err == nil {
			val := float32(value)
			*paymentParams.Value = val
		}
	}

	if deliveryID := queryParams.Get("DeliveryID"); deliveryID != "" {
		*paymentParams.DeliveryID = deliveryID
	}

	if userID := queryParams.Get("UserID"); userID != "" {
		*paymentParams.UserID = userID
	}

	// Make the request on the repo
	count, err := c.repository.DeleteAll(paymentParams)
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

func (c *PaymentController) DeleteOne(w http.ResponseWriter, r *http.Request) {
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

func (c *PaymentController) Update(w http.ResponseWriter, r *http.Request) {
	// There is no update method, but this needs to be included since the controller is implementing an interface (IController)
}
