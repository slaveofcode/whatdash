package wa

import (
	"fmt"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type MsgHandler struct {
	whatsapp.Handler
}

// func (*MsgHandler) HandleError(err error) {
// 	fmt.Println("Error:")
// }

func (*MsgHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println("MSG:", message.Text)
}
func (*MsgHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	fmt.Println("IMG:", message.Thumbnail)
}
func (*MsgHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	fmt.Println("VID:", message.Thumbnail)
}
func (*MsgHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	fmt.Println("AUDIO:", message.Type)
}
func (*MsgHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	fmt.Println("DOC:", message.Title)
}
func (*MsgHandler) HandleJsonMessage(message string) {
	fmt.Println("JSON:", message)
}
