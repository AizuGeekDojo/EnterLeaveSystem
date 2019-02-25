<template>
    <div id='welcome' class="container align-middle">
        <h1 class="contents align-middle"> Welcome To Geek Dojo {{ username }} </h1>
    </div>
</template>

<script>
// import router from "../router";
import util from "../util.js";

export default {
  name: "welcome",
  data() {
    return {
      message: "Now Reading...",
      username: " ",
      isEnter: true,
      sid: " "
    };
  },
  mounted: function() {
    var userinfo = this.$route.params.userinfo;
    this.username = userinfo["UserName"]
    var sid = userinfo["SID"]
    var IsEnter = userinfo["IsEnter"]
    util.addLog(sid,true,"")
    var self = this
    setTimeout(function() {
      self.$router.push({ name: "top" });
    }, 5000);
  }
};
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
