package api

import (
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

func (c *WhatsApp) LoginOnExisting(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

func (c *WhatsApp) Logout(w http.ResponseWriter, r *http.Request) {

}

func (c *WhatsApp) SendText(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	waMgr, err := c.GetManager(number)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	waMgr.SendMessage("6285716114426", "Ngapain aja bebas!")

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}
