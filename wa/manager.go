package wa

import (
	"fmt"
	"io"
	"time"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

func Connect() (*whatsapp.Conn, error) {
	// create new connection
	wac, err := whatsapp.NewConn(60 * time.Second)
	wac.SetClientName("WhatDash Dashboard", "WhatDash")

	if err != nil {
		return nil, err
	}

	return wac, nil
}

type Manager struct {
	Conn        *whatsapp.Conn
	OwnerNumber string
}

func (w *Manager) SendMessage(jId, message string) error {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: jId,
		},
		Text: message,
	}
	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending message: %v", err)
	}

	return nil
}

func (w *Manager) SendImage(toNumber string, img io.Reader, caption string) error {
	msg := whatsapp.ImageMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: toNumber + "@s.whatsapp.net",
		},
		Type:    "image/jpeg",
		Caption: caption,
		Content: img,
	}

	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending image: %v\n", err)
	}

	return nil
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
	newSession, err := w.Conn.RestoreSession(session)

	if err != nil {
		return nil, fmt.Errorf("Error during restoring session: %v\n", err)
	}

	return &newSession, nil
}

func (w *Manager) DisconnectSocket() error {
	return w.Conn.CloseSocket()
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

func (w *Manager) TriggerLoadMessage(jId string, count int) error {
	_, err := w.Conn.LoadMessages(jId, count)
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
