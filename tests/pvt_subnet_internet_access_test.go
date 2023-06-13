package test

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"testing"
	"time"
	//"os/exec"
	//"bytes"
	"golang.org/x/crypto/ssh"
	//"golang.org/x/crypto/ssh/knownhosts"
	"bufio"
	"io"
)

func pingServer(subjectIp string) *ping.Statistics {
	pinger, err := ping.NewPinger(subjectIp)
	if err != nil {
		panic(err)
	}
	pinger.Timeout = 5 * time.Second
	pinger.Count = 3

	pinger.Run() // blocks until finished

	return pinger.Statistics()
}

func expectPong(t *testing.T, subjectIp string) bool {
	stats := pingServer(subjectIp)

	if stats.PacketsRecv > 0 {
		return true
	}

	errString := fmt.Sprintf("Error: %s no pong returned", subjectIp)
	t.Errorf(errString)

	return false
}

func unexpectPong(t *testing.T, subjectIp string) bool {
	stats := pingServer(subjectIp)

	if stats.PacketsRecv == 0 {
		return true
	}

	errString := fmt.Sprintf("Error: %s pong has returned", subjectIp)
	t.Errorf(errString)

	return false
}

func getAccessableDNS(t *testing.T, opts *terraform.Options) func(string) string {
	return func(outputVar string) string {
		return terraform.OutputRequired(t, opts, outputVar)
	}
}

func TestPvtSubnetInternetAccess(t *testing.T) {
	opts := &terraform.Options{
		TerraformDir: "../examples/networking/pvt_subnet_internet_access",
	}

	defer terraform.Destroy(t, opts)

	terraform.Init(t, opts)

	terraform.Apply(t, opts)

	getDNS := getAccessableDNS(t, opts)

	pingServerPubDNS := getDNS("ping_server-public_DNS")
	pub1aServerPubDNS := getDNS("pub_server-public_DNS")

	pvt1aServerPubDNS := getDNS("pvt_server-net_access-public_DNS")
	pvt1bServerPubDNS := getDNS("pvt_server-no_access-public_DNS")

	pvt1aServerPvtDNS := getDNS("pvt_server-net_access-private_DNS")
	//pvt1bServerPvtDNS := getDNS("pvt_server-no_access-private_DNS")

	expectPong(t, pingServerPubDNS)
	expectPong(t, pub1aServerPubDNS)

	unexpectPong(t, pvt1aServerPubDNS)
	unexpectPong(t, pvt1bServerPubDNS)

	//	cmd := exec.Command("whoami")
	//    cmd.Stdin = os.Stdin
	//    cmd.Stderr = os.Stderr
	//
	//	var out bytes.Buffer
	//	cmd.Stdout = &out
	//
	//	err := cmd.Run()
	//	if err != nil {
	//	    panic(err)
	//	}
	//
	//    t.Log(out.String())

	// execOverSSH(pub1aServerPubDNS, pvt1aServerPvtDNS, pingServerPubDNS)
	execOverSSH(pub1aServerPubDNS, pvt1aServerPvtDNS)
}

func execOverSSH(serverDNS string, pvtDNS string) {
	var err error
	var signer ssh.Signer

	host := fmt.Sprintf("%s:22", serverDNS)

	fmt.Println(host)

	pKey, err := os.ReadFile("../examples/networking/.ssh/id_rsa") // just pass the file name
	if err != nil {
		fmt.Println(err.Error())
	}
	user := "ubuntu"
	pwd := ""

	signer, err = ssh.ParsePrivateKey(pKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	//	var hostkeyCallback ssh.HostKeyCallback
	//	hostkeyCallback, err = knownhosts.New("~/.ssh/known_hosts")
	//	if err != nil {
	//		fmt.Println(err.Error())
	//	}

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

	wr := make(chan []byte, 10)

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

	//	fmt.Println("pwd")

	//	scanner := bufio.NewScanner(os.Stdin)
	//	scanner.Scan()
	//	text := scanner.Text()

	wr <- []byte("pwd\n")
	wr <- []byte("ls -l ./.ssh\n")
	wr <- []byte("whoami\n")

    time.Sleep(5 * time.Second)

    // host key varification failed
    sshIntoPvt := "ssh -tt -i \"./.ssh/id_rsa\" " + "ubuntu@"+ pvtDNS + "\n"
    
	wr <- []byte(sshIntoPvt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	fmt.Println(text)

    time.Sleep(5 * time.Second)
}

// func execOverSSH(gatewayIp string, privateIp string, externalIp string) {
// func execOverSSH(gatewayIp string) {
//     // sshString := fmt.Sprintf("-t -i \"../examples/networking/.ssh/id_rsa\" ubuntu@%s ssh -t -i \"../examples/networking/.ssh/id_rsa\" ubuntu@%s ping %s", gatewayIp, privateIp, externalIp)
//
//     //pvtKey := fmt.Sprintf("-i \"../examples/networking/.ssh/id_rsa\"")
//     //sshDNS := fmt.Sprintf("ubuntu@%s", gatewayIp)
//     // sshString := fmt.Sprintf("-i \"../examples/networking/.ssh/id_rsa\" ubuntu@%s", gatewayIp)
//     pvtKeyLocation := "-i ../examples/networking/.ssh/id_rsa"
//
//     cmd := exec.Command("ssh", pvtKeyLocation)
// 	//cmd := exec.Command("ssh", pvtKey, sshDNS)
//     cmd.Stdin = os.Stdin
//     cmd.Stderr = os.Stderr
//
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
//
// 	err := cmd.Run()
// 	if err != nil {
// 	    panic(err)
// 	}
// }
