package shopping_cart_test

import (
	"fmt"
	"log"
	"regexp"
	"sipub-test/internal/shopping_cart"
	testhelper "sipub-test/pkg/test_helper"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShoppingCart(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &shopping_cart.MySQLShoppingCartRepository{}
		repo.SetDB(db)

		params := shopping_cart.ShoppingCartParams{
			UserID:        testhelper.StringPointer("user-123"),
			ProductID:     testhelper.StringPointer("product-456"),
			ProductAmount: testhelper.UintPointer(5),
		}

		mock.ExpectExec(`INSERT INTO shopping_cart`).
			WithArgs(sqlmock.AnyArg(), *params.UserID, *params.ProductID, *params.ProductAmount).
			WillReturnResult(sqlmock.NewResult(1, 1))

		cart, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		assert.Equal(t, *params.UserID, cart.ToDTO().UserID)
		assert.Equal(t, *params.ProductID, cart.ToDTO().ProductID)
		assert.Equal(t, *params.ProductAmount, cart.ToDTO().ProductAmount)
	})
}

func TestGetAllShoppingCarts(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &shopping_cart.MySQLShoppingCartRepository{}
		repo.SetDB(db)

		rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_amount"}).
			AddRow("cart-123", "user-123", "product-456", 5)

		mock.ExpectQuery(`SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE 1=1 AND userID = ?`).
			WithArgs("user-123").
			WillReturnRows(rows)

		filter := shopping_cart.ShoppingCartParams{UserID: testhelper.StringPointer("user-123")}
		carts, err := repo.GetAll(filter)

		assert.NoError(t, err, "Should have no errors")
		assert.Len(t, carts, 1, "Length should be 1")
		assert.Equal(t, "cart-123", carts[0].ToDTO().Id, "Id should be equal to cart-123")
	})

	t.Run("ShouldReturnAnErrorIfDidntFindAny", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &shopping_cart.MySQLShoppingCartRepository{}
		repo.SetDB(db)

		mock.ExpectQuery(`SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE 1=1 AND userID = ?`).
			WithArgs("nonexistent-user").
			WillReturnError(fmt.Errorf("failed to get shopping_cart"))

		filter := shopping_cart.ShoppingCartParams{UserID: testhelper.StringPointer("nonexistent-user")}
		carts, err := repo.GetAll(filter)

		assert.Error(t, err, "Should have an error")
		assert.Len(t, carts, 0, "Length should be 0")
	})
}

func TestGetShoppingCartByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &shopping_cart.MySQLShoppingCartRepository{}
	repo.SetDB(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_amount"}).
		AddRow("cart-123", "user-123", "product-456", 5)

	mock.ExpectQuery(`SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE id = ?`).
		WithArgs("cart-123").
		WillReturnRows(rows)

	cart, err := repo.GetOne("cart-123")

	assert.NoError(t, err, "Should have no errors")
	assert.Equal(t, "cart-123", cart.ToDTO().Id, "Id should be the same")
	assert.Equal(t, "user-123", cart.ToDTO().UserID, "User ID should be the same")
	assert.Equal(t, "product-456", cart.ToDTO().ProductID, "Product ID should be the same")
	assert.Equal(t, uint(5), cart.ToDTO().ProductAmount, "Product amount should be the same")
}

func TestDeleteShoppingCart(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &shopping_cart.MySQLShoppingCartRepository{}
	repo.SetDB(db)

	mock.ExpectExec(`DELETE FROM shopping_cart WHERE id = ?`).
		WithArgs("cart-123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("cart-123")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
}

func TestUpdateShoppingCart(t *testing.T) {
	t.Run("ValidUpdateWithNonZeroAmount", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &shopping_cart.MySQLShoppingCartRepository{}
		repo.SetDB(db)

		newParams := shopping_cart.ShoppingCartParams{
			ProductAmount: testhelper.UintPointer(5),
		}

		// Create a new row for existing shopping cart
		sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_amount"}).
			AddRow("123", "user-1", "product-1", 3)
		// UPDATE query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE shopping_cart SET product_amount = ? WHERE id = ?`)).
			WithArgs(5, "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Final SELECT for updated shopping cart
		updatedRows := sqlmock.NewRows([]string{"id", "user_id", "product_id", "product_amount"}).
			AddRow("123", "user-1", "product-1", 5)

		mock.ExpectQuery(`SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedRows)

		shoppingCart, err := repo.Update("123", newParams)

		log.Println(err)
		assert.NoError(t, err, "Should contain no errors")
		assert.Equal(t, "123", shoppingCart.ToDTO().Id, "ID should remain the same")
		assert.Equal(t, uint(5), shoppingCart.ToDTO().ProductAmount, "Product amount should be updated")
	})

	t.Run("ValidUpdateWithZeroAmount", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &shopping_cart.MySQLShoppingCartRepository{}
		repo.SetDB(db)

		newParams := shopping_cart.ShoppingCartParams{
			ProductAmount: testhelper.UintPointer(0),
		}

		// DELETE query when ProductAmount is zero
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM shopping_cart WHERE id = ?`)).
			WithArgs("123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		shoppingCart, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should not return an error")
		assert.Equal(t, shopping_cart.ShoppingCartModel{}, shoppingCart, "Should return an empty ShoppingCartModel")
	})
}
