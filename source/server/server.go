package main

import (
	"crypto/tls"
	"log"
	"net"
)

type Server struct {
	si       ServerInterface
	Listener net.Listener
}

func Listen(network, address string) (*Server, error) {
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen(network, address, config)
	if err != nil {
		return nil, err
	}
	return &Server{Listener: listener}, nil
}

func (t *Server) Run(si ServerInterface) error {
	t.si = si
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			return err
		}
		go t.accept(conn)

	}
}
func (t *Server) accept(conn net.Conn) {
	defer conn.Close()

	if !t.si.OnAccept(conn) {
		return
	}

	defer t.si.OnClientClose(conn)

	for {
		cmd, ext, data, err := Recv(conn)
		if err != nil {
			t.si.OnRecvError(conn, err)
			break
		}
		if err := t.si.OnData(conn, cmd, ext, data); err != nil {
			break
		}
	}
}

func (t *Server) Send(conn net.Conn, cmd, ext uint16, data []byte) (int, error) {
	return Send(conn, cmd, ext, data)
}
