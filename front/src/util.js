export default {
  getUserInfo (sid) {
    return fetch(`http://localhost:3000/api/user?sid=${sid}`, {
      mode: 'cors',
      method: 'GET',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      }
    })
      .then(response => {
        if (!response.ok) {
          throw response
        }
        if (response.status === 204) {
          return {
            'SID': sid,
            'UserName': '',
            'IsEnter': true
          }
        }
        return response.json()
      })
      .catch(response => {
        console.error(response)
      })
  },
  registCardInfo (cardid, sid) {
    return fetch('http://localhost:3000/api/user', {
      mode: 'cors',
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        SID: sid,
        CardID: cardid
      })
    })
      .then(response => {
        return response.json()
      })
      .catch(function (error) {
        console.error(error)
      })
  },
  addLog (sid, isenter, ext) {
    console.log({
      SID: sid,
      IsEnter: (isenter ? 1 : 0),
      Ext: ext
    })
    return fetch('http://localhost:3000/api/log', {
      mode: 'cors',
      method: 'POST',
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        SID: sid,
        IsEnter: (isenter ? 1 : 0),
        Ext: ext
      })
    })
      .then(response => {
        if (!response.ok) {
          throw response
        }
        return response
      })
      .catch(function (error) {
        console.error(error)
        return null
      })
  },
  roomName () {
    var roomname = process.env.VUE_APP_ROOMNAME
    if (roomname === undefined) {
      roomname = 'University of Aizu'
    }
    return roomname
  },
  isShowQuestion () {
    var isshow = process.env.VUE_APP_SHOWQUESTION
    if (isshow === undefined) {
      isshow = 'false'
    }
    return isshow === 'true'
  }
}
