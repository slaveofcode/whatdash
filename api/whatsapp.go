package api

import (
	"net/http"
	"time"
	"whatdash/wa"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsApp struct {
	WA *wa.ActiveConnections
}

func (c *WhatsApp) Login(w http.ResponseWriter, r *http.Request) {
	number := "6287886837648"

	conn := c.WA.Get(number)

	if conn == nil {
		// create new connection
		wac, err := whatsapp.NewConn(8 * time.Second)
		if err != nil {
			ResponseJSON(w, 200, []byte(`{"status": "error", "err": "`+err.Error()+`"}`))
		}

		stringQr := make(chan string)

		wawa := wa.WA{Conn: wac}
		go func(number string, wawa *wa.WA, wac *whatsapp.Conn, c *WhatsApp, stringQr chan string) {
			sess, _ := wawa.LoginAccount(number, stringQr)

			c.WA.Add(number, wac, sess)
		}(number, &wawa, wac, c, stringQr)

		ResponseJSON(w, 200, []byte(`{"status": "create", "qr": "`+<-stringQr+`"}`))

		return
	}

	// return existing connection
	ResponseJSON(w, 200, []byte(`{"status": "exist"}`))
	return
}

func (c *WhatsApp) LoginOnExisting(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

func (c *WhatsApp) Logout(w http.ResponseWriter, r *http.Request) {

}

func (c *WhatsApp) SendMsg(w http.ResponseWriter, r *http.Request) {
	wrapper := c.WA.Get("6287886837648")

	if wrapper == nil {
		ResponseJSON(w, 404, []byte(`{"status": "fail"}`))
		return
	}

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "6285716114426@s.whatsapp.net",
		},
		Text: "Hello from API",
	}

	wrapper.Conn.Send(msg)

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}
