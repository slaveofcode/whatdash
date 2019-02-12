<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <b-jumbotron header="WhatsApp Collaboration Dashboard " lead="Start by register or reconnect your existing account" bg-variant="transparent" text-variant="#333">
      </b-jumbotron>
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
      msg: "Welcome to WhatDash"
    };
  }
};
</script>

<style lang="scss">
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
