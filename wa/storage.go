package wa

import whatsapp "github.com/slaveofcode/go-whatsapp"

type ConnWrapper struct {
	IsFilled bool
	Conn     *whatsapp.Conn
	Sess     *whatsapp.Session
}

type Storage struct {
	Connections map[string]ConnWrapper
}

func (c *Storage) Add(number string, conn *whatsapp.Conn, sess *whatsapp.Session) {
	c.Connections[number] = ConnWrapper{
		IsFilled: true,
		Conn:     conn,
		Sess:     sess,
	}
}

func (c *Storage) IsExist(number string) bool {
	return c.Connections[number].IsFilled
}

func (c *Storage) Remove(number string) {
	delete(c.Connections, number)
}

func (c *Storage) Get(number string) *ConnWrapper {
	wrapper := c.Connections[number]
	if !wrapper.IsFilled || wrapper.Conn == nil || wrapper.Sess == nil {
		return nil
	}

	return &wrapper
}

func (c *Storage) Reset() {
	for k := range c.Connections {
		delete(c.Connections, k)
	}
}
