package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"whatdash/wa"

	"gopkg.in/mgo.v2/bson"
)

type Dashboard struct {
	SessionHandler
}

// func (c *Dashboard) NewAccount(w http.ResponseWriter, r *http.Request) {

// }

func (c *Dashboard) ListConnectedAccounts(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, 200, []byte(`{"list": [1,2,3]}`))
}

// func (c *Dashboard) ReconnectAccount(w http.ResponseWriter, r *http.Request) {

// }

func (c *Dashboard) LoadChatHistory(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number string `json:"number"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	mgoSession := c.Bucket.MgoSession.Copy()
	defer mgoSession.Close()

	pipe := []bson.M{
		bson.M{"$match": bson.M{"ownerNumber": params.Number}},
		bson.M{"$group": bson.M{
			"_id": bson.M{
				"jid": "$wamsg.info.remoteJid",
			},
			"msgCount": bson.M{"$sum": 1},
		}},
	}

	var results []struct {
		ID struct {
			JIDNumber string `bson:"jid" json:"jid"`
		} `bson:"_id" json:"wa"`
		Count     int    `bson:"msgCount" json:"msgCount"`
		Timestamp uint64 `bson:"timestamp" json:"lastChatTime"`
	}

	err = mgoSession.DB(wa.DBName()).
		C(wa.WaMsgCollName).
		Pipe(pipe).
		All(&results)

	if err != nil {
		ShowError(w, "Unable to fetch chat history")
	}

	for key, val := range results {
		var timestamp struct {
			OwnerNumber string   `bson:"ownerNumber"`
			Wamsg       wa.WaMsg `bson:"wamsg"`
		}

		err = mgoSession.DB(wa.DBName()).
			C(wa.WaMsgCollName).
			Find(bson.M{"ownerNumber": params.Number, "wamsg.info.remoteJid": val.ID.JIDNumber}).
			Sort("-wamsg.info.timestamp").
			Limit(1).
			One(&timestamp)

		results[key].Timestamp = timestamp.Wamsg.Info.Timestamp
	}

	sort.Slice(results, func(a, b int) bool {
		return results[a].Timestamp > results[b].Timestamp
	})

	data, _ := json.Marshal(results)

	ResponseJSON(w, 200, data)
	return
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

	mgoSession := c.Bucket.MgoSession.Copy()
	defer mgoSession.Close()

	var query bson.M
	if params.RemoteJid != "" {
		query = bson.M{
			"ownerNumber": params.Number,
			"$or": []bson.M{
				bson.M{"wamsg.info.remoteJid": params.RemoteJid},
				bson.M{"wamsg.info.remoteJid": params.Number + "s.whatsapp.net"},
			},
		}
	} else {
		query = bson.M{"ownerNumber": params.Number}
	}

	var results []map[string]interface{}
	err = mgoSession.DB(wa.DBName()).
		C(wa.WaMsgCollName).
		Find(query).
		Sort("-wamsg.info.timestamp").
		Limit(params.Count).
		All(&results)

	data, _ := json.Marshal(results)

	ResponseJSON(w, 200, data)
	return
}
