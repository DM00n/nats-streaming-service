package nats

import (
	"github.com/nats-io/stan.go"
)

func NewStanConnect() (stan.Conn, error) {
	sc, err := stan.Connect("test-cluster", "111")
	if err != nil {
		return nil, err
	}
	return sc, nil
}

func Sub(sc stan.Conn) (stan.Subscription, error) {
	sub, err := sc.Subscribe("foo",
		func(m *stan.Msg) {
			println(string(m.Data))
			m.Ack()
		},
		stan.StartWithLastReceived())
	if err != nil {
		return nil, err
	}
	return sub, nil
}
