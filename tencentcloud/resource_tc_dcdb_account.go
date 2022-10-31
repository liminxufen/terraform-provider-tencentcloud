/*
Provides a resource to create a dcdb account

Example Usage

```hcl
resource "tencentcloud_dcdb_account" "account" {
  instance_id = ""
  user_name = ""
  host = ""
  password = ""
  read_only = ""
  description = ""
  max_user_connections = ""
}

```
Import

dcdb account can be imported using the id, e.g.
```
$ terraform import tencentcloud_dcdb_account.account account_id
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

func resourceTencentCloudDcdbAccount() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudDcdbAccountRead,
		Create: resourceTencentCloudDcdbAccountCreate,

		Update: resourceTencentCloudDcdbAccountUpdate,
		Delete: resourceTencentCloudDcdbAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance id.",
			},

			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "account name.",
			},

			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "db host.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "password.",
			},

			"read_only": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "is it a readonly account.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "description for account.",
			},

			"max_user_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "max user connections.",
			},
		},
	}
}

func resourceTencentCloudDcdbAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCreateAccountRequest()
		response   *dcdb.CreateAccountResponse
		instanceId string
		userName   string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("password"); ok {
		request.Password = helper.String(v.(string))
	}

	if v, ok := d.GetOk("read_only"); ok {
		request.ReadOnly = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("max_user_connections"); ok {
		request.MaxUserConnections = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CreateAccount(request)
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
		log.Printf("[CRITAL]%s create dcdb account failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId

	d.SetId(instanceId + FILED_SP + userName)
	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	accounts, err := service.DescribeDcdbAccount(ctx, instanceId)

	if err != nil {
		return err
	}

	if accounts == nil {
		d.SetId("")
		return fmt.Errorf("resource `account` %s does not exist", userName)
	}

	if len(accounts.Users) == 0 {
		d.SetId("")
		return fmt.Errorf("resource `account` %s does not exist", userName)
	}

	account := &dcdb.DBAccount{}
	for _, ac := range accounts.Users {
		if *ac.UserName == userName {
			account = ac
		}
	}

	if accounts.InstanceId != nil {
		_ = d.Set("instance_id", accounts.InstanceId)
	}

	if account.UserName != nil {
		_ = d.Set("user_name", account.UserName)
	}

	if account.Host != nil {
		_ = d.Set("host", account.Host)
	}

	if account.ReadOnly != nil {
		_ = d.Set("read_only", account.ReadOnly)
	}

	if account.Description != nil {
		_ = d.Set("description", account.Description)
	}

	return nil
}

func resourceTencentCloudDcdbAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)

	request := dcdb.NewModifyAccountDescriptionRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	request.InstanceId = &instanceId
	request.UserName = &userName

	if d.HasChange("instance_id") {

		return fmt.Errorf("`instance_id` do not support change now.")

	}

	if d.HasChange("user_name") {

		return fmt.Errorf("`user_name` do not support change now.")

	}

	if d.HasChange("host") {

		return fmt.Errorf("`host` do not support change now.")

	}

	if d.HasChange("password") {

		return fmt.Errorf("`password` do not support change now.")

	}

	if d.HasChange("read_only") {

		return fmt.Errorf("`read_only` do not support change now.")

	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

	}

	if d.HasChange("max_user_connections") {

		return fmt.Errorf("`max_user_connections` do not support change now.")

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyAccountDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dcdb account failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudDcdbAccountRead(d, meta)
}

func resourceTencentCloudDcdbAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_account.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	instanceId := idSplit[0]
	userName := idSplit[1]

	if err := service.DeleteDcdbAccountById(ctx, instanceId, userName); err != nil {
		return err
	}

	return nil
}
