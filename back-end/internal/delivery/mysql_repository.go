package delivery

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"
	"sipub-test/pkg/nilcheck"
	"time"

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
	CREATE TABLE IF NOT EXISTS deliveries (
		id CHAR(36) NOT NULL,
		isActive BOOLEAN NOT NULL DEFAULT TRUE,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
        createdAt CHAR(19) NOT NULL,
        user_id CHAR(36) NOT NULL,
        address_id char(36) NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
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

func (r *MySQLDeliveryRepository) Create(params DeliveryParams) (DeliveryModel, error) {
	id := uuid.NewString()

	timeCreated := time.Now().Format("2006-01-02 15:04:05")

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields
	model := DeliveryModel{
		id:        id,
		isActive:  nilcheck.NotNilBool(params.IsActive, false),
		isDeleted: nilcheck.NotNilBool(params.IsDeleted, true),
		createdAt: timeCreated,
		userID:    *params.UserID,
		addressID: *params.AddressID,
	}
	fmt.Println(model)

	query := `INSERT INTO deliveries (id, isActive, isDeleted, createdAt, user_id, address_id) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, id, model.isActive, model.isDeleted, timeCreated, model.userID, model.addressID)
	if err != nil {
		return DeliveryModel{}, fmt.Errorf("failed to create delivery: %w", err)
	}

	return model, nil
}

func (r *MySQLDeliveryRepository) GetAll(filter DeliveryParams) ([]DeliveryModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, user_id, address_id FROM deliveries WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.UserID != nil {
		query += " AND user_id = ?"
		args = append(args, *filter.UserID)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []DeliveryModel
	for rows.Next() {
		var delivery DeliveryModel
		err := rows.Scan(&delivery.id,
			&delivery.isActive,
			&delivery.isDeleted,
			&delivery.createdAt,
			&delivery.userID,
			&delivery.addressID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

func (r *MySQLDeliveryRepository) GetOne(id string) (DeliveryModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, user_id, address_id FROM deliveries WHERE id = ?`

	var delivery DeliveryModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&delivery.id,
		&delivery.isActive,
		&delivery.isDeleted,
		&delivery.createdAt,
		&delivery.userID,
		&delivery.addressID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DeliveryModel{}, fmt.Errorf("delivery not found")
		}
		return DeliveryModel{}, fmt.Errorf("failed to get delivery: %w", err)
	}
	return delivery, nil
}

func (r *MySQLDeliveryRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM deliveries WHERE id = ?`
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

func (r *MySQLDeliveryRepository) DeleteAll(filter DeliveryParams) (uint, error) {
	query := `DELETE FROM deliveries WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.UserID != nil {
		query += " AND user_id = ?"
		args = append(args, *filter.UserID)
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete deliveries: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLDeliveryRepository) Update(id string, newDelivery DeliveryParams) (DeliveryModel, error) {
	previousDelivery, err := r.GetOne(id)
	if err != nil {
		return DeliveryModel{}, err
	}
	// This will check nil arguments and change only the non-nil ones. CANNOT UPDATE USERID
	updatedDelivery := DeliveryModel{
		isActive:  nilcheck.NotNilBool(newDelivery.IsActive, previousDelivery.isActive),
		isDeleted: nilcheck.NotNilBool(newDelivery.IsDeleted, previousDelivery.isDeleted),
		addressID: nilcheck.NotNilString(newDelivery.AddressID, previousDelivery.addressID),
	}
	query := `UPDATE deliveries SET isActive = ?, isDeleted = ?, address_id = ? WHERE id = ?`

	_,
		err = r.db.Exec(query,
		updatedDelivery.isActive,
		updatedDelivery.isDeleted,
		updatedDelivery.addressID,
		id)
	if err != nil {
		return DeliveryModel{}, fmt.Errorf("failed to update deliveries: %w", err)
	}
	return r.GetOne(id)
}
