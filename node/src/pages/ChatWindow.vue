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
            @click.native="poolConversation(contact).catch(err => console.log(err))"
          ></ContactItem>
        </div>
        <div class="section-messages">
          <div class="section-chat-header" v-show="conversationTitle">
            <p>{{conversationTitle}}</p>
          </div>
          <div class="section-chat">
            <span class="empty-message" :class="[conversationTitle ? 'hide' : '']">Select some chat to start conversation.</span>
            <span class="empty-message loading" v-show="conversationId && !conversations[conversationId]">Loading messages...</span>
            <div class="chat-display" v-show="key === conversationId" v-for="(conversation, key) in conversations" :key="key">
              <MessageItem
                v-for="(msg, idx) in conversation.displayMessages"
                :key="idx"
                v-bind:me="msg.isMe"
                v-bind:messages="msg.messages"
              ></MessageItem>
            </div>
          </div>
          <div class="section-input">
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
  width: 230px;
  height: 100%;
  overflow-y: auto;
  border-bottom: 1px solid #dad8d8;
  border-top: 1px solid #dad8d8;
}

.section-messages {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
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
        console.log("Sending message");
        this.chatInput = "";
      }
    },
    async initPage() {
      this.detailAccount = await this.loadAccountDetail(this.$route.params.id);
      this.contacts = await this.loadContacts(this.detailAccount.number);
      this.chatHistory = await this.loadHistory(this.detailAccount.number);
      this.contactChats = this.parseHistoryWithContact(
        this.contacts,
        this.chatHistory
      );

      this.pageTitle = `Active Chat on [${this.detailAccount.number}]`
    },
    async loadAccountDetail(accId) {
      const acc = await Req.get(`/account/detail/${accId}`);
      return acc.status === 200 ? acc.data : null;
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
    parseMessageItem(item) {
      const waMsg = item.wamsg;
      const waInfo = waMsg.info;
      if (waMsg.type === 'text') {
        return {
          id: waInfo.id,
          msg: item.text,
          me: waInfo.fromMe,
          stat: waInfo.msgStatus,
        };
      }

      return {
        id: waInfo.id,
        msg: 'non text message',
        me: waInfo.fromMe,
        stat: waInfo.msgStatus,
      };
    },
    parseMessages(messages) {
      if (!messages || messages.length <= 0) return []
      return messages.reverse().map(item => this.parseMessageItem(item))
    },
    parseMessageDisplay(messages) {
      if (!messages || messages.length <= 0) return []
      let isSectionMe = null

      const displayMsgs = []
      for (const msg of messages) {
        if (msg.me && !isSectionMe) isSectionMe = true
        if (!msg.me && isSectionMe) isSectionMe = false

        if (displayMsgs.length == 0) {
          displayMsgs.push({
            isMe: msg.me,
            messages: [msg],
          })
        } else {
          if (!isSectionMe) {
            displayMsgs.push({
              isMe: msg.me,
              messages: [msg],
            })
          } else if (isSectionMe) {
            
            displayMsgs[displayMsgs.length-1].messages.push(msg)
          }
        }
      }

      return displayMsgs
    },
    async poolConversation(contact) {
      if (this.activeConversation) clearInterval(this.activeConversation);

      this.activeConversation = setInterval(async () => {
        const conversation = this.conversations[contact.id];

        if (!this.activeConversationPool) {
          this.activeConversationPool = true;
          
          const cancelToken = axios.CancelToken
          const source = cancelToken.source()

          try {
            if (!conversation) {
              // load fresh message
              const m = await Req.post("/chat/pool", {
                number: this.detailAccount.number,
                remoteJid: contact.id,
                first: true
              }, { 
                cancelToken: source.token, 
                timeout: 1000 * 30, // 20s
              });

              if (m.status === 200) {
                const msgs = this.parseMessages(m.data.messages)
                this.$set(this.conversations, contact.id, {
                  messages: msgs,
                  displayMessages: this.parseMessageDisplay(msgs),
                  lastCount: m.data.totalCount
                })

              }

              this.activeConversationPool = false
            } else {
              // load existing message + new message on server
              const m = await Req.post("/chat/pool", {
                number: this.detailAccount.number,
                remoteJid: contact.id,
                first: false,
                lastCount: conversation.lastCount
              }, { 
                cancelToken: source.token,
                timeout: 1000 * 30, // 20s
              });

              if (m.status === 200) {
                const newMsgs = this.parseMessages(m.data.messages)
                const cvsMsgs = conversation.messages.concat(newMsgs)
                const disMsgs = conversation.displayMessages.concat(newMsgs)

                this.$set(this.conversations, contact.id, {
                  messages: cvsMsgs,
                  displayMessages: disMsgs,
                  lastCount: m.data.totalCount,
                })
              }

              this.activeConversationPool = false
            }
          } catch (err) {
            console.log("Error Pool:", err);
            source.cancel('Any active pool canceled.')
            this.activeConversationPool = false
          }
        }
      }, 300);

      // check if chat container already created
      // hide existing container if exist and show selected chat
      this.conversationId = contact.id
      this.conversationTitle = contact.name
    }
  }
};
</script>
