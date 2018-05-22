<template>
    <div id='goodbye' class="container align-middle">
        <h1>{{ message }}</h1>
    </div>
</template>

<script>
import router from '../router'
export default {
  name: 'goodbye',
  data () {
    return {
      message: 'Good bye'
    }
  },
  mounted: function () {
    setTimeout(function () {
      router.push({name: 'top'})
    }, 3000)
  }
}
</script>

<style scoped>
h1, h2 {
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
div #welcome{
  display: table;
  text-align: center;
  vertical-align: middle;
}
</style>
