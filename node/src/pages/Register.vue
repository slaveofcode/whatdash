<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <h1>Register new Account</h1>
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
      msg: "Register New Account"
    };
  },
};
</script>