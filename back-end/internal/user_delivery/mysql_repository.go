package user_delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"

	"github.com/google/uuid"
)

type MySQLUserDeliveryRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLUserDeliveryRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLUserDeliveryRepository) createNewUserDeliveryTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS user_delivery (
		id CHAR(36) NOT NULL,
        delivery_id char(36) NOT NULL,
        user_id char(36) NOT NULL,
        FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
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

func NewMySQLUserDeliveryRepository() *MySQLUserDeliveryRepository {
	repo := &MySQLUserDeliveryRepository{db: db.GetDB()}
	repo.createNewUserDeliveryTableIfNoneExists()
	return repo
}

func (r *MySQLUserDeliveryRepository) Create(params UserDeliveryParams) (UserDeliveryModel, error) {
	id := uuid.NewString()

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields
	model := UserDeliveryModel{
		id:         id,
		deliveryID: *params.DeliveryID,
		userID:     *params.UserID,
	}

	query := `INSERT INTO user_delivery (id, delivery_id, user_id) VALUES (?, ?, ?)`

	_, err := r.db.Exec(query, id, model.deliveryID, model.userID)
	if err != nil {
		fmt.Println(err)
		return UserDeliveryModel{}, fmt.Errorf("failed to create delivery: %w", err)
	}

	return model, nil
}

func (r *MySQLUserDeliveryRepository) GetAll(filter UserDeliveryParams) ([]UserDeliveryModel, error) {
	query := `SELECT id, delivery_id, user_id FROM user_delivery WHERE 1=1`
	args := []interface{}{}

	// Only looks for userid
	if filter.UserID != nil {
		query += " AND user_id = ?"
		args = append(args, *filter.DeliveryID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get userDelivery: %w", err)
	}
	defer rows.Close()

	var userDelivery []UserDeliveryModel
	for rows.Next() {
		var delivery UserDeliveryModel
		err := rows.Scan(&delivery.id, &delivery.deliveryID, &delivery.userID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan userDelivery: %w", err)
		}
		userDelivery = append(userDelivery, delivery)
	}

	return userDelivery, nil
}

func (r *MySQLUserDeliveryRepository) GetOne(id string) (UserDeliveryModel, error) {
	query := `SELECT id, delivery_id, user_id FROM user_delivery WHERE id = ?`

	var delivery UserDeliveryModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&delivery.id, &delivery.deliveryID, &delivery.userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDeliveryModel{}, fmt.Errorf("delivery not found")
		}
		return UserDeliveryModel{}, fmt.Errorf("failed to get delivery: %w", err)
	}
	return delivery, nil
}

func (r *MySQLUserDeliveryRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM user_delivery WHERE id = ?`
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

func (r *MySQLUserDeliveryRepository) DeleteAll(filter UserDeliveryParams) (uint, error) {
	query := `DELETE FROM user_delivery WHERE 1=1 AND user_id = ?`
	args := []interface{}{*filter.UserID}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user_delivery: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}
