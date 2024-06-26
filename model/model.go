package model

// Keep - модель заказа
type Keep struct {
	ID        int64  `json:"id"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	OtherPets string `json:"other_pets"`
	Feed      string `json:"feed"`
	PickUp    string `json:"pick_up"`
	Transfer  string `json:"transfer"`
	Status    string `json:"status"`
	Owner     int64  `json:"owner"`
	Sitter    int64  `json:"sitter"`
}

// Owner - модель владельца
type Owner struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Patronym  string `json:"patronym"`
	TgNick    string `json:"tg_nick"`
	TgID      string `json:"tg_id"`
	PhoneNum  string `json:"phone_num"`
	City      string `json:"city"`
	User      int    `json:"user"`
}

// Sitter - модель ситтера
type Sitter struct {
	ID   int `json:"id"`
	User struct {
		Username string `json:"username"`
	} `json:"user"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Patronym  string `json:"patronym"`
	TgNick    string `json:"tg_nick"`
	TgID      string `json:"tg_id"`
	PhoneNum  string `json:"phone_num"`
}

// Order - модель заказа
type Order struct {
	ID       int64
	OwnerId  int64
	SitterId int64
	Status   string
}

// OrderMap - хеш-таблица хранит заказы, где ключом является номер заказа.
type OrderMap map[int64]Order

// UserMap - хеш-таблица хранит пользователей. Ключ - телеграм ник, значение - телеграм ID.
// Необходима для быстрого поиска адресата сообщения в чате передержки.
type UserMap map[string]int64

// OrderPair - хеш-таблица хранит пару пользователей в паре. Ключ - телеграм ID одного из пользователей, значение - телеграм ID
// второго пользователя в паре.
type OrderPair map[int64]int64
