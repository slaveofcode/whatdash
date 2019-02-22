<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
      <b-row class="d-flex justify-content-around align-content-center flex-wrap">
        <AccountCard v-for="account in accounts" :key="account.number" :account="account"></AccountCard>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import PageTitle from "../components/PageTitle.vue";
import TopNav from "../components/TopNav.vue";
import AccountCard from "../components/AccountCard.vue";
import Req from '../req'

export default {
  components: {
    PageTitle,
    TopNav,
    AccountCard,
  },
  data() {
    return {
      pageTitle: "Reconnect Account",
      loading: false,
      accounts: null,
      error: null,
    };
  },
  watch: {
    '$route': 'loadAccounts',
  },
  created() {
    this.loadAccounts()
  },
  methods: {
    async loadAccounts() {
      this.error = this.accounts = null
      this.loading = true

      const res = await Req.get('/connected-accounts')
      if (res.status === 200) {
        this.accounts = res.data
      } else {
        this.error = 'Couldn\'t load connected accounts'
      }

      this.loading = false
    }
  }
};
</script>