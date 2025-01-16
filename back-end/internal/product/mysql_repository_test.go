package product_test

import (
	"fmt"
	"math"
	"regexp"
	"sipub-test/internal/product"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	t.Run("ValidCreate", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)

		params := product.ProductParams{
			IsActive:    boolPointer(true),
			IsDeleted:   boolPointer(false),
			WeightGrams: floatPointer(100),
			Price:       floatPointer(19.99),
			Name:        stringPointer("Test Product"),
		}

		mock.ExpectExec(`INSERT INTO products`).
			WithArgs(sqlmock.AnyArg() /* id determined at function */, true, false, sqlmock.AnyArg() /*time determined at function*/, 100.0, 19.99, "Test Product").
			WillReturnResult(sqlmock.NewResult(1, 1))

		product, err := repo.Create(params)

		assert.NoError(t, err, "Shouldn't contain any errors")
		// Won't check for id since it is created in the repository
		assert.Equal(t, *params.IsActive, product.GetIsActive())
		assert.Equal(t, *params.Price, product.GetPrice())
	})

	t.Run("InvalidCreateWithNullParams", func(t *testing.T) {
		// Although controller validates incomming request, the repo adds a new
		// protection layer to validate if the fields are nil so it doesn't
		// insert into the database
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)

		params := product.ProductParams{} // All params are nil

		_, err = repo.Create(params)

		assert.Error(t, err, "Should return error for invalid product creation")
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("ValidGetAll", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)

		// Setting the id to 123 is unreallistic but it works for a testing environment
		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", 100.0, 19.99, "Test Product")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products`).
			WillReturnRows(rows)

		filter := product.ProductParams{}
		products, err := repo.GetAll(filter)

		assert.NoError(t, err, "Should have no errors")
		assert.Len(t, products, 1, "Lenght should be 1")
		assert.Equal(t, "123", products[0].ToDTO().Id, "Id should be equal to 123")
	})
	t.Run("ShouldReturnAnErrorIfDindntFindAny", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)
		// Create one with weight 100
		sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", 100.0, 19.99, "Test Product")

			// Should return "failed to get products"
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products`).
			WillReturnError(fmt.Errorf("failed to get products"))

		// Will search for one with weight 10 and should return 0 found
		filter := product.ProductParams{WeightGrams: floatPointer(10)}
		products, err := repo.GetAll(filter)

		assert.Error(t, err, "Should have no errors")
		assert.Len(t, products, 0, "Lenght should be 0")
	})
}

func TestGetProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &product.MySQLProductRepository{}
	repo.SetDB(db)

	// Setting the id to 123 is unreallistic but it works for a testing environment
	rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
		AddRow("123", true, false, "2025-01-15 12:00:00", 100.0, 19.99, "Test Product")

	mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
		WithArgs("123").
		WillReturnRows(rows)

	product, err := repo.GetOne("123")

	// Rounded to fix floating point innacuracy
	roundedWeight := math.Round(float64(product.ToDTO().WeightInGrams)*100) / 100
	roundedPrice := math.Round(float64(product.ToDTO().Price)*100) / 100

	assert.NoError(t, err, "Should have no errors")
	assert.Equal(t, "123", product.ToDTO().Id, "Id should be the same")
	assert.Equal(t, true, product.GetIsActive(), "IsActive should be the same")
	assert.Equal(t, "2025-01-15 12:00:00", product.ToDTO().CreatedAt, "Date created should be the same")
	assert.Equal(t, 100.0, roundedWeight, "Weigth should be the same")
	assert.Equal(t, 19.99, roundedPrice, "Price should be the same")
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock DB: %v", err)
	}
	defer db.Close()

	repo := &product.MySQLProductRepository{}
	repo.SetDB(db)

	mock.ExpectExec(`DELETE FROM products WHERE id = ?`).
		WithArgs("123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	count, err := repo.DeleteOne("123")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), count)
}

func TestUpdateProduct(t *testing.T) {
	t.Run("ValidUpdateWithAllParams", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)

		newParams := product.ProductParams{
			IsActive:    boolPointer(false),
			IsDeleted:   boolPointer(true),
			WeightGrams: floatPointer(200),
			Price:       floatPointer(29.99),
			Name:        stringPointer("Updated Product"),
		}

		// Create a new row
		rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", 150.0, 19.99, "Original Product")

		// Initial SELECT for GetOne
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(rows)

		// UPDATE query
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`)).
			WithArgs(false, true, 200.0, 29.99, "Updated Product", "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Final SELECT for updated product
		updatedRows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", false, true, "2025-01-15 12:00:00", 200.0, 29.99, "Updated Product")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedRows)

		product, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should contain no errors")
		assert.Equal(t, "123", product.ToDTO().Id, "Id should remain the same")
		assert.Equal(t, true, product.GetIsDeleted(), "IsDeleted should be set to true")
	})
	t.Run("Update with Null Params", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("failed to create mock DB: %v", err)
		}
		defer db.Close()

		repo := &product.MySQLProductRepository{}
		repo.SetDB(db)

		// Creating an existing product that will be retrieved from the database
		existingProduct := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", true, false, "2025-01-15 12:00:00", 100.0, 19.99, "Original Product")

		// Expect the `GetOne` call to return the existing product
		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(existingProduct)

		// Params with only some fields updated, others are nil
		newParams := product.ProductParams{
			IsActive: boolPointer(false),                         // This should update
			Name:     stringPointer("Partially Updated Product"), // This should update
		}

		// Expect the `UPDATE` query with values including the updated fields and the unchanged fields
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`)).
			WithArgs(false, false, 100.0, 19.99, "Partially Updated Product", "123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect the `GetOne` call after the update to return the updated product
		updatedProduct := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
			AddRow("123", false, false, "2025-01-15 12:00:00", 100.0, 19.99, "Partially Updated Product")

		mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
			WithArgs("123").
			WillReturnRows(updatedProduct)

		product, err := repo.Update("123", newParams)

		assert.NoError(t, err, "Should not fail when updating with partial fields")
		assert.Equal(t, "123", product.ToDTO().Id, "Id should remain the same")
		assert.Equal(t, false, product.GetIsActive(), "IsActive should match the updated value")
		assert.Equal(t, "Partially Updated Product", product.ToDTO().Name, "Name should be updated")
		assert.Equal(t, float32(19.99), product.ToDTO().Price, "Price should remain unchanged")
		assert.Equal(t, float32(100.0), product.ToDTO().WeightInGrams, "Weight should remain unchanged")
	})
}

func boolPointer(b bool) *bool {
	return &b
}

func floatPointer(f float32) *float32 {
	return &f
}

func stringPointer(s string) *string {
	return &s
}
