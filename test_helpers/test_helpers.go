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

func WriteableSSHConnSpecPort(port int) func(serverDNS string, pvtKeyPath string, wr chan []byte, done chan struct{}) {
	 return func(serverDNS string, pvtKeyPath string, wr chan []byte, done chan struct{}) {
		writeableSSHConn(port, serverDNS, pvtKeyPath, wr, done)
	}
}

func writeableSSHConn(port int, serverDNS string, pvtKeyPath string, wr chan []byte, done chan struct{}) {
    conn := makeSSHConnection(port , serverDNS , pvtKeyPath )

    session, err := conn.NewSession()

    go func() {
        select {
        case <-done:
            conn.Close()
            session.Close()
        }
    }()

	if err != nil {
		fmt.Println(err.Error())
	}

	interactWithSession(session, wr, done)

    fmt.Println("starting shell")
	session.Shell()
}

func makeSSHConnection(port int, serverDNS string, pvtKeyPath string) *ssh.Client {
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

	conn, err = ssh.Dial("tcp", host, conf)
	if err != nil {
		fmt.Println(err.Error())
	}

    return conn
}

func interactWithSession(session *ssh.Session, wr <-chan []byte, done <-chan struct{}) {
    var err error
	var stdin io.WriteCloser
	var stdout, stderr io.Reader

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
                case <-done:
                return
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
            select {
            case <-done:
                return
            }
		}
	}()
}
