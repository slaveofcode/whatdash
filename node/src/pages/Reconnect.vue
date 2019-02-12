<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <h1>Reconnect Account</h1>
    </b-container>
  </div>
</template>

<script>
import TopNav from "../components/TopNav.vue";
export default {
  components: {
    TopNav
  },
  data() {
    return {
      msg: "Reconnect Account"
    };
  }
};
</script>