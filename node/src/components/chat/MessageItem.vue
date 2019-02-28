<template>
  <div class="chat-display-item" :class="[msg.isMe ? 'from-me' : '']">
    <div class="chat-wrapper">
      <div v-for="m in msg.messages" :key="m.id">
        <p class="chat-display-msg">{{ m.msg }}</p>
        <p class="chat-loading-msg" v-show="m.sendingOnTheFly">sending...</p>
      </div>
    </div>
  </div>
</template>

<style>
.chat-display-item {
  font-size: 14px;
  padding: 0 30px;
}

.chat-display-item > .chat-wrapper {
  float: left;
  margin-bottom: 10px;
}

.chat-display-item.from-me  > .chat-wrapper{
  float: right;
}

.chat-display-item:after {
  display: table;
  content: "";
  clear: both;
}

.chat-display-item .chat-display-msg {
  max-width: 400px;
  background: #bbedc5;
  padding: 7px 12px;
  margin-bottom: 5px;
  color: #444;
  border-radius: 1em;
  display: block;
}

.chat-display-item.from-me .chat-display-msg {
  background: #cccccc;
}

.chat-loading-msg {
  font-size: 12px;
  color: #cf7400;
  padding-left: 10px;
}
</style>

<script>
export default {
  props: ['msg']
}
</script>
