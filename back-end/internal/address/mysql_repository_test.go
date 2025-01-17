package address_test

import (
	"regexp"
	"sipub-test/internal/address"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)

		params := address.AddressParams{
			IsActive:     testhelper.BoolPointer(true),
			IsDeleted:    testhelper.BoolPointer(false),
			Street:       testhelper.StringPointer("Main Street"),
			Number:       testhelper.StringPointer("123"),
			Neighborhood: testhelper.StringPointer("Downtown"),
			City:         testhelper.StringPointer("Metropolis"),
			State:        testhelper.StringPointer("NY"),
			Country:      testhelper.StringPointer("USA"),
		}

		mock.ExpectExec(`INSERT INTO addresses`).
			WithArgs(sqlmock.AnyArg(), true, false, sqlmock.AnyArg(), "Main Street", "123", "Downtown", "", "Metropolis", "NY", "USA", float64(0), float64(0), "").
			WillReturnResult(sqlmock.NewResult(1, 1))

		addr, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, *params.IsActive, addr.GetIsActive())
		assert.Equal(t, *params.Street, addr.GetStreet())
		assert.Equal(t, *params.City, addr.GetCity())
	})
}

func TestGetAllAddresses(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", "Main Street", "123", "Downtown", "", "Metropolis", "NY", "USA", 0, 0, "")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE 1=1`).
			WillReturnRows(rows)

		filter := address.AddressParams{}
		addresses, err := repo.GetAll(filter)

		assert.NoError(t, err, "Should have no errors")
		assert.Len(t, addresses, 1, "Length should be 1")
		assert.Equal(t, "123", addresses[0].GetID(), "ID should match")
	})
}

func TestGetAddressByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &address.MySQLAddressRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name"}).
		AddRow("123", true, false, "2025-01-15 12:00:00", "Main Street", "123", "Downtown", 0, "Metropolis", "NY", "USA", 0, 0, "")

	mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`).
		WithArgs("123").
		WillReturnRows(rows)

	address, err := repo.GetOne("123")

	assert.NoError(t, err, "Should have no errors")
	assert.Equal(t, "123", address.GetID(), "ID should match")
	assert.Equal(t, "Main Street", address.GetStreet(), "Street should match")
	assert.Equal(t, "Metropolis", address.GetCity(), "City should match")
}

func TestDeleteAddress(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &address.MySQLAddressRepository{}
	repo.SetDB(db)

	mock.ExpectExec(`DELETE FROM addresses WHERE id = ?`).
		WithArgs("123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("123")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
}

func TestUpdateAddress(t *testing.T) {
	t.Run("ValidUpdateWithAllParams", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &address.MySQLAddressRepository{}
		repo.SetDB(db)

		newParams := address.AddressParams{
			IsActive:  testhelper.BoolPointer(false),
			IsDeleted: testhelper.BoolPointer(true),
			Street:    testhelper.StringPointer("Updated Street"),
			City:      testhelper.StringPointer("Gotham"),
			Country:   testhelper.StringPointer("USA"),
		}

		// Initial SELECT for GetOne
		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", "Main Street", "123", "Downtown", "", "Metropolis", "NY", "USA", 0, 0, "")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(rows)

		// UPDATE query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE addresses SET isActive = ?, isDeleted = ?, street = ?, number = ?, neighborhood = ?, complement = ?, city = ?, state = ?, country = ?, latitude = ?, longitude = ?, name = ? WHERE id = ?`)).
			WithArgs(false, true, "Updated Street", "123", "Downtown", "", "Gotham", "NY", "USA", float64(0), float64(0), "", "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Final SELECT for updated address
		updatedRows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "street", "number", "neighborhood", "complement", "city", "state", "country", "latitude", "longitude", "name"}).
			AddRow("123", false, true, "2025-01-15 12:00:00", "Updated Street", "123", "Downtown", float64(0), "Gotham", "NY", "USA", float64(0), float64(0), "")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedRows)

		address, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should contain no errors")
		assert.Equal(t, "123", address.GetID(), "ID should remain the same")
		assert.Equal(t, "Updated Street", address.GetStreet(), "Street should be updated")
		assert.Equal(t, "Gotham", address.GetCity(), "City should be updated")
	})
}
