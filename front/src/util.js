export default {
  getUserInfo (sid) {
    return fetch(`http://localhost:3000/api/user?sid=${sid}`, {
      mode: "cors",
      method: "GET",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json"
      }//,
      // body: JSON.stringify({
      //   SID: sid
      // })
    })
      .then(response => {
        if (!response.ok) {
          throw response
        }
        return response.json();
      })
      .catch(function(error) {
        console.error(error);
      });
  },
  registCardInfo (cardid, sid) {
    return fetch("http://localhost:3000/api/user", {
      mode: "cors",
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        SID: sid,
        CardID: cardid
      })
    })
      .then(response => {
        if (!response.ok) {
          throw response
        }
        return response.json();
      })
      .catch(function(error) {
        console.error(error);
      });
  },
  addLog (sid, isenter, ext) {
    console.log({
      SID: sid,
      IsEnter: (isenter ? 1 : 0),
      Ext: ext
    })
    return fetch("http://localhost:3000/api/log", {
      mode: "cors",
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json"
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
        return response;
      })
      .catch(function(error) {
        console.error(error);
      });
  }
}