package model

// Сущность заказ
type Keep struct {
	ID        int    `json:"id"`
	FromDate  string `json:"from_date"`
	ToDate    string `json:"to_date"`
	OtherPets string `json:"other_pets"`
	Feed      string `json:"feed"`
	PickUp    string `json:"pick_up"`
	Transfer  string `json:"transfer"`
	Status    string `json:"status"`
	Owner     int    `json:"owner"`
	Sitter    int    `json:"sitter"`
}

// Сущность владелец
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

// Сущность ситтер
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

type OwnerSitter map[int64]int64