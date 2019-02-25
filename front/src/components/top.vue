<template>
  <div id='top' ref="message" class="container align-middle">
    <h1 class="contents align-middle">Please hold the card over the reader</h1>
  </div>
</template>

<script>
// import router from "../router";
import util from "../util.js";

export default {
  name: "top",
  // data() {
  //   return {
  //     message: ""
  //   };
  // },
  // created: function() {
  //   console.log("Created");
  // },
  destroyed: function() {
    console.log(this.ws)
    this.ws.close()
  },
  mounted: function() {
    const self = this;
    this.ws = new WebSocket("ws://localhost:3000/socket/readCard");
    this.ws.onopen = function(e) {
      console.log(" Web socket onopen ");
    };
    this.ws.onmessage = function(e) {
      // console.log(this)
      console.log(" Web socket onmessage ", e.data);
      var message = JSON.parse(e.data);
      console.log("Response = " + message);
      if (message["IsNew"] === false) {
        util.getUserInfo(message["SID"])
          .then(res => {
            if (res["IsEnter"]) {
              self.$router.push({ name: "question", params: { userinfo: res } });
            } else {
              self.$router.push({ name: "welcome", params: { userinfo: res } });
            }
          })
        // self.updateMsg("Now Reading ...");
        // self.getUser(message["SID"]);
      } else {
                self.$router.push({ name: "regist", params: { cardid: message["CardID"] } });

        // self.createUser(message["CardID"]);
      }
    };
    this.ws.onerror = function(e) {
      console.log(" Web socket error ");
      console.log(e);
    };
    this.ws.onclose = function(e) {
      console.log(" Web socket onclose " + e);
    };
  },
  methods: {
    updateMsg: function(text) {
      const self = this;
      self.msg = text;
    },
    getUser: function(SID) {
      setTimeout(function() {
        router.push({ name: "welcome", params: { sid: SID } });
      }, 500);
    },
    createUser: function(CardID) {
      setTimeout(function() {
        router.push({ name: "regist", params: { cardid: CardID } });
      }, 500);
    }
  }
};
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
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
div #top {
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
