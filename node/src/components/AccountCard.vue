<template>
  <b-card bg-variant="secondary" text-variant="white" :header="account.number" tag="article" class="account-card mb-2 text-center">
    <b-row>
      <b-col>
        <p><small class="text-warning">{{account.jid}}</small></p>
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
        <b-button href="#" size="m" variant="warning">Go Chat</b-button>
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
  font-size: 15px;
  padding: 0;
  margin: 0;
}
</style>

<script>
export default {
  props: ['account']
}
</script>
