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
    let refresh = sessionStorage.getItem("refresh")
    let username = sessionStorage.getItem("username")

    if (
      token !== "undefined" && token !== null && token !== "" && token &&
      refresh !== "undefined" && refresh !== null && refresh !== "" && refresh &&
      username !== "undefined" && username !== null && username !== "" && username
    ) {
      store.$patch({
        token: token as string,
        refreshToken: refresh as string,
        username : username as string,
      })
      store.scheduleRefresh()
    } else {
      sessionStorage.setItem("token", "")
      sessionStorage.setItem("refresh", "")
      sessionStorage.setItem("username", "")
    }
  })

</script>
