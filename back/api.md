# API JSON仕様
## UserInfo
### Create User
Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
Response
```json
{
    "Success": bool,
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
### Update User
Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
Response
```json
{
    "Success": bool,
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```

### Read User
Request
```json
{
    "CardID": string,
    "timestamp": timestamp,
}
```
Response
```json
{
    "SID": string,
    "CardID": string,
    "IsEnter": bool,
    "UserName": string,
    "timestamp": timestamp,
}
```

### Add Log
Requesr
```json
{
    "SID": string,
    "IsEnter": bool,
    "timestamp": timestamp,
}
```
Response
```json
{
    "SID": string,
    "timestamp": timestamp,
}
```

## Websocket Server -> Front
```json
{
    "IsCard": bool,
    "CardID": string,
    "IsEnter": bool,
    "IsNew": bool,
    "timestamp": timestamp
}
```

# API URI

## UserInfo
### create
URI : `/api/createuser/`  
method : `POST`  

### read
URI : `/api/readuser`
method : `GET`

### Update
URI : `/api/updateuser`
method : `UPDATE`

## WebSocket
URI : `/socket/readCard`
