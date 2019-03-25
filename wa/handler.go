package wa

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
		Text:      message.Text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
		Type:      message.Type,
		Caption:   message.Caption,
		Thumb:     message.Thumbnail,
		Content:   content,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
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
		Type:      message.Type,
		Caption:   message.Caption,
		Thumb:     message.Thumbnail,
		Content:   content,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
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
		Type:      message.Type,
		Content:   content,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
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
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("Error save doc", err)
	}
}

func (m *MsgHandler) HandleJsonMessage(message string) {
	fmt.Println("JSON:", message)

	if strings.Contains(message, "Msg") || strings.Contains(message, "Presence") || strings.Contains(message, "Cmd") {
		var msg []interface{}
		json.Unmarshal([]byte(message), &msg)

		if msg[0] == "Msg" {
			fmt.Println("JSON Message:", msg[1])
			parsed, _ := msg[1].(map[string]interface{})
			if parsed["cmd"] == "ack" {
				fmt.Println("JSON ACK:", parsed["from"], parsed["to"], parsed["id"], parsed["ack"])
			}
		}

		if msg[0] == "Presence" {
			fmt.Println("JSON Presence:", msg[1])
			parsed, _ := msg[1].(map[string]interface{})
			fmt.Println("Presence:", parsed["id"], parsed["type"], parsed["t"])
		}

		if msg[0] == "Cmd" {
			fmt.Println("JSON Command:", msg[1])
			parsed, _ := msg[1].(map[string]interface{})
			if parsed["type"] == "disconnect" {
				// user has disconnect the client, forcefully destroy the existing session
				// (&SessionStorage{MgoSession: m.MgoSession}).Destroy(m.OwnerNumber)
			}
		}
	}
}
