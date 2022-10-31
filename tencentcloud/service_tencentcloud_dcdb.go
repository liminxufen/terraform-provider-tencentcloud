package tencentcloud

import (
	"context"
	"log"

	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DcdbService struct {
	client *connectivity.TencentCloudClient
}

func (me *DcdbService) DescribeDcdbDbInstance(ctx context.Context, instanceId string) (dbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&instanceId}

	response, err := me.client.UseDcdbClient().DescribeDCDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	dbInstance = response.Response
	return
}

func (me *DcdbService) DeleteDcdbDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyDCDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DestroyDCDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcdbService) DescribeDcdbHourDbInstance(ctx context.Context, instanceId string) (hourdbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&instanceId}

	response, err := me.client.UseDcdbClient().DescribeDCDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	hourdbInstance = response.Response
	return
}

func (me *DcdbService) DeleteDcdbHourDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyHourDCDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DestroyHourDCDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcdbService) DescribeDcdbAccount(ctx context.Context, instanceId string) (account *dcdb.DescribeAccountsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseDcdbClient().DescribeAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	account = response.Response
	return
}

func (me *DcdbService) DeleteDcdbAccountById(ctx context.Context, instanceId, userName string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDeleteAccountRequest()

	request.InstanceId = &instanceId
	request.UserName = &userName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DeleteAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DcdbService) DescribeDcdbSecurityGroup(ctx context.Context, instanceId string) (securityGroup *dcdb.DescribeDBSecurityGroupsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDBSecurityGroupsRequest()
		product = "dcdb"
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.Product = &product

	response, err := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityGroup = response.Response
	return
}

func (me *DcdbService) DeleteDcdbSecurityGroupById(ctx context.Context, instanceId, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)
	var (
		product = "dcdb"
	)

	request := dcdb.NewDisassociateSecurityGroupsRequest()

	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId
	request.Product = &product

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
