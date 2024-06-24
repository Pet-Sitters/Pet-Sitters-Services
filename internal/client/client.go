package client

import (
	"Pet-Sitters-Services/config"
	"Pet-Sitters-Services/internal/command"
	"Pet-Sitters-Services/internal/storage"
	"Pet-Sitters-Services/keyboard"
	"bufio"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

var (
	bot *tgbotapi.BotAPI
	err error
)

// StartBot функция инициализирует телеграм бота по токену записанному в настройках.
func StartBot() {
	bot, err = tgbotapi.NewBotAPI(config.TG_TOKEN)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Отображает в консоли все взаимодействия с серверами телеграмма.
	bot.Debug = config.DEBUG

	//Переменная для получения обновлений от серверов каждые 60 секунд
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// updates - канал, который получает обновления.
	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

// receiveUpdates - получает обновления из канала и обрабатывает их.
func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// бесконечный цикл
	for {
		select {
		// выход из цикла при отмененном ctx
		case <-ctx.Done():
			return
		// update получает обновление из канала и отправляет его на обработку
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

// handleUpdate - обрабатывает обновления. Если обновление пришло в виде сообщения, то вызывается обработчик сообщений.
// Если обновление в виде активированной кнопки, то вызывается обработчик нажатия клавиш.
func handleUpdate(update tgbotapi.Update) {
	switch {

	case update.Message != nil:
		handleMessage(update.Message)
		break

	case update.CallbackQuery != nil:
		keyboard.HandleButton(bot, update.CallbackQuery)
		break
	}
}

// handleMessage - обработчик сообщений.
func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if message.IsCommand() {
		err = handleCommand(message)
	} else if strings.HasPrefix(text, "*chat") {
		command.Chat(message, bot)
	} else if strings.HasPrefix(text, "*startorder") {
		command.StartOrder(message, bot)
	} else if strings.HasPrefix(text, "*stoporder") {
		command.StopOrder(message, bot)
	} else if message.Photo != nil {
		command.SendPhoto(message, bot)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(message *tgbotapi.Message) error {
	var err error

	command := message.Command()

	switch command {

	case "menu":
		err = keyboard.SendMenu(bot, message.Chat.ID)
		break
	case "help":
		sendHelp(message)
		break
	case "admin":
		callAdmin(message)
		break
	case "open":
		keyboard.OpenKeyboard(bot, message)
		break
	case "close":
		keyboard.KeyboardClose(bot, message)
		break
	case "start":
		sendHello(message)
		break
	case "faq":
		sendFAQ(message)
		break
	}

	return err
}

func sendFAQ(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, command.FAQ)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func callAdmin(message *tgbotapi.Message) {

}

func sendHello(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, command.MSGHELLO)
	storage.CreateUser(message.Chat.ID, message.Chat.UserName)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func sendHelp(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, command.MSGHELP)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}
