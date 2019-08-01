package azure

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/dns/mgmt/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
	"strings"
	"time"

	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
)

type Builder struct {
	Zone string

	ClientID     string
	ClientSecret string
	TenantID     string

	SubscriptionID string
	ResourceGroup  string

	HTTPClient *http.Client
}

// AzUpdater holds the information necessary to successfully run update requests
type AzUpdater struct {
	Builder
	client *dns.RecordSetsClient
}

// DNSUpdater defines an interface to communicate with DNS Server via update commands
type DNSUpdater interface {
	RemoveRR(name, recordType string) (err error)
	AddRR(record hookTypes.DNSRecord, ttl time.Duration) (err error)
	UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (err error)
}

// New constructs a new AzUpdater instance from environment variables
func (b *Builder) New() (*AzUpdater, error) {
	result := &AzUpdater{Builder: *b}

	if succ, errs := result.check(); !succ {
		return nil, fmt.Errorf("Errors encountered:\n\t%v", strings.Join(errs, "\n\t"))
	}
	if b.HTTPClient == nil {
		b.HTTPClient = http.DefaultClient
	}

	authorizer, err := getAuthorizer(b)
	if err != nil {
		return nil, err
	}

	// just one instance
	rsc := dns.NewRecordSetsClient(b.SubscriptionID)
	rsc.Authorizer = authorizer
	result.client = &rsc

	return result, nil
}

func getAuthorizer(config *Builder) (autorest.Authorizer, error) {
	if config.ClientID != "" && config.ClientSecret != "" && config.TenantID != "" {
		oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, config.TenantID)
		if err != nil {
			return nil, err
		}

		spt, err := adal.NewServicePrincipalToken(*oauthConfig, config.ClientID, config.ClientSecret, azure.PublicCloud.ResourceManagerEndpoint)
		if err != nil {
			return nil, err
		}

		spt.SetSender(config.HTTPClient)
		return autorest.NewBearerAuthorizer(spt), nil
	}
	return auth.NewAuthorizerFromEnvironment()
}

// RemoveRR removes a Resource Record
func (azu *AzUpdater) RemoveRR(name, recordType string) (err error) {
	err = azu.checkName(name)
	relative := toRelativeRecord(name, ToFqdn(azu.Zone))
	_, err = azu.client.Delete(context.Background(), azu.ResourceGroup, azu.Zone, relative, dns.RecordType(recordType), "")
	if err != nil {
		err = fmt.Errorf("azure: %v", err)
		return
	}
	return
}

// AddRR adds a Resource Record
func (azu *AzUpdater) AddRR(record hookTypes.DNSRecord, ttl time.Duration) (err error) {
	return azu.createOrUpdate(record, ttl)
}

// UpdateRR updates a DNS Resource Record
func (azu *AzUpdater) UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (err error) {
	return azu.createOrUpdate(record, ttl)
}

func (azu *AzUpdater) createOrUpdate(record hookTypes.DNSRecord, ttl time.Duration) (err error) {
	err = azu.checkName(record.Name)
	if err != nil {
		return
	}
	relative := toRelativeRecord(record.Name, ToFqdn(azu.Zone))
	recordSetProperties, err := recordSetProperties(record)
	if err != nil {
		err = fmt.Errorf("azure: %v", err)
		return
	}
	recordSetProperties.TTL = to.Int64Ptr(int64(ttl.Seconds()))
	rec := dns.RecordSet{
		Name:                &relative,
		RecordSetProperties: recordSetProperties,
	}

	_, err = azu.client.CreateOrUpdate(context.Background(), azu.ResourceGroup, azu.Zone, relative, dns.RecordType(record.Type), rec, "", "")
	if err != nil {
		err = fmt.Errorf("azure: %v", err)
		return
	}
	return
}

func recordSetProperties(record hookTypes.DNSRecord) (*dns.RecordSetProperties, error) {
	var properties *dns.RecordSetProperties
	switch record.Type {
	case "A":
		properties = &dns.RecordSetProperties{
			ARecords: &[]dns.ARecord{{Ipv4Address: &record.Value}},
		}
	case "AAAA":
		properties = &dns.RecordSetProperties{
			AaaaRecords: &[]dns.AaaaRecord{{Ipv6Address: &record.Value}},
		}
	case "CNAME":
		properties = &dns.RecordSetProperties{
			CnameRecord: &dns.CnameRecord{Cname: &record.Value},
		}
	case "MX":
		properties = &dns.RecordSetProperties{
			MxRecords: &[]dns.MxRecord{{Exchange: &record.Value}},
		}
	case "NS":
		properties = &dns.RecordSetProperties{
			NsRecords: &[]dns.NsRecord{{Nsdname: &record.Value}},
		}
	case "PTR":
		properties = &dns.RecordSetProperties{
			PtrRecords: &[]dns.PtrRecord{{Ptrdname: &record.Value}},
		}
	case "TXT":
		properties = &dns.RecordSetProperties{
			TxtRecords: &[]dns.TxtRecord{{Value: &[]string{record.Value}}},
		}
	default:
		return nil, fmt.Errorf("record type %s not supported", record.Type)
	}
	return properties, nil
}

// Returns the relative record to the domain
func toRelativeRecord(domain, zone string) string {
	return UnFqdn(strings.TrimSuffix(domain, zone))
}
