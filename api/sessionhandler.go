package api

import (
	"fmt"
	"whatdash/wa"
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

		// handle existing connection to be replaced
		if wrapper.Conn != nil && wrapper.Conn.IsSocketConnected() {
			wrapper.Conn.Logout()
			s.Bucket.Remove(number)
		}

		newConn, err := wa.Connect()

		if err != nil {
			return waMgr, err
		}

		sessStorage := wa.SessionStorage{MgoSession: s.Bucket.MgoSession}
		session, err := sessStorage.Get(number)
		if err != nil {
			return waMgr, fmt.Errorf("Error during fetching stored session: %v", err)
		}

		waMgr = wa.Manager{Conn: newConn, OwnerNumber: number}
		newSession, err := waMgr.ReloginAccount(session)

		if err == nil {
			// re-store session to file
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
