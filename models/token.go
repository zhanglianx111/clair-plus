package models

type Token struct {
	ExpiresIn int    `json:"expires_in"`
	IssuedAt  string `json:"issued_at"`
	Token     string `json:"token"`
}

type Image struct {
	Repo string `json:"repo"`
	Tag  string `json:"tag"`
}
