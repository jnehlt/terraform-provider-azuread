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

func testCheckADApplicationPasswordExists(name string) resource.TestCheckFunc { //nolint unparam
	return func(s *terraform.State) error {
		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ApplicationsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		id, err := graph.ParseCredentialId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing Application Password Credential ID: %v", err)
		}
		resp, err := client.Get(ctx, id.ObjectId)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Azure AD Application  %q does not exist", id.ObjectId)
			}
			return fmt.Errorf("Bad: Get on Azure AD applicationsClient: %+v", err)
		}

		credentials, err := client.ListPasswordCredentials(ctx, id.ObjectId)
		if err != nil {
			return fmt.Errorf("Error Listing Password Credentials for Application %q: %+v", id.ObjectId, err)
		}

		cred := graph.PasswordCredentialResultFindByKeyId(credentials, id.KeyId)
		if cred != nil {
			return nil
		}

		return fmt.Errorf("Password Credential %q was not found in Application %q", id.KeyId, id.ObjectId)
	}
}

func testCheckADApplicationPasswordCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		client := acceptance.AzureADProvider.Meta().(*clients.AadClient).ApplicationsClient
		ctx := acceptance.AzureADProvider.Meta().(*clients.AadClient).StopContext

		if rs.Type != "azuread_application_password" {
			continue
		}

		id, err := graph.ParseCredentialId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error parsing Application Password Credential ID: %v", err)
		}

		resp, err := client.Get(ctx, id.ObjectId)
		if err != nil {
			if ar.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Azure AD Application Password Credential still exists:\n%#v", resp)
	}

	return nil
}

func TestAccAzureADApplicationPassword_basic(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADObjectPasswordApplication_basic(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADApplicationPasswordExists(resourceName),
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

func TestAccAzureADApplicationPassword_basicOld(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADObjectPasswordApplication_basicOld(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADApplicationPasswordExists(resourceName),
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

func TestAccAzureADApplicationPassword_requiresImport(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADObjectPasswordApplication_basic(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckADApplicationPasswordExists(resourceName),
				),
			},
			{
				Config:      testAccADApplicationPassword_requiresImport(applicationId, value),
				ExpectError: acceptance.RequiresImportError("azuread_application_password"),
			},
		},
	})
}

func TestAccAzureADApplicationPassword_customKeyId(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	keyId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADApplicationPassword_customKeyId(applicationId, keyId, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckADApplicationPasswordExists(resourceName),
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

func TestAccAzureADApplicationPassword_description(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADApplicationPassword_description(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					testCheckADApplicationPasswordExists(resourceName),
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

func TestAccAzureADApplicationPassword_relativeEndDate(t *testing.T) {
	resourceName := "azuread_application_password.tests"
	applicationId := uuid.New().String()
	value := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckADApplicationPasswordCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccADApplicationPassword_relativeEndDate(applicationId, value),
				Check: resource.ComposeTestCheckFunc(
					// can't assert on Value since it's not returned
					testCheckADApplicationPasswordExists(resourceName),
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

func testAccADApplicationPassword_template(applicationId string) string {
	return fmt.Sprintf(`
resource "azuread_application" "tests" {
  name = "acctestApp-%s"
}
`, applicationId)
}

func testAccADObjectPasswordApplication_basic(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "tests" {
  application_object_id = azuread_application.tests.id
  value                 = "%s"
  end_date              = "2099-01-01T01:02:03Z"
}
`, testAccADApplicationPassword_template(applicationId), value)
}

func testAccADObjectPasswordApplication_basicOld(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "tests" {
  application_id = azuread_application.tests.id
  value          = "%s"
  end_date       = "2099-01-01T01:02:03Z"
}
`, testAccADApplicationPassword_template(applicationId), value)
}

func testAccADApplicationPassword_requiresImport(applicationId, value string) string {
	template := testAccADObjectPasswordApplication_basic(applicationId, value)
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "import" {
  application_object_id = azuread_application_password.tests.application_object_id
  key_id                = azuread_application_password.tests.key_id
  value                 = azuread_application_password.tests.value
  end_date              = azuread_application_password.tests.end_date
}
`, template)
}

func testAccADApplicationPassword_customKeyId(applicationId, keyId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "tests" {
  application_object_id = azuread_application.tests.id
  key_id                = "%s"
  value                 = "%s"
  end_date              = "2099-01-01T01:02:03Z"
}
`, testAccADApplicationPassword_template(applicationId), keyId, value)
}

func testAccADApplicationPassword_description(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "tests" {
  application_object_id = azuread_application.tests.id
  description           = "terraform"
  value                 = "%s"
  end_date              = "2099-01-01T01:02:03Z"
}
`, testAccADApplicationPassword_template(applicationId), value)
}

func testAccADApplicationPassword_relativeEndDate(applicationId, value string) string {
	return fmt.Sprintf(`
%s

resource "azuread_application_password" "tests" {
  application_object_id = azuread_application.tests.id
  value                 = "%s"
  end_date_relative     = "8760h"
}
`, testAccADApplicationPassword_template(applicationId), value)
}
