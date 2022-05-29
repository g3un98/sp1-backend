# Documents

## Index

### User

| Method | Path                                             | Description            |
|--------|--------------------------------------------------|------------------------|
| POST   | [/user](#post-user)                              | 앱 계정 생성           |
| DELETE | [/user](#delete-user)                            | 앱 계정 삭제           |
| PUT    | [/user](#put-user)                               | 앱 계정 수정           |
| POST   | [/login](#post-login)                            | 로그인                 |

### Group

| Method | Path                                             | Description            |
|--------|--------------------------------------------------|------------------------|
| GET    | [/group/:groupId](#get-groupgroupid)             | 앱 그룹 조회           |
| POST   | [/group](#post-group)                            | 앱 그룹 생성           |
| DELETE | [/group/:groupId](#delete-groupgroupid)          | 앱 그룹 삭제           |
| PUT    | [/group/:groupId](#put-groupgroupid)             | 앱 그룹 수정           |

### Netflix

| Method | Path                                             | Description            |
|--------|--------------------------------------------------|------------------------|
| POST   | [/netflix](#post-netflix)                        | Netflix 계정 정보 조회 |
| DELETE | [/netflix/membership](#delete-netflixmembership) | Netflix 구독 해지      |

### Wavve

| Method | Path                                             | Description            |
|--------|--------------------------------------------------|------------------------|
| POST   | [/wavve](#post-wave)                             | Wavve 계정 정보 조회   |

## POST /user

### Request

```json
{
    "app_id": "계정 아이디",
    "app_pw": "계정 비밀번호"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 201 Created               | 계정 생성 성공                                                     |
| 400 Bad Request           | app_id 혹은 app_pw가 유효하지 않음                                 |
| 401 Unauthorized          | 이미 존재하는 app_id 사용                                          |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## DELETE /user

### Request

```json
{
    "app_id": "계정 아이디",
    "app_pw": "계정 비밀번호"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 계정 삭제 성공                                                     |
| 400 Bad Request           | app_id 혹은 app_pw가 유효하지 않음                                 |
| 401 Unauthorized          | 존재하지 않은 app_id 입력                                          |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## PUT /user

### Request

```json
{
    "app_id": "계정 아이디",
    "app_pw": "계정 비밀번호",
    "app_email": "계정 이메일"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 계정 수정 성공                                                     |
| 400 Bad Request           | app_id 혹은 app_pw가 유효하지 않음                                 |
| 401 Unauthorized          | 존재하지 않은 app_id 입력                                          |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## POST /login

### Request

```json
{
    "app_id": "계정 아이디",
    "app_pw": "계정 비밀번호"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 로그인 성공                                                        |
| 400 Bad Request           | app_id 혹은 app_pw가 유효하지 않음                                 |
| 404 Not Found             | 존재하지 않은 app_id 혹은 app_pw 입력                              |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## GET /group/:groupId

### Request

groudId는 12자리 문자열

### Response

```json
{
    "group_id": "12자리 문자열",
    "ott": "OTT 서비스명",
    "account": {
        "id": "OTT 계정 아이디",
        "pw": "OTT 계정 비밀번호",
        "payment": {
            "type": "결제 수단",
            "detail": "결제 수단 정보",
            "next": "다음 결제일"
        },
        "membership": {
            "type": "멤버십 타입 상수",
            "cost": "멤버십 가격"
        }
    },
    "updatetime": "마지막 수정 시간 (Unix time)",
    "members": [
        { "app_id": "계정 아이디", "is_admin": "0: 그룹원 || 1: 그룹장" },
        ...
    ]
}
```

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 그룹 정보 조회                                                     |
| 404 Not Found             | 존재하지 않은 group_id 입력                                        |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## POST /group

전달받은 인자를 기반으로 그룹 생성.
이미 존재하는 그룹이면 해당 그룹에 참가.

### Request

```json
{
    "app_id": "계정 아이디",
    "ott": "OTT 서비스명",
    "ott_id": "OTT 계정 아이디",
    "ott_pw": "OTT 계정 비밀번호"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 그룹 생성 성공                                                     |
| 400 Bad Request           | 유효하지 않은 데이터 입력                                          |
| 401 Unauthorized          | 이미 참가한 그룹                                                   |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## DELETE /group/:groupId

### Request

```json
{
    "app_id": "그룹장 아이디"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 그룹 삭제 성공                                                     |
| 400 Bad Request           | 유효하지 않은 app_id 입력                                          |
| 404 Not Found             | 유효하지 않은 grupdId 혹은 그룹장이 아닌 app_id 입력               |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## PUT /group/:groupId

### Request

```json
{
    "ott_pw": "OTT 계정 비밀번호",
    "payment": {
        "type": "결제 수단",
        "detail": "결제 수단 정보",
        "next": "다음 결제일"
    },
    "membership": {
        "type": "멤버십 타입 상수",
        "cost": "멤버십 가격"
    }
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 그룹 수정 성공                                                     |
| 400 Bad Request           | 유효하지 않은 데이터 입력                                          |
| 404 Not Found             | 유효하지 않은 grupdId 입력                                         |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |


## POST /netflix

### Request

```json
{
    "id": "Netflix 계정 아이디",
    "pw": "Netflix 계정 비밀번호"
}
```

### Response

```json
{
    "account": {
        "id": "Netflix 계정 아이디",
        "pw": "Netflix 계정 비밀번호",
        "payment": {
            "type": "결제 수단",
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

## DELETE /netflix/membership

### Request

```json
{
    "id": "Netflix 계정 아이디",
    "pw": "Netflix 계정 비밀번호"
}
```

### Response

상태 메시지 혹은 오류 메시지 반환

### Status code

| Status                    | Note                                                               |
|---------------------------|--------------------------------------------------------------------|
| 200 OK                    | 구독 해지 성공                                                     |
| 400 Bad Request           | id 혹은 pw가 유효하지 않음                                         |
| 401 Unauthorized          | id 혹은 pw가 틀림, 오류 메시지가 함께 반환                         |
| 405 Method Not Allowed    | 유효하지 않은 메소드 호출                                          |
| 500 Internal Server Error | 타임 아웃 혹은 예상하지 못한 오류 발생 시, 오류 메시지와 함께 반환 |

## POST /wavve

### Request

```json
{
    "id": "Wavve 계정 아이디",
    "pw": "Wavve 계정 비밀번호"
}
```

### Response

```json
{
    "account": {
        "id": "Wavve 계정 아이디",
        "pw": "Wavve 계정 비밀번호",
        "payment": {
            "type": "결제 수단",
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
