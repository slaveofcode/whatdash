package wa

import (
	"fmt"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type MsgHandler struct {
	whatsapp.Handler
	MgoSession  *mgo.Session
	OwnerNumber string
}

func (*MsgHandler) HandleError(err error) {
	fmt.Println("Error:", err.Error())
}
func (m *MsgHandler) HandleTextMessage(message whatsapp.TextMessage) {
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveText(&MsgText{
		ID:          bson.NewObjectId(),
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			Type: "text",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				FromMe:          message.Info.FromMe,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Text: message.Text,
	})

	if err != nil {
		fmt.Println("Error save text", err)
	}
}
func (m *MsgHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	content, _ := message.Download()
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		ID:          bson.NewObjectId(),
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			Type: "image",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				FromMe:          message.Info.FromMe,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Type:    message.Type,
		Caption: message.Caption,
		Thumb:   message.Thumbnail,
		Content: content,
	})
	if err != nil {
		fmt.Println("Error save img", err)
	}
}
func (m *MsgHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	content, _ := message.Download()
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		ID:          bson.NewObjectId(),
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{

			Type: "video",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				FromMe:          message.Info.FromMe,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Type:    message.Type,
		Caption: message.Caption,
		Thumb:   message.Thumbnail,
		Content: content,
	})

	if err != nil {
		fmt.Println("Error save video", err)
	}
}
func (m *MsgHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	content, _ := message.Download()
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		ID:          bson.NewObjectId(),
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			Type: "audio",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				FromMe:          message.Info.FromMe,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Type:    message.Type,
		Content: content,
	})

	if err != nil {
		fmt.Println("Error save audio", err)
	}
}
func (m *MsgHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	content, _ := message.Download()
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveDocument(&MsgDoc{
		ID:          bson.NewObjectId(),
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			Type: "document",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				FromMe:          message.Info.FromMe,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Type:      message.Type,
		Title:     message.Title,
		PageCount: message.PageCount,
		Thumb:     message.Thumbnail,
		Content:   content,
	})
	if err != nil {
		fmt.Println("Error save doc", err)
	}
}

func (*MsgHandler) HandleStickerMessage(message whatsapp.StickerMessage) {
	fmt.Println("Sticker message")
	fmt.Println("Sticker here")
}

func (*MsgHandler) HandleJsonMessage(message string) {
	fmt.Println("JSON:", message)
	// if strings.Contains(data, "Msg") || strings.Contains(data, "Presence") {
	//   var msg []interface{}
	//   json.Unmarshal([]byte(data), &msg)

	//   // its now your move to do what you wanted here
	//   if msg[0] == "Msg" { }
	//   else if acknowledgements["cmd"] == "acks" {
	// }
}
