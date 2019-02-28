<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
      <div class="chat-container">
        <div class="section-contacts">
          <ContactItem
            v-for="contact in contactChats"
            :key="contact.id"
            v-bind:contact="contact"
            @click.native="clickOnContact(contact).catch(err => console.log(err))"
          ></ContactItem>
        </div>
        <div class="section-messages">
          <div class="section-chat-header" v-show="conversationTitle">
            <p>{{conversationTitle}}</p>
            <div v-b-tooltip title="Sync socket and messages" class="resync-button socket" @click="resyncSocket(conversationId).catch(err => console.log(err))">
              <i class="fas fa-sync-alt" :class="[requestSyncSocket ? 'fa-spin' : '']"></i>
            </div>
            <div v-b-tooltip title="Sync messages" class="resync-button messages" @click="resyncConversation(conversationId).catch(err => console.log(err))">
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
              <MessageItem v-for="(msg, idx) in conversation.displayMessages" :key="idx" :msg="msg"></MessageItem>
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
          </div>
        </div>
      </div>
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
  min-width: 230px;
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
import Req from "../req";
import { axios } from "../req";

export default {
  components: {
    TopNav,
    PageTitle,
    MessageItem,
    ContactItem
  },
  data() {
    return {
      pageTitle: "Active Chat",
      detailAccount: null,
      contacts: [],
      chatHistory: [],
      contactChats: [],
      chatInput: null,
      conversationId: null,
      conversationTitle: null,
      conversations: {},
      activeConversation: null,
      activeConversationPool: false,
      requestSyncSocket: false,
      poolContactId: null,
      onPoolContacts: false,
    };
  },
  watch: {
    $route: "initPage"
  },
  created() {
    this.initPage();
  },
  methods: {
    onKeydown(evt) {
      // detecting ctrl+enter
      if (evt.keyCode === 13 && evt.ctrlKey) {
        // send message
        this.sendMessageText({
          number: this.detailAccount.number,
          jid: this.conversationId,
          text: this.chatInput,
        }).catch(err => console.log(err))
        this.chatInput = "";
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
      
      this.poolContactId = setInterval(async () => {
        if (!this.onPoolContacts) {
          this.onPoolContacts = true
          await this.syncContacts()
          this.onPoolContacts = false
        }        
      }, 5000)

      this.pageTitle = `Active Chat on [${this.detailAccount.number}]`;
    },
    async syncContacts() {
      this.detailAccount = await this.loadAccountDetail(this.$route.params.id);
      
      const contacts = await this.loadContacts(this.detailAccount.number);
      this.contacts = contacts ? contacts : []

      if (!contacts) {
        await this.requestContacts(this.detailAccount.number)
      }
      
      const history = await this.loadHistory(this.detailAccount.number);
      this.chatHistory = history ? history : []

      this.contactChats = this.parseHistoryWithContact(
        this.contacts,
        this.chatHistory
      );
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
      return c.status === 200 ? c.data : [];
    },
    async loadHistory(number) {
      const c = await Req.post("/chat/history", { number });
      return c.status === 200 ? c.data : [];
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
      const waMsg = item.wamsg;
      const waInfo = waMsg.info;
      if (waMsg.type === "text") {
        return {
          id: waInfo.id,
          type: 'text',
          msg: item.text,
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
          sendingOnTheFly: onTheFly,
        };
      } else if (waMsg.type === 'image') {
        return {
          id: waInfo.id,
          type: 'image',
          caption: item.caption,
          image: {
            thumb: item.thumb,
            content: item.content,
          },
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
          sendingOnTheFly: onTheFly,
        };
      }

      return {
        id: waInfo.id,
        type: 'unknown',
        msg: "*unknown message, see on device*",
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

      const displayMsgs = [];
      for (const msg of messages) {
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
      this.scrollDownChat()
      await this.poolConversation(contact)
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
                  timeout: 1000 * 30 // 20s
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
          } catch (err) {
            console.log("Error Pool:", err);
            source.cancel("Any active pool canceled.");
            this.activeConversationPool = false;
          }
        }
      }, 300);

      // check if chat container already created
      // hide existing container if exist and show selected chat
      this.conversationId = contact.id;
      this.conversationTitle = contact.name;
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

      await Req.post(
        "/wa/send/text",
        {
          from: number,
          to: jid,
          message: text,
        }
      );

      newMsg.sendingOnTheFly = false
    },
    async resyncSocket(jid) {
      this.requestSyncSocket = true
      await Req.post(
        "/wa/connection/terminate",
        {
          number: this.detailAccount.number,
        }
      );

      setTimeout(() => {
        this.requestSyncSocket = false
        this.resyncConversation(jid)
      }, 3000)
    },

    async resyncConversation(jid) {
      this.$delete(this.conversations, jid)
      await poolConversation(jid)
    }
  }
};
</script>
