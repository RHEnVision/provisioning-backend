package clients

type AccountIdentity struct {
	AWSDetails *AccountDetailsAWS `json:"aws,omitempty" yaml:"aws"`
}

type AccountDetailsAWS struct {
	AccountID string `json:"account_id" yaml:"account_id"`
}
