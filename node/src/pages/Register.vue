<template>
  <div>
    <TopNav></TopNav>
    <b-container>
      <PageTitle :text="pageTitle"></PageTitle>
    </b-container>
  </div>
</template>

<script>
import PageTitle from "../components/PageTitle.vue";
import TopNav from "../components/TopNav.vue";
export default {
  components: {
    PageTitle,
    TopNav
  },
  data() {
    return {
      pageTitle: "Register New Account"
    };
  },
};
</script>