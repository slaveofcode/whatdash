<template>
  <div class="main-title">
    <div class="row">
      <div class="col">
        <h1><span v-if="conversationWindow" class="fas fa-phone-volume"></span>{{text}}</h1>
        <p v-if="conversationWindow"><span class="phone">+{{number}}</span></p>
      </div>
      <div class="address-sync col col-md-1">
        <button v-if="conversationWindow" class="btn btn-md btn-info" @click="syncContact(number)" v-b-tooltip title="Sync contacts">
          <span class="fas" :class="[requestContacts ? 'fa-sync-alt fa-spin' : 'fa-address-card']"></span>
        </button>
      </div>
    </div>
  </div>
</template>

<style>
.main-title {
  padding: 20px 0;
  margin: 0 0 20px 0;
  color: #515151;
  border-bottom: 1px solid #c5c5c5;
}

.main-title h1 {
  font-size: 27px;
}
.main-title p {
  margin-bottom: 0;
  margin-left: 17px;
}
.main-title p span.phone {
  margin-left: .5em;
  font-size: 15px;
  color: #17a2b8;
  font-weight: bold
}

.main-title h1 span.fas {
  margin-right: 5px;
}

.address-sync .fas {
  font-size: 25px;
}
</style>


<script>
import Req from "../req";

export default {
  props: ['conversationWindow', 'text', 'number'],
  data() {
    return {
      requestContacts: false,
    }
  },
  methods: {
    async syncContact(number) {
      if (!this.requestContacts) {
        this.requestContacts = true
        await Req.post('/wa/contact/load', {
          number,
          reloadSocket: true,
        });

        this.requestContacts = false
      }
    }
  }
}
</script>