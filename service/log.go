// 로그 구분을 위한 변수 선언

package service

import (
	"log"
	"os"
)

var (
	// 접속 기록 저장을 위한 구조체 변수
	LogStdout = log.New(os.Stdout, "", log.LstdFlags)

	// 에러 내용 저장을 위한 구조체 변수
	LogStderr = log.New(os.Stderr, "", log.LstdFlags)
)
