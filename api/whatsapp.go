package api

import (
	"net/http"

	whatsapp "github.com/Rhymen/go-whatsapp"
)

type WhatsApp struct {
	session *whatsapp.Session
}

func (c *WhatsApp) Login(w http.ResponseWriter, r *http.Request) {

}

func (c *WhatsApp) LoginOnExisting(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

func (c *WhatsApp) Close(w http.ResponseWriter, r *http.Request) {

}
