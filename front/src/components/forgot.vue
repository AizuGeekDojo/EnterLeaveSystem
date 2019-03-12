<template>
    <div id='forgot' style="width: 100%;">
      <div class="studentNUM">
        <h1>Please Input Your Student Number</h1>
        <h1>Input: {{ sid }}</h1>
        <input v-on:keyup.enter="forgot" v-model="sid" placeholder="s120000">
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
            console.log('The ID is not found.')
            alert('The ID is not found.')
          } else {
            if (res['IsEnter']) {
              self.$router.push({ name: 'question', params: { userinfo: res } })
            } else {
              self.$router.push({ name: 'welcome', params: { userinfo: res } })
            }
            self.$router.push({ name: 'question', params: { userinfo: res } })
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
  font-size: 72px;
  width: 100%;
  font-weight: normal;
}
h2 {
  font-size: 50px;
  width: 100%;
  font-weight: normal;
}
input {
  width: 300px;
  height: 100px;
  font-size: 26px;
  border: 2px, #42b983, double;
}
div .studentNUM {
  width: 100%;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}
</style>
