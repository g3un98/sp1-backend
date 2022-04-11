// Netflix API 유닛 테스트 구현

package main_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/g3un/sp1-backend/service"
)

var n = service.NewService("Netflix").(*service.Netflix)

type TestSet struct {
    req Request
    res Response
}

type Request struct {
    method string
    target string
    body io.Reader
}

type Response struct {
    statusCode int
    status string
}

var testSets = []TestSet{
    // POST를 제외한 메소드 요청 시, 상태 코드 405 반환
    {
        Request{http.MethodGet, "/netflix/info", nil},
        Response{405, "405 Method Not Allowed"},
    },
    {
        Request{http.MethodPut, "/netflix/info", nil},
        Response{405, "405 Method Not Allowed"},
    },
    {
        Request{http.MethodDelete, "/netflix/info", nil},
        Response{405, "405 Method Not Allowed"},
    },

    // id 혹은 pw를 입력하지 않았을 시, 상태 코드 400 반환
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{}`)},
        Response{400, "400 Bad Request"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "jujujujusttetetetest@gmail.com" }`)},
        Response{400, "400 Bad Request"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "pw": "jujujujusttetetetest1!" }`)},
        Response{400, "400 Bad Request"},
    },

    // id 혹은 pw의 길이 검사 실패 시, 상태 코드 400 반환
    // 5 <= id <= 50
    // 4 <= pw <= 60
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "1234", "pw": "jujujujusttetetetest1!" }`)},
        Response{400, "400 Bad Request"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "123456789a123456789b123456789c123456789d123456789e1", "pw": "jujujujusttetetetest1!" }`)},
        Response{400, "400 Bad Request"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "jujujujusttetetetest@gmail.com", "pw": "123" }`)},
        Response{400, "400 Bad Request"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "jujujujusttetetetest@gmail.com", "pw": "123456789a123456789b123456789c123456789d123456789e123456789f1" }`)},
        Response{400, "400 Bad Request"},
    },

    // id 혹은 pw가 틀릴 시,
    // 넷플릭스 사이트 오류 발생 시, 상태 코드 401 반환
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "jujujujusttetetetest@gmail.com", "pw": "1234" }`)},
        Response{401, "401 Unauthorized"},
    },
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "1jujujujusttetetetest@gmail.com", "pw": "jujujujusttetetetest1!" }`)},
        Response{401, "401 Unauthorized"},
    },

    // id와 pw가 맞을 시, 상태 코드 200 반환
    {
        Request{http.MethodPost, "/netflix/info", strings.NewReader(`{ "id": "jujujujusttetetetest@gmail.com", "pw": "jujujujusttetetetest1!" }`)},
        Response{200, "200 OK"},
    },
}

func TestInfo(t *testing.T) {
    for _, tt := range testSets {
        w, req := httptest.NewRecorder(), httptest.NewRequest(tt.req.method, tt.req.target, tt.req.body)
        n.Info(w, req)
        res := w.Result()
        defer res.Body.Close()
        data := Response{res.StatusCode, res.Status}

        if data != tt.res {
            // 넷플릭스 사이트 에러 예외 처리
            // 에러 발생 시, 5분 정도 기다린 후 다시 테스트
            msg, _ := ioutil.ReadAll(res.Body)
            if string(msg) == `{"message":"현재 기술적인 문제가 있어 수정 작업 중에 있습니다. 잠시 후 다시 시도해 주시기 바랍니다."}` {
                t.Errorf("netflix.com: We are having technical difficulties and are actively working on a fix. Please try again in a few minutes.")
                return
            }

            t.Errorf("expected %v, got %v", tt.res, data)
        }
    }
}
