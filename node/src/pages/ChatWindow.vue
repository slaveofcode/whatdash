<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle" :number="pageTitleNumber" :conversationWindow="true"></PageTitle>
      <div class="chat-container">
        <div class="section-contacts">
          <ContactItem
            v-for="contact in contactChats"
            :key="contact.id"
            v-bind:contact="contact"
            @click.native="clickOnContact(contact)"
          ></ContactItem>
        </div>
        <div class="section-messages">
          <div class="section-chat-header" v-show="conversationTitle">
            <p>{{conversationTitle}} <span v-show="!isGroupNumber(conversationId)" class="peer-number">{{conversationPeerNumber}}</span><span v-show="isGroupNumber(conversationId)" class="peer-number">Group Chat</span></p>
            <div v-b-tooltip title="Sync socket and messages" class="resync-button socket" @click="resyncSocket()">
              <i class="fas fa-sync-alt" :class="[requestSyncSocket ? 'fa-spin' : '']"></i>
            </div>
            <div v-b-tooltip title="Sync messages" class="resync-button messages" @click="resyncConversation()">
              <i class="fas fa-comments"></i>
            </div>
          </div>
          <div class="section-chat">
            <span
              class="empty-message"
              :class="[conversationTitle ? 'hide' : '']"
            >Select some chat to start conversation.</span>
            <span
              class="empty-message loading"
              v-show="conversationId && !conversations[conversationId]"
            >Loading messages...</span>
            <div
              class="chat-display"
              v-show="key === conversationId"
              v-for="(conversation, key) in conversations"
              :key="key"
              :class="hashCode(key)"
            >
              <MessageItem 
                v-for="(msg, idx) in conversation.displayMessages" 
                :key="idx" 
                :msg="msg" 
                :conv="conversationContact.id" 
                :getContactName="getNameByGroup"
                :download="downloadMedia"
                >
              </MessageItem>
              <MediaLoading :isLoading="inUploadMedia" />
            </div>
          </div>
          <div class="section-input" v-show="conversationTitle && conversations[conversationId]">
            <textarea-autosize
              class="input-chat"
              placeholder="Type something and press (CTRL+ENTER) to send immediately"
              v-model="chatInput"
              :min-height="70"
              @keydown.native="onKeydown"
            ></textarea-autosize>
            <ChatButton :sendTextAction="sendText" :sendImageAction="sendImage"></ChatButton>
          </div>
        </div>
      </div>

      <b-modal v-model="modalErrConnShow" centered title="Server Disconnected" @hide="modalErrConnClosed()">
        Server was disconnected for a while, please wait for a minute to get it running again.
        <div slot="modal-footer" class="w-100">
          <b-button size="sm" class="float-right" variant="primary" @click="modalErrConnClosed()">Close</b-button>
        </div>
      </b-modal>
      <b-modal v-model="modalMsg.show" centered :title="modalMsg.title" @hide="modalMsg.show = false">
        {{modalMsg.message}}
        <div slot="modal-footer" class="w-100">
          <b-button size="sm" class="float-right" variant="primary" @click="modalMsg.show = false">Close</b-button>
        </div>
      </b-modal>
    </b-container>
  </div>
</template>

<style>
.chat-container {
  display: flex;
  flex-direction: row;

  width: 100%;
  height: 70vh;
}

.section-contacts {
  max-width: 200px;
  height: 100%;
  overflow-y: auto;
  border-bottom: 1px solid #dad8d8;
  border-top: 1px solid #dad8d8;
}

.section-messages {
  position: relative;
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
}

.resync-button {
  position: absolute;
  background: #ffffff;
  padding: 5px 10px;
  border-radius: 5px;
  cursor: pointer;
}

.resync-button.socket {
  right: 65px;
  top: 12px;
  color: #ea7038;
  border: 1px solid #f9b799;
}

.resync-button.messages {
  right: 20px;
  top: 12px;
  color: #007bff;
  border: 1px solid #b6d9ff;
}

.resync-button:hover {
  background: #eeeeee;
  border: 1px solid #cccccc;
}

.resync-button:active {
  background: #cccccc;
  border: 1px solid #aaaaaa;
  color: #555;
}

.section-messages .section-chat-header {
  font-weight: bold;
  padding: 15px 15px 0;
  margin: 0;
  background: #eee;
  border-top: 1px solid #dad8d8;
  border-right: 1px solid #dad8d8;
}

.section-chat-header .peer-number {
  font-weight: normal;
  font-size: 13px;
  color: #007bff;
  background: #fff;
  border-radius: 7px;
  padding: 3px 8px;
  border: 1px solid #ccc;
}

.section-chat {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
  background: #eee;
  overflow-y: auto;
}

.section-chat .empty-message {
  font-size: 14px;
  color: #aaa;
  display: flex;
  justify-content: center;
  padding-top: 25vh;
}

.section-chat .empty-message.loading {
  color: #075419;
  padding-top: 23vh !important;
}

.section-chat .empty-message.hide,
.section-chat .chat-display.hide,
.section-chat-header.hide {
  display: none;
}

.section-chat .chat-display {
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
  background: #ffffff;
  border-right: 1px solid #e0e0e0;
  border-top: 1px solid #e9e7e7;
  padding-top: 15px;
}

.section-input {
  width: 100%;
  height: 70px;
  align-self: flex-end;
  position: relative;
}

.section-input .input-chat {
  font-size: 14px;
  width: 100%;
  background: #ffffff;
  border: none;
  color: #333333;
  padding: 10px 15px;
  border-bottom: 1px solid #dad8d8;
  border-right: 1px solid #e0e0e0;
  border-top: 1px solid #e9e7e7;
}
</style>

<script>
import PageTitle from "../components/PageTitle.vue";
import TopNav from "../components/TopNav.vue";
import ContactItem from "../components/chat/ContactItem.vue";
import MessageItem from "../components/chat/MessageItem.vue";
import ChatButton from "../components/chat/Button.vue";
import MediaLoading from "../components/chat/MediaLoading.vue";
import Req from "../req";
import { axios } from "../req";

export default {
  components: {
    TopNav,
    PageTitle,
    MessageItem,
    ContactItem,
    ChatButton,
    MediaLoading,
  },
  data() {
    return {
      pageTitle: "Active Chat",
      pageTitleNumber: "",
      detailAccount: null,
      contacts: [],
      chatHistory: [],
      contactChats: [],
      chatInput: null,
      conversationId: null,
      conversationTitle: null,
      conversationPeerNumber: null,
      conversationContact: null,
      conversations: {},
      activeConversation: null,
      activeConversationPool: false,
      requestSyncSocket: false,
      poolContactId: null,
      onPoolContacts: false,
      failedPoolAttempt: 0,
      maxFailedPoolAttempt: 10,
      modalErrConnShow: false,
      modalMsg: {},
      inUploadMedia: false,
    };
  },
  watch: {
    $route: "initPage"
  },
  created() {
    this.initPage();
  },
  beforeDestroy() {
    if (this.poolContactId) clearInterval(this.poolContactId)
    if (this.activeConversation) clearInterval(this.activeConversation);
  },
  methods: {
    modalErrConnClosed() {
      this.modalErrConnShow = false
      setTimeout(() => {
        this.clickOnContact(this.conversationContact)
      }, 1000 * 10) // 10 secs
    },
    onKeydown(evt) {
      // detecting ctrl+enter
      if (evt.keyCode === 13 && evt.ctrlKey) {
        // send message
        this.sendText()
      }
    },
    base64toBlob(base64Data, contentType) {
      contentType = contentType || "";
      var sliceSize = 1024;
      var byteCharacters = atob(base64Data);
      var bytesLength = byteCharacters.length;
      var slicesCount = Math.ceil(bytesLength / sliceSize);
      var byteArrays = new Array(slicesCount);

      for (var sliceIndex = 0; sliceIndex < slicesCount; ++sliceIndex) {
        var begin = sliceIndex * sliceSize;
        var end = Math.min(begin + sliceSize, bytesLength);

        var bytes = new Array(end - begin);
        for (var offset = begin, i = 0; offset < end; ++i, ++offset) {
          bytes[i] = byteCharacters[offset].charCodeAt(0);
        }
        byteArrays[sliceIndex] = new Uint8Array(bytes);
      }
      return new Blob(byteArrays, { type: contentType });
    },
    hashCode(str) {
      var hash = 0,
        i,
        chr;
      if (str.length === 0) return 'ss' + hash;
      for (i = 0; i < str.length; i++) {
        chr = str.charCodeAt(i);
        hash = (hash << 5) - hash + chr;
        hash |= 0; // Convert to 32bit integer
      }
      return 'ss' + Math.abs(hash);
    },
    scrollDownChat() {
      const el = document.querySelector('.section-chat')
      if (el) {
        el.scrollTop = el.scrollHeight
      }
    },
    async initPage() {
      try {
        this.onPoolContacts = true
        await this.syncContacts()
        this.onPoolContacts = false  
      } catch (err) {
        this.onPoolContacts = false  
      }

      this.pageTitle = 'Active Chat'
      this.pageTitleNumber = this.detailAccount.number

      if (this.contactChats.length > 0) {
        this.clickOnContact(this.contactChats[0])
      }
    },
    async syncContacts() {
      if (!this.detailAccount) {
        const detailAccount = await this.loadAccountDetail(this.$route.params.id);
        if (detailAccount) {
          this.detailAccount = detailAccount
        }
      }

      if (this.detailAccount) {

        const contacts = await this.loadContacts(this.detailAccount.number);
        if (contacts && contacts.length > 0) {
          this.contacts = contacts
        }
        
        const history = await this.loadHistory(this.detailAccount.number);
        if (history && history.length > 0) {
          this.chatHistory = history
        }

        if (this.contacts && this.chatHistory) {
          this.contactChats = this.parseHistoryWithContact(
            this.contacts,
            this.chatHistory
          );
        }

      }
      
    },
    async loadAccountDetail(accId) {
      const acc = await Req.get(`/account/detail/${accId}`);
      return acc.status === 200 ? acc.data : null;
    },
    async requestContacts(number) {
      const c = await Req.post('/wa/contact/load', { number });
      return c.status === 200 ? true : false
    },
    async loadContacts(number) {
      const c = await Req.post("/wa/contact/list", { number });
      return c.status === 200 ? c.data : null;
    },
    async loadHistory(number) {
      const c = await Req.post("/chat/history", { number });
      return c.status === 200 ? c.data : null;
    },
    async sendMessageSignal(number, jid, loadType = 'load', messageCount = 100) {
      let uri = '/wa/messages/load'

      if (loadType === 'next') {
        uri = '/wa/messages/load-next'
      } else if (loadType === 'prev') {
        uri = '/wa/messages/load-prev'
      }

      const c = await Req.post(uri, { number, jid, messageCount, });
      const stat = c.status === 200 ? true : false
      if (!stat) return false
      return c.data.status === 'requested'
    },
    parseHistoryWithContact(contacts, chatHistory) {
      const parsed = [];

      for (const history of chatHistory) {
        const foundContact = contacts.find(c => c.jid == history.wa.jid);
        const contactNumber = history.wa.jid.split("@")[0];
        let contactName = history.wa.jid.split("@")[0];

        if (foundContact) {
          contactName = !!foundContact.contact.name
            ? foundContact.contact.name
            : foundContact.jid.split("@")[0];
        }

        parsed.push({
          id: history.wa.jid,
          number: contactNumber,
          name: contactName,
          time: new Date(history.lastChatTime * 1000),
          msgCount: history.msgCount
        });
      }

      return parsed;
    },
    parseMessageItem(item, onTheFly = false) {

      const getFileExt = filename => {
        if (item.filename && item.filename.length > 0) {
          const len = item.filename.split('.')
          return len.length > 1 ? `.${len[1]}` : undefined
        }
        return undefined
      }

      const waMsg = item.wamsg;
      const waInfo = waMsg.info;
      if (waMsg.type === "text") {
        return {
          id: waInfo.id,
          owner: waInfo.remoteJid,
          type: 'text',
          msg: item.text,
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
          sendingOnTheFly: onTheFly,
          sendingError: null,
        };
      } else if (['image', 'video'].includes(waMsg.type)) {
        return {
          id: waInfo.id,
          owner: waInfo.remoteJid,
          type: waMsg.type,
          caption: item.caption,
          image: {
            thumb: item.thumb,
            content: item.content,
          },
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
          sendingOnTheFly: onTheFly,
          sendingError: null,
          ext: getFileExt(item.filename),
        };
      } else if (['document', 'audio'].includes(waMsg.type)) {
        return {
          id: waInfo.id,
          owner: waInfo.remoteJid,
          type: waMsg.type,
          caption: item.caption,
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
          sendingOnTheFly: onTheFly,
          sendingError: null,
          ext: getFileExt(item.filename),
        };
      }

      return {
        id: waInfo.id,
        owner: waInfo.remoteJid,
        type: 'unknown',
        msg: '*unsupported message, see on device*',
        me: waInfo.fromMe,
        stat: waInfo.msgStatus,
        sendingOnTheFly: onTheFly,
      };
    },
    parseMessages(messages) {
      if (!messages || messages.length <= 0) return [];
      return messages.reverse().map(item => this.parseMessageItem(item));
    },
    parseMessageDisplay(messages) {
      if (!messages || messages.length <= 0) return [];
      let isLastSectionMe = null;

      const cleanMsgs = []
      for (const m of messages) {
        if (!m.id)  {
          cleanMsgs.push(m) // for sending msg
        } else if (!cleanMsgs.find(item => item.id === m.id)) {
          cleanMsgs.push(m)
        }
      }

      const displayMsgs = [];
      for (const msg of cleanMsgs) {
        if (displayMsgs.length == 0) {
          displayMsgs.push({
            isMe: msg.me,
            messages: [msg]
          });
        } else {
          if ((!isLastSectionMe && !msg.me) || (isLastSectionMe && msg.me)) {
            displayMsgs[displayMsgs.length - 1].messages.push(msg);
          } else if (
            (!isLastSectionMe && msg.me) ||
            (isLastSectionMe && !msg.me)
          ) {
            displayMsgs.push({
              isMe: msg.me,
              messages: [msg]
            });
          }
        }

        if (msg.me && !isLastSectionMe) isLastSectionMe = true;
        if ((!msg.me && isLastSectionMe) || (!msg.me && !isLastSectionMe))
          isLastSectionMe = false;
      }

      return displayMsgs;
    },
    async clickOnContact(contact) {
      this.failedPoolAttempt = 0

      this.scrollDownChat()
      this.requestContacts(this.detailAccount.number)
        .catch(err => console.log(err))
      await this.poolConversation(contact)

      if (this.poolContactId) clearInterval(this.poolContactId)

      this.poolContactId = setInterval(async () => {
        if (!this.onPoolContacts && this.failedPoolAttempt < this.maxFailedPoolAttempt) {
          this.onPoolContacts = true

          if (this.$route.params.id) {
            await this.syncContacts()
          }

          // detect when on other pages
          if (!this.$route.params.id) clearInterval(this.poolContactId)

          this.onPoolContacts = false
        }

        if (this.failedPoolAttempt >= this.maxFailedPoolAttempt) {
          clearInterval(this.poolContactId)
        }
      }, 8000)
    },
    async poolConversation(contact) {
      if (this.activeConversation) clearInterval(this.activeConversation);

      this.activeConversation = setInterval(async () => {
        const conversation = this.conversations[contact.id];

        if (!this.activeConversationPool) {
          this.activeConversationPool = true;

          const cancelToken = axios.CancelToken;
          const source = cancelToken.source();

          try {
            if (!conversation) {
              // load fresh message
              const m = await Req.post(
                "/chat/pool",
                {
                  number: this.detailAccount.number,
                  remoteJid: contact.id,
                  first: true
                },
                {
                  cancelToken: source.token,
                  timeout: 1000 * 30 // 30s
                }
              );

              if (m.status === 200) {
                const msgs = this.parseMessages(m.data.messages);
                this.$set(this.conversations, contact.id, {
                  messages: msgs,
                  displayMessages: this.parseMessageDisplay(msgs),
                  lastCount: m.data.totalCount
                });
                setTimeout(() => {
                  this.scrollDownChat()
                }, 200)
                
              }

              this.activeConversationPool = false;
            } else {
              // load existing message + new message on server
              const m = await Req.post(
                "/chat/pool",
                {
                  number: this.detailAccount.number,
                  remoteJid: contact.id,
                  first: false,
                  lastCount: conversation.lastCount
                },
                {
                  cancelToken: source.token,
                  timeout: 1000 * 30 // 20s
                }
              );

              if (m.status === 200) {
                const newMsgs = this.parseMessages(m.data.messages);
                const cvsMsgs = conversation.messages.concat(newMsgs);
                const disMsgs = this.parseMessageDisplay(cvsMsgs);

                this.$set(this.conversations, contact.id, {
                  messages: cvsMsgs,
                  displayMessages: disMsgs,
                  lastCount: m.data.totalCount
                });

                setTimeout(() => {
                  this.scrollDownChat();
                }, 200)
              }

              this.activeConversationPool = false;
            }

            if (this.failedPoolAttempt > 0) this.failedPoolAttempt--
          } catch (err) {
            console.log("Error Pool:", err);
            source.cancel("Any active pool canceled.");
            this.activeConversationPool = false;

            this.failedPoolAttempt++
            if (this.failedPoolAttempt >= this.maxFailedPoolAttempt) {
              this.modalErrConnShow = true
            }
            
          }
        }

        if (this.failedPoolAttempt >= this.maxFailedPoolAttempt) {
          clearInterval(this.activeConversation)
        }
      }, 300);

      this.conversationId = contact.id;
      this.conversationTitle = contact.name;
      this.conversationPeerNumber = `+${contact.number}`;

      this.conversationContact = Object.assign({}, contact)
    },
    async sendMessageText({number, jid, text}) {
      // push message into chat window
      const conversation = this.conversations[jid];
      const msgItem = {
        text: text,
        wamsg: {
          info: {
            fromMe: true,
          },
          type: 'text'
        }
      }
      const newMsg = this.parseMessageItem(msgItem, true)
      const cvsMsgs = conversation.messages.concat([newMsg]);
      const disMsgs = this.parseMessageDisplay(cvsMsgs);

      this.$set(this.conversations, jid, {
        messages: cvsMsgs,
        displayMessages: disMsgs,
        lastCount: conversation.lastCount
      });
      
      setTimeout(() => {
        this.scrollDownChat()
      }, 200)

      try {
        await Req.post(
          "/wa/send/text",
          {
            from: number,
            to: jid,
            message: text,
          }
        );

        newMsg.sendingOnTheFly = false
      } catch (err) {
        newMsg.sendingOnTheFly = false
        newMsg.sendingError = true

        this.showMessage("Send Failed", "Send message failed")
      }

      
    },
    async resyncSocket() {
      this.requestSyncSocket = true
      await Req.post(
        "/wa/connection/terminate",
        {
          number: this.detailAccount.number,
        }
      );

      setTimeout(() => {
        this.requestSyncSocket = false
        this.resyncConversation()
      }, 3000)
    },

    async resyncConversation() {
      this.$delete(this.conversations, this.conversationId)
      await this.poolConversation(this.conversationContact)
    },

    isGroupNumber(jid) {
      return jid && jid.split('@')[1] === 'g.us'
    },

    getNameByGroup(groupNumber) {
      if (!this.contacts) return groupNumber 
        ? (groupNumber.split('@')[0] 
          ? groupNumber.split('@')[0].split('-')[0] 
          : 'No Info') 
        : 'No Info'
      const cleanNumber = groupNumber.split('@')[0].split('-')[0]
      const foundContact = this.contacts
        .filter(c => c.jid.split('@')[1] === 's.whatsapp.net')
        .find(c => c.jid.split('@')[0] === cleanNumber)

      return foundContact ? foundContact.contact.name : 'No Info'
    },
    sendText() {
      this.sendMessageText({
        number: this.detailAccount.number,
        jid: this.conversationId,
        text: this.chatInput,
      })
      .then(res => {
        this.sendMessageSignal(this.detailAccount.number, this.conversationId, 'load')
          .catch(err => console.log(err))
      }).catch(err => console.log(err))
      this.chatInput = "";
    },
    async sendImage(images) {
      for (const img of images) {
        const data = new FormData()
        
        data.set('number', this.detailAccount.number)
        data.set('receipentJid', this.conversationId)
        data.append('imageFile', img.file, img.file.name)
        try {
          this.inUploadMedia = true
          setTimeout(() => this.scrollDownChat(), 1000)
          const res = await Req.post(
            "/wa/send/media",
            data,
            {
              headers: {
                'Accept': 'application/json',
                'Content-Type': `multipart/form-data`,
              },
            }
          );

          if (res.status !== 200) {
            this.showMessage("Send Failed", (res.data && res.data.status) ? res.data.status : "Send media failed")
          }

          this.inUploadMedia = false
        } catch (err) {
          this.showMessage("Send Failed", "Send media failed")
          this.inUploadMedia = false
        }
      }      
    },
    showMessage(title, message) {
      this.$set(this.modalMsg, 'show', true)
      this.$set(this.modalMsg, 'title', title)
      this.$set(this.modalMsg, 'message', message)
    },
    async downloadMedia(mid, ext) {
      try {
        const res = await Req.post(
            "/wa/download",
            {
              "number": this.detailAccount.number,
              "messageId": mid,
            },
            {
              headers: {
                'Accept': 'application/json',
                'Content-Type': `multipart/form-data`,
              },
              responseType: 'blob',
            }
          );
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'download' + (ext ? ext : ''));
        document.body.appendChild(link);
        link.click();
      } catch (err) {
        showMessage("Download Failed", "Sorry,  we failed to fetch the file content")
      }
    }
  }
};
</script>
