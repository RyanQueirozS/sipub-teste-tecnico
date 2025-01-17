package address

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

type MySQLAddressRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLAddressRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLAddressRepository) createNewAddressTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS addresses (
		id CHAR(36) NOT NULL,
		isActive BOOLEAN NOT NULL DEFAULT TRUE,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
        createdAt CHAR(19) NOT NULL,
        street VARCHAR(255) NOT NULL,
        number VARCHAR(50) NOT NULL,
        neighborhood VARCHAR(255) NOT NULL,
        complement VARCHAR(255),
        city VARCHAR(255) NOT NULL,
        state CHAR(10) NOT NULL,
        country VARCHAR(100) NOT NULL,
        latitude DECIMAL(10, 8),
        longitude DECIMAL(11, 8),
        name VARCHAR(255),
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

func NewMySQLAddressRepository() *MySQLAddressRepository {
	repo := &MySQLAddressRepository{db: db.GetDB()}
	repo.createNewAddressTableIfNoneExists()
	return repo
}

func (r *MySQLAddressRepository) Create(params AddressParams) (AddressModel, error) {
	id := uuid.NewString()

	timeCreated := time.Now().Format("2006-01-02 15:04:05")

	// Fields might be nil, but they need to be passed empty/defaulted non nil fields
	model := AddressModel{
		id:           id,
		isActive:     *params.IsActive,
		isDeleted:    *params.IsDeleted,
		createdAt:    timeCreated,
		street:       *params.Street,
		number:       *params.Number,
		neighborhood: *params.Neighborhood,
		complement:   nilcheck.NotNilString(params.Complement, ""),
		city:         *params.City,
		state:        *params.State,
		country:      *params.Country,
		latitude:     nilcheck.NotNilFloat32(params.Latitude, 0),
		longitude:    nilcheck.NotNilFloat32(params.Longitude, 0),
		name:         nilcheck.NotNilString(params.Name, ""),
	}

	query := `INSERT INTO addresses 
		(id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, id, model.isActive, model.isDeleted, timeCreated, model.street, model.number, model.neighborhood, model.complement, model.city, model.state, model.country, model.latitude, model.longitude, model.name)
	if err != nil {
		return AddressModel{}, fmt.Errorf("failed to create address: %w", err)
	}

	return model, nil
}

func (r *MySQLAddressRepository) GetAll(filter AddressParams) ([]AddressModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.Street != nil {
		query += " AND street LIKE ?"
		args = append(args, "%"+*filter.Street+"%")
	}
	if filter.City != nil {
		query += " AND city LIKE ?"
		args = append(args, "%"+*filter.City+"%")
	}
	if filter.State != nil {
		query += " AND state = ?"
		args = append(args, *filter.State)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	defer rows.Close()

	var addresses []AddressModel
	for rows.Next() {
		var address AddressModel
		err := rows.Scan(&address.id,
			&address.isActive,
			&address.isDeleted,
			&address.createdAt,
			&address.street,
			&address.number,
			&address.neighborhood,
			&address.complement,
			&address.city,
			&address.state,
			&address.country,
			&address.latitude,
			&address.longitude,
			&address.name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan address: %w", err)
		}
		addresses = append(addresses, address)
	}

	return addresses, nil
}

func (r *MySQLAddressRepository) GetOne(id string) (AddressModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, street, number, neighborhood, complement, city, state, country, latitude, longitude, name FROM addresses WHERE id = ?`

	var address AddressModel
	row := r.db.QueryRow(query, id)
	err := row.Scan(&address.id,
		&address.isActive,
		&address.isDeleted,
		&address.createdAt,
		&address.street,
		&address.number,
		&address.neighborhood,
		&address.complement,
		&address.city,
		&address.state,
		&address.country,
		&address.latitude,
		&address.longitude,
		&address.name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AddressModel{}, fmt.Errorf("address not found")
		}
		return AddressModel{}, fmt.Errorf("failed to get address: %w", err)
	}
	return address, nil
}

func (r *MySQLAddressRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM addresses WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete address: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("no address found with the given ID")
	}
	return uint(count), nil
}

func (r *MySQLAddressRepository) DeleteAll(filter AddressParams) (uint, error) {
	query := `DELETE FROM addresses WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.Street != nil {
		query += " AND street LIKE ?"
		args = append(args, "%"+*filter.Street+"%")
	}
	if filter.Number != nil {
		query += " AND number = ?"
		args = append(args, *filter.Number)
	}
	if filter.Neighborhood != nil {
		query += " AND neighborhood LIKE ?"
		args = append(args, "%"+*filter.Neighborhood+"%")
	}
	if filter.City != nil {
		query += " AND city LIKE ?"
		args = append(args, "%"+*filter.City+"%")
	}
	if filter.State != nil {
		query += " AND state = ?"
		args = append(args, *filter.State)
	}
	if filter.Country != nil {
		query += " AND country LIKE ?"
		args = append(args, "%"+*filter.Country+"%")
	}
	if filter.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete addresses: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLAddressRepository) Update(id string, newAddress AddressParams) (AddressModel, error) {
	previousAddress, err := r.GetOne(id)
	if err != nil {
		return AddressModel{}, err
	}
	// This will check nil arguments and change only the non-nil ones
	updatedAddress := AddressModel{
		isActive:     nilcheck.NotNilBool(newAddress.IsActive, previousAddress.isActive),
		isDeleted:    nilcheck.NotNilBool(newAddress.IsDeleted, previousAddress.isDeleted),
		street:       nilcheck.NotNilString(newAddress.Street, previousAddress.street),
		number:       nilcheck.NotNilString(newAddress.Number, previousAddress.number),
		neighborhood: nilcheck.NotNilString(newAddress.Neighborhood, previousAddress.neighborhood),
		complement:   nilcheck.NotNilString(newAddress.Complement, previousAddress.complement),
		city:         nilcheck.NotNilString(newAddress.City, previousAddress.city),
		state:        nilcheck.NotNilString(newAddress.State, previousAddress.state),
		country:      nilcheck.NotNilString(newAddress.Country, previousAddress.country),
		latitude:     nilcheck.NotNilFloat32(newAddress.Latitude, previousAddress.latitude),
		longitude:    nilcheck.NotNilFloat32(newAddress.Longitude, previousAddress.longitude),
		name:         nilcheck.NotNilString(newAddress.Name, previousAddress.name),
	}
	query := `UPDATE addresses 
		SET isActive = ?, isDeleted = ?, street = ?, number = ?, neighborhood = ?, complement = ?, city = ?, state = ?, country = ?, latitude = ?, longitude = ?, name = ? 
		WHERE id = ?`

	_,
		err = r.db.Exec(query,
		updatedAddress.isActive,
		updatedAddress.isDeleted,
		updatedAddress.street,
		updatedAddress.number,
		updatedAddress.neighborhood,
		updatedAddress.complement,
		updatedAddress.city,
		updatedAddress.state,
		updatedAddress.country,
		updatedAddress.latitude,
		updatedAddress.longitude,
		updatedAddress.name,
		id)
	if err != nil {
		return AddressModel{}, fmt.Errorf("failed to update address: %w", err)
	}
	return r.GetOne(id)
}
