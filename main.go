// API 서버 최상단에서 각 서비스에 따라 파싱

package main

import (
	"fmt"
	"log"
	"net/http"
)

// API 서버가 동작하는지 확인하기 위한 함수
func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello API!")
	log.Println("Hit endpoint: helloWorld")
}

// localhost:8080으로 API 서버 시작
func parseRequests() {
	http.HandleFunc("/", helloWorld)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	log.Println("Start API server")
	parseRequests()
}
