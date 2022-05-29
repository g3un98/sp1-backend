package main

type payment struct {
	Type   string `json:"type"`
	Detail string `json:"detail,omitempty"`
	Next   int64  `json:"next"`
}
