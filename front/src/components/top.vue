<template>
  <div id='top' ref="message">
    <h1>{{ msg }}</h1>
  </div>
</template>

<script>
const ws = new WebSocket('ws://localhost:3000/socket')
export default {
  name: 'top',
  data () {
    return {
      msg: 'カードをリーダーにかざしてください',
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
      if (self.message['message'] === 'True') {
        self.updateMsg('読取中')
        self.getUser()
        console.log('updated at ', self.msg)
      }
    }
    ws.onerror = function (e) {
      console.log(' Web socket error ')
      console.log(e)
    }
    ws.onclose = function (e) {
      console.log(' Web socket onclose ')
    }
  },
  methods: {
    updateMsg: function (text) {
      this.msg = text
    },
    getUser: function () {
      fetch('http://localhost:3000/api/getUser')
        .then(response => {
          return response.json()
        }).then(res => {
          if (res['isExist'] === true) {
            this.updateMsg('Welcome' + ' ' + res['Name'] + '(' + res['StudentId'] + ')')
          }
        }).catch(function (error) {
          console.log(error)
        })
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
