package api

import (
	"net/http"
	"whatdash/wa"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type WhatsApp struct {
	Storage *wa.Storage
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

		c.Storage.Add(number, wac, sess)
	}(number, &waMgr, wac, c, stringQr)

	ResponseJSON(w, 200, []byte(`{"status": "create", "qr": "`+<-stringQr+`"}`))

	return
}

func (c *WhatsApp) CheckSession(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"
	wrapper := c.Storage.Get(number)

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
	// wrapper := c.Storage.Get(number)

	// if wrapper == nil {
	// 	ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
	// 	return
	// }

	// waMgr := wa.Manager{Conn: wrapper.Conn}

	// handle closed connection
	// if !waMgr.IsConnected() {
	newConn, err := wa.Connect()

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "fail to reload session connection"}`))
		return
	}

	waMgr := wa.Manager{Conn: newConn}
	succ, _ := waMgr.ReloginAccount(number)
	if !succ {
		ResponseJSON(w, 200, []byte(`{"status": "restore session fail"}`))
	}
	// }

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "6285716114426@s.whatsapp.net",
		},
		Text: "Hello from API",
	}

	newConn.Send(msg)

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}
