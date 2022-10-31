package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbSecurityGroup_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_security_group.security_group", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_security_group.securityGroup",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbSecurityGroup = `

resource "tencentcloud_dcdb_security_group" "security_group" {
  product = "dcdb"
  security_group_id = ""
  instance_ids = ""
}

`
