package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipub-test/internal/delivery_product"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeliveryProductController_Create(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`INSERT INTO delivery_product`).
			WithArgs(sqlmock.AnyArg(), "order-id", "product-id", 10).
			WillReturnResult(sqlmock.NewResult(1, 1))

		requestBody := `{
			"OrderID": "order-id",
			"ProductID": "product-id",
			"ProductAmount": 10
		}`
		r := httptest.NewRequest(http.MethodPost, "/delivery-product", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response delivery_product.DeliveryProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "order-id", response.OrderID)
		assert.Equal(t, "product-id", response.ProductID)
		assert.Equal(t, 10, response.ProductAmount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidRequest", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		requestBody := `{"OrderID": "order-id"}` // Missing required fields
		r := httptest.NewRequest(http.MethodPost, "/delivery-product", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeliveryProductController_GetAll(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "product_amount"}).
			AddRow("id-1", "order-id", "product-id", 10)

		mock.ExpectQuery(`SELECT id, order_id, product_id, product_amount FROM delivery_product WHERE 1=1`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "/delivery-product", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []delivery_product.DeliveryProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, "order-id", response[0].OrderID)
		assert.Equal(t, "product-id", response[0].ProductID)
		assert.Equal(t, 10, response[0].ProductAmount)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidQueryParam", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		r := httptest.NewRequest(http.MethodGet, "/delivery-product?invalid_param=10", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestDeliveryProductController_GetOne(t *testing.T) {
	t.Run("ShouldReturnDeliveryIfExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		row := sqlmock.NewRows([]string{"id", "order_id", "product_id", "product_amount"}).
			AddRow(id, "order-id", "product-id", 10)

		mock.ExpectQuery(`SELECT id, order_id, product_id, product_amount FROM delivery_product WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(row)

		r := httptest.NewRequest(http.MethodGet, "/delivery-product/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response delivery_product.DeliveryProductDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "order-id", response.OrderID)
	})

	t.Run("ShouldReturnNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "non-existent-id"
		mock.ExpectQuery(`SELECT id, order_id, product_id, product_amount FROM delivery_product WHERE id = ?`).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		r := httptest.NewRequest(http.MethodGet, "/delivery-product/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestDeliveryProductController_DeleteOne(t *testing.T) {
	t.Run("ShouldDeleteSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		mock.ExpectExec(`DELETE FROM delivery_product WHERE id = \?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := httptest.NewRequest(http.MethodDelete, "/delivery-product/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "non-existent-id"

		mock.ExpectExec(`DELETE FROM delivery_product WHERE id = \?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := httptest.NewRequest(http.MethodDelete, "/delivery-product/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeliveryProductController_DeleteAll(t *testing.T) {
	t.Run("ShouldDeleteAllSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`DELETE FROM delivery_product WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 10)) // Assume 10 rows deleted

		r := httptest.NewRequest(http.MethodDelete, "/delivery-product", nil)
		w := httptest.NewRecorder()

		controller.DeleteAll(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldHandleEmptyTable", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &delivery_product.DeliveryProductController{}
		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`DELETE FROM delivery_product WHERE 1=1`).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows to delete

		r := httptest.NewRequest(http.MethodDelete, "/delivery-product", nil)
		w := httptest.NewRecorder()

		controller.DeleteAll(w, r)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
