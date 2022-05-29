# Documents

## Index

| Method | Path                | Parameters           | Description            |
|--------|---------------------|----------------------|------------------------|
| POST   | /netflix/info       | id:string, pw:strnig | Netflix 계정 정보 조회 |
| DELETE | /netflix/membership | id:string, pw:strnig | Netflix 구독 해지      |
| POST   | /wavve/info         | id:string, pw:strnig | Wavve 계정 정보 조회   |

## POST /netlifx/info

### Request

```json
{
    "id": "계정 아이디",
    "pw": "계정 비밀번호"
}
```

### Response

```json
{
    "account": {
        "id": "계정 아이디",
        "pw": "계정 비밀번호",
        "payment": {
            "type": "결제 수단 방식",
            "detail": "결제 수단 정보",
            "next": "다음 결제일"
        },
        "membership": {
            "type": "멤버십 타입 상수",
            "cost": "멤버십 가격"
        }
    }
}
```

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 계정 정보 반환                                                     |
| 400 Bad Request           | id 혹은 pw가 유효하지 않음                                         |
| 401 Unauthorized          | id 혹은 pw가 틀림, 오류 메시지가 함께 반환                         |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## DELETE /netlifx/membership

### Request

```json
{
    "id": "계정 아이디",
    "pw": "계정 비밀번호"
}
```

### Response

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 구독 해지 성공                                                     |
| 400 Bad Request           | id 혹은 pw가 유효하지 않음                                         |
| 401 Unauthorized          | id 혹은 pw가 틀림, 오류 메시지가 함께 반환                         |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## POST /wavve/info

### Request

```json
{
    "id": "계정 아이디",
    "pw": "계정 비밀번호"
}
```

### Response

```json
{
    "account": {
        "id": "계정 아이디",
        "pw": "계정 비밀번호",
        "payment": {
            "type": "결제 수단 방식",
            "detail": "결제 수단 정보",
            "next": "다음 결제일"
        },
        "membership": {
            "type": "멤버십 타입 상수",
            "cost": "멤버십 가격"
        }
    }
}
```


### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 계정 정보 반환                                                     |
| 400 Bad Request           | id 혹은 pw가 유효하지 않음                                         |
| 401 Unauthorized          | id 혹은 pw가 틀림, 오류 메시지가 함께 반환                         |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |
