package dc1server

import (
	"bufio"
	"context"
	"net"
	"sync"
	"sync/atomic"
)

const (
	channelBufSize = 256
	delim          = '\n'
)

type Session struct {
	id     string
	conn   net.Conn
	send   chan []byte
	server *Server
	close  atomic.Bool
	on     sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

func NewSession(conn net.Conn, server *Server) *Session {
	s := &Session{
		id:     conn.RemoteAddr().String(),
		conn:   conn,
		server: server,
		send:   make(chan []byte, channelBufSize),
	}
	s.ctx, s.cancel = context.WithCancel(server.baseCtx)
	return s
}

func (c *Session) Conn() net.Conn {
	return c.conn
}
func (c *Session) SessionID() string {
	return c.id
}

func (c *Session) SendMessage(message []byte) {
	if !c.close.Load() {
		select {
		case c.send <- message:
		}
	}
}

func (c *Session) write() {
	defer c.Close()
	for {
		select {
		case msg := <-c.send:
			if c.conn == nil {
				return
			}
			c.server.logger.Infof("send msg to [%s] : %s", c.id, msg)
			_, err := c.conn.Write(msg)
			if err != nil {
				c.server.logger.Errorf("[tcp] write message error: %v", err)
				return
			}
		case <-c.ctx.Done():
			return
		}

	}
}

func (c *Session) Close() {
	c.on.Do(func() {
		c.close.Store(true)
		c.server.unregister <- c
		c.cancel()
		err := c.conn.Close()
		if err != nil {
			return
		}
	})
}

func (c *Session) Listen() {
	c.server.register <- c
	go c.write()
	go c.read()
}

func (c *Session) read() {
	defer c.Close()

	reader := bufio.NewReader(c.conn)

	for {
		if c.conn == nil {
			return
		}
		readString, err := reader.ReadString(delim)
		if err != nil {
			c.server.logger.Errorf("[tcp] read message error: %v", err)
			return
		}
		err = c.server.MessageHandler(c, readString)
		if err != nil {
			c.server.logger.Warnf("MessageHandler,Error:%s", err.Error())
			return
		}
	}
}
