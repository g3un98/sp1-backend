// Netflix 결제 방식 구조체 정의

package service

type Payment struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Next   string `json:"next"`
}
