package delivery_product

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"

	"github.com/google/uuid"
)

type MySQLDeliveryRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLDeliveryRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLDeliveryRepository) createNewDeliveryTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS delivery_product (
		id CHAR(36) NOT NULL,
        delivery_id char(36) NOT NULL,
        product_id char(36) NOT NULL,
        product_amount INT UNSIGNED NOT NULL,
        FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
        FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
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

func NewMySQLDeliveryRepository() *MySQLDeliveryRepository {
	repo := &MySQLDeliveryRepository{db: db.GetDB()}
	repo.createNewDeliveryTableIfNoneExists()
	return repo
}

func (r *MySQLDeliveryRepository) Create(params DeliveryProductParams) (DeliveryProductModel, error) {
	id := uuid.NewString()

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields
	model := DeliveryProductModel{
		id:            id,
		deliveryID:    *params.DeliveryID,
		productID:     *params.ProductID,
		productAmount: *params.ProductAmount,
	}

	query := `INSERT INTO delivery_product (id, delivery_id, product_id, product_amount) VALUES (?, ?, ?, ?)`

	_, err := r.db.Exec(query, id, model.deliveryID, model.productID, model.productAmount) // todo
	if err != nil {
		fmt.Println(err)
		return DeliveryProductModel{}, fmt.Errorf("failed to create delivery: %w", err)
	}

	return model, nil
}

func (r *MySQLDeliveryRepository) GetAll(filter DeliveryProductParams) ([]DeliveryProductModel, error) {
	query := `SELECT id, delivery_id, product_id, product_amount FROM delivery_product WHERE 1=1`
	args := []interface{}{}

	// Only looks for deliveryid
	if filter.DeliveryID != nil {
		query += " AND delivery_id = ?"
		args = append(args, *filter.DeliveryID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveryProduct: %w", err)
	}
	defer rows.Close()

	var deliveryProduct []DeliveryProductModel
	for rows.Next() {
		var delivery DeliveryProductModel
		err := rows.Scan(&delivery.id, &delivery.deliveryID, &delivery.productID, &delivery.productAmount) // todo
		if err != nil {
			return nil, fmt.Errorf("failed to scan deliveryProduct: %w", err)
		}
		deliveryProduct = append(deliveryProduct, delivery)
	}

	return deliveryProduct, nil
}

func (r *MySQLDeliveryRepository) GetOne(id string) (DeliveryProductModel, error) {
	query := `SELECT id, delivery_id, product_id, product_amount FROM delivery_product WHERE id = ?`

	var delivery DeliveryProductModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&delivery.id, &delivery.deliveryID, &delivery.productID, &delivery.productAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeliveryProductModel{}, fmt.Errorf("delivery not found")
		}
		return DeliveryProductModel{}, fmt.Errorf("failed to get delivery: %w", err)
	}
	return delivery, nil
}

func (r *MySQLDeliveryRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM delivery_product WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete delivery: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("no delivery found with the given ID")
	}
	return uint(count), nil
}

func (r *MySQLDeliveryRepository) DeleteAll(filter DeliveryProductParams) (uint, error) {
	query := `DELETE FROM delivery_product WHERE 1=1`
	args := []interface{}{}

	if filter.DeliveryID != nil {
		query += " AND user_id = ?"
		args = append(args, *filter.DeliveryID)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete delivery_product: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}
