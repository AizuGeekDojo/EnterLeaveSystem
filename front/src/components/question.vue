<template>
    <div id='ques' style="width: 100%;">
      <div class="question">
        <h2>What is your purpose?</h2>
        <div class="checkboxes alignment" data-toggle="buttons-radio">
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="3DPrinter" v-model="checkedUse">
            <label for="3DPrinter" style="font-size: 20px">3DPrinter</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="LaserCutter" v-model="checkedUse">
            <label for="LaserCutter" style="font-size: 20px">LaserCutter</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="Training session" v-model="checkedUse">
            <label for="Training" style="font-size: 20px">Training</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="Other" v-model="checkedUse">
            <label for="Other" style="font-size: 20px">Other</label>
          </div>
        </div>
        <h3>If you have any request please fill in.</h3>
        <br>
        <textarea v-model="message" placeholder=""></textarea>
        <br>
        <button ref="sendbtn" class="btn btn-info" v-on:click="send">send</button>
      </div>
    </div>
</template>

<script>
import util from '../util.js'

export default {
  name: 'ques',
  data () {
    return {
      checkedUse: [],
      message: ''
    }
  },
  mounted: function () {
    if (util.isShowQuestion()){
      util.addLog(userinfo['SID'], false, '')
      .then(res => {
        this.$router.push({ name: 'goodbye' })
      })
    }
  },
  methods: {
    send: function () {
      this.$refs.sendbtn.disabled = true
      var userinfo = this.$route.params.userinfo
      var answer = JSON.stringify({
        'Use': this.checkedUse,
        'message': this.message
      })
      if (userinfo === undefined) {
        this.$router.push({ name: 'top' })
        return
      }
      util.addLog(userinfo['SID'], false, answer)
        .then(res => {
          this.$router.push({ name: 'goodbye' })
        })
    }
  }
}
</script>

<style scoped>
h1, h2 {
  font-size: 60px;
  width: 100%;
  font-weight: normal;
}
input {
  font-size: 26px;
}
.checkbox input {
  width: 50px;
  height:50px;
}
.checkbox {
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}
.checkboxes {
  width: 500px;
  height: 200px;
  display: flex;
  justify-content: center;
  align-items: center;
}
div .question textarea{
  resize: none;
  width:500px;
  height:200px;
}
div .question {
  width: 100%;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}
</style>
