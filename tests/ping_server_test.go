package test

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/go-ping/ping"
	"testing"
)

func TestPingServer(t *testing.T) {
	opts := &terraform.Options{
		TerraformDir: "../ping_server",
	}

	defer terraform.Destroy(t, opts)

	terraform.Init(t, opts)

	terraform.Apply(t, opts)

	pingServerDNS := terraform.OutputRequired(t, opts, "public_dns")

    pinger, err := ping.NewPinger(pingServerDNS)
    if err != nil {
        panic(err)
    }
    pinger.Count = 3

    pinger.Run() // blocks until finished

    stats := pinger.Statistics()

    if stats.PacketsRecv == 0 {
        t.Errorf("No packets recieved")
    }
}
