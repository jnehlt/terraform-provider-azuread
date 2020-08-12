package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/clients"
)

func TestAccAzureADGroup_basic(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroup_basic(id),
				Check:  testCheckAzureAdGroupBasic(id, "0", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_complete(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroup_complete(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "1", "1"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_owners(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupWithThreeOwners(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "0", "3"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_members(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupWithThreeMembers(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "3", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_membersAndOwners(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupWithOwnersAndMembers(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "2", "1"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_membersDiverse(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupWithDiverseMembers(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "3", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_ownersDiverse(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADGroupWithDiverseOwners(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "0", "2"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADGroup_membersUpdate(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			// Empty group with 0 members
			{
				Config: testAccAzureADGroup_basic(id),
				Check:  testCheckAzureAdGroupBasic(id, "0", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Group with 1 member
			{
				Config: testAccAzureADGroupWithOneMember(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "1", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Group with multiple members
			{
				Config: testAccAzureADGroupWithThreeMembers(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "3", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Group with a different member
			{
				Config: testAccAzureADGroupWithServicePrincipalMember(id),
				Check:  testCheckAzureAdGroupBasic(id, "1", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Empty group with 0 members
			{
				Config: testAccAzureADGroup_basic(id),
				Check:  testCheckAzureAdGroupBasic(id, "0", "0"),
			},
		},
	})
}

func TestAccAzureADGroup_ownersUpdate(t *testing.T) {
	rn := "azuread_group.tests"
	id := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureADGroupDestroy,
		Steps: []resource.TestStep{
			// Empty group with 0 owners
			{
				Config: testAccAzureADGroup_basic(id),
				Check:  testCheckAzureAdGroupBasic(id, "0", "0"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Group with multiple owners
			{
				Config: testAccAzureADGroupWithThreeOwners(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "0", "3"),
			},
			// Group with 1 owners
			{
				Config: testAccAzureADGroupWithOneOwners(id, pw),
				Check:  testCheckAzureAdGroupBasic(id, "0", "1"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Group with a different owners
			{
				Config: testAccAzureADGroupWithServicePrincipalOwner(id),
				Check:  testCheckAzureAdGroupBasic(id, "0", "1"),
			},
			{
				ResourceName:      rn,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Empty group with 0 owners is not possible
		},
	})
}

func TestAccAzureADGroup_preventDuplicateNames(t *testing.T) {
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureADGroup_duplicateName(ri),
				ExpectError: regexp.MustCompile("existing Azure Active Directory Group .+ was found"),
			},
		},
	})
}

func testCheckAzureADGroupExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).GroupsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Group %q does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("Bad: Get on Azure AD GroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureADGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuread_group" {
			continue
		}

		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).GroupsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD group still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureAdGroupBasic(id int, memberCount, ownerCount string) resource.TestCheckFunc {
	resourceName := "azuread_group.tests"

	return resource.ComposeTestCheckFunc(
		testCheckAzureADGroupExists(resourceName),
		resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("acctestGroup-%d", id)),
		resource.TestCheckResourceAttrSet(resourceName, "object_id"),
		resource.TestCheckResourceAttr(resourceName, "members.#", memberCount),
		resource.TestCheckResourceAttr(resourceName, "owners.#", ownerCount),
	)
}

func testAccAzureADGroup_basic(id int) string {
	return fmt.Sprintf(`
resource "azuread_group" "tests" {
  name    = "acctestGroup-%d"
  members = []
}
`, id)
}

func testAccAzureADGroup_complete(id int, password string) string {
	return fmt.Sprintf(`
%s

resource "azuread_group" "tests" {
  name        = "acctestGroup-%d"
  description = "Please delete me as this is a tests AD group!"
  members     = [azuread_user.tests.object_id]
  owners      = [azuread_user.tests.object_id]
}
`, testAccADUser_basic(id, password), id)
}

func testAccAzureADDiverseDirectoryObjects(id int, password string) string {
	return fmt.Sprintf(`
data "azuread_domains" "tenant_domain" {
  only_initial = true
}

resource "azuread_application" "tests" {
  name = "acctestApp-%[1]d"
}

resource "azuread_service_principal" "tests" {
  application_id = azuread_application.tests.application_id
}

resource "azuread_group" "member" {
  name = "acctestGroup-%[1]d-Member"
}

resource "azuread_user" "tests" {
  user_principal_name = "acctestUser.%[1]d@${data.azuread_domains.tenant_domain.domains.0.domain_name}"
  display_name        = "acctestUser-%[1]d"
  password            = "%[2]s"
}
`, id, password)
}

func testAccAzureADGroupWithDiverseMembers(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name    = "acctestGroup-%[2]d"
  members = [azuread_user.tests.object_id, azuread_group.member.object_id, azuread_service_principal.tests.object_id]
}
`, testAccAzureADDiverseDirectoryObjects(id, password), id)
}

func testAccAzureADGroupWithDiverseOwners(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name   = "acctestGroup-%[2]d"
  owners = [azuread_user.tests.object_id, azuread_service_principal.tests.object_id]
}
`, testAccAzureADDiverseDirectoryObjects(id, password), id)
}

func testAccAzureADGroupWithOneMember(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name    = "acctestGroup-%[2]d"
  members = [azuread_user.tests.object_id]
}
`, testAccADUser_basic(id, password), id)
}

func testAccAzureADGroupWithOneOwners(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name   = "acctestGroup-%[2]d"
  owners = [azuread_user.tests.object_id]
}
`, testAccADUser_basic(id, password), id)
}

func testAccAzureADGroupWithThreeMembers(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name    = "acctestGroup-%[2]d"
  members = [azuread_user.testA.object_id, azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, testAccADUser_threeUsersABC(id, password), id)
}

func testAccAzureADGroupWithThreeOwners(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name   = "acctestGroup-%[2]d"
  owners = [azuread_user.testA.object_id, azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, testAccADUser_threeUsersABC(id, password), id)
}

func testAccAzureADGroupWithOwnersAndMembers(id int, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "azuread_group" "tests" {
  name    = "acctestGroup-%[2]d"
  owners  = [azuread_user.testA.object_id]
  members = [azuread_user.testB.object_id, azuread_user.testC.object_id]
}
`, testAccADUser_threeUsersABC(id, password), id)
}

func testAccAzureADGroupWithServicePrincipalMember(id int) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctestApp-%[1]d"
}

resource "azuread_service_principal" "tests" {
  application_id = azuread_application.tests.application_id
}

resource "azuread_group" "tests" {
  name    = "acctestGroup-%[1]d"
  members = [azuread_service_principal.tests.object_id]
}
`, id)
}

func testAccAzureADGroupWithServicePrincipalOwner(id int) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctestApp-%[1]d"
}

resource "azuread_service_principal" "tests" {
  application_id = azuread_application.tests.application_id
}

resource "azuread_group" "tests" {
  name   = "acctestGroup-%[1]d"
  owners = [azuread_service_principal.tests.object_id]
}
`, id)
}

func testAccAzureADGroup_duplicateName(id int) string {
	return fmt.Sprintf(`
%s

resource "azuread_group" "duplicate" {
  name                    = azuread_group.tests.name
  prevent_duplicate_names = true
}
`, testAccAzureADGroup_basic(id))
}