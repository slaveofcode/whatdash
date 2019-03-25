package wa

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

func Connect() (*whatsapp.Conn, error) {
	// create new connection
	wac, err := whatsapp.NewConn((60 * 3) * time.Second)
	wac.SetClientName("WhatDash Dashboard", "WhatDash")

	if err != nil {
		return nil, err
	}

	return wac, nil
}

var randomChars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func getIDLength() int {
	min := 20
	max := 30
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func generateMessageID() string {
	c := make([]rune, getIDLength())
	for i := range c {
		c[i] = randomChars[rand.Intn(len(randomChars))]
	}
	return string(c)
}

type Manager struct {
	Conn        *whatsapp.Conn
	OwnerNumber string
}

func (w *Manager) SendMessage(jId, message string) (error, string) {
	msgID := generateMessageID()

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			Id:        msgID,
			RemoteJid: jId,
		},
		Text: message,
	}
	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending message: %v", err), ""
	}

	return nil, msgID
}

func (w *Manager) SendImage(jId string, img io.Reader, imgType, caption string) (error, string) {
	msgID := generateMessageID()
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			Id:        msgID,
			RemoteJid: jId,
		},
		Type:    imgType,
		Caption: caption,
		Content: img,
	}

	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending image: %v\n", err), ""
	}

	return nil, msgID
}

func (w *Manager) SendVideo(jId string, vid io.Reader, vidType, caption string) (error, string) {
	msgID := generateMessageID()
	msg := whatsapp.VideoMessage{
		Info: whatsapp.MessageInfo{
			Id:        msgID,
			RemoteJid: jId,
		},
		Type:    vidType,
		Caption: caption,
		Content: vid,
	}

	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending video: %v\n", err), ""
	}

	return nil, msgID
}

func (w *Manager) SendDocument(jId string, doc io.Reader, docType, title string) (error, string) {
	msgID := generateMessageID()

	msg := whatsapp.DocumentMessage{
		Info: whatsapp.MessageInfo{
			Id:        msgID,
			RemoteJid: jId,
		},
		Type:    docType,
		Title:   title,
		Content: doc,
	}

	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending document: %v\n", err), ""
	}

	return nil, msgID
}

func (w *Manager) LoginAccount(number string, qrStorage chan string) (*whatsapp.Session, error) {
	session, err := w.Conn.Login(qrStorage)
	if err != nil {
		return nil, fmt.Errorf("Error during login WhatsApp: %v\n", err)
	}

	if err != nil {
		return nil, fmt.Errorf("Error on saving session: %v\n", err)
	}

	return &session, nil
}

func (w *Manager) ReloginAccount(session whatsapp.Session) (*whatsapp.Session, error) {
	newSession, err := w.Conn.RestoreWithSession(session)

	if err != nil {
		return nil, fmt.Errorf("Error during restoring session: %v\n", err)
	}

	return &newSession, nil
}

func (w *Manager) DisconnectSocket() error {
	_, err := w.Conn.Disconnect()
	return err
}

func (w *Manager) LogoutAccount() error {
	return w.Conn.Logout()
}

func (w *Manager) SetupHandler(msgHandler *MsgHandler) {
	w.Conn.AddHandler(msgHandler)
}

func (w *Manager) LoadContacts() error {
	_, err := w.Conn.Contacts()
	return err
}

func (w *Manager) GetContacts() map[string]whatsapp.Contact {
	return w.Conn.Store.Contacts
}

func (w *Manager) TriggerLoadMessage(jId, msgId string, count int) error {
	_, err := w.Conn.LoadMessages(jId, msgId, count)
	return err
}

func (w *Manager) TriggerLoadNextMessage(jId, msgId string, count int) error {
	_, err := w.Conn.LoadMessagesAfter(jId, msgId, count)
	return err
}

func (w *Manager) TriggerLoadPrevMessage(jId, msgId string, count int) error {
	_, err := w.Conn.LoadMessagesBefore(jId, msgId, count)
	return err
}
