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
      <v-list-item v-if="store.userAccess !== '' " link title="Users" @click="go('/user')" :prepend-icon="icons.USER"></v-list-item>
      <v-list-item v-if="store.groupAccess !== '' " link title="Groups" @click="go('/group')" :prepend-icon="icons.GROUP"></v-list-item>
      <v-list-item v-if="store.claimsAccess !== '' " link title="Claims" @click="go('/claim')" :prepend-icon="icons.CLAIM"></v-list-item>
      <v-list-item v-if="store.monitorAccess" link title="Monitor" @click="go('/monitor')" :prepend-icon="icons.MONITOR"></v-list-item>
    </v-list>

    <template #append>
      <v-list>
        <v-list-item>
          <template #prepend>
            <v-icon :icon="icons.LOGOUT" color="info" @click="store.logout"></v-icon>
          </template>
          <template #title>
            <span class="text-h7 text-weight-thin"> {{ store.username }} </span>
          </template>
        </v-list-item>
        <v-list-item>
          <template #prepend>
            <v-icon :icon="pinned? icons.PIN_OFF : icons.PIN_ON" size="small" @click="pinned = !pinned"></v-icon>
          </template>
        </v-list-item>
      </v-list>
    </template>

  </v-navigation-drawer>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue"
import { useRouter } from "vue-router"
import { useAppStore } from "../stores/app"
import icons from '../util/icons'

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

onMounted(async () => {
  await store.fetchCapabilites()
})

</script>
