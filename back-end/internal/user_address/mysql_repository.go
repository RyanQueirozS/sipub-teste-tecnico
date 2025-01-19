package user_address

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sipub-test/db"

	"github.com/google/uuid"
)

type MySQLUserAddressRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLUserAddressRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLUserAddressRepository) createNewUserTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS user_address (
        id CHAR(36) NOT NULL,
        user_id char(36) NOT NULL,
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

func NewMySQLUserAddressRepository() *MySQLUserAddressRepository {
	repo := &MySQLUserAddressRepository{db: db.GetDB()}
	repo.createNewUserTableIfNoneExists()
	return repo
}

func (r *MySQLUserAddressRepository) Create(params UserAddressParams) (UserAddressModel, error) {
	id := uuid.NewString()

	// Round price to 2 decimal places, if not, there will be floating number
	// innacuracy
	query := `INSERT INTO user_address (id, user_id, address_id ) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, id, params.UserID, params.AddressID)
	if err != nil {
		return UserAddressModel{}, fmt.Errorf("failed to create userAddress: %w", err)
	}

	return UserAddressModel{
		id:        id,
		UserID:    params.UserID,
		AddressID: params.AddressID,
	}, nil
}

func (r *MySQLUserAddressRepository) GetAll(filter UserAddressParams) ([]UserAddressModel, error) {
	if filter.UserID == "" {
		return nil, fmt.Errorf("Invalid userId")
	}

	query := `SELECT id, user_id, address_id FROM user_address WHERE 1=1`
	args := []interface{}{}
	{ // Add the user id
		query += " AND user_id LIKE ?"
		args = append(args, "%"+*&filter.UserID+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get user_address: %w", err)
	}
	defer rows.Close()

	var userAddresses []UserAddressModel
	for rows.Next() {
		var userAddress UserAddressModel
		if err := rows.Scan(&userAddress.id, &userAddress.UserID, &userAddress.AddressID); err != nil {
			return nil, fmt.Errorf("failed to scan userAddress: %w", err)
		}
		userAddresses = append(userAddresses, userAddress)
	}

	return userAddresses, nil
}

func (r *MySQLUserAddressRepository) GetOne(id string) (UserAddressModel, error) {
	query := `SELECT id, user_id, address_id FROM user_address WHERE id = ?`
	var userAddress UserAddressModel
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&userAddress.id, &userAddress.UserID, &userAddress.AddressID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserAddressModel{}, fmt.Errorf("userAddress not found")
		}
		return UserAddressModel{}, fmt.Errorf("failed to get userAddress: %w", err)
	}
	return userAddress, nil
}

func (r *MySQLUserAddressRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM user_address WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete userAddress: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("failed to delete userAddress: %w", err)
	}
	return uint(count), nil
}

func (r *MySQLUserAddressRepository) DeleteAll(filter UserAddressParams) (uint, error) {
	if filter.UserID != "" {
		return 0, fmt.Errorf("Invalid UserID")
	}
	query := `DELETE FROM user_address WHERE 1=1`
	args := []interface{}{}
	{ // Process the user id
		query += " AND user_id LIKE ?"
		args = append(args, "%"+*&filter.UserID+"%")
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user_address : %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}
