package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/clients"
)

func TestAccAzureADServicePrincipal_basic(t *testing.T) {
	resourceName := "azuread_service_principal.tests"
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipal_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
					resource.TestCheckResourceAttr(resourceName, "oauth2_permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "app_role_assignment_required", "false"),
					resource.TestCheckResourceAttr(resourceName, "oauth2_permissions.0.admin_consent_description", fmt.Sprintf("Allow the application to access %s on behalf of the signed-in user.", fmt.Sprintf("acctestApp-%s", id))),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADServicePrincipal_complete(t *testing.T) {
	resourceName := "azuread_service_principal.tests"
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipal_complete(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "app_role_assignment_required", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureADServicePrincipal_update(t *testing.T) {
	resourceName := "azuread_service_principal.tests"
	id := uuid.New().String()
	updatedId := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipal_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
					resource.TestCheckResourceAttr(resourceName, "app_role_assignment_required", "false"),
				),
			},
			{
				Config: testAccADServicePrincipal_complete(updatedId),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "object_id"),
					resource.TestCheckResourceAttr(resourceName, "app_role_assignment_required", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccADServicePrincipal_basic(id),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
					resource.TestCheckResourceAttr(resourceName, "app_role_assignment_required", "false"),
				),
			},
		},
	})
}

func testCheckADServicePrincipalExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ServicePrincipalsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Service Principal %q does not exist", rs.Primary.ID)
			}
			return fmt.Errorf("Bad: Get on Azure AD ServicePrincipalsClient: %+v", err)
		}

		return nil
	}
}

func testCheckADServicePrincipalDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuread_service_principal" {
			continue
		}

		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ServicePrincipalsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext
		resp, err := client.Get(ctx, rs.Primary.ID)

		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD Service Principal still exists:\n%#v", resp)
	}

	return nil
}

func testAccADServicePrincipal_basic(id string) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctestApp-%s"
}

resource "azuread_service_principal" "tests" {
  application_id = azuread_application.tests.application_id
}
`, id)
}

func testAccADServicePrincipal_complete(id string) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctesttApp-%s"
}

resource "azuread_service_principal" "tests" {
  application_id               = azuread_application.tests.application_id
  app_role_assignment_required = true

  tags = ["tests", "multiple", "CapitalS"]
}
`, id)
}
