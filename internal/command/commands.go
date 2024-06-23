package command

import (
	"Pet-Sitters-Services/config"
	"Pet-Sitters-Services/internal/storage"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
	"time"
)

// StartOrder позволяет создать пару между владельцем питомца и ситтером для их взаимодействия.
// Без этой команды взаимодействие невозможно. Взаимодействие - переписка и отправка фотоотчетов.
func StartOrder(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var msgText string
	text := message.Text
	spaceIndex := strings.Index(text, " ")
	numStr := text[spaceIndex+1:]
	num, _ := strconv.Atoi(numStr)
	//order, err := s.GetInfo(int64(num))
	order, err := storage.GetOrderInfo(int64(num), message.Chat.ID, message.Chat.UserName)

	if err != nil || order == nil {
		msgText = fmt.Sprintf("Заказ %v не найден", num)
	} else {
		err = storage.CreatePair(order)
		if err != nil {
			msgText = fmt.Sprintf("%v", err)
		} else {
			msgText = fmt.Sprintf("Чат %v успешно создан!", num)
		}
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	bot.Send(msg)

}

// SendPhoto позволяет пересылать фотоотчеты между владельцем и ситтером.
// Пересылаемый объект должен быть именно фотографией. Одно фото - одно сообщение.
// Подписи к фотографии не пересылаются. Данная функция автоматически вызывает таймер - startTimer.
func SendPhoto(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Photo == nil {
		return
	}
	var receiver int64
	sender := message.Chat.ID
	receiver, err := storage.IsExists(sender)
	if err != nil {
		msg := tgbotapi.NewMessage(sender, fmt.Sprintf("%v Сообщение не отправлено!", err))
		bot.Send(msg)
	} else {

		photoArray := message.Photo
		lastIndex := len(photoArray) - 1
		photo := photoArray[lastIndex]

		var msg tgbotapi.Chattable

		msg = tgbotapi.NewPhoto(receiver, tgbotapi.FileID(photo.FileID))
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}

		startTimer(message, bot)
	}
}

// startTimer таймер, который предупреждает пользователей о необходимости фотоотчета.
// Время для таймера задаётся в config
func startTimer(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	t := time.NewTimer(config.TIMER * time.Second)
	<-t.C
	msg := tgbotapi.NewMessage(message.Chat.ID, "Пришло время для фотоотчета!")
	bot.Send(msg)
}

// Chat позволяет пересылать текстовые сообщения между владельцем и ситтером.
func Chat(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var receiver int64
	sender := message.Chat.ID
	folderName, logPair := storage.GetLogPairs(sender, receiver)
	date := int64(message.Date)

	receiver, err := storage.IsExists(sender)
	if err != nil {
		msg := tgbotapi.NewMessage(sender, fmt.Sprintf("%v Сообщение не отправлено!", err))
		bot.Send(msg)
	} else {
		msgText := fmt.Sprintf("%v", message.Text)
		msg := tgbotapi.NewMessage(receiver, msgText)
		msgReply := tgbotapi.NewMessage(sender, fmt.Sprint("Сообщение отправлено!"))
		storage.Logging(folderName, logPair[len(logPair)-1], sender, receiver, date, msgText)

		bot.Send(msg)
		bot.Send(msgReply)
	}

}
