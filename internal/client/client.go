package client

import (
	"Pet-Sitters-Services/internal/command"
	"Pet-Sitters-Services/internal/storage"
	"Pet-Sitters-Services/keyboard"
	"bufio"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	Id         int64
	ConsumerId int64
	SitterId   int64
}

//var order1 = Order{Id: 1, ConsumerId: 241621664, SitterId: 6048355505}

var (
	bot *tgbotapi.BotAPI
)

var s, _ = storage.New()

var err error

func StartDB() {
	s, err = storage.New()
	if err != nil {
		log.Fatal(err)
	}
}

func StartBot() {
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
		chat(message)
	} else if strings.HasPrefix(text, "*startorder") {
		startOrder(message)
	} else if message.Photo != nil {
		sendPhoto(message)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func startOrder(message *tgbotapi.Message) {
	text := message.Text
	spaceIndex := strings.Index(text, " ")
	numStr := text[spaceIndex+1:]
	num, _ := strconv.Atoi(numStr)
	//msg := tgbotapi.NewMessage(message.Chat.ID, numStr)
	order, _ := s.GetInfo(int64(num))
	msgText := fmt.Sprintf("Chat %v succesfull created!", order)
	msg := tgbotapi.NewMessage(message.Chat.ID, msgText)
	bot.Send(msg)
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
	case "starttimer":
		startTimer(message)
		break
	case "open":
		keyboard.OpenKeyboard(bot, message)
		break
	case "close":
		keyboard.KeyboardClose(bot, message)
		break
	}

	return err
}

func startTimer(message *tgbotapi.Message) {
	t := time.NewTimer(3 * time.Second)
	<-t.C
	msg := tgbotapi.NewMessage(message.Chat.ID, "Пришло время для фотоотчета!")
	bot.Send(msg)
}

func callAdmin(message *tgbotapi.Message) {

}

func sendPhoto(message *tgbotapi.Message) {
	if message.Photo == nil {
		return
	}
	var receiver int64
	sender := message.Chat.ID
	receiver, err := s.IsExists(sender)
	if err != nil {
		panic(err)
	}

	photoArray := message.Photo
	lastIndex := len(photoArray) - 1
	photo := photoArray[lastIndex]

	var msg tgbotapi.Chattable

	msg = tgbotapi.NewPhoto(receiver, tgbotapi.FileID(photo.FileID))
	if _, err := bot.Send(msg); err != nil {
		panic(err)
	}

	startTimer(message)

}

func chat(message *tgbotapi.Message) {
	var receiver int64
	sender := message.Chat.ID
	receiver, err := s.IsExists(sender)
	if err != nil {
		panic(err)
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
