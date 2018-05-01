<template>
    <div id='snum'>
        <p>学籍を入力してください: </p>
        <p>Input: {{ snum }}</p>
        <input v-on:keyup.enter="submit" v-model="snum" placeholder="s120000" style="border: 2px, #42b983, double;">
    </div>
</template>

<script>
export default {
  name: 'snum',
  data () {
    return {
      snum: ''
    }
  },
  methods: {
    submit: function () {
      console.log('submit', this.snum)
      fetch('http://localhost:3000/api/register', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({'snum': this.snum})
      }).then(response => {
        return response.json()
      }).then(res => {
        console.log(res)
      }).catch(function (error) {
        console.log(error)
      })
      location.replace('http://localhost:8080/#/welcome')
    }
  }
}
</script>

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
