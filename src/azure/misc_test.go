package azure

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"testing"
)

func TestAzUpdater_check(t *testing.T) {
	errMsg := `The "%v" must be specified`
	errorMsgRg := fmt.Sprintf(errMsg, "ResourceGroup")
	errorMsgSubscription := fmt.Sprintf(errMsg, "SubscriptionID")
	errorMsgDnsZone := fmt.Sprintf(errMsg, "DNS zone")

	type returnValue struct {
		success bool
		errs    []string
	}
	testCases := []struct {
		name      string
		azUpdater AzUpdater
		expected  returnValue
	}{
		{
			"all OK",
			AzUpdater{Builder{SubscriptionID: "sub-value", ResourceGroup: "rg-value", Zone: "test.com"}, nil},
			returnValue{true, []string{}},
		},
		{
			"all required fields",
			AzUpdater{Builder{}, nil},
			returnValue{false, []string{errorMsgDnsZone, errorMsgSubscription, errorMsgRg}},
		},
		{
			"subscription required",
			AzUpdater{Builder{ResourceGroup: "rg-value", Zone: "test.com"}, nil},
			returnValue{false, []string{errorMsgSubscription}},
		},
		{
			"resource group required",
			AzUpdater{Builder{SubscriptionID: "sub-value", Zone: "test.com"}, nil},
			returnValue{false, []string{errorMsgRg}},
		},
		{
			"DNS zone required",
			AzUpdater{Builder{SubscriptionID: "sub-value", ResourceGroup: "rg-value"}, nil},
			returnValue{false, []string{errorMsgDnsZone}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			success, errs := test.azUpdater.check()
			if success != test.expected.success {
				t.Errorf("It was expected success=false but returned true")
			}
			if len(errs) != len(test.expected.errs) {
				t.Errorf("The error array length must be %d but got %d", len(test.expected.errs), len(errs))
				t.FailNow()
			}
			for i, err := range test.expected.errs {
				if errs[i] != err {
					t.Errorf("Expected message was %s but got %s", err, errs[i])
				}
			}
		})
	}
}

func TestAzUpdater_checkName(t *testing.T) {
	azUpdater := AzUpdater{Builder{Zone: "test.com."}, nil}
	errorMsg := "the record name '%s' is not allowed. Must obey the following pattern: '<subdomain>.%s'"

	testCases := []struct {
		name     string
		expected error
	}{
		{"teste.io.", types.BadRequestError(fmt.Sprintf(errorMsg, "teste.io.", azUpdater.Zone), nil)},
		{".test.com", types.BadRequestError(fmt.Sprintf(errorMsg, ".test.com", azUpdater.Zone), nil)},
		{"test.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "test.com.", azUpdater.Zone), nil)},
		{"subdomain.test.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.test.com", azUpdater.Zone), nil)},
		{"subdomain.test.com.br", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.test.com.br", azUpdater.Zone), nil)},
		{"subdomain.subdomain.test.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.subdomain.test.com", azUpdater.Zone), nil)},
		{"subdomain.teste.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com", azUpdater.Zone), nil)},
		{"subdomain.teste.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com.", azUpdater.Zone), nil)},
		{"subdomain.etest.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.etest.com", azUpdater.Zone), nil)},
		{"subdomain.etest.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.etest.com.", azUpdater.Zone), nil)},
		{"subdomain.teste.com.br.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com.br.", azUpdater.Zone), nil)},
		{"subdomain.subdomain.test.com.", nil},
		{"subdomain.test.com.", nil},
		{"a.test.com.", nil},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := azUpdater.checkName(test.name)
			if test.expected == nil {
				if err != nil {
					t.Errorf("got = %v, want %v", err, test.expected)
				}
			} else {
				if err == nil || err.Error() != test.expected.Error() {
					t.Errorf("got = %v, want %v", err, test.expected)
				}
			}
		})
	}
}

func TestToFqdn(t *testing.T) {
	testCases := []struct {
		desc     string
		domain   string
		expected string
	}{
		{
			desc:     "simple",
			domain:   "foo.bar.com",
			expected: "foo.bar.com.",
		},
		{
			desc:     "already FQDN",
			domain:   "foo.bar.com.",
			expected: "foo.bar.com.",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.domain, func(t *testing.T) {
			if got := ToFqdn(tt.domain); got != tt.expected {
				t.Errorf("ToFqdn() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnFqdn(t *testing.T) {
	testCases := []struct {
		desc     string
		fqdn     string
		expected string
	}{
		{
			desc:     "simple",
			fqdn:     "foo.bar.com.",
			expected: "foo.bar.com",
		},
		{
			desc:     "already domain",
			fqdn:     "foo.bar.com",
			expected: "foo.bar.com",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.fqdn, func(t *testing.T) {
			if got := UnFqdn(tt.fqdn); got != tt.expected {
				t.Errorf("UnFqdn() = %v, want %v", got, tt.expected)
			}
		})
	}
}
