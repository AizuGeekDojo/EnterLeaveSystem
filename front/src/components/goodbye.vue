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
  mounted: function() {
    const self = this
    return new Promise((resolve, reject) => {
      const date = new Date()
      const sid = self.$route.params.res["SID"]
      fetch("http://localhost:3000/api/log", {
        mode: "cors",
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
        .then(() => {
          setTimeout(() => {
            router.push({ name: "top" })
          }, 3000)
          resolve(true)
        })
        .catch(err => {
          console.log(err)
          reject(err)
        })
    })
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
