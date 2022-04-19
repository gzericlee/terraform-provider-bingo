package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCommandResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testPreCheck(t) },
		ProtoV6ProviderFactories: testProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testCommandResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("bingo_cmp_command.test", "host_type", "1"),
					resource.TestCheckResourceAttr("bingo_cmp_command.test", "content", "pwd"),
					resource.TestCheckResourceAttr("bingo_cmp_command.test", "status", "new"),
					resource.TestCheckResourceAttr("bingo_cmp_command.test", "instance_ids", "c0dea473-cfc0-49a7-830e-a7edc8f1125d"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testCommandResourceConfig() string {
	return fmt.Sprintf(`
provider "bingo" {
  sso_endpoint      = "##"
  cmp_access_token	= "##"
  #cmp_client_secret	= "clientSecret1"
  cmp_endpoint		= "##"
}

resource "bingo_cmp_command" "test" {
  host_type   	= "1" 
  content     	= "pwd"
  instance_ids	= "c0dea473-cfc0-49a7-830e-a7edc8f1125d"
}
`)
}
