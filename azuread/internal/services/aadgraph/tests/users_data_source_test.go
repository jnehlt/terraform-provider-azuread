package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
)

func TestAccAzureADUsersDataSource_byUserPrincipalNames(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byUserPrincipalNames(id, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
					resource.TestCheckResourceAttr(dsn, "users.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_byUserPrincipalNamesIgnoreMissing(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byUserPrincipalNamesIgnoreMissing(id, password, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "3"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "3"),
					resource.TestCheckResourceAttr(dsn, "users.#", "3"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_byObjectIds(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byObjectIds(id, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
					resource.TestCheckResourceAttr(dsn, "users.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_byObjectIdsIgnoreMissing(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byObjectIdsIgnoreMissing(id, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
					resource.TestCheckResourceAttr(dsn, "users.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_byMailNicknames(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byMailNicknames(id, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
					resource.TestCheckResourceAttr(dsn, "mail_nicknames.#", "2"),
					resource.TestCheckResourceAttr(dsn, "users.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_byMailNicknamesIgnoreMissing(t *testing.T) {
	dsn := "data.azuread_users.tests"
	id := tf.AccRandTimeInt()
	password := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_byMailNicknamesIgnoreMissing(id, password, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "2"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "2"),
					resource.TestCheckResourceAttr(dsn, "mail_nicknames.#", "2"),
					resource.TestCheckResourceAttr(dsn, "users.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureADUsersDataSource_noNames(t *testing.T) {
	dsn := "data.azuread_users.tests"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureADUsersDataSource_noNames(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dsn, "user_principal_names.#", "0"),
					resource.TestCheckResourceAttr(dsn, "object_ids.#", "0"),
					resource.TestCheckResourceAttr(dsn, "mail_nicknames.#", "0"),
				),
			},
		},
	})
}

func testAccAzureADUsersDataSource_byUserPrincipalNames(id int, password string) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  user_principal_names = [azuread_user.testA.user_principal_name, azuread_user.testB.user_principal_name]
}
`, testAccADUser_threeUsersABC(id, password))
}

func testAccAzureADUsersDataSource_byUserPrincipalNamesIgnoreMissing(id int, password string, ri int) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  ignore_missing = true

  user_principal_names = [
    azuread_user.testA.user_principal_name,
    azuread_user.testB.user_principal_name,
    azuread_user.testC.user_principal_name,
    "not-a-real-user-%d${data.azuread_domains.tenant_domain.domains.0.domain_name}",
  ]
}
`, testAccADUser_threeUsersABC(id, password), ri)
}

func testAccAzureADUsersDataSource_byObjectIds(id int, password string) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  object_ids = [azuread_user.testA.object_id, azuread_user.testB.object_id]
}
`, testAccADUser_threeUsersABC(id, password))
}

func testAccAzureADUsersDataSource_byObjectIdsIgnoreMissing(id int, password string) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  ignore_missing = true

  object_ids = [
    azuread_user.testA.object_id,
    azuread_user.testB.object_id,
    "00000000-0000-0000-0000-000000000000"
  ]
}
`, testAccADUser_threeUsersABC(id, password))
}

func testAccAzureADUsersDataSource_byMailNicknames(id int, password string) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  mail_nicknames = [azuread_user.testA.mail_nickname, azuread_user.testB.mail_nickname]
}
`, testAccADUser_threeUsersABC(id, password))
}

func testAccAzureADUsersDataSource_byMailNicknamesIgnoreMissing(id int, password string, ri int) string {
	return fmt.Sprintf(`
%s

data "azuread_users" "tests" {
  ignore_missing = true

  mail_nicknames = [
    azuread_user.testA.mail_nickname,
    azuread_user.testB.mail_nickname,
    "not-a-real-user-%d${data.azuread_domains.tenant_domain.domains.0.domain_name}"
  ]
}
`, testAccADUser_threeUsersABC(id, password), ri)
}

func testAccAzureADUsersDataSource_noNames() string {
	return `
data "azuread_users" "tests" {
  user_principal_names = []
}
`
}
