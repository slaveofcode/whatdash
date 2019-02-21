package api

import (
	"encoding/json"
	"net/http"
	"whatdash/wa"

	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type WhatsApp struct {
	SessionHandler
}

func (c *WhatsApp) CreateSession(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	wac, err := wa.Connect()

	if err != nil {
		ResponseJSON(w, 200, []byte(`{"status": "error", "err": "`+err.Error()+`"}`))
		return
	}

	stringQr := make(chan string)

	waMgr := wa.Manager{Conn: wac}
	go func(number string, waMgr *wa.Manager, wac *whatsapp.Conn, c *WhatsApp, stringQr chan string) {
		sess, _ := waMgr.LoginAccount(number, stringQr)
		c.Bucket.Save(number, wac, sess)
	}(params.Number, &waMgr, wac, c, stringQr)

	ResponseJSON(w, 200, []byte(`{"status": "created", "qr": "`+<-stringQr+`"}`))

	return
}

func (c *WhatsApp) CheckSession(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	wrapper := c.Bucket.Get(params.Number)

	if wrapper == nil {
		ResponseJSON(w, 400, []byte(`{"status": "unregistered"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "registered"}`))
	return
}

func (c *WhatsApp) Destroy(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
		Force  bool   `json:"force"`
	}

	err := decoder.Decode(&params)
	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	err = c.CloseManager(params.Number, params.Force)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"destroyed": false}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"destroyed": true}`))
	return
}

func (c *WhatsApp) SendText(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		From    string `json:"from"`
		To      string `json:"to"`
		GroupID string `json:"groupId"`
		Message string `json:"message"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.From, false)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.SendMessage(params.To, params.Message, params.GroupID)

	if err != nil {
		ResponseJSON(w, 200, []byte(`{"status": "fail"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "sent"}`))
	return
}

func (c *WhatsApp) GetContacts(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, false)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	contacts := waMgr.GetContacts()

	data, _ := json.Marshal(contacts)

	ResponseJSON(w, 200, data)
	return
}

func (c *WhatsApp) TriggerLoadMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number       string `json:"number"`
		Jid          string `json:"jid"`
		MsgCount     int    `json:"messageCount"`
		ReloadSocket bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.TriggerLoadMessage(params.Jid, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}

func (c *WhatsApp) TriggerLoadNewMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number        string `json:"number"`
		Jid           string `json:"jid"`
		LastMessageID string `json:"lastMessageId"`
		MsgCount      int    `json:"messageCount"`
		ReloadSocket  bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.TriggerLoadNextMessage(params.Jid, params.LastMessageID, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}

func (c *WhatsApp) TriggerLoadOldMessage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number       string `json:"number"`
		Jid          string `json:"jid"`
		MessageID    string `json:"messageId"`
		MsgCount     int    `json:"messageCount"`
		ReloadSocket bool   `json:"reloadSocket"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	waMgr, err := c.GetManager(params.Number, params.ReloadSocket)
	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	err = waMgr.TriggerLoadPrevMessage(params.Jid, params.MessageID, params.MsgCount)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed"}`))
		return
	}

	ResponseJSON(w, 200, []byte(`{"status": "requested"}`))
	return
}
