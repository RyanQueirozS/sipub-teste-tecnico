package integration_test

// This file mainly does integration testing for the product creation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipub-test/internal/product"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerCreate(t *testing.T) {
	t.Run("ShouldCreateAValidProductWhenPassingAllValidParams", func(t *testing.T) {
		controller := product.NewProductController()
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewReader([]byte(`{
		"Name": "Test Product",
		"WeightGrams": 500,
		"Price": 25.50
	}`)))
		w := httptest.NewRecorder()

		controller.Create(w, r)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response product.ProductDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Product", response.Name)
		assert.EqualValues(t, 500, response.WeightInGrams)
		assert.EqualValues(t, 25.50, response.Price)
	})

	t.Run("ShouldSendInvalidRequestIfNotPassingTheRequiredFields", func(t *testing.T) {
		controller := product.NewProductController()
		// Will not pass the `Name`
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewReader([]byte(`{
		"WeightGrams": 500,
		"Price": 25.50
	}`)))
		w := httptest.NewRecorder()

		controller.Create(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should send bad request")
	})

	t.Run("ShouldSendInvalidRequestIfUsingAnInvalidField", func(t *testing.T) {
		controller := product.NewProductController()
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewReader([]byte(`{
		"invalid_field": "value",
	}`)))
		w := httptest.NewRecorder()

		controller.Create(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code, "Should send bad request")
	})
}

func TestControllerGetAll(t *testing.T) {
	// It returns status 200 with 0 found because the instruction did go
	// through, but no query was found
	t.Run("ShouldReturn200StatusEvenIfNoneExists", func(t *testing.T) {
		controller := product.NewProductController()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products?name=Test+Product", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)
		assert.Equal(t, http.StatusOK, w.Code)

		var response []product.ProductDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0, "Lenght should be zero")
	})

	t.Run("ShouldReturnBadRequestIfInvalidField", func(t *testing.T) {
		controller := product.NewProductController()
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products?invalid_pparam=10", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
