package token

type Pool struct {
	ID string `json:"id"`
}

type Volume struct {
	Volume string `json:"volume"`
}

type TokenDayData struct {
	VolumeUSD string `json:"volumeUSD"graphql:"volumeUSD"`
}
