package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int64) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	connStr := os.Getenv("DB_CONN_STR")

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {

	query := `CREATE TABLE IF NOT EXISTS Account(
		ID serial primary key,
		FIRST_NAME varchar(50),
		LAST_NAME varchar(50),
		ACCOUNT_NUMBER serial,
		ENCRYPTED_PASSWORD varchar(60),
		BALANCE real,
		CREATED_AT timestamp
	);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) (*Account, error) {

	query := `INSERT INTO Account (FIRST_NAME, LAST_NAME, ACCOUNT_NUMBER, ENCRYPTED_PASSWORD, BALANCE, CREATED_AT) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *`

	rows, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.EncryptedPassword,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanRowIntoAccount(rows)
	}

	return nil, fmt.Errorf("could not create account")
}

func (s *PostgresStore) DeleteAccount(id int) error {

	_, err := s.db.Query(`DELETE FROM Account WHERE ID = $1`, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM Account a WHERE a.ID = $1`, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanRowIntoAccount(rows)
	}

	return nil, fmt.Errorf("account id %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(accountNumber int64) (*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM Account a WHERE a.ACCOUNT_NUMBER = $1`, accountNumber)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return ScanRowIntoAccount(rows)
	}

	return nil, fmt.Errorf("account not found")
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM Account`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {

		account, err := ScanRowIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
