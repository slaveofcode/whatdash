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
		return waMgr, fmt.Errorf("Number not registered")
	}

	if wrapper.Conn != nil && wrapper.Conn.IsSocketConnected() {
		waMgr = wa.Manager{Conn: wrapper.Conn}
	} else {
		newConn, err := wa.Connect()

		if err != nil {
			return waMgr, err
		}

		waMgr = wa.Manager{Conn: newConn}
		succ, err := waMgr.ReloginAccount(number)
		s.Bucket.RenewConn(number, newConn)

		if !succ {
			return waMgr, err
		}
	}

	return waMgr, nil
}
