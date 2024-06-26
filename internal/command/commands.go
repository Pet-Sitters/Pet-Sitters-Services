package command

import (
	"Pet-Sitters-Services/config"
	"Pet-Sitters-Services/internal/storage"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
)

// StartOrder позволяет создать пару между владельцем питомца и ситтером для их взаимодействия.
// Без этой команды взаимодействие невозможно. Взаимодействие - переписка и отправка фотоотчетов.
// Вызывается нестандартной командой /startorder где после команды указывается номер передержки.
// Вызывает функции storage.GetOrderInfo и storage.CreatePair.
// Принимает на вход сообщение(message) и экземляр бота(bot)
func StartOrder(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var msgText string
	text := message.Text
	spaceIndex := strings.Index(text, " ")
	numStr := text[spaceIndex+1:]
	num, _ := strconv.Atoi(numStr)
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

// StopOrder позволяет завершить пару по окончанию передержки.
// Вызывается нестандартной командой /stoporder. После данной команды номер заказа указывать не надо.
// Вызывает функции storage.IsExists и storage.DeletePair.
// Принимает на вход сообщение(message) и экземляр бота(bot)
func StopOrder(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var msgText string
	var receiver int64
	sender := message.Chat.ID
	receiver, err := storage.IsExists(sender)
	if err != nil {
		log.Printf("An error occured: %s", err.Error())
		msgText = fmt.Sprint("Чат не создан")
	} else {
		storage.DeletePair(sender, receiver)
		msgText = fmt.Sprint("Чат успешно удален!")
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	bot.Send(msg)
}

// SendPhoto позволяет пересылать фотоотчеты между владельцем и ситтером.
// Пересылаемый объект должен быть именно фотографией. Одно фото - одно сообщение.
// Подписи к фотографии не пересылаются.
// Вызывает функции storage.IsExists и startTimer.
// Принимает на вход сообщение(message) и экземляр бота(bot).
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

		go startTimer(message, bot)
	}
}

// startTimer таймер, который предупреждает пользователей о необходимости фотоотчета.
// Время для таймера задаётся в config.
// Принимает на вход сообщение(message) и экземляр бота(bot)
func startTimer(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	t := time.NewTimer(config.TIMER * time.Second)
	<-t.C
	msg := tgbotapi.NewMessage(message.Chat.ID, "Пришло время для фотоотчета!")
	bot.Send(msg)
}

// Chat позволяет пересылать текстовые сообщения между владельцем и ситтером.
// Вызывает функции storage.IsExists и storage.GetLogPairs.
// Принимает на вход сообщение(message) и экземляр бота(bot).
func Chat(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	var receiver int64
	sender := message.Chat.ID

	//получение директории в которой хранятся все передержки соответствующих владельца и ситтера, а также всех
	//заказов данной пары
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

		//Запись чата в файл
		storage.Logging(folderName, logPair[len(logPair)-1], sender, receiver, date, msgText)

		bot.Send(msg)
		bot.Send(msgReply)
	}

}
