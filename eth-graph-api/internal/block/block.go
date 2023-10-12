package block

type Token struct {
	Symbol string `json:"symbol"`
}

type Pool struct {
	Token0 Token `json:"token0"`
	Token1 Token `json:"token1"`
}

type Swap struct {
	ID   string `json:"id"`
	Pool Pool   `json:"pool"`
}
