package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
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
		BALANCE real,
		CREATED_AT timestamp
	);`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {

	query := `INSERT INTO Account (FIRST_NAME, LAST_NAME, ACCOUNT_NUMBER, BALANCE, CREATED_AT) 
		VALUES ($1, $2, $3, $4, $5)`

	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		log.Printf("Could not insert to db: %+v\n", err.Error())
		return err
	}

	fmt.Printf("%+v\n", resp)

	resp.Close()

	return nil
}

func (s *PostgresStore) DeleteAccount(int) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(int) (*Account, error) {
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {

	rows, err := s.db.Query(`SELECT * FROM Account`)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {

		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
