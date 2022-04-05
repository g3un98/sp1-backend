// Servicer 인터페이스 정의

package service

// 각 서비스는 아래 메소드가 필수
type Servicer interface {
	GetName() (name string)
	Handler()
}
