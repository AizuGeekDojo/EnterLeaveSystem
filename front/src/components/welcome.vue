<template>
    <div id='welcome'>
        <p> Welcome Geek Dojo {{ message }} </p>
    </div>
</template>

<script>
import router from '../router'
export default {
  name: 'welcome',
  data () {
    return {
      message: ''
    }
  },
  mounted: function () {
    const self = this
    let cardid = this.$route.params.cardid
    let date = new Date()
    fetch('http://localhost:3000/api/readuser', {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(
        {
          'CardID': cardid,
          'timestamp': date.getTime()
        }
      )
    }).then(response => {
      return response.json()
    }).then(res => {
      self.message = res['UserName']
    }).catch(function (error) {
      alert('Error ' + error + ' ' + self.message)
    })
    setTimeout(router.push({name: 'top'}), 4000)
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
