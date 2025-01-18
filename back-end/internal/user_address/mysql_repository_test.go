package user_address_test

import (
	"fmt"
	"sipub-test/internal/user_address"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserAddress(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_address.MySQLUserAddressRepository{}
		repo.SetDB(db)

		params := user_address.UserAddressParams{
			UserID:    "user-123",    // UserID of an existing user
			AddressID: "address-456", // AddressID of an existing address
		}

		mock.ExpectExec(`INSERT INTO user_address`).
			WithArgs(sqlmock.AnyArg() /* id determined at function */, "user-123", "address-456").
			WillReturnResult(sqlmock.NewResult(1, 1))

		userAddress, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, params.UserID, userAddress.UserID)
		assert.Equal(t, params.AddressID, userAddress.AddressID)
	})
}

func TestGetAllUserAddresses(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_address.MySQLUserAddressRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "user_id", "address_id"}).
			AddRow("123", "user-123", "address-456")

		mock.ExpectQuery(`SELECT id, user_id, address_id FROM user_address`).
			WillReturnRows(rows)

		filter := user_address.UserAddressParams{UserID: "user-123"}
		userAddresses, err := repo.GetAll(filter)

		assert.NoError(t, err, "Should have no errors")
		assert.Len(t, userAddresses, 1, "Length should be 1")
		assert.Equal(t, "123", userAddresses[0].GetID(), "Id should be equal to 123")
	})
	t.Run("ShouldReturnAnErrorIfDidntFindAny", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_address.MySQLUserAddressRepository{}
		repo.SetDB(db)

		// Should return "failed to get user_address"
		mock.ExpectQuery(`SELECT id, user_iD, address_id FROM user_address`).
			WillReturnError(fmt.Errorf("failed to get user_address"))

		filter := user_address.UserAddressParams{UserID: "nonexistent-user-id"}
		userAddresses, err := repo.GetAll(filter)

		assert.Error(t, err, "Should have an error")
		assert.Len(t, userAddresses, 0, "Length should be 0")
	})
}

func TestGetUserAddressByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user_address.MySQLUserAddressRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "addressID"}).
		AddRow("123", "user-123", "address-456")

	mock.ExpectQuery(`SELECT id, user_id, address_id FROM user_address WHERE id = ?`).
		WithArgs("123").
		WillReturnRows(rows)

	userAddress, err := repo.GetOne("123")

	assert.NoError(t, err, "Should have no errors")
	assert.Equal(t, "user-123", userAddress.UserID, "UserID should be the same")
	assert.Equal(t, "address-456", userAddress.AddressID, "AddressID should be the same")
}

func TestDeleteUserAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user_address.MySQLUserAddressRepository{}
	repo.SetDB(db)

	mock.ExpectExec(`DELETE FROM user_address WHERE id = ?`).
		WithArgs("123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("123")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
}
