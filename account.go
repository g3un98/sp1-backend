package main

type account struct {
	Id         string     `json:"id"`
	Pw         string     `json:"pw"`
	Payment    payment    `json:"payment"`
	Membership membership `json:"membership"`
}
