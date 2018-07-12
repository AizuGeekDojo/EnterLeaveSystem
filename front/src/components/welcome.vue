<template>
    <div id='welcome' class="container align-middle">
        <h1 class="contents align-middle"> {{ message }}{{ user }} </h1>
    </div>
</template>

<script>
import router from "../router"
export default {
  name: "welcome",
  data() {
    return {
      message: "Now Reading...",
      user: " ",
      isEnter: true,
      sid: " ",
      cardid: " "
    }
  },
  mounted: function() {
    const self = this
    this.cardid = this.$route.params.cardid
    self.cardid = this.cardid
    let date = new Date()
    fetch("http://localhost:3000/api/readuser", {
      mode: "cors",
      // credentials: "include",
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        CardID: this.cardid,
        timestamp: date.getTime()
      })
    })
      .then(response => {
        // .log(response);
        return response.json()
      })
      .then(res => {
        console.log(res)
        if (res["IsEnter"] === true) {
          router.push({ name: "question", params: { res: res } })
        } else {
          self.message = "Welcome To Geek Dojo "
          self.IsEnter = false
          self.sid = res["SID"]
          self.user = res["UserName"]
          this.user = self.user
          this.sid = self.sid
          this.IsEnter = self.IsEnter
          this.message = self.message
          self.push_log()
        }
      })
      .catch(function(error) {
        console.error(error)
      })
  },
  methods: {
    push_log: function() {
      const self = this
      let date = new Date()
      fetch("http://localhost:3000/api/log", {
        mode: "cors",
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          SID: self.sid,
          IsEnter: self.IsEnter,
          timestamp: date.getTime()
        })
      }).catch(function(error) {
        console.error(error)
      }),
        setTimeout(function() {
          router.push({ name: "top" })
        }, 5000)
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
div #welcome {
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
