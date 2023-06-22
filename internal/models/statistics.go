package models

type Statistics struct {
	Usage24h []*UsageStat
	Usage28d []*UsageStat
}

type UsageStat struct {
	Provider ProviderType `db:"provider"`
	Result   string       `db:"result"`
	Count    int64        `db:"count"`
}
