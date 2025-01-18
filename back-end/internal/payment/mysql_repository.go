package payment

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"
	"time"

	"github.com/google/uuid"
)

type MySQLPaymentRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLPaymentRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLPaymentRepository) createNewPaymentTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS payments (
		id CHAR(36) NOT NULL,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
        createdAt CHAR(19) NOT NULL,
        delivery_id CHAR(16) NOT NULL,
        value FLOAT NOT NULL,
        FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE,
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

func NewMySQLPaymentRepository() *MySQLPaymentRepository {
	repo := &MySQLPaymentRepository{db: db.GetDB()}
	repo.createNewPaymentTableIfNoneExists()
	return repo
}

func (r *MySQLPaymentRepository) Create(params PaymentParams) (PaymentModel, error) {
	id := uuid.NewString()

	timeCreated := time.Now().Format("2006-01-02 15:04:05")

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields
	model := PaymentModel{
		id:         id,
		isDeleted:  *params.IsDeleted,
		createdAt:  timeCreated,
		deliveryID: *params.DeliveryID,
		value:      *params.Value,
	}

	query := `INSERT INTO payments (id, isDeleted, createdAt, delivery_id, value) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, id, model.isDeleted, timeCreated, model.deliveryID, model.value)
	if err != nil {
		return PaymentModel{}, fmt.Errorf("failed to create payment: %w", err)
	}

	return model, nil
}

func (r *MySQLPaymentRepository) GetAll(filter PaymentParams) ([]PaymentModel, error) {
	query := `
		SELECT 
			p.id, p.isDeleted, p.createdAt, p.deliveryID, p.value
		FROM 
			userDelivery ud
		JOIN 
			deliveries d ON ud.delivery_id = d.delivery_id
		JOIN 
			payments p ON p.delivery_id = d.delivery_id
		WHERE 
			1=1 AND ud.user_id = ?`
	args := []interface{}{*filter.UserID}

	// Add filtering for user_id if provided

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}
	defer rows.Close()

	var payments []PaymentModel
	for rows.Next() {
		var payment PaymentModel
		err := rows.Scan(
			&payment.id,
			&payment.isDeleted,
			&payment.createdAt,
			&payment.deliveryID,
			&payment.value,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *MySQLPaymentRepository) GetOne(id string) (PaymentModel, error) {
	query := `SELECT id, isDeleted, createdAt, delivery_id, value FROM payments WHERE id = ?`

	var payment PaymentModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&payment.id, &payment.isDeleted, &payment.createdAt, &payment.deliveryID, &payment.value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PaymentModel{}, fmt.Errorf("payment not found")
		}
		return PaymentModel{}, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

func (r *MySQLPaymentRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM payments WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete payment: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("no payment found with the given ID")
	}
	return uint(count), nil
}

func (r *MySQLPaymentRepository) DeleteAll(filter PaymentParams) (uint, error) {
	// SQL query to delete payments associated with a specific user
	query := `
		DELETE p
		FROM payments p
		JOIN deliveries d ON p.deliveryID = d.deliveryID
		JOIN userDelivery ud ON ud.deliveryID = d.deliveryID
		WHERE 1=1 AND ud.user_id = ?`
	args := []interface{}{*filter.UserID}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete payments: %w", err)
	}

	// Get the number of rows affected
	count, _ := res.RowsAffected()
	return uint(count), nil
}
