package storage

import (
	"Pet-Sitters-Services/config"
	"Pet-Sitters-Services/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	// userMap - переменная для хранения пользователей
	userMap = make(model.UserMap)

	// ownerSitter - переменная для хранения ID пользователей в паре
	ownerSitter = make(model.OrderPair)

	// client - клиент пакета "http" для отправки запросов
	client = &http.Client{}

	// orderMap - переменная для хранения заказов
	orderMap = make(model.OrderMap)
)

const (
	// KEEP_URL путь по которому отправляется запрос для получения информации о передержки.
	//Необходим для получения ID владельца и ситтера.
	KEEP_URL = "http://89.223.123.5/keep/keep_crud/"

	// TOKEN токен авторизации для сервиса. Необходим для возможности отправки запросов.
	TOKEN = "Token " + config.PS_TOKEN

	// OWNER_URL путь по которому отправляется запрос для получения информации о владельце питомца.
	//Необходим для получения телеграм ника владельца.
	OWNER_URL = "http://89.223.123.5/owner/owner_crud/"

	// SITTER_URL путь по которому отправляется запрос для получения информации о ситтере.
	//Необходим для получения телеграм ника ситтера.
	SITTER_URL = "http://89.223.123.5/sitter/sitter_crud/"

	// LOGGER_PATH путь, где хранятся истории чата передержек.
	LOGGER_PATH = "logger/chat"

	// DEFAULT_PERM права доступа к папке. Необходим для записи истории чата передержек.
	DEFAULT_PERM = 0774
)

// GetOrderInfo - функция для получения информации о передержке. Функция отправляет запрос по адресу KEEP_URL.
// Результатами запроса является JSON документ который затем парсится в структуру model.Keep.
// По ID владельца и ситтера, из model.Keep, находятся телеграм ники с помощью getOwnerTgNick и getSitterTgNick.
// sender - telegram ID отправителя, записывается в структуру model.Order.
// На вход принимается номер заказа, telegram ID(sender) и Username(userName).
// return Order, error.
func GetOrderInfo(num, sender int64, userName string) (*model.Order, error) {
	url := KEEP_URL + strconv.FormatInt(num, 10) + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var keep model.Keep
	var order model.Order
	order.ID = num

	if err := json.Unmarshal(body, &keep); err != nil {
		return nil, err
	}

	ownerTgNick, err := getOwnerTgNick(keep.Owner)
	if err != nil {
		return nil, err
	}

	sitterTgNick, err := getSitterTgNick(keep.Sitter)
	if err != nil {
		return nil, err
	}

	if userName == ownerTgNick {
		order.OwnerId = sender
		order.SitterId = getPair(sitterTgNick)
	}

	if userName == sitterTgNick {
		order.SitterId = sender
		order.OwnerId = getPair(ownerTgNick)
	}

	fmt.Println(order)
	fmt.Println(userMap)

	ownerSitter[order.OwnerId] = order.SitterId
	ownerSitter[order.SitterId] = order.OwnerId
	orderMap[keep.ID] = order

	return &order, nil
}

// getPair - функция для получения второго пользователя в паре передержки, если такой имеется.
// На вход принимается telegram Username(username).
// return telegram ID (int64).
func getPair(tgNick string) int64 {
	if receiver, ok := userMap[tgNick]; ok {
		return receiver
	}
	return 0
}

// getOwnerTgNick - функция для получения telegram username владельца. Функция отправляет запрос по адресу OWNER_URL.
// Результатами запроса является JSON документ который затем парсится в структуру model.Owner.
// На вход принимается ID владельца на сервисе(username).
// return telegram Username (string), err.
func getOwnerTgNick(id int64) (string, error) {
	url := OWNER_URL + strconv.FormatInt(id, 10) + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var owner model.Owner

	if err := json.Unmarshal(body, &owner); err != nil {
		return "", err
	}

	return owner.TgNick, err
}

// getSitterTgNick - функция для получения telegram username ситтера. Функция отправляет запрос по адресу SITTER_URL.
// Результатами запроса является JSON документ который затем парсится в структуру model.Sitter.
// На вход принимается ID ситтера на сервисе(username).
// return telegram Username (string), err.
func getSitterTgNick(id int64) (string, error) {
	url := SITTER_URL + strconv.FormatInt(id, 10) + "/"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var sitter model.Sitter

	if err := json.Unmarshal(body, &sitter); err != nil {
		fmt.Println(err)
	}

	return sitter.TgNick, err
}

// IsExists - функция проверяет существование получателя сообщения по отправителю. Проверка осуществляется
// в хеш-таблице содержащей ID пользователей в паре. На вход функция принимает телеграм ID отправителя.
// На выходе функция возвращает телеграм ID получателя и ошибку, при её наличии.
func IsExists(sender int64) (receiver int64, err error) {
	if receiver, ok := ownerSitter[sender]; ok {
		return receiver, nil
	}
	return 0, errors.New("Пары не существует!")
}

// CreatePair - функция проверяет завершен ли заказ. Если его статус не завершен, то вызывается функция createDir.
func CreatePair(order *model.Order) error {
	if order.Status == "done" {
		return errors.New("Заказ уже выполнен!")
	}

	createDir(order)
	return nil
}

// createDir - функция, создающая директорию по пути LOGGER_PATH. На вход функция принимает ссылку на объект
// типа *model.Order. По завершению функции по указанному пути образуется директорию имеющая название
// в формате "Tg_ID_owner-Tg_ID_sitter". В данной папке будут храниться файлы с логами всех чатов передержек
// данной пары.
// Пример, имеется пара пользователей Tg_ID_owner = 123456789 и Tg_ID_sitter = 987654321, тогда
// создастся директория logger/chat/123456789-987654321
func createDir(order *model.Order) {
	path := fmt.Sprintf("%v-%v", order.OwnerId, order.SitterId)
	fpath := filepath.Join(LOGGER_PATH, path)
	err := os.MkdirAll(fpath, DEFAULT_PERM)
	if err != nil {
		log.Fatal("Can't create dir: %w", err)
	}
	createFile(fpath, order)
}

// createFile - функция, создающая файл в директории из createDir. На вход функция принимает директорию и ссылку
// на объект типа *model.Order. По завершению функции в директории создастся файл с именем номера
// передержки и расширением log.
// Пример, имеется пара пользователей Tg_ID_owner = 123456789 и Tg_ID_sitter = 987654321. У данной пары был
// заказ под номером 2 и текущий под номером 5. Тогда в директории logger/chat/123456789-987654321 будет
// два файла - 2.log и 5.log.
func createFile(fpath string, order *model.Order) {
	filePath := fmt.Sprintf("%v/%v.log", fpath, order.ID)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, DEFAULT_PERM)
	if err != nil {
		log.Fatal("Can't create file: %w", err)
	}
	defer file.Close()

	log := fmt.Sprintf("%v: Chat has been created\n", time.Now())
	file.WriteString(log)
}

// GetLogPairs - функция для чтения директории с логами чатов. На вход функции подаётся телеграм ID отправителя и
// получателя сообщения. Функция возвращает директорию в которой хранятся логи чатов и слайс с номерами чатов
// передержек.
func GetLogPairs(sender, receiver int64) (folder string, logPair []string) {
	senderStr := strconv.Itoa(int(sender))
	receiverStr := strconv.Itoa(int(receiver))

	dir, err := os.ReadDir(LOGGER_PATH)
	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if strings.Contains(d.Name(), senderStr) && strings.Contains(d.Name(), receiverStr) {
			folder = LOGGER_PATH + "/" + d.Name()
		}
	}

	f, err := os.Open(folder)
	if err != nil {
		log.Fatal(err)
	}

	logPair, _ = f.Readdirnames(0)

	return folder, logPair
}

// Logging - функция для записи текстового чата в файл. Файл расположен по пути logger/chat.
// Входные переменные:
// - путь к папке, которая имеет название состоящее из телеграм ID владельца и ситтера,
// - имя документа с чатом последней передержки,
// - Телеграм ID отправителя и получателя сообщения,
// - дата сообщения,
// - текст сообщения.
func Logging(folderName, s string, sender, receiver int64, date int64, text string) {
	file, err := os.OpenFile(folderName+"/"+s, os.O_APPEND|os.O_WRONLY, DEFAULT_PERM)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log := fmt.Sprintf("%v from %v to %v: %v\n", time.Unix(date, 0), sender, receiver, text)
	file.WriteString(log)
}

// CreateUser - функция записывает в хеш-таблицу пользователей для быстрого поиска телеграм ID.
// На вход подаётся телеграм ID, и телеграм userName.
func CreateUser(id int64, name string) {
	userMap[name] = id
}

// DeletePair - функция удаляет из хеш-таблицы ownerSitter телеграм ID пользователей пары, после завершения заказа.
func DeletePair(sender, receiver int64) {
	delete(ownerSitter, receiver)
	delete(ownerSitter, sender)
}
