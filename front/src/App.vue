<template>
  <div style="background: #EEEEEE;">
    <div id="app" class="container">
    <div style="position:fixed; height: 100px;">
      <router-link to="/" id="totopPos">←</router-link>
    </div>
    <h2 id='clockPos'>{{ clocktext }}</h2>
    <router-view/>
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data: function () {
  return {
    clocktext: '0000/00/00  00:00:00'
  }
  },
  mounted: function () {
    const self = this
    setInterval(() => {
      const da = new Date()
      const year = da.getFullYear()
      const month = da.getMonth() + 1
      const date = da.getDate()
      const hour = da.getHours()
      const minute = da.getMinutes()
      const second = da.getSeconds()

      self.clocktext = `${year}/`
      if (month < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${month}/`
      if (date < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${date}  `
      if (hour < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${hour}:`
      if (minute < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${minute}:`
      if (second < 10) {
        self.clocktext += '0'
      }
      self.clocktext += `${second}`
    }, 1000);
  }
}
</script>

<style>
#app {
  font-family: verdana;
  vertical-align: middle;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  overflow: hidden;
}

#clockPos{
  text-align: right;
  padding: 24px 48px;
  font-size: 2rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
}

#totopPos{
  position: absolute;
  left: 0px;
  padding: 24px 48px;
  font-size: 2rem;
  font-weight: bold;
  letter-spacing: 0.05rem;
  color: #3282ce;
}
</style>
