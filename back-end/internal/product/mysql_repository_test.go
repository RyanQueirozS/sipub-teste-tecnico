package product_test

import (
	"math"
	"regexp"
	"sipub-test/internal/product"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
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
}

func TestGetAllProducts(t *testing.T) {
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

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`)).
		WithArgs(false, true, 200.0, 29.99, "Updated Product", "123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "isActive", "isDeleted", "createdAt", "weightGrams", "price", "name"}).
		AddRow("123", false, false, "2025-01-15 12:00:00", 200.0, 29.99, "Updated Product")

	mock.ExpectQuery(`SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`).
		WithArgs("123").
		WillReturnRows(rows)

	product, err := repo.Update("123", newParams)

	assert.NoError(t, err, "Should contain no errors")
	assert.Equal(t, "123", product.ToDTO().Id, "Id should remain the same")
	assert.Equal(t, false, product.GetIsActive(), "IsActive should remain the same")
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
