<template>
    <div id='forgot' style="width: 100%;">
      <div class="studentNUM">
        <h1>Input your student number</h1>
        <input v-on:keyup.enter="forgot" v-model="sid" placeholder="s12xxxxx">
        <h4 style="color:#aaaaaa;"><br><span style="color:#88aace;">Enter</span> key to continue -></h4>
      </div>
    </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'forgot',
  data () {
    return {
      sid: ''
    }
  },
  methods: {
    forgot: function () {
      const self = this
      util.getUserInfo(self.sid)
        .then(res => {
          console.log(res)
          if (res['UserName'] === '') {
            console.log('Your ID is incorrect.')
            alert('Your ID is incorrect.')
          } else {
            if (res['IsEnter']) {
              self.$router.push({ name: 'question', params: { userinfo: res } })
            } else {
              self.$router.push({ name: 'welcome', params: { userinfo: res } })
            }
          }
        })
        .catch(function (error) {
          console.log(error)
          self.$router.push({ name: 'top' })
        })
    }
  }
}
</script>

<style scoped>
h1 {
  font-size: 4rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
}
h2 {
  font-size: 50px;
  width: 100%;
  font-weight: normal;
}
input {
  width: 325px;
  height: 125px;
  font-size: 2.2rem;
  font-weight: lighter;
  letter-spacing: 0.05rem;
  margin: 15px;
}
div .studentNUM {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  width: 80%;
  height: 200px;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}
</style>
