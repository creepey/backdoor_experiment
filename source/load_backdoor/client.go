package main

import (
	"crypto/tls"
	"net"
)

type Client struct {
	ci   ClientInterface
	Conn net.Conn
	T    chan bool
	off  int
}

func Dial(network, adress string, ci ClientInterface) (*Client, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial(network, adress, conf)
	if err != nil {
		return nil, err
	}
	cl := &Client{Conn: conn, ci: ci, T: make(chan bool)}
	cl.ci.OnSuccess(conn)
	go cl.clinet()
	return cl, nil
}

func (t *Client) clinet() {
	defer t.ci.OnClose()

	for {
		cmd, ext, data, err := Recv(t.Conn)
		if err != nil {
			t.T <- true
			t.ci.OnRecvError(err)
			break
		}
		if err := t.ci.OnData(cmd, ext, data); err != nil {
			break
		}
	}

}

func (t *Client) Close() {
	t.Conn.Close()
}
func (t *Client) Send(cmd, ext uint16, data []byte) (int, error) {
	return Send(t.Conn, cmd, ext, data)
}
