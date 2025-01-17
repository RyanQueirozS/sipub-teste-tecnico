package main

// This file mainly does integration testing for the product creation

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sipub-test/internal/product"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Creating through the controller itself (making an http request)
func TestControllerCreate(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`INSERT INTO products`).
			WithArgs(sqlmock.AnyArg(), true, false, sqlmock.AnyArg(), 500.0, 25.50, "Test Product").
			WillReturnResult(sqlmock.NewResult(1, 1))

		requestBody := `{
			"Name": "Test Product",
			"WeightGrams": 500,
			"Price": 25.50,
			"IsActive": true,
			"IsDeleted": false
		}`
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response product.ProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Product", response.Name)
		assert.EqualValues(t, 500, response.WeightGrams)
		assert.EqualValues(t, 25.50, response.Price)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidRequest", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		requestBody := `{
			"WeightGrams": 500,
			"Price": 25.50
		}` // Missing required Name field
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/products", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestControllerGetAll(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("1", true, false, "2023-01-01 12:00:00", 500.0, 25.50, "Test Product")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []product.ProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, "Test Product", response[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnEmptyResult", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"})
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []product.ProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidQueryParam", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products?invalid_param=10", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestControllerGetOne(t *testing.T) {
	t.Run("ShouldReturnProductIfExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		rows := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted",
			"createdAt", "weightGrams", "price", "name",
		}).
			AddRow(id, true, false, "2023-01-01 00:00:00", 500, 25.50, "Test Product")
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt,
        weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet,
			"http://localhost:8080/products/"+id, nil)
		// Set the path value of the ID parameter
		r.SetPathValue("id", id)
		w := httptest.NewRecorder()

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response product.ProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, id, response.Id)
		assert.Equal(t, "Test Product", response.Name)
	})
	t.Run("ShouldReturn404IfProductNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/products/"+id, nil)
		w := httptest.NewRecorder()

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestControllerDeleteOne(t *testing.T) {
	t.Run("ShouldDeleteProductSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/products/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ShouldReturn404IfProductNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/products/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestControllerDeleteAll(t *testing.T) {
	t.Run("ShouldDeleteAllMatchingProducts", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller := &product.ProductController{}
		controller.SetRepository(repo)

		mock.ExpectExec(`DELETE FROM products WHERE 1=1 AND name LIKE ?`).
			WithArgs("%Test%").
			WillReturnResult(sqlmock.NewResult(0, 3)) // Simulate 3 rows deleted

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/products?Name=Test", nil)
		w := httptest.NewRecorder()

		controller.DeleteAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response uint
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(3), response) // Assert that 3 rows were deleted
	})
}

func TestControllerUpdate(t *testing.T) {
	t.Run("ShouldUpdateProductSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &product.ProductController{}
		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mock previous product fetch
		rowsBeforeUpdate := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow(id, true, false, "2023-01-01 00:00:00", 500.0, 25.50, "Old Product")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`)).
			WithArgs(id).
			WillReturnRows(rowsBeforeUpdate)

		// Mock update query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`)).
			WithArgs(true, false, 600.0, 25.50, "Updated Name", id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Mock updated product fetch
		rowsAfterUpdate := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow(id, true, false, "2023-01-01 00:00:00", 600.0, 25.50, "Updated Name")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`)).
			WithArgs(id).
			WillReturnRows(rowsAfterUpdate)

		body := `{"Name": "Updated Name", "WeightGrams": 600}`
		r := httptest.NewRequest(http.MethodPut, "http://localhost:8080/products/"+id, bytes.NewReader([]byte(body)))
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.Update(w, r)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)

		var response product.ProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		assert.Equal(t, "Updated Name", response.Name)
		assert.Equal(t, float32(600.0), response.WeightGrams)
	})
}

func boolPointer(b bool) *bool {
	return &b
}

func floatPointer(f float32) *float32 {
	return &f
}

func stringPointer(s string) *string {
	return &s
}
