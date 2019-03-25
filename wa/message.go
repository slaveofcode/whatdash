package wa

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
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
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
}

type MsgMedia struct {
	WaMsg
	ID          bson.ObjectId `bson:"_id"`
	OwnerNumber string        `bson:"ownerNumber"`
	Type        string        `bson:"type"`
	Caption     string        `bson:"caption"`
	Thumb       []byte        `bson:"thumb"`
	Content     []byte        `bson:"content"`
	FileName    string        `bson:"filename"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
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
	FileName    string        `bson:"filename"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
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
		media.FileName = CreateFileName(media.OwnerNumber, getExtByMime(media.Type))
		mk.StoreBytes(media.FileName, media.Content)
		// cleanup
		media.Content = []byte(``)
		err = db.C(WaMsgCollName).Insert(media)
	} else {
		// check status of message
		if msg.WaMsg.Info.MessageStatus != media.Info.MessageStatus {
			// update status of message
			err = db.C(WaMsgCollName).Update(
				bson.M{"_id": msg.ID},
				bson.M{"$set": bson.M{
					"wamsg.info.msgStatus": media.Info.MessageStatus,
					"thumb":                media.Thumb,
					"updatedAt":            time.Now()},
				},
			)
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
		doc.FileName = CreateFileName(doc.OwnerNumber, getExtByMime(doc.Type))
		mk.StoreBytes(doc.FileName, doc.Content)
		// cleanup
		doc.Content = []byte(``)

		err = db.C(WaMsgCollName).Insert(doc)
	} else {
		if msg.WaMsg.Info.MessageStatus != doc.Info.MessageStatus {
			// update status of message
			err = db.C(WaMsgCollName).Update(bson.M{"_id": msg.ID}, bson.M{"$set": bson.M{"wamsg.info.msgStatus": doc.Info.MessageStatus, "updatedAt": time.Now()}})
		}
	}

	return err
}

func (mk *MessageKeeper) DestroyMessages(number string) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	_, err := sess.DB(DBName()).
		C(WaMsgCollName).
		RemoveAll(bson.M{"ownerNumber": number})

	return err
}

func (mk *MessageKeeper) DestroyFiles(number string) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	var itemBinds []struct {
		FileName string `bson:"filename"`
	}
	err := sess.DB(DBName()).
		C(WaMsgCollName).
		Find(bson.M{"ownerNumber": number, "filename": bson.M{"$exists": true}}).
		All(&itemBinds)

	for _, item := range itemBinds {
		sess.DB(DBName()).GridFS("files").Remove(item.FileName)
	}

	return err
}

func (mk *MessageKeeper) StoreFile(fileName string, fileToStore io.Reader) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	file, err := sess.DB(DBName()).GridFS("files").Create(fileName)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, fileToStore)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (mk *MessageKeeper) StoreBytes(fileName string, bytes []byte) error {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	file, err := sess.DB(DBName()).GridFS("files").Create(fileName)
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}

	file.Close()

	return nil
}

func (mk *MessageKeeper) GetStoredFile(filename string) (*mgo.GridFile, *mgo.Session, error) {
	sess := mk.MgoSession.Copy()

	file, err := sess.DB(DBName()).GridFS("files").Open(filename)

	if err != nil {
		sess.Close()
		return nil, nil, err
	}
	return file, sess, nil
}

func (mk *MessageKeeper) IsMessageOwnerExist(ownerNumber, messageId string) bool {
	sess := mk.MgoSession.Copy()
	defer sess.Close()
	count, _ := sess.DB(DBName()).C(WaMsgCollName).Find(bson.M{"ownerNumber": ownerNumber, "wamsg.info.id": messageId}).Count()
	return count > 0
}

func (mk *MessageKeeper) GetMessageFile(ownerNumber, messageId string) (interface{}, string) {
	sess := mk.MgoSession.Copy()
	defer sess.Close()

	var r struct {
		WaMsg struct {
			Type string `bson:"type"`
		} `bson:"wamsg"`
	}

	db := sess.DB(DBName())
	db.C(WaMsgCollName).Find(bson.M{"ownerNumber": ownerNumber, "wamsg.info.id": messageId}).One(&r)

	switch r.WaMsg.Type {
	case "image", "video", "audio":
		var media MsgMedia
		db.C(WaMsgCollName).Find(bson.M{"ownerNumber": ownerNumber, "wamsg.info.id": messageId}).One(&media)
		return media, media.FileName
	case "document":
		var doc MsgDoc
		db.C(WaMsgCollName).Find(bson.M{"ownerNumber": ownerNumber, "wamsg.info.id": messageId}).One(&doc)
		return doc, doc.FileName
	}

	return nil, ""
}

func CreateFileName(number, ext string) string {
	token := randToken(12)
	return number + "-" + token + ext
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func getExtByMime(mime string) string {
	var ext string

	switch mime {
	case "image/jpeg", "image/jpg":
		ext = ".jpg"
		break
	case "image/gif":
		ext = ".gif"
		break
	case "image/png":
		ext = ".png"
		break
	case "video/mp4":
		ext = ".mp4"
		break
	case "application/octet-stream", "video/3gpp":
		ext = ".3gp"
		break
	case "application/pdf":
		ext = ".pdf"
		break
	}

	return ext
}
