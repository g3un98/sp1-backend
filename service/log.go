// 로그 구분을 위한 변수 선언

package service

import (
	"log"
	"os"
)

var (
	// 정보 저장을 위한 구조체 변수
	LogInfo = log.New(os.Stdout, "", log.LstdFlags)

	// 에러 저장을 위한 구조체 변수
	LogErr = log.New(os.Stderr, "", log.LstdFlags)
)
