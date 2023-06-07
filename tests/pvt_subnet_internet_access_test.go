package test

import (
	"github.com/go-ping/ping"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
    "time"
    "fmt"
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

func TestPvtSubnetInternetAccess(t *testing.T) {
	opts := &terraform.Options{
		TerraformDir: "../examples/networking/pvt_subnet_internet_access",
	}

	defer terraform.Destroy(t, opts)

	terraform.Init(t, opts)

	terraform.Apply(t, opts)

    pingServerIp := terraform.OutputRequired(t, opts, "ping_server-public_ip")
	pub1aServerIp := terraform.OutputRequired(t, opts, "pub_server-public_ip")
	pvt1aServerIp := terraform.OutputRequired(t, opts, "pvt_server-net_access-public_ip")
    pvt1bServerIp := terraform.OutputRequired(t, opts, "pvt_server-no_access-public_ip")

    expectPong(t, pingServerIp)
    expectPong(t, pub1aServerIp)
    unexpectPong(t, pvt1aServerIp)
    unexpectPong(t, pvt1bServerIp)
}

//
//func execOverSSH() {
//	cmd := exec.Command("ssh", pub1aServerIp, "bash-command")
//	var out bytes.Buffer
//	cmd.Stdout = &out
//	err := cmd.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
//}
