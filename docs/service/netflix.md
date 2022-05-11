# Netflix

## Index

| Method | Target                 | Body                   |
|--------|------------------------|------------------------|
| POST   | [/netflix/info](#info) | id: string, pw: string |

## Info

계정 정보를 받아오는 API.

### 요청

```json
{
    "id": "계정 이메일 혹은 전화번호",
    "pw": "계정 비밀번호"
}
```

### 응답

```json
{
    "account": {
        "id": "계정 이메일 혹은 전화번호",
        "pw": "계정 비밀번호",
        "payment": {
            "type": "결제 수단 방식",
            "detail": "결제 수단 정보",
            "next": "다음 결제일"
        },
        "membership": {
            "type": 멤버십 타입 상수,
            "cost": 멤버십 가격 
        }
    }
}
```

### 응답 코드

| Status                    | Note                                            |
|---------------------------|-------------------------------------------------|
| 200 OK                    | id와 pw가 맞을 시, 계정 정보와 함께 반환        |
| 400 Bad Request           | id 혹은 pw가 유효하지 않을 시, 반환             |
| 401 Unauthorized          | id 혹은 pw가 틀릴 시, 오류 메시지와 함께 반환   |
| 405 Method Not Allowed    | POST를 제외한 메소드 요청 시, 반환              |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 반환 |
