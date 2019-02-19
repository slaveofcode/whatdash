package wa

import (
	"github.com/slaveofcode/go-whatsapp"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MsgInfo struct {
	ID              string                 `bson:"id"`
	RemoteJid       string                 `bson:"remoteJid"`
	SenderJid       string                 `bson:"senderJid"`
	Timestamp       uint64                 `bson:"timestamp"`
	PushName        string                 `bson:"pushName"`
	MessageStatus   whatsapp.MessageStatus `bson:"msgStatus"`
	QuotedMessageID string                 `bson:"quotedMessageID"`
}

type WaMsg struct {
	ID   bson.ObjectId `bson:"_id"`
	Type string        `bson:"type"`
	Info MsgInfo       `bson:"info"`
}

type MsgJSON struct {
	WaMsg
	JSON interface{} `bson:"json"`
}

type MsgText struct {
	WaMsg
	OwnerNumber string `bson:"ownerNumber"`
	Text        string `bson:"text"`
}

type MsgMedia struct {
	WaMsg
	OwnerNumber string `bson:"ownerNumber"`
	Type        string `bson:"type"`
	Caption     string `bson:"caption"`
	Thumb       []byte `bson:"thumb"`
	Content     []byte `bson:"content"`
}

type MsgDoc struct {
	WaMsg
	OwnerNumber string `bson:"ownerNumber"`
	Type        string `bson:"type"`
	Title       string `bson:"title"`
	PageCount   uint32 `bson:"pageCount"`
	Thumb       []byte `bson:"thumb"`
	Content     []byte `bson:"content"`
}

const WaMsgCollName = "WaMessages"

type MessageKeeper struct {
	MgoSession *mgo.Session
}

func (mk *MessageKeeper) isMsgExist(db *mgo.Database, msgID string) bool {
	var msg MsgText
	db.C(WaMsgCollName).Find(bson.M{"wamsg.info.id": msgID}).One(&msg)
	return msg != MsgText{}
}

func (mk *MessageKeeper) SaveText(text *MsgText) error {
	defer mk.MgoSession.Close()

	db := mk.MgoSession.Copy().DB(DBName())

	msgExist := mk.isMsgExist(db, text.Info.ID)

	if !msgExist {
		err := db.C(WaMsgCollName).Insert(text)
		return err
	}

	return nil
}

func (mk *MessageKeeper) SaveMedia(media *MsgMedia) error {
	defer mk.MgoSession.Close()

	db := mk.MgoSession.Copy().DB(DBName())

	msgExist := mk.isMsgExist(db, media.Info.ID)

	if !msgExist {
		err := db.C(WaMsgCollName).Insert(media)
		return err
	}

	return nil
}

func (mk *MessageKeeper) SaveDocument(doc *MsgDoc) error {
	defer mk.MgoSession.Close()

	db := mk.MgoSession.Copy().DB(DBName())

	msgExist := mk.isMsgExist(db, doc.Info.ID)

	if !msgExist {
		err := db.C(WaMsgCollName).Insert(doc)
		return err
	}

	return nil
}
