package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
)

func TestAccAzureADApplicationDataSource_byObjectId(t *testing.T) {
	dataSourceName := "data.azuread_application.tests"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(ri),
			},
			{
				Config: testAccAzureADApplicationDataSource_objectId(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "required_resource_access.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "webapp/api"),
					resource.TestCheckResourceAttr(dataSourceName, "app_roles.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_permissions.0.admin_consent_description", fmt.Sprintf("Allow the application to access %s on behalf of the signed-in user.", fmt.Sprintf("acctest-APP-%d", ri))),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAzureADApplicationDataSource_byObjectIdComplete(t *testing.T) {
	dataSourceName := "data.azuread_application.tests"
	ri := tf.AccRandTimeInt()
	pw := "p@$$wR2" + acctest.RandStringFromCharSet(7, acctest.CharSetAlphaNum)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_complete(ri, pw),
			},
			{
				Config: testAccAzureADApplicationDataSource_objectIdComplete(ri, pw),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://homepage-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.0.access_token.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.0.id_token.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "required_resource_access.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "group_membership_claims", "All"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAzureADApplicationDataSource_byApplicationId(t *testing.T) {
	dataSourceName := "data.azuread_application.tests"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(ri),
			},
			{
				Config: testAccAzureADApplicationDataSource_applicationId(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "required_resource_access.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAzureADApplicationDataSource_byName(t *testing.T) {
	dataSourceName := "data.azuread_application.tests"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplication_basic(ri),
			},
			{
				Config: testAccAzureADApplicationDataSource_name(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("https://acctest-APP-%d", ri)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "optional_claims.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "required_resource_access.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func testAccAzureADApplicationDataSource_objectId(ri int) string {
	template := testAccApplication_basic(ri)
	return fmt.Sprintf(`
%s

data "azuread_application" "tests" {
  object_id = azuread_application.tests.object_id
}
`, template)
}

func testAccAzureADApplicationDataSource_objectIdComplete(ri int, pw string) string {
	template := testAccApplication_complete(ri, pw)
	return fmt.Sprintf(`
%s

data "azuread_application" "tests" {
  object_id = azuread_application.tests.object_id
}
`, template)
}

func testAccAzureADApplicationDataSource_applicationId(ri int) string {
	template := testAccApplication_basic(ri)
	return fmt.Sprintf(`
%s

data "azuread_application" "tests" {
  application_id = azuread_application.tests.application_id
}
`, template)
}

func testAccAzureADApplicationDataSource_name(ri int) string {
	template := testAccApplication_basic(ri)
	return fmt.Sprintf(`
%s

data "azuread_application" "tests" {
  name = azuread_application.tests.name
}
`, template)
}
