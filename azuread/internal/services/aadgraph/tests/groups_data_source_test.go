package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
)

func TestAccAzureADGroupsDataSource_byUserPrincipalNames(t *testing.T) {
	dsn := "data.azuread_groups.tests"
	id := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupsDataSource_byDisplayNames(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADGroupsDataSource_byObjectIds(t *testing.T) {
	dsn := "data.azuread_groups.tests"
	id := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupsDataSource_byObjectIds(id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADGroupsDataSource_noNames(t *testing.T) {
	dsn := "data.azuread_groups.tests"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupsDataSource_noNames(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "names.#", "0"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "0"),
				),
			},
		},
	})
}

func testAccAzureADGroup_multiple(id int) string {
	return fmt.Sprintf(`
resource "azuread_group" "testA" {
  name    = "acctestGroup-%[1]d"
  members = []
}

resource "azuread_group" "testB" {
  name    = "acctestGroup-%[1]d"
  members = []
}
`, id)
}

func testAccAzureADGroupsDataSource_byDisplayNames(id int) string {
	return fmt.Sprintf(`
%s

data "azuread_groups" "tests" {
  names = [azuread_group.testA.name, azuread_group.testB.name]
}
`, testAccAzureADGroup_multiple(id))
}

func testAccAzureADGroupsDataSource_byObjectIds(id int) string {
	return fmt.Sprintf(`
%s

data "azuread_groups" "tests" {
  object_ids = [azuread_group.testA.object_id, azuread_group.testB.object_id]
}
`, testAccAzureADGroup_multiple(id))
}
func testAccAzureADGroupsDataSource_noNames() string {
	return `
data "azuread_groups" "tests" {
  names = []
}
`
}
