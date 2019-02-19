package api

import (
	"encoding/json"
	"net/http"
	"whatdash/wa"

	"gopkg.in/mgo.v2/bson"
)

type Dashboard struct {
	SessionHandler
}

func (c *Dashboard) NewAccount(w http.ResponseWriter, r *http.Request) {

}

func (c *Dashboard) ListConnectedAccounts(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

func (c *Dashboard) ReconnectAccount(w http.ResponseWriter, r *http.Request) {

}

func (c *Dashboard) LoadChats(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number    string `json:"number"`
		RemoteJid string `json:"remoteJid"`
		Count     int    `json:"count"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "please login first"}`))
		return
	}

	mgoSession := c.Bucket.MgoSession.Copy()

	defer mgoSession.Close()

	// Load Texts
	var texts []wa.MsgText
	err = mgoSession.DB(wa.DBName()).
		C(wa.WaMsgCollName).
		Find(bson.M{"ownerNumber": params.Number, "wamsg.info.remoteJid": params.RemoteJid, "wamsg.type": "text"}).
		Sort("-wamsg.info.timestamp").
		Limit(params.Count).
		All(&texts)

	data, _ := json.Marshal(texts)

	ResponseJSON(w, 200, data)
	return
}
