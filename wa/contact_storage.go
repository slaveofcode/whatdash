package wa

import (
	"fmt"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	whatsapp "github.com/slaveofcode/go-whatsapp"
)

const ContactCollName = "contacts"

type Contact struct {
	ID          string           `bson:"_id"`
	OwnerNumber string           `bson:"ownerNumber"`
	JID         string           `bson:"jid"`
	Contact     whatsapp.Contact `bson:"contact"`
}

type Contacts []Contact

type ContactStorage struct {
	MgoSession *mgo.Session
}

func (s *ContactStorage) Get(number, jid string) (error, Contact) {
	sess := s.MgoSession.Copy()
	defer sess.Close()

	var contact Contact
	err := sess.DB(DBName()).
		C(ContactCollName).
		Find(bson.M{"ownerNumber": number, "jid": jid}).
		One(&contact)

	if err != nil {
		return fmt.Errorf("Error fetch contacts: %v", err), contact
	}

	return nil, contact
}

func (s *ContactStorage) Save(contact *Contact) error {
	sess := s.MgoSession.Copy()
	defer sess.Close()

	return sess.DB(DBName()).
		C(ContactCollName).
		Insert(contact)
}

func (s *ContactStorage) FetchAll(number string) (error, []Contacts) {
	var contacts []Contacts

	sess := s.MgoSession.Copy()
	defer sess.Close()

	// var iface []map[string]interface{}

	err := sess.DB(DBName()).
		C(ContactCollName).
		Find(bson.M{"ownerNumber": number}).
		All(&contacts)

	if err != nil {
		return fmt.Errorf("Error fetch contacts: %v", err), contacts
	}

	return nil, contacts
}
