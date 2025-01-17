package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"sipub-test/internal/address"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestControllerCreate(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`INSERT INTO addresses`).
			WithArgs(sqlmock.AnyArg(), true, false, sqlmock.AnyArg(), "Test Street", "123", "Test Neighborhood", sqlmock.AnyArg(), "Test City", "NY", "USA", float64(0), float64(0), "Test Address").
			WillReturnResult(sqlmock.NewResult(1, 1))

		requestBody := `{
			"IsActive": true,
			"IsDeleted": false,
			"Street": "Test Street",
			"Number": "123",
			"Neighborhood": "Test Neighborhood",
			"City": "Test City",
			"State": "NY",
			"Country": "USA",
			"Name": "Test Address",
			"IsActive": true,
			"IsDeleted": false,
            "Latitude": 0,
            "Longitude": 0
		}`
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/addresses", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response address.AddressDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Address", response.Name)
		assert.Equal(t, "Test Street", response.Street)
		assert.Equal(t, "123", response.Number)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidRequest", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		requestBody := `{
			"Number": "123",
			"Neighborhood": "Test Neighborhood",
			"City": "Test City"
		}` // Missing required Street, State, Country, etc.
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/addresses", bytes.NewReader([]byte(requestBody)))
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

		controller := address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name",
		}).
			AddRow("1", true, false, "2023-01-01 12:00:00", "Test Street", "123", "Test Neighborhood", "", "Test City", "NY", "USA", 0, 0, "Test Address")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/addresses", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []address.AddressDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, "Test Address", response[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnEmptyResult", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name",
		})
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/addresses", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []address.AddressDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidQueryParam", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/addresses?invalid_param=10", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestControllerDeleteOne(t *testing.T) {
	t.Run("ShouldDeleteAddressSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM addresses WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/addresses/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ShouldReturn404IfAddressNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		controller := &address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM addresses WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/addresses/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestControllerGetOne(t *testing.T) {
	t.Run("ShouldReturnAddressIfExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		rows := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted",
			"createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name",
		}).
			AddRow(id, true, false, "2023-01-01 00:00:00", "Main St", "123", "Downtown", "", "City", "NY", "USA", float64(0), float64(0), "Test Address")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/addresses/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response address.AddressDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, id, response.Id)
		assert.Equal(t, "Test Address", response.Name)
	})

	t.Run("ShouldReturn404IfAddressNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/addresses/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestControllerDeleteAll(t *testing.T) {
	t.Run("ShouldDeleteAllMatchingAddresses", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller := &address.AddressController{}
		controller.SetRepository(repo)

		mock.ExpectExec(`DELETE FROM addresses WHERE 1=1 AND street LIKE ?`).
			WithArgs("%Main%").
			WillReturnResult(sqlmock.NewResult(0, 3)) // Simulate 3 rows deleted

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/addresses?Street=Main", nil)
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
	t.Run("ShouldUpdateAddressSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &address.AddressController{}
		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		// Mock previous address fetch
		rowsBeforeUpdate := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude",
			"name",
		}).
			AddRow(id, true, false, "2023-01-01 00:00:00", "Main St", "123", "Downtown", "", "City", "NY", "USA", float64(0), float64(0), "Old Address")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`)).
			WithArgs(id).
			WillReturnRows(rowsBeforeUpdate)

			// Mock update query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE addresses SET isActive = ?, isDeleted = ?, street = ?, number = ?, neighborhood = ?, complement = ?, city = ?, state = ?, country = ?, latitude = ?, longitude = ?, name = ? WHERE id = ?`)).
			WithArgs(true, false, "New St", "456", "Downtown", "", "City", "NY", "USA", float64(0), float64(0), "Updated Address", id).
			WillReturnResult(sqlmock.NewResult(1, 1))

			// Mock updated address fetch
		rowsAfterUpdate := sqlmock.NewRows([]string{
			"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name",
		}).
			AddRow(id, true, false,
				"2023-01-01 00:00:00", "New St", "456", "Downtown", "", "City", "NY", "USA", float64(0), float64(0), "Updated Address")

		mock.ExpectQuery(regexp.QuoteMeta(`
        SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?
        `)).
			WithArgs(id).
			WillReturnRows(rowsAfterUpdate)

		body := `
        {
            "Name": "Updated Address",
            "Street": "New St",
            "Number": "456"
        }
        `
		r := httptest.NewRequest(http.MethodPut,
			"http://localhost:8080/addresses/"+id, bytes.NewReader([]byte(body)))
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.Update(w, r)

		// Validate response
		assert.Equal(t, http.StatusOK, w.Code)

		var response address.AddressDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Address", response.Name)
		assert.Equal(t, "New St", response.Street)
		assert.Equal(t, "456", response.Number)
	})
}
