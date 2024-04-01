<template>
  <div class="d-flex justify-center ma-auto" height="100%">
    <v-sheet width="30vw" height="50vh" elevation="15">
      <v-alert
        class="mb-n16"
        type="error"
        text="Login Failed."
        style="z-index : 100;"
        v-if="loginError"
        @click:close="loginError=false"
        closable
      ></v-alert>
      <v-card class="h-100 w-100">
        <template #title>
          <div class="pt-5 pb-3 d-flex justify-center">
            <StokeIcon size="40"/>
            <div class="d-flex flex-column pl-1">
              <span class="text-h4 font-weight-thin">Stoke</span>
            </div>
          </div>
        </template>
        <v-form @submit.prevent="loginOrShowError" class="mx-15">
          <v-text-field
            class="py-3"
            label="Username"
            prepend-icon="mdi-account"
            :rules="[rules.usernameRequired]"
            v-model="username"
            validate-on="submit"
          > </v-text-field>
          <v-text-field
            class="pb-3"
            label="Password"
            prepend-icon="mdi-key"
            v-model="password"
            :rules="[rules.passwordRequired]"
            :type="showPass ? 'text' : 'password'"
            :append-inner-icon="showPass ? 'mdi-eye' : 'mdi-eye-off'"
            @click:append-inner="showPass = !showPass"
            validate-on="submit"
          > </v-text-field>
          <div class="d-flex justify-center mb-5">
            <v-btn type="submit" variant="tonal" color="info" rounded="lg" elevation="2" density="comfortable" size="large"> Login </v-btn>
          </div>
        </v-form>
      </v-card>
    </v-sheet>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue"
import { useAppStore } from "../stores/app"
import { useRouter } from "vue-router"

const username = ref("")
const password = ref("")
const showPass = ref(false)
const loginError = ref(false)
const rules = {
  usernameRequired: value => !!value || 'Username is required.',
  passwordRequired: value => !!value || 'Password is required.',
}

const store = useAppStore()
const router = useRouter()

async function loginOrShowError(event : Promise<SubmitEvent>) {
  try {
    const results = await event
    if ( ! results.valid ) return
    const resp = await store.login(username.value, password.value)
    router.push("/user")
  } catch (err) {
    console.error(err)
    loginError.value = true
  }
}

</script>

