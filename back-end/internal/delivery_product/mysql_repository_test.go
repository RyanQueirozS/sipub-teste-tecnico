package delivery_product_test

import (
	"regexp"
	"sipub-test/internal/delivery"
	"sipub-test/internal/delivery_product"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeliveryProduct(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)

		params := delivery_product.DeliveryProductParams{
			OrderID:       testhelper.StringPointer("order-123"),
			ProductID:     testhelper.StringPointer("product-123"),
			ProductAmount: testhelper.UintPointer(5),
		}

		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO delivery_product (id, order_id, product_id, product_amount) VALUES (?, ?, ?, ?)`)).
			WithArgs(sqlmock.AnyArg(), "order-123", "product-123", 5).
			WillReturnResult(sqlmock.NewResult(1, 1))

		result, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, "order-123", result.ToDTO().OrderID, "OrderID should match")
		assert.Equal(t, "product-123", result.ToDTO().ProductID, "ProductID should match")
		assert.Equal(t, uint(5), result.ToDTO().ProductAmount, "ProductAmount should match")
	})
}

func TestGetAllDeliveryProducts(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &delivery_product.MySQLDeliveryRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "product_amount"}).
			AddRow("delivery-123", "order-123", "product-123", 5)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, order_id, product_id, product_amount FROM delivery_product WHERE 1=1`)).
			WillReturnRows(rows)

		filter := delivery_product.DeliveryProductParams{}
		results, err := repo.GetAll(filter)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Len(t, results, 1, "Result length should be 1")
		assert.Equal(t, "order-123", results[0].ToDTO().OrderID, "OrderID should match")
		assert.Equal(t, "product-123", results[0].ToDTO().ProductID, "ProductID should match")
	})
}

func TestGetDeliveryProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &delivery_product.MySQLDeliveryRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "order_id", "product_id", "product_amount"}).
		AddRow("delivery-123", "order-123", "product-123", 5)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, order_id, product_id, product_amount FROM delivery_product WHERE id = ?`)).
		WithArgs("delivery-123").
		WillReturnRows(rows)

	result, err := repo.GetOne("delivery-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, "delivery-123", result.ToDTO().Id, "ID should match")
	assert.Equal(t, "order-123", result.ToDTO().OrderID, "OrderID should match")
	assert.Equal(t, "product-123", result.ToDTO().ProductID, "ProductID should match")
}

func TestDeleteDeliveryProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &delivery_product.MySQLDeliveryRepository{}
	repo.SetDB(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM delivery_product WHERE id = ?`)).
		WithArgs("delivery-123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("delivery-123")

	assert.NoError(t, err, "Shouldn't contain any errors")
	assert.Equal(t, uint(1), count, "Affected row count should be 1")
}

func TestUpdateDelivery(t *testing.T) {
	t.Run("ValidUpdate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &delivery.MySQLDeliveryRepository{}
		repo.SetDB(db)

		existingRows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "user_id", "address_id"}).
			AddRow("delivery-123", true, false, "2025-01-15 12:00:00", "user-123", "address-123")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, user_id, address_id FROM deliveries WHERE id = ?`).
			WithArgs("delivery-123").
			WillReturnRows(existingRows)

		newParams := delivery.DeliveryParams{
			IsActive:  testhelper.BoolPointer(false),
			AddressID: testhelper.StringPointer("new-address-123"),
		}

		mock.ExpectExec(regexp.QuoteMeta(`UPDATE deliveries SET isActive = ?, isDeleted = ?, address_id = ? WHERE id = ?`)).
			WithArgs(false, false, "new-address-123", "delivery-123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		updatedRows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "user_id", "address_id"}).
			AddRow("delivery-123", false, false, "2025-01-15 12:00:00", "user-123", "new-address-123")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, user_id, address_id FROM deliveries WHERE id = ?`).
			WithArgs("delivery-123").
			WillReturnRows(updatedRows)

		delivery, err := repo.Update("delivery-123", newParams)

		assert.NoError(t, err)
		assert.Equal(t, "delivery-123", delivery.ToDTO().Id, "ID should match")
		assert.Equal(t, "new-address-123", delivery.ToDTO().AddressID, "AddressID should match updated value")
		assert.Equal(t, false, delivery.ToDTO().IsActive, "IsActive should match updated value")
	})
}
