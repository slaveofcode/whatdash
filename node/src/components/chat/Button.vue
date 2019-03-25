<template>
  <div class="container-buttons">
    <div class="send-btn send-text">
      <button class="btn btn-warning" @click="sendTextAction()">
        <span class="fas fa-paper-plane"></span>
      </button>
    </div>
    <div class="send-btn send-image">
      <file-upload class="btn btn-warning"
        post-action="/upload/post"
        extensions="gif,jpg,jpeg,png,webp"
        accept="image/png,image/gif,image/jpeg,image/webp"
        :multiple="false"
        :size="1024 * 1024 * 10"
        ref="upload"
        v-model="files"
        @input-filter="inputFilter"
        @input="setFileUpload"
        >
        <span class="fas fa-paperclip"></span>
      </file-upload>
    </div>
  </div>
</template>

<style lang="css">
.container-buttons {
  position: absolute;
  top: 15px;
  right: 0px;
  width: 110px;
}
.send-btn {
  position: relative;
  float:left;
  margin-left: 5px;
  cursor: pointer;
}
.send-btn .btn {
  border-radius: 40%;
  color: #fff;
}

</style>


<script>
export default {
  props: [
    'sendTextAction',
    'sendImageAction'
  ],
  data() {
    return {
      files: [],
    }
  },
  methods: {
    setFileUpload(file){
      this.sendImageAction(this.files)
      this.files = []
    },
    inputFilter(newFile, oldFile, prevent) {
      if (newFile && !oldFile) {
        // Before adding a file
        // Filter system files or hide files
        if (/(\/|^)(Thumbs\.db|desktop\.ini|\..+)$/.test(newFile.name)) {
          return prevent()
        }
        // Filter php html js file
        if (/\.(php5?|html?|jsx?)$/i.test(newFile.name)) {
          return prevent()
        }
      }
    },
    inputFile(newFile, oldFile) {
      // this.$refs.upload.active = true
      console.log('inputFile')
    }
  }
}
</script>
