package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sipub-test/internal/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserControllerCreate(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		mock.ExpectExec(`INSERT INTO users`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		requestBody := `{
			"Email": "test@example.com",
			"Cpf": "12345678901",
			"Name": "Test User",
			"IsActive": true,
			"IsDeleted": false
		}`
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/u", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response user.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Test User", response.Name)
		assert.Equal(t, "test@example.com", response.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidRequest", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		// Arrange
		controller := user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		requestBody := `{
			"Cpf": "12345678901",
			"Email": "test@example.com"
		}` // Missing required Name field
		r := httptest.NewRequest(http.MethodPost, "http://localhost:8080/u", bytes.NewReader([]byte(requestBody)))
		w := httptest.NewRecorder()

		controller.Create(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserControllerGetAll(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{"id", "createdAt", "email", "cpf", "name", "isActive", "isDeleted"}).
			AddRow("1", true, false, "2023-01-01 12:00:00", "test@example.com", "12345678901", "Test User")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE 1=1`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/u", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []user.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
		assert.Equal(t, "Test User", response[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnEmptyResult", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		rows := sqlmock.NewRows([]string{"id", "createdAt", "email", "cpf", "name", "isActive", "isDeleted"})
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE 1=1`).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/u", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []user.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("ShouldReturnInvalidQueryParam", func(t *testing.T) {
		db, _, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/u?invalid_param=10", nil)
		w := httptest.NewRecorder()

		controller.GetAll(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserControllerGetOne(t *testing.T) {
	t.Run("ShouldReturnUserIfExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"

		rows := sqlmock.NewRows([]string{
			"id", "createdAt", "email", "cpf", "name", "isActive", "isDeleted",
		}).
			AddRow(id, true, false, "2023-01-01 00:00:00", "test@example.com", "12345678901", "Test User")
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
			WithArgs(id).
			WillReturnRows(rows)

		r := httptest.NewRequest(http.MethodGet,
			"http://localhost:8080/u/"+id, nil)
		r.SetPathValue("id", id)
		w := httptest.NewRecorder()

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		var response user.UserDTO
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, id, response.Id)
		assert.Equal(t, "Test User", response.Name)
	})
	t.Run("ShouldReturn404IfUserNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectQuery(`SELECT id, createdAt, email, cpf, name, isActive, isDeleted FROM users WHERE id = ?`).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/u/"+id, nil)
		w := httptest.NewRecorder()

		controller.GetOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserControllerDeleteOne(t *testing.T) {
	t.Run("ShouldDeleteUserSuccessfully", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM users WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/u/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ShouldReturn404IfUserNotFound", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		controller := &user.UserController{}
		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)
		controller.SetRepository(repo)

		id := "123e4567-e89b-12d3-a456-426614174000"
		mock.ExpectExec(`DELETE FROM users WHERE id = ?`).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		r := httptest.NewRequest(http.MethodDelete, "http://localhost:8080/u/"+id, nil)
		w := httptest.NewRecorder()
		r.SetPathValue("id", id)

		controller.DeleteOne(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
