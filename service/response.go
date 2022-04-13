// 응답을 처리하기 위한 함수 정의

package service

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, statusCode int, resp map[string]interface{}) (err error) {
	w.WriteHeader(statusCode)
	w.Header().Set("Context-type", "application/json")
	if resp["message"] == "" {
		resp["message"] = http.StatusText(statusCode)
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("An error has occurred while JSON Marshal: %s", err)
	}
	if _, err = w.Write(jsonResp); err != nil {
		return fmt.Errorf("An error has occurred while respond: %s", err)
	}
	return
}
