package test

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
    "os"
    "os/exec"
    "bytes"
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

    execOverSSH(pub1aServerPubDNS, pvt1aServerPvtDNS, pingServerPubDNS)
}

func execOverSSH(gatewayIp string, privateIp string, externalIp string) {
    sshString := fmt.Sprintf("-t -i '../.ssh/id_rsa' ubuntu@%s ssh -t -i '/home/ubuntu/.ssh/id_rsa' ubuntu@%s ping %s", gatewayIp, privateIp, externalIp)

	cmd := exec.Command("ssh", sshString)
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
	    panic(err)
	}
}
