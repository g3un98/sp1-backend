package main

type payment struct {
	Type   string `json:"type"`
	Detail string `json:"detail,omit"`
	Next   int64  `json:"next"`
}
