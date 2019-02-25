<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
      <div class="chat-container">
        <div class="section-contacts">
          <div class="contact-item" v-for="chat in contactChats" :key="chat.id">
            <p>
              <span class="item-title">{{chat.name}}</span>
              <span class="item-count">({{chat.msgCount}})</span>
            </p>
            <span class="item-time"><timeago :datetime="chat.time" :auto-update="60"></timeago></span>
          </div>
        </div>
        <div class="section-messages">
          <div class="section-chat">
            <span class="empty-message">Select some chat to start conversation.</span>
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

.section-contacts .contact-item {
  font-size: 13px;
  padding: 10px 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}

.section-contacts .contact-item:hover {
  background: rgba(0, 123, 255, .1);
  border-bottom: 1px solid rgba(0, 123, 255, .2);
}

.section-contacts .contact-item p {
  margin: 0;
  padding: 0;
  font-weight: bold;
  color: #333;
}

.section-contacts .contact-item .item-time {
  font-size: 11px;
  color: #075419;
}

.section-messages {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
}

.section-chat {
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
  background: #eee;
}

.section-chat .empty-message {
  font-size: 14px;
  color: #aaa;
  display: flex;
  justify-content: center;
  padding-top: 25vh;
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
import Req from "../req";

export default {
  components: {
    TopNav,
    PageTitle,

  },
  data() {
    return {
      pageTitle: 'Chat',
      detailAccount: null,
      contacts: [],
      chatHistory: [],
      contactChats: [],
      chatInput: null,
    }
  },
  watch: {
    '$route': 'initPage',
  },
  created() {
    this.initPage()
  },
  methods: {
    onKeydown(evt) {
      // detecting ctrl+enter
      if (evt.keyCode === 13 && evt.ctrlKey) {
        // send message
        console.log('Sending message')
        this.chatInput = ''
      }
    },
    async initPage() {
      this.detailAccount = await this.loadAccountDetail(this.$route.params.id)
      this.contacts = await this.loadContacts(this.detailAccount.number)
      this.chatHistory = await this.loadHistory(this.detailAccount.number)
      this.contactChats = this.parseHistoryWithContact(this.contacts, this.chatHistory)
    },
    async loadAccountDetail(accId) {
      const acc = await Req.get(`/account/detail/${accId}`)
      return (acc.status === 200) ? acc.data : null
    },
    async loadContacts(number) {
      const c = await Req.post('/wa/contact/list', { number, })
      return (c.status === 200) ? c.data : []
    },
    async loadHistory(number) {
      const c = await Req.post('/chat/history', { number, })
      return (c.status === 200) ? c.data : []
    },
    parseHistoryWithContact(contacts, chatHistory) {
      const parsed = []
      for (const history of chatHistory) {
        const foundContact = contacts.find(c => c.jid == history.wa.jid)
        let contactName = history.wa.jid.split('@')[0]
        
        if (foundContact) {
          contactName = !!foundContact.contact.name 
            ? foundContact.contact.name 
            : (foundContact.jid.split('@')[0])
        }

        parsed.push({
          id: history.jid,
          name: contactName,
          time: new Date(history.lastChatTime),
          msgCount: history.msgCount,
        })
      }

      return parsed
    },
    async loadMessages() {},
    async loadContactMessages() {},
    async sendMessage(){},
    async poolMessage(){},
  }
}
</script>
