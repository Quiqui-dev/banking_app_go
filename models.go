package main

import (
	"database/sql"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type TransferRequest struct {
	ToAccount int64   `json:"to_account"`
	Amount    float64 `json:"amount"`
}

type LoginRequest struct {
	AccountNumber int64  `json:"account_number"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	AccountNumber int64  `json:"account_number"`
	Token         string `json:"token"`
}

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Number            int64     `json:"account_number"`
	EncryptedPassword string    `json:"-"`
	Balance           float64   `json:"balance"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewAccount(FirstName string, LastName string, Password string) (*Account, error) {

	encrypedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         FirstName,
		LastName:          LastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(encrypedPassword),
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func ScanRowIntoAccount(rows *sql.Rows) (*Account, error) {

	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt,
	)

	return account, err
}
