// Netflix 관련 API 제공

package service

import (
	"fmt"
	"log"
	"net/http"
)

// Netflix 파싱이 동작하는지 확인하기 위한 함수
func helloNetflix(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello Netflix!")
	log.Println("[/netflix] helloNetflix")
}

// TODO
// Netflix 파싱이 동작하는지 확인하기 위한 함수
func loginNetflix(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "NOT IMPLEMENT")
	log.Println("[/netflix/login] loginNetflix")
}

// Netflix 핸들러를 패턴에 맞게 연결
func HandleNetflix() {
    http.HandleFunc("/netflix", helloNetflix)
    http.HandleFunc("/netflix/login", loginNetflix)
}
