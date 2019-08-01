package azure

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	azureClientID       = "azure-client-id"
	azureClientSecret   = "azure-client-secret"
	azureResourceGroup  = "azure-resource-group"
	azureSubscriptionID = "azure-subscription-id"
	azureTenantID       = "azure-tenant-id"
	managedZone         = "zone"
)

// AddFlags adds flags for Builder.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(azureClientID, "", "Client ID")
	flags.String(azureClientSecret, "", "Client secret")
	flags.String(azureResourceGroup, "", "Resource group")
	flags.String(azureSubscriptionID, "", "Subscription ID")
	flags.String(azureTenantID, "", "Tenant ID")
	flags.String(managedZone, "", "Managed zone")
}

// InitFromViper initializes Builder with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.ClientID = v.GetString(azureClientID)
	b.ClientSecret = v.GetString(azureClientSecret)
	b.ResourceGroup = v.GetString(azureResourceGroup)
	b.SubscriptionID = v.GetString(azureSubscriptionID)
	b.TenantID = v.GetString(azureTenantID)
	b.Zone = v.GetString(managedZone)
	return b
}
