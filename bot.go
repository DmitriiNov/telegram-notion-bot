package telegramnotionbot

import (
	"log"

	"github.com/DmitriiNov/telegram-notion-bot/notion"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	TelegramID = 0
)

func Run() {
	tgToken, botToken := readConfig()
	notionApi := notion.NewNotionApi(botToken)
	startBot(tgToken, notionApi)
}

func startBot(token string, notionApi *notion.NotionApi) {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.From.ID != TelegramID {
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func readConfig() (string, string) {
	return "tgToken", "botToken"
}
