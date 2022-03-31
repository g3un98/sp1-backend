// API 서버 최상단에서 각 서비스에 따라 파싱

package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"service"
)

// API 서버가 동작하는지 확인하기 위한 함수
func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello API!")
	log.Println("[/] helloWorld")
}

// localhost:8080으로 API 서버 시작
// 각 서비스 핸들러 호출
func handleRequests() {
	var (
		// 핸들러 초기화를 위해, 각 서비스를 배열에 삽입
		services = [...]service.Servicer{
			service.Netflix{},
		}
		// 동기화 작업을 위한 WaitGroup
		wg sync.WaitGroup
	)

	// "/" 경로와 helloWorld 함수를 연결
	http.HandleFunc("/", helloWorld)

	// services 각 서비스의 핸들러 호출
	for _, s := range services {
		wg.Add(1)
		go func(s service.Servicer) {
			defer wg.Done()

			log.Printf("Prepare %s APIs\n", s.GetName())
			s.Handler()
			log.Printf("%s APIs are ready\n", s.GetName())
		}(s)
	}

	// 모든 핸들러가 동작을 완료할 때까지 대기
	wg.Wait()

	// localhost:8080으로 서버 시작
	// 에러가 발생할 시, 로그 작성 및 프로그램 종료
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// API 서버 시작 로그를 남기고, 요청 핸들러 호출
func main() {
	log.Println("Start API server")
	handleRequests()
}
