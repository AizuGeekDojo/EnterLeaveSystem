<template>
    <div id='ques' style="width: 100%;">
        <p>What did you use?</p>
        <input type="checkbox" value="3DPrinter" v-model="checkedUse">
        <label for="3DPrinter">3DPrinter</label>
        <input type="checkbox" value="LaserCutter" v-model="checkedUse">
        <label for="LaserCutter">LaserCutter</label>
        <input type="checkbox" value="Other" v-model="checkedUse">
        <label for="Other">Other</label>

        <p>If you have any request please fill in</p>
        <p style="white-space: pre-line;">{{ message }}</p>
        <br>
        <textarea v-model="message" placeholder="add multiple lines"></textarea>
        <br>

        <button v-on:click="send">send</button>
        <span>Checked names: {{ checkedUse }}</span>
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
  font-size: 72px;
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
  width: 300px;
  height: 100px;
  font-size: 26px;
}
div .studentNUM {
  width: 100%;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
}
</style>
