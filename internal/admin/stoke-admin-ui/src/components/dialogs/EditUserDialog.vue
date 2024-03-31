<template>
  <v-sheet class="d-flex flex-column px-4 pt-4 mb-n5" width="40vw" height="70vh">
    <v-row>
      <v-text-field variant="solo-filled" clearable label="Username" v-model="username" no-resize disabled></v-text-field>
      <v-btn class="ml-5 h-75" variant="tonal" color="info" stacked prepend-icon="mdi-lock-open-variant" disabled> Unlock </v-btn>
      <v-btn class="ml-2 h-75" variant="tonal" color="error" stacked density="compact">
          <span>Change</span>
          <span>Password</span>
      </v-btn>
    </v-row>
    <v-row>
      <v-text-field class="" variant="solo-filled" clearable label="First Name" v-model="fname" @blur="updateCurrentUser"></v-text-field>
      <v-text-field class="ml-4" variant="solo-filled" clearable label="Last Name" v-model="lname" @blur="updateCurrentUser"></v-text-field>
    </v-row>
    <v-row>
      <v-text-field variant="solo-filled" clearable label="Email" v-model="email" no-resize @blur="updateCurrentUser"></v-text-field>
    </v-row>
    <v-row class="mb-5 d-flex flex-grow-1 overflow-auto h-100">
      <GroupList :groups="groups" showFooter :addButton="addGroup">
      </GroupList>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import { useAppStore } from "@/stores/app"

  const store = useAppStore()

  const fname = ref(store.currentUser.fname)
  const lname = ref(store.currentUser.lname)
  const username = ref(store.currentUser.username)
  const email = ref(store.currentUser.email)
  const groups = ref(store.currentGroups)

  function updateCurrentUser() {
    store.$patch({
      currentUser: {
        ...store.currentUser,
        fname: fname.value,
        lname: lname.value,
        email: email.value,
      },
    })
  }

  function addGroup() {
    console.log("components/dialogs/EditUserDialog.vue")
  }
</script>
