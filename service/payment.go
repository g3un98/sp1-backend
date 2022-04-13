// Netflix 결제 방식 구조체 정의

package service

// 결제 방식이 가지는 정보
type Payment struct {
	Type   string `json:"type"`
	Detail string `json:"detail"`
	Next   string `json:"next"`
}
