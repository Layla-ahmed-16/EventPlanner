package data

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

// UserExists checks if a username already exists
func UserExists(db *sql.DB, username string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)", username).Scan(&exists)
	return exists, err
}

// CreateUser inserts a new user with hashed password
func CreateUser(db *sql.DB, username, password string) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO users(username, password) VALUES($1, $2)", username, string(hashedPwd))
	return err
}

// ValidateUser checks if username exists and password matches
func ValidateUser(db *sql.DB, username, password string) (bool, error) {
	var dbPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&dbPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // user not found
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if err != nil {
		return false, nil // password mismatch
	}

	return true, nil
}
