package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVpcAcl_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc_acl.foo", "ingress"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_acl.foo", "egress"),
				),
			},
		},
	})
}
func TestAccTencentCloudVpcAclRulesUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#80#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#80-90#TCP"),
				),
			},
			{
				Config: testAccVpcACLConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "ingress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.0", "ACCEPT#192.168.1.0/24#800#TCP"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "egress.1", "ACCEPT#192.168.1.0/24#800-900#TCP"),
				),
			},
		},
	})
}

func testAccCheckVpcACLExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		_, _, has, err := service.DescribeNetWorkByACLID(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has > 0 {
			return nil
		}

		return fmt.Errorf("vpc network acl %s not exists", rs.Primary.ID)
	}
}

func testAccCheckVpcACLDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := VpcService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_vpc_acl" {
			continue
		}
		time.Sleep(5 * time.Second)
		_, _, has, err := service.DescribeNetWorkByACLID(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if has == 0 {
			return nil
		}

		return fmt.Errorf("vpc acl %s still exists", rs.Primary.ID)
	}

	return nil
}

const testAccVpcACLConfig = `
data "tencentcloud_vpc_instances" "default" {
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
    	"ACCEPT#192.168.1.0/24#80-90#TCP",
	]
}  
`

const testAccVpcACLConfigUpdate = `
data "tencentcloud_vpc_instances" "default" {
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.default.instance_list.0.vpc_id
    name  	= "test_acl_gogoowang"
	ingress = [
		"ACCEPT#192.168.1.0/24#800#TCP",
		"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#800#TCP",
    	"ACCEPT#192.168.1.0/24#800-900#TCP",
	]
} 
`
