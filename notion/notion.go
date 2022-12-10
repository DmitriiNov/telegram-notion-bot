package notion

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jomei/notionapi"
)

type NotionApi struct {
	Client   *notionapi.Client
	MainPage *notionapi.Page
	BlockId  notionapi.BlockID
}

func NewNotionApi(token, pageId string) (*NotionApi, error) {
	np := &NotionApi{
		Client: notionapi.NewClient(notionapi.Token(token)),
	}
	page, err := np.Client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return nil, err
	}
	np.MainPage = page
	np.BlockId = notionapi.BlockID(page.ID)

	fmt.Println(np.BlockId)
	return np, nil
}

func getTimestamp() string {
	tm := time.Now()
	return fmt.Sprintf("[%s]\n", tm.Format("02.01.06 15:04"))
}

func GetRichText(message *tgbotapi.Message) notionapi.RichText {
	rch := notionapi.RichText{
		Text:        &notionapi.Text{Content: message.Text},
		Annotations: &notionapi.Annotations{Bold: true},
	}
	for _, ent := range message.Entities {
		ent := tgbotapi.MessageEntity(ent)
		url := string([]rune(message.Text)[ent.Offset : ent.Offset+ent.Length])
		fmt.Println(url)
		if ent.Type == "url" {
			rch.Text.Link = &notionapi.Link{Url: url}
			break
		}
	}
	return rch
}

func (n *NotionApi) WriteNewNote(message *tgbotapi.Message) error {
	tm := getTimestamp()
	blocks := make([]notionapi.Block, 1)
	blocks[0] = notionapi.ParagraphBlock{
		Paragraph: notionapi.Paragraph{
			RichText: []notionapi.RichText{
				{Text: &notionapi.Text{Content: tm}},
				GetRichText(message),
			},
		},
		BasicBlock: notionapi.BasicBlock{
			Type:   notionapi.BlockTypeParagraph,
			Object: notionapi.ObjectTypeBlock,
		},
	}
	req := notionapi.AppendBlockChildrenRequest{Children: blocks}
	_, err := n.Client.Block.AppendChildren(context.Background(), n.BlockId, &req)
	if err != nil {
		return err
	}
	return nil
}
