package users

import (
	"database/sql"
	"fmt"
)

type UserRepository struct {
	sqlDatabaseConnection *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{
		sqlDatabaseConnection: database,
	}
}

func (repository *UserRepository) Create(userData *User) (*User, error) {

	const CREATE_USER_QUERY = `INSERT INTO users (email,password) VALUES (?, ?)`

	statement, statementError := repository.sqlDatabaseConnection.Prepare(CREATE_USER_QUERY)

	if statementError != nil {
		return nil, statementError
	}

	createResult, createError := statement.Exec(userData.Email, userData.Password)

	if createError != nil {
		return nil, createError
	}

	lastInsertedId, insertError := createResult.LastInsertId()

	if insertError != nil {
		return nil, insertError
	}

	defer func() {
		statementCloseError := statement.Close()

		if statementCloseError != nil {
			insertError = fmt.Errorf("failed to close statement: %w", statementCloseError)

		}
	}()

	userData.ID = lastInsertedId

	return userData, nil
}

func (repository *UserRepository) FindUserByEmail(email string) (*User, error) {
	var foundUser = &User{}

	const FIND_USER_BY_EMAIL_QUERY = `SELECT id,email,password FROM users WHERE email = ?`

	statement, statementError := repository.sqlDatabaseConnection.Prepare(FIND_USER_BY_EMAIL_QUERY)

	if statementError != nil {
		return nil, statementError
	}

	if findError := statement.QueryRow(email).Scan(&foundUser.ID, &foundUser.Email, &foundUser.Password); findError != nil {
		if findError == sql.ErrNoRows {
			return nil, fmt.Errorf("User Atau Password Tidak Sesuai")
		}
		return nil, fmt.Errorf("%v", findError)
	}

	return foundUser, nil
}
