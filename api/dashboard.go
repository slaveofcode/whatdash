package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"
	"whatdash/wa"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
)

type AccountSession struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Number    string        `bson:"number" json:"number"`
	JID       string        `json:"jid"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type Dashboard struct {
	SessionHandler
}

// func (c *Dashboard) NewAccount(w http.ResponseWriter, r *http.Request) {

// }

func (c *Dashboard) ListConnectedAccounts(w http.ResponseWriter, r *http.Request) {
	var sessions []AccountSession

	mgoSession := c.Bucket.MgoSession.Copy()
	defer mgoSession.Close()

	err := mgoSession.DB(wa.DBName()).
		C(wa.SessionCollName).
		Find(bson.M{}).
		All(&sessions)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed", "reason": "`+err.Error()+`"}`))
		return
	}

	for key, item := range sessions {
		sessions[key].JID = item.Number + "@s.whatsapp.net"
	}

	data, _ := json.Marshal(sessions)

	ResponseJSON(w, 200, data)
	return
}

func (c *Dashboard) DetailAccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	mgoSession := c.Bucket.MgoSession.Copy()
	defer mgoSession.Close()

	var account AccountSession
	err := mgoSession.DB(wa.DBName()).
		C(wa.SessionCollName).
		Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).
		One(&account)

	if err != nil {
		ResponseJSON(w, 400, []byte(`{"status": "failed", "reason": "`+err.Error()+`"}`))
		return
	}

	data, _ := json.Marshal(account)

	ResponseJSON(w, 200, data)
	return
}

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

func (c *Dashboard) PoolNewMessages(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params struct {
		Number         string `json:"number"`
		RemoteJid      string `json:"remoteJid"`
		First          bool   `json:"first"`
		FirstLoadCount int    `json:"firstLoadCount"`
		LastCount      int    `json:"lastCount"`
	}
	err := decoder.Decode(&params)

	if err != nil {
		ShowError(w, "Invalid request")
		return
	}

	mgoSession := c.Bucket.MgoSession.Copy()
	defer mgoSession.Close()

	query := bson.M{
		"ownerNumber": params.Number,
		"$or": []bson.M{
			bson.M{"wamsg.info.remoteJid": params.RemoteJid},
			bson.M{"wamsg.info.remoteJid": params.Number + "s.whatsapp.net"},
		},
	}

	if params.First {

		var results []map[string]interface{}

		count := 35
		if params.FirstLoadCount != 0 {
			count = params.FirstLoadCount
		}

		mgoSession.DB(wa.DBName()).
			C(wa.WaMsgCollName).
			Find(query).
			Sort("-wamsg.info.timestamp").
			Limit(count).
			All(&results)

		count, _ = mgoSession.DB(wa.DBName()).
			C(wa.WaMsgCollName).
			Find(query).
			Sort("-wamsg.info.timestamp").
			Count()

		r := struct {
			Msg        interface{} `json:"messages"`
			TotalCount int         `json:"totalCount"`
		}{
			Msg:        results,
			TotalCount: count,
		}
		data, _ := json.Marshal(r)

		ResponseJSON(w, 200, data)
		return
	}

	newCount := make(chan int)
	timeout := make(chan bool)

	go func() {
		for count := 0; count < 100; count++ {
			count, _ := mgoSession.DB(wa.DBName()).
				C(wa.WaMsgCollName).
				Find(query).
				Count()

			if count > params.LastCount {
				newCount <- count
				break
			}

			// keeping sock connection alive
			c.keepConnAlive(params.Number)

			// 150 * 100 / 1000 = 15 secs
			time.Sleep(time.Millisecond * 150)
		}

		timeout <- true
	}()

	select {
	case totalNewCount := <-newCount:
		// select new messages and return
		var newResults []map[string]interface{}
		mgoSession.DB(wa.DBName()).
			C(wa.WaMsgCollName).
			Find(query).
			Sort("-wamsg.info.timestamp").
			Limit(totalNewCount - params.LastCount).
			All(&newResults)

		r := struct {
			Msg        interface{} `json:"messages"`
			TotalCount int         `json:"totalCount"`
		}{
			Msg:        newResults,
			TotalCount: totalNewCount,
		}
		data, _ := json.Marshal(r)

		ResponseJSON(w, 200, data)
	case <-timeout:
		ResponseJSON(w, 204, []byte(``))
	}

	return
}
