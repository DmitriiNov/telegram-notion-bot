package telegramnotionbot

import (
	"log"

	"github.com/DmitriiNov/telegram-notion-bot/notion"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() {
	tgBot := TgBot{BotAPI: nil, Notion: nil, Token: TelegramApiCode}
	notion, err := notion.NewNotionApi(NotionApiCode, NotionPageId)
	if err != nil {
		log.Panic(err)
	}
	tgBot.Notion = notion
	tgBot.startBot()
}

type TgBot struct {
	BotAPI *tgbotapi.BotAPI
	Notion *notion.NotionApi
	Token  string
}

func (bot *TgBot) startBot() {
	botAPI, err := tgbotapi.NewBotAPI(bot.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.BotAPI = botAPI
	// bot.BotAPI.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.BotAPI.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.From.ID != TelegramID {
				continue
			}
			if len(update.Message.Text) > 0 {
				go bot.sendNoteToNotion(update.Message)
			}
		}
	}
}

func (bot *TgBot) sendNoteToNotion(message *tgbotapi.Message) {
	err := bot.Notion.WriteNewNote(message.Text)
	newMessageText := "✅"
	if err != nil {
		newMessageText = "❌"
		log.Println(err)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, newMessageText)
	msg.ReplyToMessageID = message.MessageID
	_, err = bot.BotAPI.Send(msg)
	if err != nil {
		log.Println(err)
	}
}
