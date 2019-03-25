<template>
  <div class="media-loader" v-show="isLoading">
    <p><span class="fas fa-sync fa-spin"></span>&nbsp;&nbsp;Uploading attachment...</p>
  </div>
</template>

<style>
.media-loader {
  padding: 5px 0 5px;
  font-size: 13px;
  text-align: center;
  background: #fff5d7;
}

.media-loader p {
  font-weight: bold;
  color: #b48b0f;
  padding: 0;
  margin: 0;
}
</style>


<script>
export default {
  props: ['isLoading'],
}
</script>

