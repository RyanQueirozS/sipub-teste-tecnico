package product

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"

	"github.com/google/uuid"
)

func createNewProductTableIfNoneExists() {
	db := db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id CHAR(36) NOT NULL,
		isActive BOOLEAN NOT NULL DEFAULT TRUE,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
		deletedAt DATETIME NULL DEFAULT NULL,
		createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		weightGrams FLOAT NOT NULL,
		price FLOAT NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	// https://dev.mysql.com/doc/refman/8.4/en/innodb-benefits.html
	// Mainly using because it supports foreing keys

	// Execute the query
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Table 'products' has been created successfully (if it did not already exist).")
}

type MySQLProductRepository struct {
	db *sql.DB
}

func NewMySQLproductRepository() *MySQLProductRepository {
	createNewProductTableIfNoneExists()
	return &MySQLProductRepository{db: db.GetDB()}
}

func (r *MySQLProductRepository) Create(params ProductParams) (ProductModel, error) {
	id := uuid.NewString()
	query := `INSERT INTO products (id, isActive, isDeleted, weightGrams, price, name) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, id, params.IsActive, params.IsDeleted, params.WeightGrams, params.Price, params.Name)
	if err != nil {
		return ProductModel{}, fmt.Errorf("failed to create product: %w", err)
	}

	return ProductModel{
		id:          id,
		isActive:    *params.IsActive,
		isDeleted:   *params.IsDeleted,
		weightGrams: *params.WeightGrams,
		price:       *params.Price,
		name:        *params.Name,
	}, nil
}

func (r *MySQLProductRepository) GetAll(filter ProductParams) ([]ProductModel, error) {
	query := `SELECT id, isActive, isDeleted, deletedAt, createdAt, weightGrams, price, name FROM products WHERE 1=1`
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
		return nil, fmt.Errorf("failed to fetch products: %w", err)
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
	query := `SELECT id, isActive, isDeleted, deletedAt, createdAt, weightGrams, price, name FROM products WHERE id = ?`
	var product ProductModel
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&product.id, &product.isActive, &product.isDeleted, &product.createdAt, &product.weightGrams, &product.price, &product.name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ProductModel{}, fmt.Errorf("product not found")
		}
		return ProductModel{}, fmt.Errorf("failed to fetch product: %w", err)
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
	query := `UPDATE products SET isActive = ?, isDeleted = ?, weightGrams = ?, price = ?, name = ? WHERE id = ?`
	_, err := r.db.Exec(query, newProduct.IsActive, newProduct.IsDeleted, newProduct.WeightGrams, newProduct.Price, newProduct.Name, id)
	if err != nil {
		return ProductModel{}, fmt.Errorf("failed to update product: %w", err)
	}
	return r.GetOne(id)
}
