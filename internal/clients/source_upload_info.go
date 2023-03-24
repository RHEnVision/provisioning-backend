package clients

type AccountIdentity struct {
	AWSDetails *AccountDetailsAWS `json:"aws,omitempty" yaml:"aws"`
}

type AccountDetailsAWS struct {
	AccountID string `json:"account_id" yaml:"account_id"`
}

type AzureSourceDetail struct {
	TenantID       string   `json:"tenant_id"`
	SubscriptionID string   `json:"subscription_id"`
	ResourceGroups []string `json:"resource_groups"`
}
