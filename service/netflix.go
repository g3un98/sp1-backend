// Netflix 관련 API 제공

package service

import (
	"fmt"
	"log"
	"net/http"
)

// 메소드를 구현하기 위해, 빈 구조체 선언
type Netflix struct {}

// Netflix 서비스명 반환
func (n Netflix) GetName() (name string) {
    return "Netflix"
}

// Netflix 파싱이 동작하는지 확인하기 위한 함수
func (n Netflix) Hello(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello Netflix!")
	log.Println("[/netflix] helloNetflix")
}

// TODO
// Netflix에 로그인
func (n Netflix) Login(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "NOT IMPLEMENT")
	log.Println("[/netflix/login] loginNetflix")
}

// Netflix 핸들러를 패턴에 맞게 연결
func (n Netflix) Handler() {
    http.HandleFunc("/netflix", n.Hello)
    http.HandleFunc("/netflix/login", n.Login)
}
