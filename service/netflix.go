// Netflix API 구현

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// Netflix 웹사이트 주소 상수
const (
	URI_HOME   = `https://www.netflix.com/kr`
	URI_LOGIN  = URI_HOME + `/login`
	URI_LOGOUT = URI_HOME + `/signout`
	URI_INFO   = URI_HOME + `/youraccount`
)

// Netflix 웹사이트 선택자 상수
const (
	// 로그인
	SEL_LOGIN_ID  = `input[data-uia="login-field"]`
	SEL_LOGIN_PW  = `input[data-uia="password-field"]`
	SEL_LOGIN_BTN = `button[data-uia="login-submit-button"]`
	SEL_LOGIN_ERR = `div[data-uia="error-message-container"]`

	// 계정 정보 조회
	SEL_INFO_PAYMENT_TYPE = `div[data-uia="wallet-mop"]`
	SEL_INFO_PAYMENT_NEXT = `div[data-uia="nextBillingDate-item"]`
	SEL_INFO_MEMBERSHIP   = `div[data-uia="plan-section"] > section`
)

// Netflix 멤버십 상수
const (
	// 멤버십 종류
	MEMBERSHIP_NO = iota
	MEMBERSHIP_BASIC
	MEMBERSHIP_STANDARD
	MEMBERSHIP_PREMIUM

	// 멤버십 종류별 가격
	MEMBERSHIP_COST_NO       = 0
	MEMBERSHIP_COST_BASIC    = 9_500
	MEMBERSHIP_COST_STANDARD = 13_500
	MEMBERSHIP_COST_PREMIUM  = 17_000
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

		chromedp.Sleep(1*time.Second),
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

	// 웹 접속을 위한 컨텐스트 선언
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// 1분 타임아웃 설정
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	n.ctx = ctx

	// 로그인
	if msg, err := n.Login(account); err != nil {
		log.Fatal(err)
	} else if msg != "" {
		fmt.Fprintln(w, msg)
		return
	}

	// 로그아웃
	defer func() {
		if err := n.Logout(); err != nil {
			log.Fatal(err)
		}
	}()

	// 가공 전 데이터를 저장하는 변수
	var rawPayment, rawDate, rawMembership string

	// 계정 정보 조회
	if err := chromedp.Run(
		ctx,

		chromedp.Navigate(URI_INFO),

		chromedp.WaitVisible(SEL_INFO_PAYMENT_TYPE),
		chromedp.Text(SEL_INFO_PAYMENT_TYPE, &rawPayment, chromedp.NodeVisible),

		chromedp.WaitVisible(SEL_INFO_PAYMENT_NEXT),
		chromedp.Text(SEL_INFO_PAYMENT_NEXT, &rawDate, chromedp.NodeVisible),

		chromedp.WaitVisible(SEL_INFO_MEMBERSHIP),
		chromedp.Text(SEL_INFO_MEMBERSHIP, &rawMembership, chromedp.NodeVisible),
	); err != nil {
		log.Fatal(err)
	}

	// 가공 과정에서 사용하는 변수
	var (
		payment          []string
		year, month, day int
		dummy            string
	)
	payment = strings.Split(rawPayment, "\n")
	fmt.Sscanf(rawDate, "%s %s %d%s %d%s %d%s", &dummy, &dummy, &year, &dummy, &month, &dummy, &day, &dummy)

	account.Payment.Type = payment[0]
	account.Payment.Detail = payment[1]
	account.Payment.Next = fmt.Sprintf("%d-%d-%d", year, month, day)

	// 멤버십 타입에 따라 동작
	switch strings.Split(rawMembership, "\n")[0] {
	case "스트리밍 멤버십에 가입하지 않으셨습니다.":
		account.Membership.Type = MEMBERSHIP_NO
		account.Membership.Cost = MEMBERSHIP_COST_NO
	case "베이식":
		account.Membership.Type = MEMBERSHIP_BASIC
		account.Membership.Cost = MEMBERSHIP_COST_BASIC
	case "스탠다드":
		account.Membership.Type = MEMBERSHIP_STANDARD
		account.Membership.Cost = MEMBERSHIP_COST_STANDARD
	case "프리미엄":
		account.Membership.Type = MEMBERSHIP_PREMIUM
		account.Membership.Cost = MEMBERSHIP_COST_PREMIUM
	}

	// Account 구조체를 json 형식으로 변환
	res, err := json.Marshal(account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, string(res))
}

// 패턴에 맞게 메소드 연결
func (n Netflix) Handler() {
	http.HandleFunc("/netflix", n.Hello)
	http.HandleFunc("/netflix/info", n.Info)
}
