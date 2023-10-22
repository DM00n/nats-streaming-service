package stan

import (
	"github.com/nats-io/stan.go"
)

type Conn struct {
	conn stan.Conn
}

func NewStanConnect() (*Conn, error) {
	c, err := stan.Connect("test-cluster", "111")
	if err != nil {
		return nil, err
	}
	return &Conn{c}, nil
}

func (c Conn) Sub(f func(data []byte)) (stan.Subscription, error) {
	sub, err := c.conn.Subscribe("foo",
		func(m *stan.Msg) {
			f(m.Data)
		})
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (c Conn) Close() error {
	err := c.conn.Close()
	return err
}
