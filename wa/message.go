package wa

import (
	"gopkg.in/mgo.v2/bson"
)

type WaMsgItem struct {
	ID      bson.ObjectId `bson:"_id"`
	Type    string        `bson:"type"`
	Message interface{}   `bson:"message"`
}

type MsgText struct {
	RemoteJid       string `bson:"remoteJid"`
	SenderJid       string `bson:"senderJid"`
	Message         string `bson:"message"`
	Timestamp       uint64 `bson:"timestamp"`
	PushName        string `bson:"pushName"`
	MessageStatus   int    `bson:"msgStatus"`
	QuotedMessageID string `bson:"quotedMessageID"`
}

type MsgImg struct {
	Type    string `bson:"type"`
	Caption string `bson:"caption"`
	Thumb   []byte `bson:"thumb"`
	URL     string `bson:"url"`
}

const WaMsgCollName = "WaMessages"

type MessageKeeper struct{}

func (mk *MessageKeeper) SaveText(text *MsgText) {
	// sessDB, db := ConnectionOpen()
	// defer ConnectionClose(sessDB)

	// db.C(WaMsgCollName).Find(bson.M{""})
}
