package tencentcloud

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudVpcAclBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc_acl.foo", "name", "test_acl_gogoowang"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc_acl.foo", "ingress"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc_acl.foo", "egress"),
				),
			},
		},
	})
}

/*
func TestAccTencentCloudVpcACLUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcACLConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", "test_acl_gogoowang"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "ingress.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "egress.#"),
				),
			},
			{
				Config: testAccVpcConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcACLExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),

					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", fmt.Sprintf("%s.%d", "dns_servers", hashcode.String("119.29.29.29")), "119.29.29.29"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", fmt.Sprintf("%s.%d", "dns_servers", hashcode.String("182.254.116.116")), "182.254.116.116"),
				),
			},
		},
	})
}
*/
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
    vpc_id            	= data.tencentcloud_vpc_instances.foo.0.vpc_id
    network_acl_name  	= "test_acl_gogoowang"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
		"ACCEPT#192.168.1.0/24#440,900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
    	"ACCEPT#192.168.1.0/24#80-90#TCP",
    	"ACCEPT#192.168.1.0/24#80,90#TCP",
	]
}  
`

const testAccVpcACLConfigUpdate = testAccVpcACLConfig + `
data "tencentcloud_vpc_instances" "default" {
}

resource "tencentcloud_vpc_acl" "foo" {  
    vpc_id            	= data.tencentcloud_vpc_instances.foo.0.vpc_id
    network_acl_name  	= "test_acl_gogoowang"
	ingress = [
		"ACCEPT#192.168.1.0/24#80#TCP",
		"ACCEPT#192.168.1.0/24#80-90#TCP",
		"ACCEPT#192.168.1.0/24#440,900#TCP",
	]
	egress = [
    	"ACCEPT#192.168.1.0/24#80#TCP",
    	"ACCEPT#192.168.1.0/24#80-90#TCP",
    	"ACCEPT#192.168.1.0/24#80,90#TCP",
	]
} 
`
