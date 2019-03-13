# Enter Leave System API
## Web Server
### GET /
Provide front page and etc (like *.js, *.css).

## API Server
### WebSocket /socket/readCard
Provide card reader connection.
Response
``` json
{
  "IsCard":true,
  "CardID":"Detect card ID",
  "SID":"User ID associated with the card ID",
  "IsNew":false
}
```
If the card is not registered, return info with SID="" and IsNew=true.

### GET /api/user
Get user info by user ID.
Request
http://localhost:3000/api/user?sid={{UserID}}  

Response
``` json
{
  "SID":"UserID",
  "UserName":"UserName",
  "IsEnter":true
}
```
When users enter, IsEnter=true.

### POST /api/user
Register card info with given user ID.
Request
``` json
{
  "SID":"UserID",
  "CardID":"CardID"
}
```
Response
``` json
{
  "Success":true,
  "Reason":""
}
```

### POST /api/log
Add enter/leave log.
Request
``` json
{
  "SID":"UserID",
  "IsEnter":0,
  "Ext":"{...}"
}
```
If the log contains answer of question, the data is given as json string in Ext.