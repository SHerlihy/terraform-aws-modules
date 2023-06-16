package test_helpers

import (
	"bufio"
	"fmt"
	"github.com/go-ping/ping"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"time"
)

func PingServer(subjectIp string) *ping.Statistics {
	pinger, err := ping.NewPinger(subjectIp)
	if err != nil {
		panic(err)
	}
	pinger.Timeout = 5 * time.Second
	pinger.Count = 3

	pinger.Run() // blocks until finished

	return pinger.Statistics()
}

func WriteableSSHConnSpecPort(port int) func() int {
	return func() int {
		return 7 + port
	}
}

//func WriteableSSHConnSpecPort(port int) func(string, string, chan []byte) {
//	 return func(serverDNS string, pvtKeyPath string, wr chan []byte) {
//		return writeableSSHConn(port, serverDNS, pvtKeyPath, wr)
//	}
//}

func writeableSSHConn(port int, serverDNS string, pvtKeyPath string, wr chan []byte) {
	var err error
	var signer ssh.Signer

	host := fmt.Sprintf("%s:%v", serverDNS, port)

	pKey, err := os.ReadFile(pvtKeyPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	user := "ubuntu"
	pwd := ""

	signer, err = ssh.ParsePrivateKey(pKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	hostkeyCallback := ssh.InsecureIgnoreHostKey()

	conf := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: hostkeyCallback,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
			ssh.PublicKeys(signer),
		},
	}

	var conn *ssh.Client

	var stdin io.WriteCloser
	var stdout, stderr io.Reader

	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer conn.Close()

	var session *ssh.Session
	session, err = conn.NewSession()
	if err != nil {
		fmt.Println(err.Error())
	}
	defer session.Close()

	stdin, err = session.StdinPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stdout, err = session.StdoutPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	stderr, err = session.StderrPipe()
	if err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		for {
			select {
			case d := <-wr:
				_, err := stdin.Write(d)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for {
			if tkn := scanner.Scan(); tkn {
				rcv := scanner.Bytes()
				raw := make([]byte, len(rcv))
				copy(raw, rcv)
				fmt.Println(string(raw))
			} else {
				if scanner.Err() != nil {
					fmt.Println(scanner.Err())
				} else {
					fmt.Println("io.EOF")
				}
				return
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	session.Shell()
}
