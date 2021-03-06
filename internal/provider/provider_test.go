package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	_ = os.Setenv("debug", "true")
	_ = os.Setenv("TF_LOG", "INFO")

	_ = os.Setenv("IAM_ENDPOINT", "https://sso.bingosoft.net")
	_ = os.Setenv("CMP_ENDPOINT", "https://cmp-dev.bingosoft.net")
	_ = os.Setenv("IAM_CLIENT_ID", "ajcNcUVYSmEW99qCyA9PnT")
	_ = os.Setenv("IAM_CLIENT_SECRET", "b25da097-657d-4ed0-a579-47da34ad87e1")
	_ = os.Setenv("USER_NAME", "bingo")
	_ = os.Setenv("PASSWORD", "pass@cmp#2019")
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"bingo": func() (*schema.Provider, error) {
		return New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
