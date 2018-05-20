<template>
    <div id='snum' class="container">
      <div class="box">
        <h1 class="">Please Input Your Student Number</h1>
        <h1 class="">Input: {{ snum }}</h1>
        <input class="" v-on:keyup.enter="regist" v-model="snum" placeholder="s120000" style="border: 2px, #42b983, double;">
      </div>
    </div>
</template>

<script>
import router from '../router'
export default {
  name: 'snum',
  data () {
    return {
      snum: ''
    }
  },
  methods: {
    regist: function () {
      let cardid = this.$route.params.cardid
      let date = new Date()
      fetch('http://localhost:3000/api/createuser', {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(
          {
            'SID': this.snum,
            'CardID': cardid,
            'timestamp': date.getTime()
          }
        )
      }).then(response => {
        return response.json()
      }).then(res => {
        if (res['Success'] !== true) {
          alert('Create failed')
          setTimeout(router.push({name: 'top'}), 500)
        } else {
          setTimeout(router.push({name: 'welcome', params: {cardid: cardid}}), 500)
        }
      }).catch(function (error) {
        alert(error)
      })
    }
  }
}
</script>

<style scoped>
h1, h2 {
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
.box {

}
.container{
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
