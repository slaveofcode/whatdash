<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
      <div class="chat-container">
        <div class="section-contacts">

        </div>
        <div class="section-messages">
          <div class="section-chat"></div>
          <div class="section-input"></div>
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
  background: #333333;
}

.section-contacts {
  width: 230px;
  background: blue;
  height: 100%;
}

.section-messages {
  display: flex;
  flex-direction: column;
  background: red;
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
}

.section-chat {
  flex-grow: 1;
  flex-shrink: 1;
  flex-basis: auto;
  background: yellow;
}

.section-input {
  width: 100%;
  height: 70px;
  background: green;
  align-self: flex-end;
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
      pageTitle: 'Chat'
    }
  },
  watch: {
    '$route': 'initPage',
  },
  created() {
    this.initPage()
  },
  methods: {
    async initPage() {
      const detailAccount = await this.loadAccountDetail(this.$route.params.id)
      const contacts = await this.loadContacts(detailAccount.number)
      const history = await this.loadHistory(detailAccount.number)
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
    async loadMessages() {},
    async loadContactMessages() {},
    async sendMessage(){},
    async poolMessage(){},
  }
}
</script>
