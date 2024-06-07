package main

import (
	command "Pet-Sitters-Services/command"
	"bufio"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strings"
)

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	// Store bot screaming status
	screaming = false
	bot       *tgbotapi.BotAPI

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

type Order struct {
	Id         int64
	ConsumerId int64
	SitterId   int64
}

var order1 = Order{Id: 1, ConsumerId: 241621664, SitterId: 6048355505}

func main() {
	var err error
	bot, err = tgbotapi.NewBotAPI("6954948262:AAFx4f8_efENBQ7CDeu0o27d_otTVnAKP4U")
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
		break
	}
}

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
		chat(message, order1)
	} else if message.Photo != nil {
		sendPhoto(message, order1)
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
		err = sendMenu(message.Chat.ID)
		break
	case "help":
		sendHelp(message)
		break
	}

	return err
}

func sendPhoto(message *tgbotapi.Message, order Order) {
	if message.Photo == nil {
		return
	}
	var receiver int64
	sender := message.Chat.ID
	if sender == order.ConsumerId {
		receiver = order.SitterId
	} else if sender == order.SitterId {
		receiver = order.ConsumerId
	} else {
		fmt.Println("Такого заказа нет!")
	}

	photoArray := message.Photo
	lastIndex := len(photoArray) - 1
	photo := photoArray[lastIndex]

	var msg tgbotapi.Chattable

	msg = tgbotapi.NewPhoto(receiver, tgbotapi.FileID(photo.FileID))
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}

}

func chat(message *tgbotapi.Message, order Order) {
	var receiver int64
	sender := message.Chat.ID
	if sender == order.ConsumerId {
		receiver = order.SitterId
	} else if sender == order.SitterId {
		receiver = order.ConsumerId
	} else {
		fmt.Println("Такого заказа нет!")
	}

	msg := tgbotapi.NewMessage(receiver, message.Text)
	msgReply := tgbotapi.NewMessage(sender, "Сообщение отправлено!")

	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}
	if _, err := bot.Send(msgReply); err != nil {
		panic(err)
	}
}

func sendHelp(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, command.MSGHELP)
	msg.ReplyToMessageID = message.MessageID
	bot.Send(msg)
}

func handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == nextButton {
		text = secondMenu
		markup = secondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = firstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}
