package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbDbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_db_instance.db_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_db_instance.dbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbDbInstance = `

resource "tencentcloud_dcdb_db_instance" "db_instance" {
  zones = ""
  shard_memory = ""
  shard_storage = ""
  shard_node_count = ""
  shard_count = ""
  period = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  security_group_id = ""
  project_id = ""
  instance_name = ""
  resource_tgas {
			tag_key = ""
			tag_value = ""

  }
}

`
