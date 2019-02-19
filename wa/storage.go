package wa

import (
	mgo "github.com/globalsign/mgo"
	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type ConnWrapper struct {
	IsFilled bool
	Conn     *whatsapp.Conn
	Sess     *whatsapp.Session
}

type BucketSession struct {
	Items      map[string]ConnWrapper
	MgoSession *mgo.Session
}

func (c *BucketSession) Sync() {
	// sync with existing session
	var storedSessions WASessions
	(&SessionStorage{MgoSession: c.MgoSession.Copy()}).FetchAll(&storedSessions)

	if len(storedSessions) > 0 {
		for _, item := range storedSessions {
			c.Items[item.number] = ConnWrapper{
				IsFilled: true,
				Conn:     nil,
				Sess:     &item.session,
			}
		}
	}
}

func (c *BucketSession) Save(number string, conn *whatsapp.Conn, sess *whatsapp.Session) {
	c.Items[number] = ConnWrapper{
		IsFilled: true,
		Conn:     conn,
		Sess:     sess,
	}

	// store session to file
	(&SessionStorage{MgoSession: c.MgoSession.Copy()}).Save(number, *sess)
}

func (c *BucketSession) IsExist(number string) bool {
	return c.Items[number].IsFilled
}

func (c *BucketSession) Remove(number string) {
	delete(c.Items, number)
	(&SessionStorage{MgoSession: c.MgoSession.Copy()}).Destroy(number)
}

func (c *BucketSession) Get(number string) *ConnWrapper {
	var wrapperExist bool
	wrapper := c.Items[number]
	if wrapper.IsFilled && wrapper.Sess != nil {
		wrapperExist = true
	}

	if !wrapperExist {
		session, err := (&SessionStorage{MgoSession: c.MgoSession.Copy()}).Get(number)

		if err == nil {
			wrapper = ConnWrapper{
				IsFilled: true,
				Conn:     nil,
				Sess:     &session,
			}
			c.Items[number] = wrapper
		}
	}

	return &wrapper
}

func (c *BucketSession) Reset() {
	sess := SessionStorage{MgoSession: c.MgoSession.Copy()}
	for number := range c.Items {
		delete(c.Items, number)
		sess.Destroy(number)
	}
}
