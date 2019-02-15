package wa

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"

	whatsapp "github.com/slaveofcode/go-whatsapp"
	"gopkg.in/mgo.v2/bson"
)

type WASession struct {
	number  string
	session whatsapp.Session
}

type WASessions []WASession

const SessionCollName = "savedSessions"

type SavedSession struct {
	ID        bson.ObjectId `bson:"_id"`
	Number    string        `bson:"number`
	Session   []byte        `bson:"session"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}

type SessionStorage struct{}

func (s *SessionStorage) storePath() string {
	return os.TempDir() + "/whatdash/"
}

func (s *SessionStorage) FetchAll(storedSessions *WASessions) {
	dbSess, db := ConnectionOpen()
	defer ConnectionClose(dbSess)

	var savedSessions []SavedSession
	err := db.C(SessionCollName).Find(bson.M{}).All(&savedSessions)

	if err != nil {
		fmt.Println("Fetching session error", err)
	}

	for _, sess := range savedSessions {
		session := whatsapp.Session{}

		reader := bytes.NewReader(sess.Session)
		decoder := gob.NewDecoder(reader)
		err = decoder.Decode(&session)
		*storedSessions = append(*storedSessions, WASession{
			number:  sess.Number,
			session: session,
		})
	}

}

func (s *SessionStorage) Save(number string, session whatsapp.Session) error {
	dbSess, db := ConnectionOpen()
	defer ConnectionClose(dbSess)

	var dummyBuff bytes.Buffer
	encoder := gob.NewEncoder(&dummyBuff)
	err := encoder.Encode(session)
	if err != nil {
		return err
	}

	var existingSession SavedSession
	err = db.C(SessionCollName).Find(bson.M{"number": number}).One(&existingSession)
	if err == nil && existingSession.ID != "" {
		err = db.C(SessionCollName).Update(bson.M{"number": number}, bson.M{"$set": bson.M{"session": dummyBuff.Bytes(), "updatedAt": time.Now()}})
	} else {
		err = db.C(SessionCollName).Insert(&SavedSession{
			ID:        bson.NewObjectId(),
			Number:    number,
			Session:   dummyBuff.Bytes(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})

		if err != nil {
			return fmt.Errorf("Something wrong when Store Session: %v", err)
		}
	}

	return nil
}

func (s *SessionStorage) Get(number string) (whatsapp.Session, error) {
	session := whatsapp.Session{}
	dbSess, db := ConnectionOpen()
	defer ConnectionClose(dbSess)

	var savedSession SavedSession
	err := db.C(SessionCollName).Find(bson.M{
		"number": number,
	}).One(&savedSession)

	if err != nil {
		return session, err
	}

	reader := bytes.NewReader(savedSession.Session)
	decoder := gob.NewDecoder(reader)
	err = decoder.Decode(&session)

	if err != nil {
		return session, err
	}

	return session, nil
}

func (s *SessionStorage) Destroy(number string) error {
	dbSess, db := ConnectionOpen()
	defer ConnectionClose(dbSess)

	err := db.C(SessionCollName).Remove(bson.M{"number": number})

	return err
}

func (s *SessionStorage) Reset() error {
	dbSess, db := ConnectionOpen()
	defer ConnectionClose(dbSess)

	_, err := db.C(SessionCollName).RemoveAll(bson.M{})

	return err
}
