package api

import (
	"fmt"
	"whatdash/wa"
)

type SessionHandler struct {
	Bucket *wa.BucketSession
}

func (s *SessionHandler) GetManager(number string) (wa.Manager, error) {
	wrapper := s.Bucket.Get(number)

	var waMgr wa.Manager

	if wrapper == nil {
		return waMgr, fmt.Errorf("Session number not registered")
	}

	if wrapper.Conn != nil && wrapper.Conn.IsSocketConnected() {
		waMgr = wa.Manager{Conn: wrapper.Conn, OwnerNumber: number}
	} else {
		newConn, err := wa.Connect()

		if err != nil {
			return waMgr, err
		}

		sessStorage := wa.SessionStorage{MgoSession: s.Bucket.MgoSession.Copy()}
		session, err := sessStorage.Get(number)
		if err != nil {
			return waMgr, fmt.Errorf("Error during fetching stored session: %v", err)
		}

		waMgr = wa.Manager{Conn: newConn, OwnerNumber: number}
		newSession, err := waMgr.ReloginAccount(session)

		// re-store session to file
		s.Bucket.Save(number, newConn, newSession)

		// added message handler
		waMgr.SetupHandler(s.Bucket.MgoSession.Copy())
	}

	return waMgr, nil
}

func (s *SessionHandler) CloseManager(number string) error {
	waMgr, err := s.GetManager(number)
	if err != nil {
		return err
	}

	waMgr.LogoutAccount()
	s.Bucket.Remove(number)

	return nil
}
