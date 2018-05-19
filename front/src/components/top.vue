<template>
  <div id='top' ref="message">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
import router from '../router'
const ws = new WebSocket('ws://localhost:3000/socket')
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
    ws.onopen = function (e) {
      console.log(' Web socket onopen ')
    }
    ws.onmessage = function (e) {
      console.log(' Web socket onmessage ', e.data)
      self.message = JSON.parse(e.data)
      if (self.message['IsCard'] === true) {
        self.updateMsg('now reading ...')
        self.getUser()
      } else {
        self.createUser(self.message['CardID'])
      }
    }
    ws.onerror = function (e) {
      console.log(' Web socket error ')
      console.log(e)
    }
    ws.onclose = function (e) {
      console.log(' Web socket onclose ')
      setTimeout(self.mounted(), 5000)
    }
  },
  methods: {
    updateMsg: function (text) {
      const self = this
      self.msg = text
    },
    getUser: function () {
      setTimeout(router.push({name: 'welcome', params: {cardid: this.$route.params.CardID}}), 2000)
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
  font-weight: normal;
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
</style>
