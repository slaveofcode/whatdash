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