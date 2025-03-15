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
          :prepend-icon="icons.USER"
          :rules="[require('Username')]"
          v-model="username"
          validate-on="submit"
        > </v-text-field>
        <v-text-field
          class="pb-3"
          label="Password"
          :prepend-icon="icons.PASSWORD"
          v-model="password"
          :rules="[require('Password')]"
          :type="showPass ? 'text' : 'password'"
          :append-inner-icon="showPass ? icons.HIDE : icons.SHOW"
          @click:append-inner="showPass = !showPass"
          validate-on="submit"
        > </v-text-field>
        <div class="d-flex justify-center mb-5">
          <v-btn-group
              rounded="md"
              divided
          >
            <v-btn
              type="submit"
              variant="elevated"
              color="blue-lighten-1"
              density="comfortable"
              size="large"
            > Login </v-btn>
            <v-menu
                location="bottom start"
                open-on-hover
            >
              <template #activator="{ props }">
                <v-btn
                  size="small"
                  :icon="icons.MENU_DOWN"
                  color="blue-lighten-1"
                  v-bind="props"
                  variant="elevated"
                  density="comfortable"
                  > </v-btn>
              </template>
              <v-list
                variant="tonal"
                elevation="5"
                density="compact"
              >
                <v-list-item
                    density="compact"
                    v-for="(prov, i) in openIDProviders()"
                    :value="i"
                    :key="i"
                    @click="handleOIDCLogin(prov)"
                >
                  <v-list-item-title>
                    <v-icon
                      color="orange-lighten-1"
                      size="large"
                      class="mr-3"
                      :icon="icons.OPENID"
                    ></v-icon>
                    <span class="text-button">{{ prov.name }}</span>
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-btn-group>
        </div>
      </v-form>
    </v-card>
  </div>
</template>

<script setup lang="ts">
import icons from "../util/icons"
import { ref, onMounted } from "vue"
import { useAppStore } from "../stores/app"
import { useRouter } from "vue-router"
import { require } from "../util/rules"

const username = ref("")
const password = ref("")
const showPass = ref(false)
const loginError = ref(false)
const loading = ref(false)
const formValid = ref(false)

const store = useAppStore()
const router = useRouter()

function openIDProviders() {
  return store.availableProviders.filter((p) => p.provider_type == "OIDC")
}

// TODO bring in provider type
function handleOIDCLogin(prov) {
  let u = new URL(store.api_url + "/oidc/" + prov.name + "?xfer=window&next=" + window.location.origin, window.location.origin)
  console.log(u)
  addEventListener("message", async (event: MessageEvent) => {
    if ( event.origin !== u.origin ) return;
    let result = JSON.parse(event.data)
    try {
      if ( result.id_token && result.access_code ) {
        await store.login(result.id_token, result.access_code, prov.name, () => router.push("/user"))
      }
    } catch (err) {
      console.error(err)
      loginError.value = true
    }
  })
  window.open(u.toString(), prov.name + " Login", "popup")
}

async function loginOrShowError(event : Promise<SubmitEvent>) {
  loading.value = true
  try {
    await event
    if ( ! formValid.value ) return
    await store.login(username.value, password.value, "LOCAL", () => router.push("/user"))
  } catch (err) {
    console.error(err)
    loginError.value = true
  }
  loading.value = false
}


onMounted(() => {
  store.fetchAvailableProviders()
})

</script>

