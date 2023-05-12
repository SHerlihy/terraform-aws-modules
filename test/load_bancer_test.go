package test

import (
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/gruntwork-io/terratest/modules/http-helper"
    "testing"
    "fmt"
    "time"
)

func TestLB(t *testing.T){
    opts := &terraform.Options{
        TerraformDir: "../load_balancer",

        Vars: map[string]interface{}{
            // "lb_security_group_ids": [1]string{"default",},
            "http_open": 8080,
        },
    }

    defer terraform.Destroy(t, opts)

    terraform.Init(t, opts)
    terraform.Apply(t, opts)

    lbDNS := terraform.OutputRequired(t, opts, "load_balancer_dns")
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
