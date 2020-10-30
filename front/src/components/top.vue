<template>
  <div id='top' ref="message">
    <div id="message">
      <h1>Put your card over the reader<br/></h1>
      <router-link to="/forgot" id='forgetPos'>Forgot card?</router-link>
    </div>
    <h4 id="roomPos" style="position:fixed;bottom:0;">{{roomname}}</h4>
  </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'top',
  data: function () {
    return {
      closeflg: false,
      roomname: ''
    }
  },
  destroyed: function () {
    clearTimeout(this.reconnecttimer)
    clearTimeout(this.clocktimer)
    this.closeflg = true
    this.ws.close()
  },
  mounted: function () {
    this.roomname = util.roomName()
    this.connectCardReader()
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
        if (message['IsCard'] === true) {
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
  font-size: 4rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
}

#forgetPos{
  font-size: 45px;
  color: #3282ce;
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

#roomPos {
  padding: 24px 48px;
  font-size: 1.5rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
  color: #8b8b8b;
}
</style>
