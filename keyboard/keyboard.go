package keyboard

import (
	"Pet-Sitters-Services/internal/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Menu texts
	firstMenu  = command.MSGHELLO
	secondMenu = command.MSGHELP
	faqMenu    = command.FAQ

	// Button texts
	helpButton = "Help"
	backButton = "Back"
	faqButton  = "FAQ"

	numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/help"),
			tgbotapi.NewKeyboardButton("/faq"),
			tgbotapi.NewKeyboardButton("/close"),
		),
	)

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(helpButton, helpButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
			tgbotapi.NewInlineKeyboardButtonData(faqButton, faqButton),
		),
	)

	faqMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
	)
)

func HandleButton(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	switch query.Data {
	case helpButton:
		text = secondMenu
		markup = secondMenuMarkup
	case backButton:
		text = firstMenu
		markup = firstMenuMarkup
	case faqButton:
		text = faqMenu
		markup = faqMenuMarkup
	default:
		text = helpButton
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func SendMenu(bot *tgbotapi.BotAPI, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

func OpenKeyboard(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	msg.ReplyMarkup = numericKeyboard
	bot.Send(msg)
}

func KeyboardClose(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}