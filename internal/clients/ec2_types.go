package clients

type EC2Tag struct {
	Key   *string
	Value *string
}

type EC2KeyPairInfo struct {
	KeyPairId      *string
	KeyName        *string
	KeyFingerprint *string
	Tags           []EC2Tag
}
