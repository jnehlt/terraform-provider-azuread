package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/clients"
)

func testCheckADServicePrincipalPasswordExists(name string) resource.TestCheckFunc { //nolint unparam
	return func(s *terraform.State) error {
		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ServicePrincipalsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		id, err := graph.ParseCredentialId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error Service Principal Password Credential ID: %v", err)
		}

		resp, err := client.Get(ctx, id.ObjectId)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Service Principal %q does not exist", id.ObjectId)
			}
			return fmt.Errorf("Bad: Get on Azure AD ServicePrincipalsClient: %+v", err)
		}

		credentials, err := client.ListPasswordCredentials(ctx, id.ObjectId)
		if err != nil {
			return fmt.Errorf("Error Listing Password Credentials for Service Principal %q: %+v", id.ObjectId, err)
		}

		cred := graph.PasswordCredentialResultFindByKeyId(credentials, id.KeyId)
		if cred != nil {
			return nil
		}

		return fmt.Errorf("Password Credential %q was not found in Service Principal %q", id.KeyId, id.ObjectId)
	}
}

func testCheckADServicePrincipalPasswordCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ApplicationsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext

		if rs.Type != "azuread_service_principal_password" {
			continue
		}

		id, err := graph.ParseCredentialId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing Service Principal Password Credential ID: %v", err)
		}

		resp, err := client.Get(ctx, id.ObjectId)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD Service Principal Password Credential still exists:\n%#v", resp)
	}

	return nil
}

func TestAccAzureADServicePrincipalPassword_basic(t *testing.T) {
	resourceName := "azuread_service_principal_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipalPassword_basic(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2099-01-01T01:02:03Z"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAzureADServicePrincipalPassword_requiresImport(t *testing.T) {
	resourceName := "azuread_service_principal_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipalPassword_basic(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckADServicePrincipalPasswordExists(resourceName),
				),
			},
			{
				Config:      testAccADServicePrincipalPassword_requiresImport(applicationId, value),
				ExpectError: acceptance.RequiresImportError("azuread_service_principal_password"),
			},
		},
	})
}

func TestAccAzureADServicePrincipalPassword_customKeyId(t *testing.T) {
	resourceName := "azuread_service_principal_password.tests"
	applicationId := uuid.New().String()
	keyId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipalPassword_customKeyId(applicationId, keyId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttr(resourceName, "key_id", keyId),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2099-01-01T01:02:03Z"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAzureADServicePrincipalPassword_description(t *testing.T) {
	resourceName := "azuread_service_principal_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipalPassword_description(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "end_date", "2099-01-01T01:02:03Z"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccAzureADServicePrincipalPassword_relativeEndDate(t *testing.T) {
	resourceName := "azuread_service_principal_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADServicePrincipalPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADServicePrincipalPassword_relativeEndDate(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADServicePrincipalPasswordExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "start_date"),
					resource.TestCheckResourceAttrSet(resourceName, "key_id"),
					resource.TestCheckResourceAttrSet(resourceName, "end_date"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"end_date_relative", "value"},
			},
		},
	})
}

func testAccADServicePrincipalPassword_template(applicationId string) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctestApp-%s"
}

resource "azuread_service_principal" "tests" {
  application_id = azuread_application.tests.application_id
}
`, applicationId)
}

func testAccADServicePrincipalPassword_basic(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_service_principal_password" "tests" {
  service_principal_id = azuread_service_principal.tests.id
  value                = "%s"
  end_date             = "2099-01-01T01:02:03Z"
}
`, testAccADServicePrincipalPassword_template(applicationId), value)
}

func testAccADServicePrincipalPassword_requiresImport(applicationId, value string) string {
	template := testAccADServicePrincipalPassword_basic(applicationId, value)
	return fmt.Sprintf(`
%s

resource "azuread_service_principal_password" "import" {
  key_id               = azuread_service_principal_password.tests.key_id
  service_principal_id = azuread_service_principal_password.tests.service_principal_id
  value                = azuread_service_principal_password.tests.value
  end_date             = azuread_service_principal_password.tests.end_date
}
`, template)
}

func testAccADServicePrincipalPassword_customKeyId(applicationId, keyId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_service_principal_password" "tests" {
  service_principal_id = azuread_service_principal.tests.id
  key_id               = "%s"
  value                = "%s"
  end_date             = "2099-01-01T01:02:03Z"
}
`, testAccADServicePrincipalPassword_template(applicationId), keyId, value)
}

func testAccADServicePrincipalPassword_description(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_service_principal_password" "tests" {
  service_principal_id = azuread_service_principal.tests.id
  description          = "terraform"
  value                = "%s"
  end_date             = "2099-01-01T01:02:03Z"
}
`, testAccADServicePrincipalPassword_template(applicationId), value)
}

func testAccADServicePrincipalPassword_relativeEndDate(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_service_principal_password" "tests" {
  service_principal_id = azuread_service_principal.tests.id
  value                = "%s"
  end_date_relative    = "8760h"
}
`, testAccADServicePrincipalPassword_template(applicationId), value)
}
