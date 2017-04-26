package dto

//Currency represents two currency codes, a from and a to, to aid in conversion
type Currency struct {
	CodeFrom string `json:"from"`
	CodeTo   string `json:"to"`
}

//Rate represents an exchange rate of a currency
type Rate struct {
	Conversion Currency `json:"conversion"`
	Value      float64  `json:"value"`
}

//Represente a con and an amount
type WalletItem struct {
	coin  string  `json: "coin"`
	value float64 `json: "value"`
}

type WalletItems []WalletItem
