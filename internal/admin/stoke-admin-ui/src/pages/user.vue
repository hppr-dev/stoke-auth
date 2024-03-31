<template>
  <Viewport>
    <v-row class="h-screen pb-10 ml-0">
      <v-col class="ml-n3 h-100" md="8" sm="5">
        <UserList showSearch showFooter :addButton="addUser"/>
      </v-col>
      <v-col class="h-100" md="4" offset-sm="1" offset-md="0" sm="5">
        <UserInfo />
      </v-col>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  function addUser() {
    console.log("pages/user.vue")
  }
</script>
