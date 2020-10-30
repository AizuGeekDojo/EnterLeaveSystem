<template>
    <div id='welcome'>
        <h1> Welcome to {{roomname}} {{ username }}! </h1>
    </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'welcome',
  data () {
    return {
      username: '',
      roomname: ''
    }
  },
  mounted: function () {
    this.roomname = util.roomName()
    var userinfo = this.$route.params.userinfo
    this.username = userinfo['UserName']
    var sid = userinfo['SID']
    util.addLog(sid, true, '')
    var self = this
    setTimeout(function () {
      self.$router.push({ name: 'top' })
    }, 5000)
  }
}
</script>

<style scoped>
h1{
  font-size: 4rem;
}
div #welcome {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  width: 80%;
  height: 200px;
  text-align: center;
}
</style>
