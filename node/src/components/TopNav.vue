<template>
  <b-navbar toggleable="md" type="dark" variant="dark">
    <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>

    <b-navbar-brand>
      <b-link :to="{name: 'home'}">WhatDash</b-link>
    </b-navbar-brand>

    <b-collapse is-nav id="nav_collapse">
      <b-navbar-nav>
        <b-nav-item :to="{name: 'reconnect'}">Reconnect Account</b-nav-item>
      </b-navbar-nav>

      <b-navbar-nav class="ml-auto">
        <b-nav-item :to="{name: 'register'}">Register Account</b-nav-item>
      </b-navbar-nav>
    </b-collapse>
  </b-navbar>
</template>
