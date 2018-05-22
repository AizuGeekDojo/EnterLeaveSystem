<template>
    <div id='ques' style="width: 100%;">
      <div class="question">
        <h2>What did you use?</h2>
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
            <input class="btn btn-info" type="checkbox" value="Other" v-model="checkedUse">
            <label for="Other" style="font-size: 20px">Other</label>
          </div>
        </div>
        <h3>If you have any request please fill in.</h3>
        <br>
        <textarea v-model="message" placeholder=""></textarea>
        <br>
        <button class="btn btn-info" v-on:click="send">send</button>
      </div>
    </div>
</template>

<script>
import router from '../router'
export default {
  name: 'ques',
  data () {
    return {
      checkedUse: [],
      message: ''
    }
  },
  methods: {
    send: function () {
      const self = this
      let date = new Date()
      fetch('http://localhost:3000/api/log', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(
          {
            'SID': self.sid,
            'IsEnter': self.IsEnter,
            'Ext': {
              'Use': self.checkedUse,
              'message': self.message
            },
            'timestamp': date.getTime()
          }
        )
      })
      setTimeout(function () {
        router.push({name: 'goodbye'})
      }, 5000)
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
