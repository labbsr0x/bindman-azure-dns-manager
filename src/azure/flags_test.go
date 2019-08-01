package azure

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBingFlags(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	clientIdValue := "client-id-value"
	clientSecretValue := "client-secret-value"
	resourceGroupValue := "resource-group-value"
	subscriptionIdValue := "subscription-id-value"
	tenantIdValue := "tenant-id-value"
	managedZoneValue := "managed-zone-value"

	err := command.ParseFlags([]string{
		fmt.Sprintf("--%s=%s", azureClientID, clientIdValue),
		fmt.Sprintf("--%s=%s", azureClientSecret, clientSecretValue),
		fmt.Sprintf("--%s=%s", azureResourceGroup, resourceGroupValue),
		fmt.Sprintf("--%s=%s", azureSubscriptionID, subscriptionIdValue),
		fmt.Sprintf("--%s=%s", azureTenantID, tenantIdValue),
		fmt.Sprintf("--%s=%s", managedZone, managedZoneValue),
	})
	require.NoError(t, err)

	b := &Builder{}
	b.InitFromViper(v)

	assert.Equal(t, clientIdValue, b.ClientID)
	assert.Equal(t, clientSecretValue, b.ClientSecret)
	assert.Equal(t, resourceGroupValue, b.ResourceGroup)
	assert.Equal(t, subscriptionIdValue, b.SubscriptionID)
	assert.Equal(t, tenantIdValue, b.TenantID)
	assert.Equal(t, managedZoneValue, b.Zone)
}

func TestDefaultValues(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	b := &Builder{}
	b.InitFromViper(v)

	assert.Equal(t, "", b.ClientID)
	assert.Equal(t, "", b.ClientSecret)
	assert.Equal(t, "", b.ResourceGroup)
	assert.Equal(t, "", b.SubscriptionID)
	assert.Equal(t, "", b.TenantID)
	assert.Equal(t, "", b.Zone)
}
