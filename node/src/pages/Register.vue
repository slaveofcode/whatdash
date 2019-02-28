<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
      <b-row class="form-whatsapp justify-content-md-center">
        <b-col md="6">
          <b-form @submit="requestQrCode">
            <b-form-group
              id="group-phone"
              label="WhatsApp Phone Number:"
              label-for="phone"
              description="This action will add your whatsapp session to start communicating">
              <b-form-input
                id="phone"
                type="text"
                v-model="form.phone"
                required
                placeholder="Enter phone with country code without ( + ) e.g: 6285765746576" />
            </b-form-group>
            <b-button block type="submit" variant="primary">{{form.submitText}}</b-button>
          </b-form>
        </b-col>
      </b-row>
       <b-row v-show="qrcode" class="qr-whatsapp justify-content-md-center">
         <b-col md="3">
           <p>Scan image below with your WhatsApp Scanner</p>
           <qr-code :text="qrcode"></qr-code>
         </b-col>
       </b-row>
    </b-container>


    <b-modal v-model="modalSuccessShow" centered title="Register new account">
      Successfully register new account, will be redirected to list accounts..
      <div slot="modal-footer" class="w-100">
        <b-button size="sm" class="float-right" variant="primary" @click="modalSuccessShow=false">Close</b-button>
      </div>
    </b-modal>
    <b-modal v-model="modalFailedShow" centered title="Register new account">
      Register new account failed
      <div slot="modal-footer" class="w-100">
        <b-button size="sm" class="float-right" variant="primary" @click="modalFailedShow=false">Close</b-button>
      </div>
    </b-modal>
  </div>
</template>

<style>
.form-whatsapp {
  margin: 50px 0 30px;
}
.qr-whatsapp {
  text-align: center;
  font-size: 14px;
}
</style>


<script>
import PageTitle from "../components/PageTitle.vue";
import TopNav from "../components/TopNav.vue";
import Req from '../req'

export default {
  components: {
    PageTitle,
    TopNav
  },
  data() {
    return {
      pageTitle: "Register New Account",
      qrcode: '',
      intervalCheckId: null,
      onCheckInterval: false,
      checkIntervalCount: 0,
      form: {
        phone: '',
        submitText: 'Register Number',
        submitEnable: true,
      },
      modalSuccessShow: false,
      modalFailedShow: false,
    };
  },
  methods: {
    async isSessionRegistered(phone) {
      const res = await Req.post('/wa/session/check', {
        number: phone
      })
      if (res.status != 200) return false
      const data = res.data

      return (data.status === 'registered')
    },
    async createSession(phone){
      const res = await Req.post('/wa/session/create', {
        number: phone
      })
      if (res.status != 200) return false
      const data = res.data

      return data.qr
    },
    async requestQrCode() {
      if (this.form.submitEnable) {
        const phone = this.form.phone
        const tmpText = this.form.submitText
        
        this.form.submitText = 'Loading...'
        this.form.submitEnable = false
        
        const registered = await this.isSessionRegistered(phone)
        if (!registered) {
          const qrCodeStr = await this.createSession(phone)
          this.qrcode = qrCodeStr
          this.form.submitText = tmpText

          if (this.intervalCheckId) clearInterval(this.intervalCheckId)
          this.intervalCheckId = setInterval(async () => {
            this.onCheckInterval = true
            this.checkIntervalCount++
            const registered = await this.isSessionRegistered(phone)
            if (registered) {
              // show modal
              this.modalSuccessShow = true
              clearInterval(this.intervalCheckId)
              setTimeout(() => {
                this.$router.push({ path: 'reconnect' })
              }, 4000)
            }
            this.onCheckInterval = false

            // max attempt to checking the register status
            if (this.checkIntervalCount >= 20) {
              clearInterval(this.intervalCheckId)
              this.modalFailedShow = true
              this.qrcode = ''
            }
          }, 1000)
        } else {
          // show modal
        }
      }
    }
  }
};
</script>