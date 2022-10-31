/*
Provides a resource to create a dcdb hourdb_instance

Example Usage

```hcl
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

```
Import

dcdb hourdb_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_hourdb_instance.hourdb_instance hourdbInstance_id
```
*/

package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbHourDbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbHourDbInstanceRead,
		Create: resourceTencentCloudDcdbHourDbInstanceCreate,

		Update: resourceTencentCloudDcdbHourDbInstanceUpdate,
		Delete: resourceTencentCloudDcdbHourDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "available zone of instance.",
			},

			"shard_memory": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "shard memory.",
			},

			"shard_storage": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "shard storage.",
			},

			"shard_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "node count for each shard.",
			},

			"shard_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "instance shard count.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "vpc id for this instance.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "subnet id for this instance, it&amp;#39;s required when vpcId is set.",
			},

			"db_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "db engine version for this instance, default to Percona 5.7.17.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "security group id.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "project id.",
			},

			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of this instance.",
			},

			"resource_tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "resource tags.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag key.",
						},
						"tag_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "tag value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDcdbHourDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCreateHourDCDBInstanceRequest()
		response   *dcdb.CreateHourDCDBInstanceResponse
		instanceId string
	)

	if v, ok := d.GetOk("zones"); ok {
		zonesSet := v.(*schema.Set).List()
		for i := range zonesSet {
			zones := zonesSet[i].(string)
			request.Zones = append(request.Zones, &zones)
		}
	}

	if v, ok := d.GetOk("shard_memory"); ok {
		request.ShardMemory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_storage"); ok {
		request.ShardStorage = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_node_count"); ok {
		request.ShardNodeCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("shard_count"); ok {
		request.ShardCount = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("db_version_id"); ok {
		request.DbVersionId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("security_group_id"); ok {
		request.SecurityGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		request.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource_tags"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			resourceTag := dcdb.ResourceTag{}
			if v, ok := dMap["tag_key"]; ok {
				resourceTag.TagKey = helper.String(v.(string))
			}
			if v, ok := dMap["tag_value"]; ok {
				resourceTag.TagValue = helper.String(v.(string))
			}

			request.ResourceTags = append(request.ResourceTags, &resourceTag)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateHourDCDBInstance(request)
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
		log.Printf("[CRITAL]%s create dcdb hourdbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]

	d.SetId(instanceId)
	return resourceTencentCloudDcdbHourDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	hourdbInstances, err := service.DescribeDcdbDbInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if hourdbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `hourdbInstance` %s does not exist", instanceId)
	}

	if len(hourdbInstances.Instances) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `hourdbInstance` %s does not exist", instanceId)
	}
	hourdbInstance := hourdbInstances.Instances[0]

	if hourdbInstance.Zone != nil {
		_ = d.Set("zones", []*string{hourdbInstance.Zone})
	}

	if hourdbInstance.Memory != nil {
		_ = d.Set("shard_memory", hourdbInstance.Memory)
	}

	if hourdbInstance.Storage != nil {
		_ = d.Set("shard_storage", hourdbInstance.Storage)
	}

	if hourdbInstance.ShardCount != nil {
		_ = d.Set("shard_count", hourdbInstance.ShardCount)
	}

	if hourdbInstance.VpcId != nil {
		_ = d.Set("vpc_id", hourdbInstance.VpcId)
	}

	if hourdbInstance.SubnetId != nil {
		_ = d.Set("subnet_id", hourdbInstance.SubnetId)
	}

	if hourdbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", hourdbInstance.DbVersionId)
	}

	if hourdbInstance.ProjectId != nil {
		_ = d.Set("project_id", hourdbInstance.ProjectId)
	}

	if hourdbInstance.InstanceName != nil {
		_ = d.Set("instance_name", hourdbInstance.InstanceName)
	}

	if hourdbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range hourdbInstance.ResourceTags {
			resourceTagsMap := map[string]interface{}{}
			if resourceTags.TagKey != nil {
				resourceTagsMap["tag_key"] = resourceTags.TagKey
			}
			if resourceTags.TagValue != nil {
				resourceTagsMap["tag_value"] = resourceTags.TagValue
			}

			resourceTagsList = append(resourceTagsList, resourceTagsMap)
		}
		_ = d.Set("resource_tags", resourceTagsList)
	}

	return nil
}

func resourceTencentCloudDcdbHourDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := dcdb.NewModifyDBInstanceNameRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("zones") {

		return fmt.Errorf("`zones` do not support change now.")

	}

	if d.HasChange("shard_memory") {

		return fmt.Errorf("`shard_memory` do not support change now.")

	}

	if d.HasChange("shard_storage") {

		return fmt.Errorf("`shard_storage` do not support change now.")

	}

	if d.HasChange("shard_node_count") {

		return fmt.Errorf("`shard_node_count` do not support change now.")

	}

	if d.HasChange("shard_count") {

		return fmt.Errorf("`shard_count` do not support change now.")

	}

	if d.HasChange("vpc_id") {

		return fmt.Errorf("`vpc_id` do not support change now.")

	}

	if d.HasChange("subnet_id") {

		return fmt.Errorf("`subnet_id` do not support change now.")

	}

	if d.HasChange("db_version_id") {

		return fmt.Errorf("`db_version_id` do not support change now.")

	}

	if d.HasChange("security_group_id") {

		return fmt.Errorf("`security_group_id` do not support change now.")

	}

	if d.HasChange("project_id") {

		return fmt.Errorf("`project_id` do not support change now.")

	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOk("instance_name"); ok {
			request.InstanceName = helper.String(v.(string))
		}

	}

	if d.HasChange("resource_tags") {

		return fmt.Errorf("`resource_tags` do not support change now.")

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb hourdbInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbHourDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbHourDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_hourdb_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteDcdbHourDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
