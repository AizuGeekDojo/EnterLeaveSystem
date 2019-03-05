<template>
  <div id='top' ref="message">
    <h2>{{clocktext}}</h2>
    <div id="message">
      <h1>Please hold the card over the reader</h1>
    </div>
  </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'top',
  data: function () {
    return {
      closeflg: false,
      clocktext: '----/--/--  --:--:--'
    }
  },
  destroyed: function () {
    clearTimeout(this.reconnecttimer)
    clearTimeout(this.clocktimer)
    this.closeflg = true
    this.ws.close()
  },
  mounted: function () {
    this.connectCardReader()
    const self = this
    setInterval(() => {
      const da = new Date()
      const year = da.getFullYear()
      const month = da.getMonth() + 1
      const date = da.getDate()
      const hour = da.getHours()
      const minute = da.getMinutes()
      const second = da.getSeconds()

      self.clocktext = `${year}/`
      if (month < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${month}/`
      if (date < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${date}  `
      if (hour < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${hour}:`
      if (minute < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${minute}:`
      if (second < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${second}`
    }, 1000)
  },
  methods: {
    connectCardReader: function () {
      const self = this
      this.ws = new WebSocket('ws://localhost:3000/socket/readCard')
      this.ws.onopen = function (e) {
        console.log('Card reader standby')
      }
      this.ws.onmessage = function (e) {
        var message = JSON.parse(e.data)
        console.log('Read card data:', message)
        if (message['IsNew'] === false) {
          util.getUserInfo(message['SID'])
            .then(res => {
              if (res['IsEnter']) {
                self.$router.push({ name: 'question', params: { userinfo: res } })
              } else {
                self.$router.push({ name: 'welcome', params: { userinfo: res } })
              }
            })
        } else {
          self.$router.push({ name: 'regist', params: { cardid: message['CardID'] } })
        }
      }
      this.ws.onerror = function (e) {
        console.log('Card reader communication error', e)
      }
      this.ws.onclose = function (e) {
        console.log('Card reader stopped')
        if (!self.closeflg) {
          self.reconnecttimer = setTimeout(() => {
            self.connectCardReader()
          }, 3000)
        }
      }
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
h1{
  font-size: 72px;
}
h2{
  font-size: 50px;
}
div #message {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  width: 80%;
  height: 200px;
}
div #top {
  text-align: center;
}
</style>
