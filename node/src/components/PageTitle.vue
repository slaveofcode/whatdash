<template>
  <div class="main-title">
    <h1>{{text}}</h1>
  </div>
</template>

<style>
.main-title {
  padding: 20px 0;
  margin: 0 0 20px 0;
  color: #515151;
  border-bottom: 1px solid #c5c5c5;
}
</style>


<script>
export default {
  props: ['text'],
}
</script>