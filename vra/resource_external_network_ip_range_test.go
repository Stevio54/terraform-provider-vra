package vra

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/vmware/vra-sdk-go/pkg/client/network_ip_range"
)

func TestAccVRAExternalNetworkIPRangeBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVRAExternalNetworkIPRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckVRAExternalNetworkIPRangeConfig(rInt),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckVRAExternalNetworkIPRangeExists("vra_network_ip_range.this"),
					resource.TestMatchResourceAttr(
						"vra_network_ip_range.this", "name", regexp.MustCompile("^my-vra-network-ip-range-"+strconv.Itoa(rInt))),
				),
			},
		},
	})
}

func testAccCheckVRAExternalNetworkIPRangeExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no network ip range ID is set")
		}
		return nil
	}
}

func testAccCheckVRAExternalNetworkIPRangeDestroy(s *terraform.State) error {
	apiClient := testAccProviderVRA.Meta().(*Client).apiClient

	for _, rs := range s.RootModule().Resources {
		_, err := apiClient.NetworkIPRange.GetExternalNetworkIPRange(network_ip_range.NewGetExternalNetworkIPRangeParams().WithID(rs.Primary.ID))
		if err == nil {
			return fmt.Errorf("Resource 'vra_external_network_ip_range' still exists with id %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckVRAExternalNetworkIPRangeConfig(rInt int) string {

	fabricNetworkName := os.Getenv("VRA_FABRICNETWORK_NAME")
	externalNetworkId := os.Getenv("VRA_EXTERNALNETWORK_ID")

	return fmt.Sprintf(`
	data "vra_fabric_network" "this" {
		filter = "name eq '%s'"
	}
	  
	resource "vra_external_network_ip_range" "this" {
		external_id		  = "%s"
		fabric_network_id = data.vra_fabric_network.this.id
	}	  
	 
`, fabricNetworkName, externalNetworkId)
}
