package api

import (
	"fmt"
	"whatdash/wa"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type SessionHandler struct {
	Bucket *wa.BucketSession
}

func (s *SessionHandler) GetManager(number string, forceNewSession bool) (wa.Manager, error) {
	wrapper := s.Bucket.Get(number)

	var waMgr wa.Manager

	if wrapper == nil {
		return waMgr, fmt.Errorf("Session number not registered")
	}

	if wrapper.Conn != nil && wrapper.Conn.IsSocketConnected() && !forceNewSession {
		waMgr = wa.Manager{Conn: wrapper.Conn, OwnerNumber: number}
		return waMgr, nil
	} else {
		newConn, err := wa.Connect()

		if err != nil {
			return waMgr, err
		}

		// Session fetch
		sessStorage := wa.SessionStorage{MgoSession: s.Bucket.MgoSession}
		session, err := sessStorage.Get(number)
		if err != nil {
			return waMgr, fmt.Errorf("Error during fetching stored session: %v", err)
		}

		waMgr = wa.Manager{Conn: newConn, OwnerNumber: number}
		newSession, err := waMgr.ReloginAccount(session)

		if err == nil {
			// re-store new session
			s.Bucket.Save(number, newConn, newSession)

			// added message handler
			waMgr.SetupHandler(&wa.MsgHandler{
				MgoSession:  s.Bucket.MgoSession,
				OwnerNumber: number,
			})

			// collecting contacts
			go collectContacts(&waMgr, s.Bucket.MgoSession)
		}

		return waMgr, err
	}
}

func (s *SessionHandler) keepConnAlive(number string) {
	wrapper := s.Bucket.Get(number)

	if wrapper.Conn == nil || !wrapper.Conn.IsSocketConnected() {
		s.GetManager(number, true)
	}
}

func (s *SessionHandler) CloseManager(number string, force bool) error {
	waMgr, err := s.GetManager(number, false)
	if err != nil && !force {
		return err
	}

	if err == nil {
		waMgr.LogoutAccount()
		s.Bucket.Remove(number)
	} else if err != nil && force {
		s.Bucket.Remove(number)
	}

	return nil
}

func (s *SessionHandler) TerminateConn(number string) error {
	waMgr, err := s.GetManager(number, false)
	if err != nil {
		return err
	}
	return waMgr.DisconnectSocket()
}

func collectContacts(waMgr *wa.Manager, mgoSession *mgo.Session) {
	waMgr.LoadContacts()

	cs := wa.ContactStorage{MgoSession: mgoSession}

	for {
		contacts := waMgr.GetContacts()
		if len(contacts) > 0 {
			// save contact
			for jid, contact := range contacts {
				err, _ := cs.Get(waMgr.OwnerNumber, jid)
				if err != nil {
					// means contact not found
					errSaving := cs.Save(&wa.Contact{
						ID:          bson.NewObjectId(),
						OwnerNumber: waMgr.OwnerNumber,
						JID:         jid,
						Contact: &wa.WaContact{
							Jid:    contact.Jid,
							Notify: contact.Notify,
							Name:   contact.Name,
							Short:  contact.Short,
						},
					})

					if errSaving != nil {
						fmt.Println("Error saving contact:", jid)
					}
				}
			}

			// exist loop
			break
		}
	}

}
