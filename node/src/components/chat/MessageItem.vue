<template>
  <div class="chat-display-item" :class="[fromMe ? 'from-me' : '']">
    <div class="chat-wrapper">
      <span class="chat-display-msg">Hello</span>
      <span class="chat-display-msg">How r you?</span>
    </div>
  </div>
</template>

<style>
.chat-display-item {
  font-size: 14px;
  padding: 0 10px;
}

.chat-display-item > .chat-wrapper {
  float: left;
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
  background: #cccccc;
  padding: 7px 12px;
  margin-bottom: 2px;
  color: #444;
  border-radius: 1em;
  display: block;
}

.chat-display-item.from-me .chat-display-msg {
  background: #bbedc5;
}
</style>

<script>
export default {
  props: ['fromMe']
}
</script>
