package wa

import whatsapp "github.com/Rhymen/go-whatsapp"

type WAConnectionWrapper struct {
	IsFilled bool
	Conn     *whatsapp.Conn
	Sess     *whatsapp.Session
}

type ActiveConnections struct {
	Connections map[string]WAConnectionWrapper
}

func (c *ActiveConnections) Add(number string, conn *whatsapp.Conn, sess *whatsapp.Session) {
	c.Connections[number] = WAConnectionWrapper{
		IsFilled: true,
		Conn:     conn,
		Sess:     sess,
	}
}

func (c *ActiveConnections) IsExist(number string) bool {
	return c.Connections[number].IsFilled
}

func (c *ActiveConnections) Remove(number string) {
	delete(c.Connections, number)
}

func (c *ActiveConnections) Get(number string) *WAConnectionWrapper {
	wrapper := c.Connections[number]
	if !wrapper.IsFilled || wrapper.Conn == nil || wrapper.Sess == nil {
		return nil
	}

	return &wrapper
}

func (c *ActiveConnections) Reset() {
	for k := range c.Connections {
		delete(c.Connections, k)
	}
}
