<template>
    <div id='goodbye'>
        <h1>bye!</h1>
    </div>
</template>

<script>
import router from '../router'
export default {
  name: 'goodbye',
  mounted: function () {
    setTimeout(function () {
      router.push({ name: 'top' })
    }, 3000)
  }
}
</script>

<style scoped>
h1{
  font-size: 5rem;
  font-weight: normal;
  letter-spacing: 0.05rem;
}
div #goodbye {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  width: 80%;
  height: 200px;
  text-align: center;
}
</style>
