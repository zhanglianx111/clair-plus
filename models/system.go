package models

type Memary struct {
	Totle uint64 `json: total`
	Available uint64 `json:available`
	UsedPercent float64 `json:"usedPercent"`
}

type CPU struct {
	Cores int32 `json:"cores"`
	Mhz   float64 `json:"mhz"`
	UsedPercent float64 `json:"usedPercent"`
}