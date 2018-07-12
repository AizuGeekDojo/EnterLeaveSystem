<template>
    <div id='goodbye' class="container align-middle">
        <h1 class="contents align-middle">{{ message }} </h1>
    </div>
</template>

<script>
import router from "../router"
export default {
  name: "goodbye",
  data() {
    return {
      message: "Good bye"
    }
  },
  mounted: () => {
    this.send()
    setTimeout(function() {
      router.push({ name: "top" })
    }, 3000)
  },
  methods: {
    send: () => {
      return new Promise(function(resolve, reject) {
        const self = this
        let date = new Date()
        let sid = this.$route.params.res["SID"]
        // console.log('question res', this.$route.params.res)
        fetch("http://localhost:3000/api/log", {
          mode: "no-cors",
          credentials: "include",
          method: "POST",
          headers: {
            Accept: "application/json",
            "Content-Type": "application/json"
          },
          body: JSON.stringify({
            SID: sid,
            IsEnter: true,
            timestamp: date.getTime()
          })
        })
      })
    }
  }
}
</script>

<style scoped>
h1,
h2 {
  font-size: 72px;
  display: table-cell;
  height: 100%;
  width: 100%;
  font-weight: normal;
  text-align: center;
  vertical-align: middle;
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
div #goodbye {
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
