<template>
  <div id='top' ref="message" class="container align-middle">
    <h1 class="contents align-middle">Please hold the card over the reader</h1>
  </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'top',
  destroyed: function () {
    this.ws.close()
  },
  mounted: function () {
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
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
h1,
h2 {
  font-size: 72px;
  display: table-cell;
  height: 100%;
  width: 100%;
  font-weight: normal;
  text-align: center;
  vertical-align: middle;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
div #top {
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
