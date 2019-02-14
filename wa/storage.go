package wa

import (
	whatsapp "github.com/slaveofcode/go-whatsapp"
)

type ConnWrapper struct {
	IsFilled bool
	Conn     *whatsapp.Conn
	Sess     *whatsapp.Session
}

type BucketSession struct {
	Items map[string]ConnWrapper
}

func (c *BucketSession) Sync() {
	// sync with existing session
	var storedSessions WASessions
	(&SessionStorage{}).FetchAll(&storedSessions)

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
	(&SessionStorage{}).Save(number, *sess)
}

// func (c *BucketSession) RenewConn(number string, conn *whatsapp.Conn) {
// 	wrapper := c.Items[number]
// 	wrapper.Conn = conn
// 	c.Items[number] = wrapper
// }

func (c *BucketSession) IsExist(number string) bool {
	return c.Items[number].IsFilled
}

func (c *BucketSession) Remove(number string) {
	delete(c.Items, number)
	(&SessionStorage{}).Destroy(number)
}

func (c *BucketSession) Get(number string) *ConnWrapper {
	wrapper := c.Items[number]
	if !wrapper.IsFilled || wrapper.Sess == nil {
		return nil
	}

	return &wrapper
}

func (c *BucketSession) Reset() {
	sess := SessionStorage{}
	for number := range c.Items {
		delete(c.Items, number)
		sess.Destroy(number)
	}
}
