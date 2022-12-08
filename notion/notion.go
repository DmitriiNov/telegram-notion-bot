package notion

import (
	"fmt"

	"github.com/jomei/notionapi"
)

type NotionApi struct {
	Client *notionapi.Client
}

func NewNotionApi(token string) *NotionApi {
	return &NotionApi{
		Client: notionapi.NewClient(notionapi.Token(token)),
	}
}

func (n *NotionApi) WriteNewNote(note string) {
	fmt.Println(note)
}
