// Netflix 계정 구조체 정의

package netflix

// Netflix 계정이 가지는 정보
type Account struct {
    Id, Pw string
    Membership membership
}
