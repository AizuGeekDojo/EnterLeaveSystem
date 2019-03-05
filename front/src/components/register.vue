<template>
    <div id="regist" style="width: 100%;">
      <div class="studentNUM">
        <h1>Please Input Your Student Number</h1>
        <h2>Input: {{ sid }}</h2>
        <input v-on:keyup.enter="regist" v-model="sid" placeholder="s120000">
      </div>
    </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'regist',
  data () {
    return {
      sid: ''
    }
  },
  methods: {
    regist: function () {
      const self = this
      const cardid = this.$route.params.cardid
      util.registCardInfo(cardid, this.sid)
        .then(res => {
          if (res['Success'] !== true) {
            console.log('Card register failed')
            alert('The ID is not found.')
          } else {
            util.getUserInfo(self.sid)
              .then(res => {
                if (res['IsEnter']) {
                  self.$router.push({ name: 'question', params: { userinfo: res } })
                } else {
                  self.$router.push({ name: 'welcome', params: { userinfo: res } })
                }
              })
          }
        })
        .catch(function (error) {
          console.error(error)
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
