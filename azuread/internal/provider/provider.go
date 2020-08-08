package provider

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/services/aadgraph"

	"github.com/hashicorp/go-azure-helpers/authentication"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/internal/clients"
)

// Provider returns a terraform.ResourceProvider.
func AzureADProvider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			// TODO: remove subscription_id field at next major version
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_SUBSCRIPTION_ID", ""),
			},

			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_ID", ""),
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_TENANT_ID", ""),
			},

			"environment": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_ENVIRONMENT", "public"),
			},

			// Client Certificate specific fields
			"client_certificate_password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PASSWORD", ""),
			},

			"client_certificate_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_CERTIFICATE_PATH", ""),
			},

			// Client Secret specific fields
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_CLIENT_SECRET", ""),
			},

			// Managed Service Identity specific fields
			"use_msi": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_USE_MSI", false),
			},
			"msi_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ARM_MSI_ENDPOINT", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"azuread_application":       aadgraph.DataApplication(),
			"azuread_domains":           aadgraph.DataDomains(),
			"azuread_client_config":     aadgraph.DataClientConfig(),
			"azuread_group":             aadgraph.DataGroup(),
			"azuread_groups":            aadgraph.DataGroups(),
			"azuread_service_principal": aadgraph.DataServicePrincipal(),
			"azuread_user":              aadgraph.DataUser(),
			"azuread_users":             aadgraph.DataUsers(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"azuread_application":                   aadgraph.ResourceApplication(),
			"azuread_application_certificate":       aadgraph.ResourceApplicationCertificate(),
			"azuread_application_password":          aadgraph.ResourceApplicationPassword(),
			"azuread_group":                         aadgraph.ResourceGroup(),
			"azuread_group_member":                  aadgraph.ResourceGroupMember(),
			"azuread_service_principal":             aadgraph.ResourceServicePrincipal(),
			"azuread_service_principal_certificate": aadgraph.ResourceServicePrincipalCertificate(),
			"azuread_service_principal_password":    aadgraph.ResourceServicePrincipalPassword(),
			"azuread_user":                          aadgraph.ResourceUser(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		// TODO: drop subscription_id in v1.0
		// When constructing the Builder, we default to using the tenant ID for the subscription ID.
		// Although this has no effect since we never consume it, this practise mimics
		// the Azure CLI and it seems the most sensible value to use after a nonsense string.
		// However, if subscription_id _is_ configured for the provider, we'll use that since it's
		// currently exposed via data.azuread_client_config.
		subscriptionId := d.Get("subscription_id").(string)
		if subscriptionId == "" {
			subscriptionId = d.Get("tenant_id").(string)
		}

		builder := &authentication.Builder{
			ClientID:           d.Get("client_id").(string),
			ClientSecret:       d.Get("client_secret").(string),
			SubscriptionID:     subscriptionId,
			TenantID:           d.Get("tenant_id").(string),
			Environment:        d.Get("environment").(string),
			MsiEndpoint:        d.Get("msi_endpoint").(string),
			ClientCertPassword: d.Get("client_certificate_password").(string),
			ClientCertPath:     d.Get("client_certificate_path").(string),

			// Feature Toggles
			SupportsClientCertAuth:         true,
			SupportsClientSecretAuth:       true,
			SupportsManagedServiceIdentity: d.Get("use_msi").(bool),
			SupportsAzureCliToken:          true,
		}

		config, err := builder.Build()
		if err != nil {
			return nil, fmt.Errorf("Error building AzureAD Client: %s", err)
		}

		client, err := clients.GetAadClient(config, p.TerraformVersion, p.StopContext())
		if err != nil {
			return nil, err
		}

		client.StopContext = p.StopContext()

		// replaces the context between tests
		p.MetaReset = func() error { //nolint unparam
			client.StopContext = p.StopContext()
			return nil
		}

		return client, nil
	}
}
