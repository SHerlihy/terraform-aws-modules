package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/http-helper"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"testing"
	"time"
)

func TestLB(t *testing.T) {
	var defaultVPC = aws.GetDefaultVpc(
		t,
		"eu-west-2",
	)

	defaultSubnetIds := []string{}

	for i := range defaultVPC.Subnets {
		defaultSubnetIds = append(defaultSubnetIds, defaultVPC.Subnets[i].Id)
	}

	fmt.Println(defaultSubnetIds)

	opts := &terraform.Options{
		TerraformDir: "../../../wrapped_modules/load_balancer/single_ingress_all_egress",

		Vars: map[string]interface{}{
			"subnet_ids": defaultSubnetIds,
			"open_port":  80,
		},
	}

	defer terraform.Destroy(t, opts)

	terraform.Init(t, opts)
	terraform.Apply(t, opts)

	lbDNS := terraform.OutputRequired(t, opts, "lb_dns_name")
	url := fmt.Sprintf("http://%s", lbDNS)

	expectedStatus := 404
	expectedBody := "404: page not found"
	maxRetries := 10
	timeBetweenRetries := 10 * time.Second

	http_helper.HttpGetWithRetry(
		t,
		url,
		nil,
		expectedStatus,
		expectedBody,
		maxRetries,
		timeBetweenRetries,
	)
}
