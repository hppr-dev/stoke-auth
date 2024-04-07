<template>
  <div class="d-flex justify-center ma-auto pb-10" height="100%">
    <v-alert
      width="30vw"
      type="error"
      text="Login Failed."
      style="z-index : 100;"
      v-if="loginError"
      @click:close="loginError=false"
      closable
      position="absolute"
    ></v-alert>
    <v-card width="30vw" elevation="15" class="py-10" :loading="loading">
      <template #title>
        <div class="pt-5 pb-3 d-flex justify-center">
          <StokeIcon size="40"/>
          <div class="d-flex flex-column pl-1">
            <span class="text-h4 font-weight-thin">Stoke</span>
          </div>
        </div>
      </template>
      <v-form v-model="formValid" @submit.prevent="loginOrShowError" class="mx-15">
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
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useAppStore } from "../stores/app"
import { useRouter } from "vue-router"

const username = ref("")
const password = ref("")
const showPass = ref(false)
const loginError = ref(false)
const loading = ref(false)
const formValid = ref(false)

const rules = {
  usernameRequired: ( value : string ) => !!value || 'Username is required.',
  passwordRequired: ( value : string ) => !!value || 'Password is required.',
}

const store = useAppStore()
const router = useRouter()

async function loginOrShowError(event : Promise<SubmitEvent>) {
  loading.value = true
  try {
    await event
    if ( ! formValid.value ) return
    await store.login(username.value, password.value, () => router.push("/user"))
  } catch (err) {
    console.error(err)
    loginError.value = true
  }
  loading.value = false
}

</script>

