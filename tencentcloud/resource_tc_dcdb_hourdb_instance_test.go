package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDcdbHourDbInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbHourdbInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dcdb_hourdb_instance.hourdb_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dcdb_hourdb_instance.hourdbInstance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbHourdbInstance = `

resource "tencentcloud_dcdb_hourdb_instance" "hourdb_instance" {
  zones = ""
  shard_memory = ""
  shard_storage = ""
  shard_node_count = ""
  shard_count = ""
  vpc_id = ""
  subnet_id = ""
  db_version_id = ""
  security_group_id = ""
  project_id = ""
  instance_name = ""
  resource_tags {
			tag_key = ""
			tag_value = ""

  }
}

`
