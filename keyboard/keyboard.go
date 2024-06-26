package keyboard

import (
	"Pet-Sitters-Services/internal/command"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Текст меню
	firstMenu  = command.MSGHELLO
	secondMenu = command.MSGHELP
	faqMenu    = command.FAQ

	// Кнопки с текстом
	helpButton = "Help"
	backButton = "Back"
	faqButton  = "FAQ"

	// Кнопки встроенной клавиатуры снизу
	numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/help"),
			tgbotapi.NewKeyboardButton("/faq"),
			tgbotapi.NewKeyboardButton("/close"),
		),
	)

	// Раскладка клавиатуры для первого меню. Одна кнопка, по одной в ряду
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(helpButton, helpButton),
		),
	)

	// Раскладка клавиатуры для второго меню. Две кнопки, по одной в ряду
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
			tgbotapi.NewInlineKeyboardButtonData(faqButton, faqButton),
		),
	)

	// Раскладка клавиатуры для первого меню. Одна кнопка, по одной в ряду
	faqMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
	)
)

// HandleButton - функция для обработки действия с кнопками
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

	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

// SendMenu - функция для отправки меню
func SendMenu(bot *tgbotapi.BotAPI, chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)
	return err
}

// OpenKeyboard - функция для открытия встроенной клавиатуры
func OpenKeyboard(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)

	msg.ReplyMarkup = numericKeyboard
	bot.Send(msg)
}

// KeyboardClose - функция для закрытия встроенной клавиатуры
func KeyboardClose(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	bot.Send(msg)
}
