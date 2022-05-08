// Servicer 인터페이스 정의

package service

import "context"

// 각 서비스가 사용하는 변수
type Service struct {
	ctx context.Context
}

// 각 서비스는 아래 메소드가 필수
type Servicer interface {
	GetName() (name string)
	Handler()
}

// 인자로 들어온 타입에 따라 각 서비스를 반환
func NewService(t string) Servicer {
	switch t {
	case "Netflix":
		return &Netflix{}
	default:
		// 정의되지 않은 서비스 입력 시, 에러 로그 작성 및 프로그램 종료
		LogErr.Fatalf("An error has occurred while create new service: %s\n", t)
		return nil
	}
}
