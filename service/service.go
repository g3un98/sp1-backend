// 기본 구조체와 인터페이스 정의

package service

// 각 서비스는 아래 메소드들이 필수로 정의되어야 함
type Servicer interface {
    GetName() (name string)
    Handler()
}
