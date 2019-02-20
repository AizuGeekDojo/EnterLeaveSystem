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
    "SID": string,
    "timestamp": timestamp,
}
```
Response
```json
{
    "SID": string,
    "IsEnter": bool,
    "UserName": string,
    "timestamp": timestamp,
}
```

### Check User
Request
```json
{
    "SID": string,
    "timestamp": timestamp,
}
```
Response
```json
{
    "isValid": string,
    "SID": string,
    "timestamp": timestamp,
}
```

### Add Log
Requesr
```json
{
    "SID": string,
    "IsEnter": bool,
    "Ext": {
        "Use" :"Purpose of using Aizu Geek Dojo",
        "message" :"impression",
    },
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
    "SID": string,
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

### Check
URI : `/api/checkuser`
method : `GET`

## WebSocket
URI : `/socket/readCard`
