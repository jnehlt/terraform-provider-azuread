package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/acceptance"
)

func TestAccDataSourceAzureADDomains_basic(t *testing.T) {
	dataSourceName := "data.azuread_domains.tests"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azuread_domains" "tests" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.authentication_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_initial"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_verified"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureADDomains_onlyDefault(t *testing.T) {
	dataSourceName := "data.azuread_domains.tests"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azuread_domains" "tests" {
					only_default = true
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.domain_name"),
					resource.TestCheckResourceAttr(dataSourceName, "domains.0.is_default", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_verified"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureADDomains_onlyInitial(t *testing.T) {
	dataSourceName := "data.azuread_domains.tests"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azuread_domains" "tests" {
					only_initial = true
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_default"),
					resource.TestCheckResourceAttr(dataSourceName, "domains.0.is_initial", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.is_verified"),
				),
			},
		},
	})
}
