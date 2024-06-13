package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type Order struct {
	id        int64
	owner     int64
	sitter    int64
	completed bool
}

type Storage struct {
	db *sql.DB
}

//	type Storage interface {
//		Update(ctx context.Context, order *Order) error
//		GetInfo(ctx context.Context, order int64) (*Order, error)
//	}

func New() (*Storage, error) {

	//Open db.
	connStr := "user=postgres password=password dbname=order sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	//Checking connection with file with db.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Update(ctx context.Context, order *Order) error {

	_, err := s.db.Exec("UPDATE orders SET completed = $1 WHERE id = $2", true, order.id)
	if err != nil {
		return (err)
	}
	return nil
}

func (s *Storage) GetInfo(ctx context.Context, orderId int64) (*Order, error) {
	order := Order{}

	row := s.db.QueryRow("SELECT * FROM orders WHERE id = $1", orderId)

	err := row.Scan(&order.id, &order.owner, &order.sitter, &order.completed)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *Storage) GetAll() {
	q := "SELECT * FROM orders"
	_, err := s.db.Exec(q)
	if err != nil {
		return
	}
}
