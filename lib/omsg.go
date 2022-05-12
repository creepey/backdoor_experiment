package main

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type ServerInterface interface {
	OnAccept(conn net.Conn) bool
	OnData(conn net.Conn, cmd, ext uint16, data []byte) error
	OnRecvError(conn net.Conn, err error)
	OnClientClose(conn net.Conn)
}

type ClientInterface interface {
	OnSuccess(conn net.Conn)
	OnData(cmd, ext uint16, data []byte) error
	OnRecvError(err error)
	OnClose()
}

type head struct {
	Sign uint16
	CRC  uint16
	Cmd  uint16
	Ext  uint16
	Size uint32
}

const signWord = 0x1234

var headSize = binary.Size(head{})

func Send(conn net.Conn, cmd, ext uint16, data []byte) (int, error) {
	buffer := make([]byte, headSize+len(data))

	binary.LittleEndian.PutUint16(buffer, signWord)
	binary.LittleEndian.PutUint16(buffer[2:], crc(data))
	binary.LittleEndian.PutUint16(buffer[4:], cmd)
	binary.LittleEndian.PutUint16(buffer[6:], ext)
	binary.LittleEndian.PutUint32(buffer[8:], uint32(len(data)))

	copy(buffer[headSize:], data)
	n, err := conn.Write(buffer)
	if err != nil {
		return n, err
	}
	return n, nil
}
func Recv(conn net.Conn) (cmd, ext uint16, data []byte, err error) {

	header := make([]byte, headSize)
	if _, err := io.ReadFull(conn, header); err != nil {
		return 0, 0, nil, err
	}
	//fmt.Println(hex.Dump(header))
	if signWord != binary.LittleEndian.Uint16(header) {
		return 0, 0, nil, errors.New("signWord err")
	}

	icrc := binary.LittleEndian.Uint16(header[2:])
	cmd = binary.LittleEndian.Uint16(header[4:])
	ext = binary.LittleEndian.Uint16(header[6:])
	dataLen := binary.LittleEndian.Uint32(header[8:])

	data = make([]byte, int(dataLen))

	if _, err := io.ReadFull(conn, data); err != nil {
		return 0, 0, nil, err
	}
	now_crc := crc(data)
	if now_crc != icrc {
		return 0, 0, nil, errors.New("crc err")
	}
	return
}

func crc(data []byte) uint16 {
	crc := uint16(0xFFFF)
	length := len(data)
	for i := 0; i < length; i++ {
		crc = (crc >> 8) ^ uint16(data[i])
		for j := 0; j < 8; j++ {
			flag := crc & 0x0001
			crc >>= 1
			if flag == 1 {
				crc ^= 0xA001
			}
		}
	}
	return crc
}
