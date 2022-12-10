package telegramnotionbot

import (
	"log"
	"os"
	"strconv"

	"github.com/DmitriiNov/telegram-notion-bot/notion"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() {
	TelegramApiCode := os.Getenv("TELEGRAM_API_CODE")
	NotionApiCode := os.Getenv("NOTION_API_CODE")
	TelegramID, err := strconv.Atoi(os.Getenv("TELEGRAM_ID"))
	if err != nil {
		log.Panic(err)
	}
	NotionPageId := os.Getenv("NOTION_PAGE_ID")
	tgBot := TgBot{BotAPI: nil, Notion: nil, Token: TelegramApiCode, Id: TelegramID}
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
	Id     int
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
			if update.Message.From.ID != int64(bot.Id) {
				continue
			}
			// jsn, _ := json.Marshal(update.Message)
			// fmt.Println(string(jsn))
			if len(update.Message.Text) > 0 {
				go bot.sendNoteToNotion(update.Message)
			}
		}
	}
}

func (bot *TgBot) sendNoteToNotion(message *tgbotapi.Message) {
	err := bot.Notion.WriteNewNote(message)
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
