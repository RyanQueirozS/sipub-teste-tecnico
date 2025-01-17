package user_test

import (
	"fmt"
	"regexp"
	"sipub-test/internal/user"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)

		params := user.UserParams{
			IsActive:  testhelper.BoolPointer(true),
			IsDeleted: testhelper.BoolPointer(false),
			Email:     testhelper.StringPointer("testuser@example.com"),
			Cpf:       testhelper.StringPointer("12345678901"),
			Name:      testhelper.StringPointer("Test User"),
		}

		mock.ExpectExec(`INSERT INTO users`).
			WithArgs(sqlmock.AnyArg() /* id determined at function */, true, false, sqlmock.AnyArg() /*time determined at function*/, "testuser@example.com", "12345678901", "Test User").
			WillReturnResult(sqlmock.NewResult(1, 1))

		user, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		// Won't check for id since it is created in the repository
		assert.Equal(t, *params.IsActive, user.GetIsActive())
		assert.Equal(t, *params.Email, user.ToDTO().Email)
		assert.Equal(t, *params.Cpf, user.ToDTO().Cpf)
		assert.Equal(t, *params.Name, user.ToDTO().Name)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", "testuser@example.com", "12345678901", "Test User")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users`).
			WillReturnRows(rows)

		filter := user.UserParams{}
		users, err := repo.GetAll(filter)

		assert.NoError(t, err, "Should have no errors")
		assert.Len(t, users, 1, "Length should be 1")
		assert.Equal(t, "123", users[0].ToDTO().Id, "Id should be equal to 123")
	})
	t.Run("ShouldReturnAnErrorIfDidntFindAny", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)

		// Should return "failed to get users"
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users`).
			WillReturnError(fmt.Errorf("failed to get users"))

		filter := user.UserParams{Email: testhelper.StringPointer("nonexistent@example.com")}
		users, err := repo.GetAll(filter)

		assert.Error(t, err, "Should have an error")
		assert.Len(t, users, 0, "Length should be 0")
	})
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user.MySQLUserRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
		AddRow("123", true, false, "2025-01-15 12:00:00", "testuser@example.com", "12345678901", "Test User")

	mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
		WithArgs("123").
		WillReturnRows(rows)

	user, err := repo.GetOne("123")

	assert.NoError(t, err, "Should have no errors")
	assert.Equal(t, "123", user.ToDTO().Id, "Id should be the same")
	assert.Equal(t, "testuser@example.com", user.ToDTO().Email, "Email should be the same")
	assert.Equal(t, "12345678901", user.ToDTO().Cpf, "Cpf should be the same")
	assert.Equal(t, "2025-01-15 12:00:00", user.ToDTO().CreatedAt, "Date created should be the same")
	assert.Equal(t, "Test User", user.ToDTO().Name, "Name should be the same")
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user.MySQLUserRepository{}
	repo.SetDB(db)

	mock.ExpectExec(`DELETE FROM users WHERE id = ?`).
		WithArgs("123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("123")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
}

func TestUpdateUser(t *testing.T) {
	t.Run("ValidUpdateWithAllParams", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)

		newParams := user.UserParams{
			IsActive:  testhelper.BoolPointer(false),
			IsDeleted: testhelper.BoolPointer(true),
			Email:     testhelper.StringPointer("updateduser@example.com"),
			Cpf:       testhelper.StringPointer("10987654321"),
			Name:      testhelper.StringPointer("Updated User"),
		}

		// Create a new row
		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", "testuser@example.com", "12345678901", "Original User")

		// Initial SELECT for GetOne
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(rows)

		// UPDATE query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET isActive = ?, isDeleted = ?, email = ?, cpf = ?, name = ? WHERE id = ?`)).
			WithArgs(false, true, "updateduser@example.com", "10987654321", "Updated User", "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Final SELECT for updated user
		updatedRows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
			AddRow("123", false, true, "2025-01-15 12:00:00", "updateduser@example.com", "10987654321", "Updated User")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedRows)

		user, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should contain no errors")
		assert.Equal(t, "123", user.ToDTO().Id, "Id should remain the same")
		assert.Equal(t, "updateduser@example.com", user.ToDTO().Email, "Email should be updated")
		assert.Equal(t, "10987654321", user.ToDTO().Cpf, "Cpf should be updated")
		assert.Equal(t, "Updated User", user.ToDTO().Name, "Name should be updated")
	})

	t.Run("Update with Null Params", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user.MySQLUserRepository{}
		repo.SetDB(db)

		existingUser := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", "testuser@example.com", "12345678901", "Original User")

		// Expect the `GetOne` call to return the existing user
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(existingUser)

		// Params with only some fields updated, others are nil
		newParams := user.UserParams{
			IsActive: testhelper.BoolPointer(false),
			Name:     testhelper.StringPointer("Partially Updated User"),
		}

		// Expect the `UPDATE` query with values including the updated fields and the unchanged fields
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET isActive = ?, isDeleted = ?, email = ?, cpf = ?, name = ? WHERE id = ?`)).
			WithArgs(false, false, "testuser@example.com", "12345678901", "Partially Updated User", "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect the `GetOne` call after the update to return the updated user
		updatedUser := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "email", "cpf", "name"}).
			AddRow("123", false, false, "2025-01-15 12:00:00", "testuser@example.com", "12345678901", "Partially Updated User")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedUser)

		user, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should not fail when updating with partial fields")
		assert.Equal(t, "123", user.ToDTO().Id, "Id should remain the same")
		assert.Equal(t, false, user.GetIsActive(), "IsActive should match the updated value")
		assert.Equal(t, "Partially Updated User", user.ToDTO().Name, "Name should be updated")
	})
}
