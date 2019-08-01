package azure

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"strings"
)

// check tests if a NSUpdate setup is ok; returns a set of error strings in case something is not right
func (azu *AzUpdater) check() (success bool, errs []string) {
	errMsg := `The "%v" must be specified`
	if strings.TrimSpace(azu.Zone) == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "DNS zone"))
	}
	if strings.TrimSpace(azu.SubscriptionID) == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "SubscriptionID"))
	}
	if strings.TrimSpace(azu.ResourceGroup) == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "ResourceGroup"))
	}
	return len(errs) == 0, errs
}

// checkName checks if the name is in the expected format: subdomain.zone
func (azu *AzUpdater) checkName(name string) (err error) {
	if !strings.HasSuffix(name, "."+azu.Zone) {
		err = types.BadRequestError(fmt.Sprintf("the record name '%s' is not allowed. Must obey the following pattern: '<subdomain>.%s'", name, azu.Zone), nil)
	}
	return
}

// ToFqdn converts the name into a fqdn appending a trailing dot.
func ToFqdn(name string) string {
	n := len(name)
	if n == 0 || name[n-1] == '.' {
		return name
	}
	return name + "."
}

// UnFqdn converts the fqdn into a name removing the trailing dot.
func UnFqdn(name string) string {
	n := len(name)
	if n != 0 && name[n-1] == '.' {
		return name[:n-1]
	}
	return name
}
