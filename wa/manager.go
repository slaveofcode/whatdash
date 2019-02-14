package wa

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
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
	Conn *whatsapp.Conn
}

func (w *Manager) SendMessage(toNumber, message string) error {
	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: toNumber + "@s.whatsapp.net",
		},
		Text: message,
	}
	err := w.Conn.Send(msg)
	if err != nil {
		return fmt.Errorf("Error during sending message: %v\n", err)
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

func (w *Manager) IsConnected() bool {
	return w.Conn.IsSocketConnected()
}

func (w *Manager) LoginAccount(number string, qrStorage chan string) (*whatsapp.Session, error) {
	session, err := w.Conn.Login(qrStorage)
	if err != nil {
		return nil, fmt.Errorf("Error during login WhatsApp: %v\n", err)
	}

	err = w.SaveLoginSession(number, session)
	if err != nil {
		return nil, fmt.Errorf("Error on saving session: %v\n", err)
	}

	return &session, nil
}

func (w *Manager) ReloginAccount(number string) (bool, error) {
	session, err := w.FetchSavedSession(number)
	if err != nil {
		return false, fmt.Errorf("Error during fetching stored session: %v\n", err)
	}

	newSession, err := w.Conn.RestoreSession(session)
	if err != nil {
		return false, fmt.Errorf("Error during restoring session: %v\n", err)
	}

	err = w.SaveLoginSession(number, newSession)
	if err != nil {
		return false, fmt.Errorf("Error on saving session: %v\n", err)
	}

	return true, nil
}

func (w *Manager) SaveLoginSession(number string, session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/wa-" + number + ".gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}

func (w *Manager) FetchSavedSession(number string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/wa-" + number + ".gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}
