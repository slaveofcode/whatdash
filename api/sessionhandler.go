package api

import (
	"fmt"
	"time"
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

	if wrapper.Conn != nil && !forceNewSession {
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
		}

		return waMgr, err
	}
}

func (s *SessionHandler) CloseManager(number string, force bool) error {
	removeAccount := func(number string) {
		// remove account
		s.Bucket.Remove(number)

		// remove contacts
		(&wa.ContactStorage{MgoSession: s.Bucket.MgoSession}).DestroyAll(number)

		mk := wa.MessageKeeper{MgoSession: s.Bucket.MgoSession}
		// remove files
		mk.DestroyFiles(number)
		// remove messages
		mk.DestroyMessages(number)
	}

	if !force {
		waMgr, err := s.GetManager(number, false)
		if err != nil && !force {
			return err
		}
		waMgr.LogoutAccount()
		removeAccount(number)
	} else {
		removeAccount(number)
	}

	return nil
}

func (s *SessionHandler) TerminateConn(number string) error {
	waMgr, err := s.GetManager(number, true)
	if err != nil {
		return err
	}
	return waMgr.DisconnectSocket()
}

func CollectContacts(waMgr *wa.Manager, mgoSession *mgo.Session) {
	waMgr.LoadContacts()

	cs := wa.ContactStorage{MgoSession: mgoSession}
	tryCount := 0
	maxTry := 5

	for tryCount < maxTry {
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
		} else {
			time.Sleep(time.Second * 3)
			tryCount++
		}
	}

}
