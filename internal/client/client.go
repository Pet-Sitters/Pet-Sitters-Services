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

func StartBot() {
	bot, err = tgbotapi.NewBotAPI(config.TG_TOKEN)
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
		keyboard.HandleButton(bot, update.CallbackQuery)
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
		command.Chat(message, bot)
	} else if strings.HasPrefix(text, "*startorder") {
		command.StartOrder(message, bot)
	} else if strings.HasPrefix(text, "*stoporder") {
		//stopOrder(message)
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
