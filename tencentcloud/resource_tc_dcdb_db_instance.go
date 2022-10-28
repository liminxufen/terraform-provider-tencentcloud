/*
Provides a resource to create a dcdb db_instance

Example Usage

```hcl
resource "tencentcloud_dcdb_db_instance" "db_instance" {
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
  period = ""
  instance_name = ""
  resource_tags {
	tag_key = ""
	tag_value = ""

  }
}

```
Import

dcdb db_instance can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_db_instance.db_instance dbInstance_id
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

func resourceTencentCloudDcdbDbInstance() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbDbInstanceRead,
		Create: resourceTencentCloudDcdbDbInstanceCreate,

		Update: resourceTencentCloudDcdbDbInstanceUpdate,
		Delete: resourceTencentCloudDcdbDbInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zones": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
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

			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "subscribe months.",
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

func resourceTencentCloudDcdbDbInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCreateDCDBInstanceRequest()
		response   *dcdb.CreateDCDBInstanceResponse
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

	if v, ok := d.GetOk("period"); ok {
		request.Period = helper.IntInt64(v.(int))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateDCDBInstance(request)
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
		log.Printf("[CRITAL]%s create dcdb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceIds[0]

	d.SetId(instanceId)
	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	dbInstances, err := service.DescribeDcdbDbInstance(ctx, instanceId)

	if err != nil {
		return err
	}

	if dbInstances == nil {
		d.SetId("")
		return fmt.Errorf("resource `dbInstance` %s does not exist", instanceId)
	}

	if len(dbInstances.Instances) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `dbInstance` %s does not exist", instanceId)
	}

	dbInstance := dbInstances.Instances[0]

	if dbInstance.Zone != nil {
		_ = d.Set("zones", []string{*dbInstance.Zone})
	}

	if dbInstance.Memory != nil {
		_ = d.Set("shard_memory", dbInstance.Memory)
	}

	if dbInstance.Storage != nil {
		_ = d.Set("shard_storage", dbInstance.Storage)
	}

	if dbInstance.ShardCount != nil {
		_ = d.Set("shard_count", dbInstance.ShardCount)
	}

	if dbInstance.VpcId != nil {
		_ = d.Set("vpc_id", dbInstance.VpcId)
	}

	if dbInstance.SubnetId != nil {
		_ = d.Set("subnet_id", dbInstance.SubnetId)
	}

	if dbInstance.DbVersionId != nil {
		_ = d.Set("db_version_id", dbInstance.DbVersionId)
	}

	if dbInstance.ProjectId != nil {
		_ = d.Set("project_id", dbInstance.ProjectId)
	}

	if dbInstance.InstanceName != nil {
		_ = d.Set("instance_name", dbInstance.InstanceName)
	}

	if dbInstance.ResourceTags != nil {
		resourceTagsList := []interface{}{}
		for _, resourceTags := range dbInstance.ResourceTags {
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

func resourceTencentCloudDcdbDbInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.update")()
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

	if d.HasChange("period") {

		return fmt.Errorf("`period` do not support change now.")

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
		log.Printf("[CRITAL]%s create dcdb dbInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbDbInstanceRead(d, meta)
}

func resourceTencentCloudDcdbDbInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	if err := service.DeleteDcdbDbInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
