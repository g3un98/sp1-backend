// Netflix 선택자 선언

package netflix

// 로그인 과정에서 사용하는 선택자
const (
    SEL_LOGIN_ID = `input[id="id_userLoginId"]`
    SEL_LOGIN_PW = `input[id="id_password"]`
    SEL_LOGIN_BTN = `button[class="btn login-button btn-submit btn-small"]`
    SEL_LOGIN_ERR = `div[class="ui-message-container ui-message-error"]`
)

// 계정 정보를 가져오는 과정에서 사용하는 선택자
const (
    SEL_INFO_MEMBERSHIP = `div[data-uia="plan-section"] > section`
)
