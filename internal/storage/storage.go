package storage

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

const (
	basePath    = "logger/chat"
	defaultPerm = 0774
)

func New() *Storage {

	//Open db.
	connStr := "user=postgres password=password dbname=order sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Can't open database: %w", err)
	}

	return &Storage{db: db}
}

func (s *Storage) GetInfo(orderId int64) (*Order, error) {

	order := Order{}
	row := s.db.QueryRow("SELECT * FROM orders WHERE id = $1", orderId)
	err := row.Scan(&order.id, &order.owner, &order.sitter, &order.completed)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func CreatePair(order *Order) error {
	if order.completed == true {
		return errors.New("Заказ уже выполнен!")
	}

	pair[order.id] = *order
	ownerSitter[order.owner] = order.sitter
	ownerSitter[order.sitter] = order.owner
	CreateDir(order)
	return nil
}

func IsExists(sender int64) (receiver int64, err error) {
	if receiver, ok := ownerSitter[sender]; ok {
		return receiver, nil
	}
	return 0, fmt.Errorf("Пары не существует!")
}

func DeletePair(message *tgbotapi.Message, receiver int64) {
	delete(ownerSitter, receiver)
	delete(ownerSitter, message.Chat.ID)
}

func CreateDir(order *Order) {
	path := fmt.Sprintf("%v-%v", order.owner, order.sitter)
	fpath := filepath.Join(basePath, path)
	err := os.MkdirAll(fpath, defaultPerm)
	if err != nil {
		log.Fatal("Can't create dir: %w", err)
	}
	createFile(fpath, order)
}

func createFile(fpath string, order *Order) {
	filePath := fmt.Sprintf("%v/%v.log", fpath, order.id)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, defaultPerm)
	if err != nil {
		log.Fatal("Can't create file: %w", err)
	}
	defer file.Close()

	log := fmt.Sprintf("%v: Chat has been created\n", time.Now())
	file.WriteString(log)
}

func GetLogPairs(sender, receiver int64) (folder string, logPair []string) {
	senderStr := strconv.Itoa(int(sender))
	receiverStr := strconv.Itoa(int(receiver))

	dir, err := os.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if strings.Contains(d.Name(), senderStr) && strings.Contains(d.Name(), receiverStr) {
			folder = basePath + "/" + d.Name()
		}
	}

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	logPair, _ = f.Readdirnames(0)

	return folder, logPair
}

func Logging(folderName, s string, sender, receiver int64, date int64, text string) {
	file, err := os.OpenFile(folderName+"/"+s, os.O_APPEND|os.O_WRONLY, defaultPerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log := fmt.Sprintf("%v from %v to %v: %v\n", time.Unix(date, 0), sender, receiver, text)
	file.WriteString(log)
}
