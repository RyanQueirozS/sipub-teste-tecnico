package user_delivery_test

import (
	"regexp"
	"sipub-test/internal/user_delivery"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserDelivery(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_delivery.MySQLUserDeliveryRepository{}
		repo.SetDB(db)

		params := user_delivery.UserDeliveryParams{
			DeliveryID: testhelper.StringPointer("order-123"),
			UserID:     testhelper.StringPointer("user-123"),
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_delivery (id, delivery_id, user_id) VALUES (?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), "order-123", "user-123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, "order-123", result.ToDTO().DeliveryID, "DeliveryID should match")
		assert.Equal(t, "user-123", result.ToDTO().UserID, "UserID should match")
	})
}

func TestGetAllUserDeliveries(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_delivery.MySQLUserDeliveryRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "delivery_id", "user_id"}).
			AddRow("delivery-123", "order-123", "user-123")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, delivery_id, user_id FROM user_delivery WHERE 1=1`)).
			WillReturnRows(rows)

		filter := user_delivery.UserDeliveryParams{}
		results, err := repo.GetAll(filter)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Len(t, results, 1, "Result length should be 1")
		assert.Equal(t, "order-123", results[0].ToDTO().DeliveryID, "DeliveryID should match")
		assert.Equal(t, "user-123", results[0].ToDTO().UserID, "UserID should match")
	})
}

func TestGetUserDeliveryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user_delivery.MySQLUserDeliveryRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "delivery_id", "user_id"}).
		AddRow("delivery-123", "order-123", "user-123")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, delivery_id, user_id FROM user_delivery WHERE id = ?`)).
		WithArgs("delivery-123").
		WillReturnRows(rows)

	result, err := repo.GetOne("delivery-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, "delivery-123", result.ToDTO().Id, "ID should match")
	assert.Equal(t, "order-123", result.ToDTO().DeliveryID, "DeliveryID should match")
	assert.Equal(t, "user-123", result.ToDTO().UserID, "UserID should match")
}

func TestDeleteUserDelivery(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &user_delivery.MySQLUserDeliveryRepository{}
	repo.SetDB(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM user_delivery WHERE id = ?`)).
		WithArgs("delivery-123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("delivery-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, uint(1), count, "Affected row count should be 1")
}

func TestDeleteAllUserDeliveries(t *testing.T) {
	t.Run("ValidDeleteAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &user_delivery.MySQLUserDeliveryRepository{}
		repo.SetDB(db)

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM user_delivery WHERE 1=1 AND user_id = ?`)).
			WithArgs("user-123").
			WillReturnResult(sqlmock.NewResult(1, 3))

		filter := user_delivery.UserDeliveryParams{
			UserID: testhelper.StringPointer("user-123"),
		}

		count, err := repo.DeleteAll(filter)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, uint(3), count, "Affected row count should be 3")
	})
}
