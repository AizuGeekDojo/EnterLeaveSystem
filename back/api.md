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
    "UserName": string,
    "timestamp": timestamp,
}
```

### Add Log
Request
```json
{
    "SID": string,
    "Category": string,
    "Ext": string,
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
