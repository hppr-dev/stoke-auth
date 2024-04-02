<template>
  <v-app>
    <NavDrawer v-if="store.authenticated"/>
    <router-view />
    <AppFooter />
  </v-app>
</template>

<script lang="ts" setup>
  import { onMounted } from "vue"
  import { useAppStore } from "./stores/app"
  import { useRouter } from "vue-router"

  const store = useAppStore()
  const router = useRouter()

  router.beforeEach( (to, _) => {
    if ( !store.authenticated && to.name != "/" ) {
      return { name: "/" }
    }
  })

  router.afterEach(store.resetSelections)

  onMounted(() => {
    let token = sessionStorage.getItem("token")
    let username = sessionStorage.getItem("username")

    if (
      token !== "undefined" && token !== null && token !== "" && token &&
      username !== "undefined" && username !== null && username !== "" && username
    ) {
      store.$patch({
        token: token as string,
        username : username as string,
      })
    } else {
      sessionStorage.setItem("token", "")
      sessionStorage.setItem("username", "")
    }
  })

</script>
