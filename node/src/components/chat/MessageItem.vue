<template>
  <div class="chat-display-item" :class="[msg.isMe ? 'from-me' : '']">
    <div class="chat-wrapper">
      <div v-for="m in msg.messages" :key="m.id">
        <!-- <p class="chat-display-name" v-if="isGroupConversation(conv) && !m.me && ((msg.messages[index-1] && msg.messages[index-1].owner != m.owner) || index === 0)">
          {{getContactName(m.owner)}}
        </p> -->
        <div v-if="m.type === 'text'">
          <p class="chat-display-msg" :class="[isFailureSent(m) ? 'unsent-msg' : '']">{{ m.msg }}</p>
        </div>
        <div v-if="['image', 'video'].includes(m.type)">
          <p v-if="m.image.content.trim().length > 0 || m.image.thumb.trim().length > 0" class="chat-display-msg pic text-center">
            <b-img v-bind:src="'data:image/jpeg;base64,'+(m.image.content || m.image.thumb)" fluid :alt="m.caption" />
            <span v-show="!m.me" class="fas fa-file-download media-download-icon" @click="download(m.id, m.ext)"></span>
          </p>
          <p v-if="m.image.thumb.trim().length === 0 && m.type === 'image'" class="chat-display-msg font-weight-bold font-italic">*See image on the device* <span class="fas fa-file-download media-download-icon" @click="download(m.id, m.ext)"></span></p>
          <p v-if="m.image.thumb.trim().length === 0 && m.type === 'video'" class="chat-display-msg font-weight-bold font-italic">*Watch video on the device* <span class="fas fa-file-download media-download-icon" @click="download(m.id, m.ext)"></span></p>
        </div>
        <div v-if="m.type === 'document'">
          <p class="chat-display-msg font-weight-bold font-italic">*See document file on the device* <span class="fas fa-file-download media-download-icon" @click="download(m.id, m.ext)"></span></p>
        </div>
        <div v-if="m.type === 'audio'">
          <p class="chat-display-msg font-weight-bold font-italic">*See audio file on the device* <span class="fas fa-file-download media-download-icon" @click="download(m.id, m.ext)"></span></p>
        </div>
        <div v-if="m.type === 'unknown'">
          <p class="chat-display-msg font-weight-bold font-italic">{{ m.msg }}</p>
        </div>
        <p class="chat-loading-msg" v-show="m.sendingOnTheFly">sending...</p>
        <p class="chat-failure-msg" v-show="isFailureSent(m)">Failed</p>
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

.chat-display-item .chat-display-name {
  color: #17a2b8;
  font-weight: bold;
  padding: 5px;
  padding-bottom: 0;
  margin-bottom: 5px;
}

.chat-display-item .chat-display-msg {
  max-width: 500px;
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

.chat-display-item.from-me .chat-display-msg.unsent-msg {
  background: #fed4d4;
}

.unsent-msg

.chat-display-msg.pic {
  padding: 15px 15px;
  background: #ffe5cf;
}

.chat-loading-msg, .chat-failure-msg {
  font-size: 12px;
  color: #cf7400;
  padding-left: 10px;
}

.chat-failure-msg {
  color: #932f2f;
  font-weight: bold;
}

.media-download-icon {
  color: #1672ad;
  font-size: 20px;
  margin-left: 5px;
  cursor: pointer;
}
.media-download-icon:hover {
  color: #047c8f;
}
</style>

<script>
export default {
  props: ['msg', 'conv', 'getContactName', 'download'],
  data() {
    return {
      lastMsgDisplay: false,
    }
  },
  methods: {
    isGroupConversation(conv) {
      return conv && conv.split('-').length > 1
    },
    isFailureSent(msg) {
      return msg.sendingError !== null && msg.sendingError
    }
  }
}
</script>
