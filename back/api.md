UserInfo
Create Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
Update Request
```json
{
    "SID": string,
    "CardID": string,
    "timestamp": timestamp,
}
```
Read Response
```json
{
    "SID": string,
    "CardID": string,
    "UserName": string,
    "timestamp": timestamp,
}
```

Websocket Server -> Front
```json
{
    "IsCard": bool,
    "CardID": string,
    "IsNew": bool,
    "timestamp": timestamp
}
```