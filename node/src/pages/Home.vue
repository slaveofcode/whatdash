<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <b-jumbotron
        header="WhatsApp Collaboration Dashboard "
        lead="Start by register or reconnect your existing account"
        bg-variant="transparent"
        text-variant="#333"
      ></b-jumbotron>
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
  },
  methods: {},
  socket: {
    events: {
      connect() {
        console.log("Websocket connected");
      },
      disconnect() {
        console.log("Websocket disconnected");
      },
      changed() {
        console.log("Something changed in Websocket");
      },
      error(err) {
        console.log("Websocket error:", err);
      }
    }
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
