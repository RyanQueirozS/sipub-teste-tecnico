package user

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

type MySQLUserRepository struct {
	db *sql.DB
}

// Mainly used for testing, but could be used elsewhere
func (r *MySQLUserRepository) SetDB(db *sql.DB) { r.db = db }

func (r *MySQLUserRepository) createNewUserTableIfNoneExists() {
	r.db = db.GetDB()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id CHAR(36) NOT NULL,
		isActive BOOLEAN NOT NULL DEFAULT TRUE,
		isDeleted BOOLEAN NOT NULL DEFAULT FALSE,
        createdAt CHAR(19) NOT NULL,
		email CHAR(100) NOT NULL,
		cpf CHAR(11) NOT NULL,
		name VARCHAR(255) NOT NULL,
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

func NewMySQLUserRepository() *MySQLUserRepository {
	repo := &MySQLUserRepository{db: db.GetDB()}
	repo.createNewUserTableIfNoneExists()
	return repo
}

func (r *MySQLUserRepository) Create(params UserParams) (UserModel, error) {
	id := uuid.NewString()

	// Round price to 2 decimal places, if not, there will be floating number
	// innacuracy
	timeCreated := time.Now().Format("2006-01-02 15:04:05")

	query := `INSERT INTO users (id, isActive, isDeleted, createdAt, email, cpf, name) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, id, *params.IsActive, *params.IsDeleted, timeCreated, *params.Email, *params.Cpf, *params.Name)
	if err != nil {
		return UserModel{}, fmt.Errorf("failed to create user: %w", err)
	}

	return UserModel{
		id:        id,
		isActive:  *params.IsActive,
		isDeleted: *params.IsDeleted,
		createdAt: timeCreated,
		email:     *params.Email,
		cpf:       *params.Cpf,
		name:      *params.Name,
	}, nil
}

func (r *MySQLUserRepository) GetAll(filter UserParams) ([]UserModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.Email != nil {
		query += " AND email LIKE ?"
		args = append(args, "%"+*filter.Email+"%")
	}
	if filter.Cpf != nil {
		query += " AND cpf = ?"
		args = append(args, *filter.Cpf)
	}
	if filter.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var user UserModel
		if err := rows.Scan(&user.id, &user.isActive, &user.isDeleted, &user.createdAt, &user.email, &user.cpf, &user.name); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *MySQLUserRepository) GetOne(id string) (UserModel, error) {
	query := `SELECT id, isActive, isDeleted, createdAt, email, cpf, name FROM users WHERE id = ?`
	var user UserModel
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&user.id, &user.isActive, &user.isDeleted, &user.createdAt, &user.email, &user.cpf, &user.name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserModel{}, fmt.Errorf("user not found")
		}
		return UserModel{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *MySQLUserRepository) DeleteOne(id string) (uint, error) {
	query := `DELETE FROM users WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return 0, fmt.Errorf("failed to delete user: %w", err)
	}
	return uint(count), nil
}

func (r *MySQLUserRepository) DeleteAll(filter UserParams) (uint, error) {
	query := `DELETE FROM users WHERE 1=1`
	args := []interface{}{}

	if filter.IsActive != nil {
		query += " AND isActive = ?"
		args = append(args, *filter.IsActive)
	}
	if filter.IsDeleted != nil {
		query += " AND isDeleted = ?"
		args = append(args, *filter.IsDeleted)
	}
	if filter.Email != nil {
		query += " AND email LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}
	if filter.Cpf != nil {
		query += " AND cpf = ?"
		args = append(args, "%"+*filter.Name+"%")
	}
	if filter.Name != nil {
		query += " AND name LIKE ?"
		args = append(args, "%"+*filter.Name+"%")
	}

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to delete users: %w", err)
	}
	count, _ := res.RowsAffected()
	return uint(count), nil
}

func (r *MySQLUserRepository) Update(id string, newUser UserParams) (UserModel, error) {
	previousUser, err := r.GetOne(id)
	if err != nil {
		return UserModel{}, err
	}
	// This will check nil arguments and change only the non-nil ones
	updatedUser := UserModel{
		isActive:  nilcheck.NotNilBool(newUser.IsActive, previousUser.isActive),
		isDeleted: nilcheck.NotNilBool(newUser.IsDeleted, previousUser.isDeleted),
		email:     nilcheck.NotNilString(newUser.Email, previousUser.email),
		cpf:       nilcheck.NotNilString(newUser.Cpf, previousUser.cpf),
		name:      nilcheck.NotNilString(newUser.Name, previousUser.name),
	}
	query := `UPDATE users SET isActive = ?, isDeleted = ?, email = ?, cpf = ?, name = ? WHERE id = ?`

	_, err = r.db.Exec(query, updatedUser.isActive, updatedUser.isDeleted, updatedUser.email, updatedUser.cpf, updatedUser.name, id)
	if err != nil {
		return UserModel{}, fmt.Errorf("failed to update user: %w", err)
	}
	return r.GetOne(id)
}
