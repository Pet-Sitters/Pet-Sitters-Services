package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

type Pair map[int64]Order

type OwnerSitter map[int64]int64

var ownerSitter = make(OwnerSitter)

var pair = make(Pair)

func New() (*Storage, error) {

	//Open db.
	connStr := "user=postgres password=password dbname=order sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
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

func (s *Storage) GetInfo(orderId int64) (*Order, error) {
	order := Order{}

	row := s.db.QueryRow("SELECT * FROM orders WHERE id = $1", orderId)

	err := row.Scan(&order.id, &order.owner, &order.sitter, &order.completed)
	if err != nil {
		return nil, err
	}

	fmt.Println(order)
	pair[orderId] = order
	ownerSitter[order.owner] = order.sitter
	ownerSitter[order.sitter] = order.owner
	return &order, nil
}

func (s *Storage) GetAll() {
	q := "SELECT * FROM orders"
	_, err := s.db.Exec(q)
	if err != nil {
		return
	}
}

func (s *Storage) IsExists(sender int64) (receiver int64, err error) {
	if receiver, ok := ownerSitter[sender]; ok {
		return receiver, nil
	}
	return 0, fmt.Errorf("sender not exists")
}
