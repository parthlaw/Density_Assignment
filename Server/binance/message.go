package binance

type Message struct {
	Kline  Kline  `json:"k"`
	Symbol string `json:"s"`
}
type Kline struct {
	Open_Price  string `json:"o"`
	Close_Price string `json:"c"`
}
