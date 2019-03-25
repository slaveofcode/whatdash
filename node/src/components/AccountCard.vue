<template>
  <b-card bg-variant="secondary" text-variant="white" :header="account.number" tag="article" class="account-card mb-2 text-center">
    <b-row>
      <b-col>
        <p><small class="text-success">{{account.jid}}</small></p>
        <p class="info">
          <strong>Logged in:</strong> <timeago :datetime="account.createdAt" :auto-update="60"></timeago>
        </p>
        <p class="info">
          <strong>Last Login:</strong> <timeago :datetime="account.updatedAt" :auto-update="60"></timeago>
        </p>
      </b-col>
    </b-row>
    <br>
    <b-row>
      <b-col>
        <b-button :to="{name: 'chat', params: { id: account.id } }" size="m" variant="warning">Go Chat</b-button>
      </b-col>
    </b-row>
    <b-row>
      <b-col>
        <b-button 
          class="float-right" 
          size="sm" 
          variant="secondary" 
          @click="logout(account.number)">
          <i class="fas fa-power-off"></i>
        </b-button>
      </b-col>
    </b-row>
  </b-card>
</template>

<style>
.card-title {
  color: #c97903
}
.account-card {
  max-width: 16rem;
}
p.info {
  font-size: 14px;
  padding: 0;
  margin: 0;
}
</style>

<script>
import Req from "../req";

export default {
  props: ['account'],
  data() {
    return {
      logoutAttempts: {},
    }
  },
  methods: {
    async logout(number) {
      if (!this.logoutAttempts[number]) {
        this.logoutAttempts[number] = 1
      } else {
        this.logoutAttempts[number] = this.logoutAttempts[number] + 1
      }

      const forceLogout = this.logoutAttempts[number] === 3
      await this.requestLogoutAccount(number, forceLogout)
      this.$router.go() 
    },
    async requestLogoutAccount(number, force = false) {
      const c = await Req.post("/wa/session/destroy", { number, force, });
      const destroyed = (c.status === 200)
      if (!destroyed) return false
      return c.data.destroyed
    }
  }
}
</script>
