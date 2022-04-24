package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestCommandResource(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testCommandResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bingo_cmp_command.dev", "host_type", "1"),
					resource.TestCheckResourceAttr("bingo_cmp_command.dev", "content", "pwd"),
					resource.TestCheckResourceAttr("bingo_cmp_command.dev", "status", "new"),
					resource.TestCheckResourceAttr("bingo_cmp_command.dev", "instance_ids", "c0dea473-cfc0-49a7-830e-a7edc8f1125d"),
				),
			},
		},
	})
}

func testCommandResourceConfig() string {
	return fmt.Sprintf(`
provider "bingo" {

}

resource "bingo_cmp_command" "dev" {
  host_type   	= "1" 
  content     	= "pwd"
  instance_ids	= "c0dea473-cfc0-49a7-830e-a7edc8f1125d"
}
`)
}
