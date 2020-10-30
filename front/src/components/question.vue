<template>
    <div id='ques' style="width: 100%;">
      <div class="question">
        <h2>What is your purpose?</h2>
        <div class="checkboxes alignment" data-toggle="buttons-radio">
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="3DPrinter" v-model="checkedUse">
            <label for="3DPrinter" style="font-size: 18px">3D Printer</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="LaserCutter" v-model="checkedUse">
            <label for="LaserCutter" style="font-size: 18px">Laser Cutter</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="Training session" v-model="checkedUse">
            <label for="Training" style="font-size: 18px">Training</label>
          </div>
          <div class="checkbox">
            <input class="btn btn-info" type="checkbox" value="Other" v-model="checkedUse">
            <label for="Other" style="font-size: 18px">Other</label>
          </div>
        </div>
        <h3>If you have any request, please fill in.</h3>
        <br>
        <textarea v-model="message" placeholder="" @keyup.ctrl.enter="send"></textarea>
        <br>
        <button ref="sendbtn" class="btn btn-info" v-on:click="send">send</button>
        <br>
        <h5 style="color:#aaaaaa;"><span style="color:#88aace;">Ctrl + Enter</span> to submit your request.</h5>
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
      message: '',
      sending: false
    }
  },
  mounted: function () {
    if (!util.isShowQuestion()) {
      var userinfo = this.$route.params.userinfo
      util.addLog(userinfo['SID'], false, '')
      this.$router.push({ name: 'goodbye' })
    }
  },
  methods: {
    send: function () {
      if (this.sending) {
        return
      }
      this.sending = true
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
      this.$router.push({ name: 'goodbye' })
    }
  }
}
</script>

<style scoped>
h2 {
  font-size: 3rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
  width: 100%;
  font-weight: normal;
}

h3 {
  font-size: 1.5rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
  width: 100%;
  font-weight: normal;
}

input {
  font-size: 26px;
}

.checkbox input {
  width: 3rem;
  height:3rem;
}

.checkbox {
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}

.checkboxes {
  width: 50%;
  height: 175px;
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
  text-align: center;
  align-items: center;
}

div #ques {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  height: 600px;
}
</style>