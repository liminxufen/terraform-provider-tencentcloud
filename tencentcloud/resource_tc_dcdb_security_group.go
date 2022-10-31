/*
Provides a resource to create a dcdb security_group

Example Usage

```hcl
resource "tencentcloud_dcdb_security_group" "security_group" {
  product = "dcdb"
  security_group_id = ""
  instance_ids = ""
}

```
Import

dcdb security_group can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_security_group.security_group securityGroup_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbSecurityGroupRead,
		Create: resourceTencentCloudDcdbSecurityGroupCreate,
		Delete: resourceTencentCloudDcdbSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "product, the value is fiexed to dcdb.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "security group id.",
			},

			"instance_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "attached instance ids.",
			},
		},
	}
}

func resourceTencentCloudDcdbSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = dcdb.NewAssociateSecurityGroupsRequest()
		response        *dcdb.AssociateSecurityGroupsResponse
		instanceId      string
		securityGroupId string
	)

	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().AssociateSecurityGroups(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb securityGroup failed, reason:%+v", logId, err)
		return err
	}

	securityGroupId = *request.SecurityGroupId

	d.SetId(instanceId + FILED_SP + securityGroupId)
	return resourceTencentCloudDcdbSecurityGroupRead(d, meta)
}

func resourceTencentCloudDcdbSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	securityGroups, err := service.DescribeDcdbSecurityGroup(ctx, instanceId)

	if err != nil {
		return err
	}

	if securityGroups == nil {
		d.SetId("")
		return fmt.Errorf("resource `securityGroup` %s does not exist", securityGroupId)
	}

	if len(securityGroups.Groups) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `securityGroup` %s does not exist", securityGroupId)
	}
	securityGroup := &dcdb.SecurityGroup{}
	for _, sg := range securityGroups.Groups {
		if *sg.SecurityGroupId == securityGroupId {
			securityGroup = sg
		}
	}

	if securityGroup.SecurityGroupId != nil {
		_ = d.Set("security_group_id", securityGroup.SecurityGroupId)
	}

	_ = d.Set("instance_ids", []string{instanceId})

	return nil
}

func resourceTencentCloudDcdbSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_security_group.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	securityGroupId := idSplit[1]

	if err := service.DeleteDcdbSecurityGroupById(ctx, instanceId, securityGroupId); err != nil {
		return err
	}

	return nil
}
