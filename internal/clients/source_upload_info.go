package clients

type AccountIdentity struct {
	AWSDetails *AccountDetailsAWS `json:"aws,omitempty" yaml:"aws"`
}

type AccountDetailsAWS struct {
	AccountID string `json:"account_id" yaml:"account_id"`
}

func (a AccountDetailsAWS) CacheKeyName() string {
	return "account_detail_aws"
}

type AzureTenantId string

func (a AzureTenantId) CacheKeyName() string {
	return "azure_tenant_id"
}

type AccountDetailsAzure struct {
	TenantID       AzureTenantId `json:"tenant_id"`
	SubscriptionID string        `json:"subscription_id"`
	ResourceGroups []string      `json:"resource_groups"`
}

type AccountDetailsGCP struct{}
