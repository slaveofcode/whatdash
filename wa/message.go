package wa

import (
	"fmt"
	"time"

	"github.com/slaveofcode/go-whatsapp"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type MsgInfo struct {
	ID              string                 `bson:"id"`
	RemoteJid       string                 `bson:"remoteJid"`
	SenderJid       string                 `bson:"senderJid"`
	FromMe          bool                   `bson:"fromMe"`
	Timestamp       uint64                 `bson:"timestamp"`
	PushName        string                 `bson:"pushName"`
	MessageStatus   whatsapp.MessageStatus `bson:"msgStatus"`
	QuotedMessageID string                 `bson:"quotedMessageID"`
}

type WaMsg struct {
	Type string  `bson:"type"`
	Info MsgInfo `bson:"info"`
}

type MsgJSON struct {
	WaMsg
	JSON interface{} `bson:"json"`
}

type MsgText struct {
	WaMsg
	ID          bson.ObjectId `bson:"_id"`
	OwnerNumber string        `bson:"ownerNumber"`
	Text        string        `bson:"text"`
}

type MsgMedia struct {
	WaMsg
	ID          bson.ObjectId `bson:"_id"`
	OwnerNumber string        `bson:"ownerNumber"`
	Type        string        `bson:"type"`
	Caption     string        `bson:"caption"`
	Thumb       []byte        `bson:"thumb"`
	Content     []byte        `bson:"content"`
}

type MsgDoc struct {
	WaMsg
	ID          bson.ObjectId `bson:"_id"`
	OwnerNumber string        `bson:"ownerNumber"`
	Type        string        `bson:"type"`
	Title       string        `bson:"title"`
	PageCount   uint32        `bson:"pageCount"`
	Thumb       []byte        `bson:"thumb"`
	Content     []byte        `bson:"content"`
}

const WaMsgCollName = "WaMessages"

type MessageKeeper struct {
	MgoSession *mgo.Session
}

func (mk *MessageKeeper) isMsgTextExist(db *mgo.Database, msgID string) (bool, MsgText) {
	var msg MsgText
	db.C(WaMsgCollName).Find(bson.M{"wamsg.info.id": msgID}).One(&msg)
	return msg != MsgText{}, msg
}

func (mk *MessageKeeper) isMsgMediaExist(db *mgo.Database, msgID string) (bool, MsgMedia) {
	var msg MsgMedia
	db.C(WaMsgCollName).Find(bson.M{"wamsg.info.id": msgID}).One(&msg)
	return msg.OwnerNumber != "", msg
}

func (mk *MessageKeeper) isMsgDocumentExist(db *mgo.Database, msgID string) (bool, MsgDoc) {
	var msg MsgDoc
	db.C(WaMsgCollName).Find(bson.M{"wamsg.info.id": msgID}).One(&msg)
	return msg.OwnerNumber != "", msg
}

func (mk *MessageKeeper) SaveText(text *MsgText) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	db := sess.DB(DBName())

	msgExist, msg := mk.isMsgTextExist(db, text.WaMsg.Info.ID)

	var err error
	if !msgExist {
		fmt.Println("Insert text:", text.Text)
		err = db.C(WaMsgCollName).Insert(text)
	} else {
		// check status of message
		if msg.WaMsg.Info.MessageStatus != text.Info.MessageStatus {
			fmt.Println("Update text:", text.Text)
			// update status of message
			err = db.C(WaMsgCollName).Update(bson.M{"_id": msg.ID}, bson.M{"$set": bson.M{"wamsg.info.msgStatus": text.Info.MessageStatus, "updatedAt": time.Now()}})
		}
	}

	return err
}

func (mk *MessageKeeper) SaveMedia(media *MsgMedia) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	db := sess.DB(DBName())

	msgExist, msg := mk.isMsgMediaExist(db, media.WaMsg.Info.ID)

	var err error
	if !msgExist {
		err = db.C(WaMsgCollName).Insert(media)
	} else {
		// check status of message
		if msg.WaMsg.Info.MessageStatus != media.Info.MessageStatus {
			// update status of message
			err = db.C(WaMsgCollName).Update(bson.M{"_id": msg.ID}, bson.M{"$set": bson.M{"wamsg.info.msgStatus": media.Info.MessageStatus, "updatedAt": time.Now()}})
		}
	}

	return err
}

func (mk *MessageKeeper) SaveDocument(doc *MsgDoc) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	db := sess.DB(DBName())

	msgExist, msg := mk.isMsgDocumentExist(db, doc.WaMsg.Info.ID)

	var err error
	if !msgExist {
		err = db.C(WaMsgCollName).Insert(doc)
	} else {
		if msg.WaMsg.Info.MessageStatus != doc.Info.MessageStatus {
			// update status of message
			err = db.C(WaMsgCollName).Update(bson.M{"_id": msg.ID}, bson.M{"$set": bson.M{"wamsg.info.msgStatus": doc.Info.MessageStatus, "updatedAt": time.Now()}})
		}
	}

	return err
}
