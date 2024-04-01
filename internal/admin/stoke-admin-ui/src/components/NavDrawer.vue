<template>
  <v-navigation-drawer class="fill-height" width="225" :rail="true" :expand-on-hover="!pinned" mobile-breakpoint="md" >
    <v-list>
      <v-list-item>
        <template #prepend>
          <StokeIcon class="ml-n1" size="x-large" />
        </template>
        <template #title>
          <div class="text-center w-50">
            <span class="text-h5 text-weight-thin">Stoke</span>
          </div>
        </template>
        <template #subtitle>
          <span class="text-caption">Auth by hppr.dev</span>
        </template>
      </v-list-item>
    </v-list>

    <v-divider></v-divider>

    <v-list>
      <v-list-item link title="Monitor" @click="go('/monitor')" prepend-icon="mdi-chart-line"></v-list-item>
      <v-list-item link title="Users" @click="go('/user')" prepend-icon="mdi-account"></v-list-item>
      <v-list-item link title="Groups" @click="go('/group')" prepend-icon="mdi-account-multiple"></v-list-item>
      <v-list-item link title="Keys" @click="go('/key')" prepend-icon="mdi-key-chain"></v-list-item>
    </v-list>

    <template #append>
      <v-list>
        <v-list-item>
          <template #prepend>
            <v-icon icon="mdi-logout" color="info" @click="logoutAndReturn"></v-icon>
          </template>
          <template #title>
            <span class="text-h7 text-weight-thin"> {{ store.username }} </span>
          </template>
        </v-list-item>
        <v-list-item>
          <template #prepend>
            <v-icon :icon="pinned? 'mdi-pin-off' : 'mdi-pin'" size="small" @click="pinned = !pinned"></v-icon>
          </template>
        </v-list-item>
      </v-list>
    </template>

  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { useAppStore } from "../stores/app"

const pinned = ref(false)
const router = useRouter()
const store = useAppStore()

function go(path : string) {
  if (store.authenticated) {
    router.push(path)
  } else {
    router.push("/")
  }
}

function logoutAndReturn() {
  store.logout()
  router.push("/")
}

</script>
