// Netflix 멤버십 상수 선언

package netflix

type membership int

// Netflix 멤버십 종류
const (
    MEMBERSHIP_NO membership = iota
    MEMBERSHIP_BASIC
    MEMBERSHIP_STANDARD
    MEMBERSHIP_PREMIUM
)

// Netflix 멤버십 종류별 가격
const (
    MEMBERSHIP_COST_BASIC = 9_500
    MEMBERSHIP_COST_STANDARD = 13_500
    MEMBERSHIP_COST_PREMIUM = 17_000
)
