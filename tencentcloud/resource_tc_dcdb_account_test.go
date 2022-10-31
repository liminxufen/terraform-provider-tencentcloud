package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbAccount_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbAccount,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_account.account", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_account.account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbAccount = `

resource "tencentcloud_dcdb_account" "account" {
  instance_id = ""
  user_name = ""
  host = ""
  passwod = ""
  read_only = ""
  description = ""
  max_user_connections = ""
}

`
