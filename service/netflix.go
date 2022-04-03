// Netflix 관련 API 제공

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
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
	log.Println("[/netflix] Netflix.Hello")
}

// Netflix 계정 정보를 가져오는 함수
func (n Netflix) Info(w http.ResponseWriter, r *http.Request) {
    // 리퀘스트로부터 계정 id, pw를 받아옴
    var account Account
    json.NewDecoder(r.Body).Decode(&account)

    ctx, cancel := chromedp.NewContext(
        context.Background(),
        //chromedp.WithDebugf(log.Printf),
    )
    defer cancel()

    var (
        // 가공하기 전 데이터를 저장하는 변수
        raw string
    )

    // TODO: 로그인, 로그아웃 분리
    // chromedp.Run 또는 []chromedp.Action을 사용하는 방식으로 구현하면 에러가 발생해서
    // 일단 구현함
    if err := chromedp.Run(
        ctx,

        // FIXME: 로그인 과정에서 에러가 발생하는 경우를 확인하지 않음 (e.g., 계정 정보 오류, 사이트 서버 문제)
        // 로그인
        chromedp.Navigate(`https://www.netflix.com/kr/login/`),
        chromedp.Click(`//input[@id="id_userLoginId"]`, chromedp.NodeVisible),
        chromedp.SendKeys(`//input[@id="id_userLoginId"]`, account.Id, chromedp.NodeVisible),
        chromedp.Click(`//input[@id="id_password"]`, chromedp.NodeVisible),
        chromedp.SendKeys(`//input[@id="id_password"]`, account.Pw, chromedp.NodeVisible),
        chromedp.Click(`//button[@class="btn login-button btn-submit btn-small"]`, chromedp.NodeVisible),
        chromedp.Sleep(1 * time.Second),

        // 로그인 과정에서 에러가 발생했을 경우
        //chromedp.Text(`//div[@class="ui-message-container ui-message-error"]`, &res, chromedp.NodeVisible),

        // 계정 정보를 가져옴
        chromedp.Navigate(`https://www.netflix.com/youraccount/`),
        chromedp.Text(`//div[@class="bd"]`, &raw, chromedp.NodeVisible),
        chromedp.Sleep(1 * time.Second),

        // FIXME: 현재 로그인 상태인지 확인하지 않음
        // 로그아웃
        chromedp.Navigate(`https://www.netflix.com/kr/signout/`),
        chromedp.Sleep(1 * time.Second),
    ); err != nil {
        log.Fatal(err)
    }

    fmt.Fprintln(w, raw)
	log.Println("[/netflix/Info] Netflix.Info")
}

// Netflix 핸들러를 패턴에 맞게 연결
func (n Netflix) Handler() {
    http.HandleFunc("/netflix", n.Hello)
    http.HandleFunc("/netflix/info", n.Info)
}
