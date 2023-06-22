package test

import (
	"fmt"
	helpers "github.com/SHerlihy/terraform-aws-modules/test_helpers"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
)

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

	t.Run("ping servers", func(t *testing.T) {
		expectPong(t, pingServerPubDNS)
		expectPong(t, pub1aServerPubDNS)

		unexpectPong(t, pvt1aServerPubDNS)
		unexpectPong(t, pvt1bServerPubDNS)
	})

	t.Run("ssh into pvt", func(t *testing.T) {
		writeableSSHConnOn22 := helpers.WriteableSSHConnSpecPort(22)

		wr := make(chan []byte)
		done := make(chan struct{})

		writeableSSHConnOn22(pub1aServerPubDNS, "../examples/networking/.ssh/id_rsa", wr, done)

		fmt.Println("writting to wr")
		sshIntoPvt := "ssh -tt -i \"./.ssh/id_rsa\" " + "ubuntu@" + pvt1aServerPvtDNS + "\n"
		pingInitial := "ping -c 3 " + pub1aServerPubDNS + "\n"

		wr <- []byte("pwd\n")
		wr <- []byte(sshIntoPvt)
		wr <- []byte("hostname -I\n")
		wr <- []byte(pingInitial)

		time.Sleep(5 * time.Second)
		close(done)

		expectPong(t, pingServerPubDNS)
		//time.Sleep(5 * time.Second)
		//wr := make(chan []byte)
		//writeableSSHConnOn22(pub1aServerPubDNS, "../examples/networking/.ssh/id_rsa", wr)

		//sshIntoPvt := "ssh -tt -i \"./.ssh/id_rsa\" " + "ubuntu@" + pvt1aServerPvtDNS + "\n"
		//pingInitial := "ping -c 3 " + pub1aServerPubDNS + "\n"

		//wr <- []byte(sshIntoPvt)
		//wr <- []byte("hostname -I\n")
		//wr <- []byte(pingInitial)

		//scanner := bufio.NewScanner(os.Stdin)
		//scanner.Scan()
		//text := scanner.Text()
		//fmt.Println(text)

		//time.Sleep(5 * time.Second)
	})
}

func expectPong(t *testing.T, subjectIp string) bool {
	stats := helpers.PingServer(subjectIp)

	if stats.PacketsRecv > 0 {
		return true
	}

	errString := fmt.Sprintf("Error: %s no pong returned", subjectIp)
	t.Errorf(errString)

	return false
}

func unexpectPong(t *testing.T, subjectIp string) bool {
	stats := helpers.PingServer(subjectIp)

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
