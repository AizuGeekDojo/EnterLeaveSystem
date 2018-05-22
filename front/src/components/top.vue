<template>
  <div id='top' ref="message" class="container align-middle">
    <h1 class="contents align-middle">{{ msg }}</h1>
  </div>
</template>

<script>
import router from '../router'
export default {
  name: 'top',
  data () {
    return {
      msg: 'Please hold the card over the reader',
      message: ''
    }
  },
  mounted: function () {
    const self = this
    const ws = new WebSocket('ws://localhost:3000/socket/readCard')
    ws.onopen = function (e) {
      console.log(' Web socket onopen ')
    }
    ws.onmessage = function (e) {
      console.log(' Web socket onmessage ', e.data)
      self.message = JSON.parse(e.data)
      console.log('Response = ' + self.message)
      if (self.message['IsNew'] === false) {
        self.updateMsg('Now Reading ...')
        self.getUser(self.message['CardID'])
      } else {
        self.createUser(self.message['CardID'])
      }
    }
    ws.onerror = function (e) {
      console.log(' Web socket error ')
      console.log(e)
    }
    ws.onclose = function (e) {
      console.log(' Web socket onclose ' + e)
    }
  },
  methods: {
    updateMsg: function (text) {
      const self = this
      self.msg = text
    },
    getUser: function (CardID) {
      setTimeout(router.push({name: 'welcome', params: {cardid: CardID}}), 500)
    },
    createUser: function (CardID) {
      setTimeout(router.push({name: 'regist', params: {cardid: CardID}}), 500)
    }
  }
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
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
div #top{
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
