// API 서버 최상단에서 각 서비스에 따라 파싱

package main

import (
	"net/http"
	"sync"

	"github.com/g3un/sp1-backend/service"
)

// localhost:8000으로 API 서버 시작
// 각 서비스 핸들러 호출
func handleRequests() {
	var (
		// 핸들러 초기화를 위해, 각 서비스를 배열에 삽입
		services = [...]service.Servicer{
            service.NewService("Netflix"),
            service.NewService("Wavve"),
		}
		// 동기화 작업을 위한 WaitGroup
		wg sync.WaitGroup
	)

	// 배열 속 각 서비스의 핸들러 호출
	for _, s := range services {
		wg.Add(1)
		go func(s service.Servicer) {
			defer wg.Done()

			service.LogInfo.Printf("Prepare %s APIs\n", s.GetName())
			s.Handler()
			service.LogInfo.Printf("%s APIs are ready\n", s.GetName())
		}(s)
	}

	// 모든 핸들러가 동작을 완료할 때까지 대기
	wg.Wait()

	// localhost:8000으로 서버 시작
	// 에러 발생 시, 로그 작성 및 프로그램 종료
    service.LogErr.Fatal(http.ListenAndServe(":8000", nil))
}

// API 서버 시작 로그를 남기고, 요청 핸들러 호출
func main() {
	service.LogInfo.Println("Start API server")
	handleRequests()
}
