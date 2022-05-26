// Wavve API 구현

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// Wavve 구조체는 Servicer 인터페이스를 만족
var _ Servicer = (*Wavve)(nil)

// Wavve 구조체 선언
type Wavve struct {
	Service
}

// Wavve 상수
const (
	// URI
	WAVVE_URI_HOME  = `https://www.wavve.com`
	WAVVE_URI_LOGIN = WAVVE_URI_HOME + `/login`
	WAVVE_URI_INFO  = WAVVE_URI_HOME + `/my/subscription_ticket`

	// 로그인/로그아웃 선택자
	WAVVE_SEL_LOGIN_ID     = `input[title="아이디"]`
	WAVVE_SEL_LOGIN_PW     = `input[title="비밀번호"]`
	WAVVE_SEL_LOGIN_SUBMIT = `a[title="로그인"]`
	WAVVE_SEL_LOGIN_ERR    = `p[class="login-error-red"]`
	WAVVE_SEL_LOGOUT       = `#app > div.body > div:nth-child(2) > header > div:nth-child(1) > div.header-nav > div > ul > li.over-parent-1depth > div > ul > li:nth-child(4) > button`

	// 계정 정보 선택자
	WAVVE_SEL_INFO_PAYMENT_TYPE    = `#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(6) > span > span`
	WAVVE_SEL_INFO_PAYMENT_NEXT    = `#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(5)`
	WAVVE_SEL_INFO_MEMBERSHIP_TYPE = `#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(2) > div > p.my-pay-tit > span:nth-child(3)`
	WAVVE_SEL_INFO_MEMBERSHIP_COST = `#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(4)`
)

// Wavve 서비스명 반환
func (wv Wavve) GetName() (name string) {
	return "Wavve"
}

// Wavve 웹사이트 로그인
func (wv *Wavve) Login(a Account) (msg string, err error) {
	var url string

	if err = chromedp.Run(
		wv.ctx,

		chromedp.Navigate(WAVVE_URI_LOGIN),

		chromedp.WaitVisible(WAVVE_SEL_LOGIN_ID),
		chromedp.Click(WAVVE_SEL_LOGIN_ID, chromedp.NodeVisible),
		chromedp.SendKeys(WAVVE_SEL_LOGIN_ID, a.Id, chromedp.NodeVisible),

		chromedp.WaitVisible(WAVVE_SEL_LOGIN_PW),
		chromedp.Click(WAVVE_SEL_LOGIN_PW, chromedp.NodeVisible),
		chromedp.SendKeys(WAVVE_SEL_LOGIN_PW, a.Pw, chromedp.NodeVisible),

		chromedp.WaitVisible(WAVVE_SEL_LOGIN_SUBMIT),
		chromedp.Click(WAVVE_SEL_LOGIN_SUBMIT, chromedp.NodeVisible),

		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return "", fmt.Errorf("An error has occurred while login to Wavve: %s\n", err)
	}

	if url == WAVVE_URI_LOGIN {
		if err = chromedp.Run(
			wv.ctx,
			chromedp.Text(WAVVE_SEL_LOGIN_ERR, &msg, chromedp.NodeVisible),
		); err != nil {
			return "", fmt.Errorf("An error has occurred while load error message from web: %s", err)
		}
		return msg, nil
	}
	return
}

// Wavve 웹사이트 로그아웃
func (n *Wavve) Logout() (err error) {
	if err = chromedp.Run(
		n.ctx,

		chromedp.Navigate(WAVVE_URI_HOME),

		chromedp.WaitVisible(WAVVE_SEL_LOGOUT),
		chromedp.Click(WAVVE_SEL_LOGOUT, chromedp.NodeVisible),
	); err != nil {
		return fmt.Errorf("An error has occurred while logout to Wavve: %s", err)
	}
	return
}

// Wavve 계정 정보 조회
func (wv *Wavve) Info(w http.ResponseWriter, r *http.Request) {
	LogInfo.Println("[/wavve/info] Wavve.Info")

	// 요청-응답 처리 과정에서 사용할 변수
	var (
		resp = make(map[string]interface{})
		err  error
	)

	// POST 메소드가 아닌 요청은 405 에러 반환
	if r.Method != "POST" {
		if err := Response(w, http.StatusMethodNotAllowed, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	// TODO
	// 타임아웃 발생시, 200이 아닌 다른 코드 반환
	// 웹사이트 접속을 위한 컨텐스트 선언 및 1분 타임아웃 설정
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		//chromedp.WithDebugf(LogInfo.Printf),
	)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	wv.ctx = ctx

	var (
		// 유저 정보를 담을 구조체 변수
		account Account

		// 가공 전 데이터를 담을 변수
		rawPaymentType, rawPaymentNext, rawMembershipType, rawMembershipCost string

		dummy string
	)

	if err = json.NewDecoder(r.Body).Decode(&account); err != nil {
		LogErr.Printf("An error has occurred while decode json from request: %s\n", err)
		if err = Response(w, http.StatusInternalServerError, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}
	defer r.Body.Close()

	// id와 pw 길이 검사
	if len(account.Id) < 1 || len(account.Pw) < 1 {
		if err = Response(w, http.StatusBadRequest, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	// 로그인
	if msg, err := wv.Login(account); err != nil {
		LogErr.Println(err)
		return
	} else if msg != "" {
		resp["message"] = msg
		if err = Response(w, http.StatusUnauthorized, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	// 로그아웃
	defer func() {
		if err := wv.Logout(); err != nil {
			LogErr.Println(err)
		}
		return
	}()

	// 계정 정보 조회
	if err := chromedp.Run(
		ctx,

		chromedp.Navigate(WAVVE_URI_INFO),

		chromedp.Text(`#contents`, &dummy, chromedp.NodeVisible),
	); err != nil {
		LogErr.Printf("An error has occurred while load account infomation: %s\n", err)
		if err = Response(w, http.StatusInternalServerError, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	if dummy == "이용권 결제 내역이 없어요." {
		account.Payment = Payment{}
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_NO
		account.Membership.Cost = MEMBERSHIP_WAVVE_COST_NO

		resp["account"] = account
		if err = Response(w, http.StatusOK, resp); err != nil {
			LogErr.Println(err)
		}

		return
	}

	if err := chromedp.Run(
		ctx,

		chromedp.Navigate(WAVVE_URI_INFO),

		chromedp.WaitVisible(WAVVE_SEL_INFO_PAYMENT_TYPE),
		chromedp.Text(WAVVE_SEL_INFO_PAYMENT_TYPE, &rawPaymentType, chromedp.NodeVisible),
		chromedp.WaitVisible(WAVVE_SEL_INFO_PAYMENT_NEXT),
		chromedp.Text(WAVVE_SEL_INFO_PAYMENT_NEXT, &rawPaymentNext, chromedp.NodeVisible),

		chromedp.WaitVisible(WAVVE_SEL_INFO_MEMBERSHIP_TYPE),
		chromedp.Text(WAVVE_SEL_INFO_MEMBERSHIP_TYPE, &rawMembershipType, chromedp.NodeVisible),
		chromedp.WaitVisible(WAVVE_SEL_INFO_MEMBERSHIP_COST),
		chromedp.Text(WAVVE_SEL_INFO_MEMBERSHIP_COST, &rawMembershipCost, chromedp.NodeVisible),
	); err != nil {
		LogErr.Printf("An error has occurred while load account infomation: %s\n", err)
		if err = Response(w, http.StatusInternalServerError, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	account.Payment = Payment{
		Type: rawPaymentType,
		Next: strings.Split(rawPaymentNext, " ")[0],
	}

	if _, err = fmt.Sscanf(rawMembershipCost, "%d%s", &account.Membership.Cost, &dummy); err != nil {
		LogErr.Printf("An error has occurred while convert datae: %s\n", err)
		if err = Response(w, http.StatusInternalServerError, resp); err != nil {
			LogErr.Println(err)
		}
		return
	}

	switch rawMembershipType {
	case "Basic":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BASIC
	case "Standard":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_STANDARD
	case "Premium":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_PREMIUM
	case "Basic X FLO 무제한":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_FLO
	case "Basic X Bugs 듣기":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BUGS
	case "Basic X KB 나라사랑카드":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_KB
	}

	resp["account"] = account
	if err = Response(w, http.StatusOK, resp); err != nil {
		LogErr.Println(err)
	}
}

// 패턴에 맞게 메소드 연결
func (wv Wavve) Handler() {
	http.HandleFunc("/wavve/info", wv.Info)
}
