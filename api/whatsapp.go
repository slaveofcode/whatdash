package api

import (
	"fmt"
	"net/http"
	"whatdash/wa"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type WhatsApp struct {
	SessionHandler
}

func (c *WhatsApp) CreateSession(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	wac, err := wa.Connect()

	if err != nil {
		ResponseJSON(w, 200, []byte(`{"status": "error", "err": "`+err.Error()+`"}`))
	}

	stringQr := make(chan string)

	waMgr := wa.Manager{Conn: wac}
	go func(number string, waMgr *wa.Manager, wac *whatsapp.Conn, c *WhatsApp, stringQr chan string) {
		sess, _ := waMgr.LoginAccount(number, stringQr)
		c.Bucket.Save(number, wac, sess)
	}(number, &waMgr, wac, c, stringQr)

	ResponseJSON(w, 200, []byte(`{"status": "create", "qr": "`+<-stringQr+`"}`))

	return
}

func (c *WhatsApp) CheckSession(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"
	wrapper := c.Bucket.Get(number)

	if wrapper == nil {
		ResponseJSON(w, 400, []byte(`{"status": "unregistered"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "registered"}`))
	return
}

func (c *WhatsApp) Destroy(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	err := c.CloseManager(number)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"destroyed": false}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"destroyed": true}`))
	return
}

func (c *WhatsApp) SendText(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	waMgr, err := c.GetManager(number)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	waMgr.SendMessage("6285716114426", "Pulang pulang...")

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}

func (c *WhatsApp) GetContacts(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	waMgr, err := c.GetManager(number)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	contacts := waMgr.GetContacts()

	fmt.Println(contacts)

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}
