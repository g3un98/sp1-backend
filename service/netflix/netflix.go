// Netflix API 구현

package netflix

import (
	"context"
    "strings"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

// Netflix 구조체 선언
type Netflix struct {
    ctx context.Context
}

// Netflix 서비스명 반환
func (n Netflix) GetName() (name string) {
    return "Netflix"
}

// Netflix 파싱이 동작하는지 확인
func (n *Netflix) Hello(w http.ResponseWriter, _ *http.Request) {
	log.Println("[/netflix] Netflix.Hello")
	fmt.Fprintln(w, "Hello Netflix!")
}

// Netflix 웹사이트 로그인
func (n *Netflix) Login(a Account) (msg string, err error) {
    var url string

    if err = chromedp.Run(
        n.ctx,

        chromedp.Navigate(URI_LOGIN),

        chromedp.WaitVisible(SEL_LOGIN_ID),
        chromedp.Click(SEL_LOGIN_ID, chromedp.NodeVisible),
        chromedp.SendKeys(SEL_LOGIN_ID, a.Id, chromedp.NodeVisible),

        chromedp.WaitVisible(SEL_LOGIN_PW),
        chromedp.Click(SEL_LOGIN_PW, chromedp.NodeVisible),
        chromedp.SendKeys(SEL_LOGIN_PW, a.Pw, chromedp.NodeVisible),

        chromedp.WaitVisible(SEL_LOGIN_BTN),
        chromedp.Click(SEL_LOGIN_BTN, chromedp.NodeVisible),

        chromedp.Sleep(1 * time.Second),
        chromedp.Location(&url),
    ); err != nil {
        return
    }

    if url == URI_LOGIN {
        if err = chromedp.Run(
            n.ctx,
            chromedp.Text(SEL_LOGIN_ERR, &msg, chromedp.NodeVisible),
        ); err != nil {
            return
        }
        return
    }
    return
}

// Netflix 웹사이트 로그아웃
func (n *Netflix) Logout() (err error) {
    return chromedp.Run(
        n.ctx,
        chromedp.Navigate(URI_LOGOUT),
    )
}

// Netflix 계정 정보 조회
func (n *Netflix) Info(w http.ResponseWriter, r *http.Request) {
	log.Println("[/netflix/info] Netflix.Info")

    // 요청으로부터 id, pw
    var account Account
    json.NewDecoder(r.Body).Decode(&account)

    ctx, cancel := chromedp.NewContext(
        context.Background(),
        //chromedp.WithDebugf(log.Printf),
    )
    defer cancel()
    ctx, cancel = context.WithTimeout(ctx, 1 * time.Minute)
    defer cancel()

    n.ctx = ctx

    var res string

    if msg, err := n.Login(account); err != nil {
        log.Fatal(err)
    } else if msg != "" {
        fmt.Fprintln(w, msg)
        return
    }
    defer func() {
        if err := n.Logout(); err != nil {
            log.Fatal(err)
        }
    }()

    // 계정 정보 조회
    if err := chromedp.Run(
        ctx,
        chromedp.Navigate(URI_INFO),
        chromedp.WaitVisible(SEL_INFO_MEMBERSHIP),
        chromedp.Text(SEL_INFO_MEMBERSHIP, &res, chromedp.NodeVisible),
    ); err != nil {
        log.Fatal(err)
    }

    switch strings.Split(res, "\n")[0] {
    case "스트리밍 멤버십에 가입하지 않으셨습니다.":
        account.Membership = MEMBERSHIP_NO
    default:
        fmt.Println(strings.Split(res, "\n")[0])
        account.Membership = MEMBERSHIP_PREMIUM
    }

    fmt.Fprintln(w, account)
}

// 패턴에 맞게 메소드 연결
func (n Netflix) Handler() {
    http.HandleFunc("/netflix", n.Hello)
    http.HandleFunc("/netflix/info", n.Info)
}
