package shopping_cart

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"

	"github.com/google/uuid"
)

type MySQLShoppingCartRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLShoppingCartRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLShoppingCartRepository) createNewShoppingCartTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS shopping_cart (
		id CHAR(36) NOT NULL,
        user_id CHAR(36) NOT NULL,
        product_id char(36) NOT NULL,
        product_amount UNSIGNED INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
        PRIMARY KEY (id)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	// https://dev.mysql.com/doc/refman/8.4/en/innodb-benefits.html
	// Mainly using InnoDB because it supports foreing keys
	// createdAt is a string because it is simpler to handle. It uses this
	// format 2006-01-02 15:04:05 (19 chars)

	if _, err := r.db.Exec(createTableQuery); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

func NewMySQLShoppingCartRepository() *MySQLShoppingCartRepository {
	repo := &MySQLShoppingCartRepository{db: db.GetDB()}
	repo.createNewShoppingCartTableIfNoneExists()
	return repo
}

func (r *MySQLShoppingCartRepository) Create(params ShoppingCartParams) (ShoppingCartModel, error) {
	id := uuid.NewString()

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields. None of the fields should be nil
	model := ShoppingCartModel{
		id:            id,
		userID:        *params.UserID,
		productID:     *params.ProductID,
		productAmount: *params.ProductAmount,
	}

	query := `INSERT INTO shopping_cart (id, user_id, product_id, product_amount) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, id, model.userID, model.productID, model.productAmount)
	if err != nil {
		return ShoppingCartModel{}, fmt.Errorf("failed to create ShoppingCart: %w", err)
	}

	return model, nil
}

func (r *MySQLShoppingCartRepository) GetAll(filter ShoppingCartParams) ([]ShoppingCartModel, error) {
	query := `SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE 1=1 AND userID = ?`
	args := []interface{}{filter.UserID}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get shoppingCart: %w", err)
	}
	defer rows.Close()

	var shopping_cart []ShoppingCartModel
	for rows.Next() {
		var shoppingCart ShoppingCartModel
		err := rows.Scan(&shoppingCart.id,
			&shoppingCart.userID,
			&shoppingCart.productID,
			&shoppingCart.productAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shoppingCart: %w", err)
		}
		shopping_cart = append(shopping_cart, shoppingCart)
	}

	return shopping_cart, nil
}

func (r *MySQLShoppingCartRepository) GetOne(id string) (ShoppingCartModel, error) {
	query := `SELECT id, user_id, product_id, product_amount FROM shopping_cart WHERE id = ?`

	var shoppingCart ShoppingCartModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&shoppingCart.id,
		&shoppingCart.userID,
		&shoppingCart.productID,
		&shoppingCart.productAmount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ShoppingCartModel{}, fmt.Errorf("ShoppingCart not found")
		}
		return ShoppingCartModel{}, fmt.Errorf("failed to get ShoppingCart: %w", err)
	}
	return shoppingCart, nil
}

func (r *MySQLShoppingCartRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM shopping_cart WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete ShoppingCart: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("no ShoppingCart found with the given ID")
	}
	return uint(count), nil
}

func (r *MySQLShoppingCartRepository) DeleteAll(filter ShoppingCartParams) (uint, error) {
	query := `DELETE FROM shopping_cart WHERE 1=1 AND user_id = ?`
	args := []interface{}{*filter.UserID}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete shopping_cart: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLShoppingCartRepository) Update(id string, newShoppingCart ShoppingCartParams) (ShoppingCartModel, error) {
	// Since the productAmount is not nil
	updatedShoppingCart := ShoppingCartModel{
		productAmount: *newShoppingCart.ProductAmount,
	}
	if *newShoppingCart.ProductAmount == 0 {
		_, err := r.DeleteOne(id)
		if err != nil {
			return ShoppingCartModel{}, err
		}
		return ShoppingCartModel{}, nil
	}
	query := `UPDATE shopping_cart SET product_amount = ? WHERE id = ?`

	_, err := r.db.Exec(query,
		updatedShoppingCart.productAmount,
		id)
	if err != nil {
		return ShoppingCartModel{}, fmt.Errorf("failed to update shoppingCart: %w", err)
	}
	return r.GetOne(id)
}
