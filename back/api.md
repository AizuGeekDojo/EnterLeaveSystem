# API JSON仕様
## UserInfo
### Create Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
### Update Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
### Read Response
```json
{
    "SID": string,
    "CardID": string,
    "UserName": string,
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


