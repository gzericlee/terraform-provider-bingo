package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func init() {
	_ = os.Setenv("TF_ACC", "1")
	_ = os.Setenv("TF_LOG", "INFO")
}

var testProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"bingo": func() (tfprotov6.ProviderServer, error) {
		return tfsdk.NewProtocol6Server(New("test")()), nil
	},
}

func testPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
