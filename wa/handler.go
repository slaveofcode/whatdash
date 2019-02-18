package wa

import (
	"fmt"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type MsgHandler struct {
	whatsapp.Handler
	OwnerNumber string
	MgoSession  *mgo.Session
}

func (*MsgHandler) HandleError(err error) {
	fmt.Println("Error:", err.Error())
}
func (m *MsgHandler) HandleTextMessage(message whatsapp.TextMessage) {
	(&MessageKeeper{MgoSession: m.MgoSession}).SaveText(&MsgText{
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			ID:   bson.NewObjectId(),
			Type: "text",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Text: message.Text,
	})
}
func (m *MsgHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	content, _ := message.Download()
	err := (&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			ID:   bson.NewObjectId(),
			Type: "image",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
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
	fmt.Println(err)
}
func (m *MsgHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	content, _ := message.Download()
	(&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			ID:   bson.NewObjectId(),
			Type: "video",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
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
}
func (m *MsgHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	content, _ := message.Download()
	(&MessageKeeper{MgoSession: m.MgoSession}).SaveMedia(&MsgMedia{
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			ID:   bson.NewObjectId(),
			Type: "audio",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
				Timestamp:       message.Info.Timestamp,
				PushName:        message.Info.PushName,
				MessageStatus:   message.Info.Status,
				QuotedMessageID: message.Info.QuotedMessageID,
			},
		},
		Type:    message.Type,
		Content: content,
	})
}
func (m *MsgHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	content, _ := message.Download()
	(&MessageKeeper{MgoSession: m.MgoSession}).SaveDocument(&MsgDoc{
		OwnerNumber: m.OwnerNumber,
		WaMsg: WaMsg{
			ID:   bson.NewObjectId(),
			Type: "document",
			Info: MsgInfo{
				ID:              message.Info.Id,
				RemoteJid:       message.Info.RemoteJid,
				SenderJid:       message.Info.SenderJid,
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
}
func (*MsgHandler) HandleJsonMessage(message string) {
	// fmt.Println("JSON:", message)
}
