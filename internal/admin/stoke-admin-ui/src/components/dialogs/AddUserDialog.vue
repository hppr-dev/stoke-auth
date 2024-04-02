<template>
  <v-sheet class="d-flex flex-column px-4 pt-4 mb-n5" width="40vw" height="40vh">
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Username"
        :rules="[require('Username')]"
        v-model="username"
      ></v-text-field>
    </v-row>
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="First Name"
        v-model="fname"
        :rules="[require('First Name')]"
        @blur="updateScratchUser"
      ></v-text-field>
      <v-text-field
        class="ml-4"
        variant="solo-filled"
        label="Last Name"
        v-model="lname"
        :rules="[require('Last Name')]"
        @blur="updateScratchUser"
      ></v-text-field>
    </v-row>
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Email"
        v-model="email"
        :rules="[require('Email')]"
        @blur="updateScratchUser"
      ></v-text-field>
    </v-row>
    <v-row class="pb-3">
      <v-text-field
        label="Password"
        v-model="password"
        :rules="[require('Password')]"
        @blur="updateScratchUser"
      > </v-text-field>
      <v-text-field
        class="ml-3"
        type="password"
        label="Repeat Password"
        v-model="passwordCheck"
        :rules="[require('Repeat password')]"
      > </v-text-field>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import { useAppStore } from "../../stores/app"
  import { UserWithCreds } from "../../util/entityTypes"
  import { require } from "../../util/rules"

  const store = useAppStore()

  const fname = ref("")
  const lname = ref("")
  const username = ref("")
  const email = ref("")
  const password = ref("")
  const passwordCheck = ref("")

  function updateScratchUser() {
    store.$patch({
      scratchUser: {
        ...store.scratchUser,
        fname: fname.value,
        lname: lname.value,
        username: username.value,
        email: email.value,
        password: password.value,
        provider: "local",
      } as UserWithCreds,
    })
  }

</script>
