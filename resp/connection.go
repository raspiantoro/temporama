package resp

import "net"

type Connection struct {
	net.Conn
	proto int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		Conn:  conn,
		proto: 2, // default protocol
	}
}
