package payment_test

import (
	"regexp"
	"sipub-test/internal/payment"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePayment(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &payment.MySQLPaymentRepository{}
		repo.SetDB(db)

		params := payment.PaymentParams{
			IsDeleted:  testhelper.BoolPointer(false),
			DeliveryID: testhelper.StringPointer("delivery-123"),
			Value:      testhelper.FloatPointer(150.50),
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO payments (id, isDeleted, createdAt, delivery_id, value) VALUES (?, ?, ?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), false, sqlmock.AnyArg(), "delivery-123", 150.50).
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, "delivery-123", result.ToDTO().DeliveryID, "DeliveryID should match")
		assert.Equal(t, float32(150.50), result.ToDTO().Value, "Payment value should match")
		assert.Equal(t, false, result.ToDTO().IsDeleted, "IsDeleted should match")
	})
}

func TestGetAllPayments(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &payment.MySQLPaymentRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "isDeleted", "createdAt", "deliveryID", "value"}).
			AddRow("payment-123", false, "2025-01-15 12:00:00", "delivery-123", 150.50)

		mock.ExpectQuery(regexp.QuoteMeta(`
			SELECT 
				p.id, p.isDeleted, p.createdAt, p.deliveryID, p.value
			FROM 
				userDelivery ud
			JOIN 
				deliveries d ON ud.delivery_id = d.delivery_id
			JOIN 
				payments p ON p.delivery_id = d.delivery_id
			WHERE 
				1=1 AND ud.user_id = ?`)).
			WithArgs("user-123").
			WillReturnRows(rows)

		filter := payment.PaymentParams{UserID: testhelper.StringPointer("user-123")}
		results, err := repo.GetAll(filter)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Len(t, results, 1, "Result length should be 1")
		assert.Equal(t, "delivery-123", results[0].ToDTO().DeliveryID, "DeliveryID should match")
		assert.Equal(t, float32(150.50), results[0].ToDTO().Value, "Payment value should match")
	})
}

func TestGetPaymentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &payment.MySQLPaymentRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "isDeleted", "createdAt", "deliveryID", "value"}).
		AddRow("payment-123", false, "2025-01-15 12:00:00", "delivery-123", 150.50)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, isDeleted, createdAt, delivery_id, value FROM payments WHERE id = ?`)).
		WithArgs("payment-123").
		WillReturnRows(rows)

	result, err := repo.GetOne("payment-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, "payment-123", result.ToDTO().Id, "Payment ID should match")
	assert.Equal(t, "delivery-123", result.ToDTO().DeliveryID, "DeliveryID should match")
	assert.Equal(t, float32(150.50), result.ToDTO().Value, "Payment value should match")
}

func TestDeletePayment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &payment.MySQLPaymentRepository{}
	repo.SetDB(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM payments WHERE id = ?`)).
		WithArgs("payment-123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("payment-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, uint(1), count, "Affected row count should be 1")
}

func TestDeleteAllPayments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &payment.MySQLPaymentRepository{}
	repo.SetDB(db)

	mock.ExpectExec(regexp.QuoteMeta(`
		DELETE p
		FROM payments p
		JOIN deliveries d ON p.deliveryID = d.deliveryID
		JOIN userDelivery ud ON ud.deliveryID = d.deliveryID
		WHERE 1=1 AND ud.user_id = ?`)).
		WithArgs("user-123").
		WillReturnResult(sqlmock.NewResult(1, 5))

	filter := payment.PaymentParams{UserID: testhelper.StringPointer("user-123")}
	count, err := repo.DeleteAll(filter)

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, uint(5), count, "Affected row count should be 5")
}
