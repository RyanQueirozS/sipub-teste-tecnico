package product

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math"
	"sipub-test/db"
	"time"

	"github.com/google/uuid"
)

type MySQLProductRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLProductRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLProductRepository) createNewProductTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id CHAR(36) NOT NULL,
		isActive BOOLEAN NOT NULL DEFAULT TRUE,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
        createdAt CHAR(19) NOT NULL,
		weightGrams FLOAT NOT NULL,
		price FLOAT NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	// https://dev.mysql.com/doc/refman/8.4/en/innodb-benefits.html
	// Mainly using InnoDB because it supports foreing keys
	// createdAt is a string because it is simpler to handle. It uses this format 2006-01-02 15:04:05 (19 chars)

	if _, err := r.db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func NewMySQLproductRepository() *MySQLProductRepository {
	repo := &MySQLProductRepository{db: db.GetDB()}
	repo.createNewProductTableIfNoneExists()
	return repo
}

func (r *MySQLProductRepository) Create(params ProductParams) (ProductModel, error) {
	if err := validateProductParams(params); err != nil {
		return ProductModel{}, err
	}
	id := uuid.NewString()

	// Round price to 2 decimal places, if not, there will be floating number
	// innacuracy
	price := math.Round(float64(*params.Price)*100) / 100.0
	timeCreated := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO products (id, isActive, isDeleted, createdAt, weightGrams, price, name) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, id, *params.IsActive, *params.IsDeleted, timeCreated, *params.WeightGrams, price, *params.Name)
	if err != nil {
		return ProductModel{}, fmt.Errorf("failed to create product: %w", err)
	}

	return ProductModel{
		id:          id,
		isActive:    *params.IsActive,
		isDeleted:   *params.IsDeleted,
		createdAt:   timeCreated,
		weightGrams: *params.WeightGrams,
		price:       float32(price),
		name:        *params.Name,
	}, nil
}

func (r *MySQLProductRepository) GetAll(filter ProductParams) ([]ProductModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []ProductModel
	for rows.Next() {
		var product ProductModel
		if err := rows.Scan(&product.id, &product.isActive, &product.isDeleted, &product.createdAt, &product.weightGrams, &product.price, &product.name); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *MySQLProductRepository) GetOne(id string) (ProductModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, weightGrams, price, name FROM products WHERE id = ?`
	var product ProductModel
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&product.id, &product.isActive, &product.isDeleted, &product.createdAt, &product.weightGrams, &product.price, &product.name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProductModel{}, fmt.Errorf("product not found")
		}
		return ProductModel{}, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (r *MySQLProductRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM products WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete product: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLProductRepository) DeleteAll(filter ProductParams) (uint, error) {
	query := `DELETE FROM products WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, filter.IsDeleted)
	}
	if filter.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete products: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLProductRepository) Update(id string, newProduct ProductParams) (ProductModel, error) {
	previousProduct, err := r.GetOne(id)
	if err != nil {
		return ProductModel{}, err
	}
	// This will check nil arguments and change only the non-nil ones
	updatedProduct := ProductModel{
		isActive:    ifNotNilBool(newProduct.IsActive, previousProduct.isActive),
		isDeleted:   ifNotNilBool(newProduct.IsDeleted, previousProduct.isDeleted),
		weightGrams: ifNotNilFloat32(newProduct.WeightGrams, previousProduct.weightGrams),
		price:       ifNotNilFloat32(newProduct.Price, previousProduct.price),
		name:        ifNotNilString(newProduct.Name, previousProduct.name),
	}
	// Remove floating point innacuracy
	roundedWeight := math.Round(float64(updatedProduct.weightGrams)*100) / 100
	roundedPrice := math.Round(float64(updatedProduct.price)*100) / 100

	query := `UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`

	_, err = r.db.Exec(query, updatedProduct.isActive, updatedProduct.isDeleted, roundedWeight, roundedPrice, updatedProduct.name, id)
	if err != nil {
		return ProductModel{}, fmt.Errorf("failed to update product: %w", err)
	}
	return r.GetOne(id)
}

// Helper functions

func validateProductParams(product ProductParams) error {
	if product.Name == nil {
		return errors.New("Invalid name")
	}
	if *product.Name == "" {
		return errors.New("invalid Name")
	}
	if product.Price == nil {
		return errors.New("Invalid Price")
	}
	if *product.Price <= 0 {
		return errors.New("Invalid Price")
	}
	if product.IsActive == nil {
		return errors.New("Invalid IsActive")
	}
	if product.IsDeleted == nil {
		return errors.New("Invalid IsDeleted")
	}
	if product.WeightGrams == nil {
		return errors.New("Invalid WeightGrams")
	}
	if *product.WeightGrams <= 0 {
		return errors.New("Invalid WeightGrams")
	}
	return nil
}

func ifNotNilBool(newVal *bool, oldVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func ifNotNilFloat32(newVal *float32, oldVal float32) float32 {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func ifNotNilString(newVal *string, oldVal string) string {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}
