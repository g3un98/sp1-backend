// 계정 구조체 정의

package service

// 계정이 가지는 정보
type Account struct {
    Id string `json:"id"`
    Pw string `json:"pw"`
    Payment `json:"payment"`
    Membership `json:"membership"`
}
