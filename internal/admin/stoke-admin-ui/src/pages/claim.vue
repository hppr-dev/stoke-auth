<template>
  <Viewport>
    <v-row class="h-screen pb-10 ml-0">
      <v-col class="ml-n3 h-100" lg="7" md="6" sm="12">
        <ClaimList
          showSearch
          showFooter
          addButton
          :claims="store.allClaims"
          :rowClick="setCurrentClaim"
          :rowProps="highlightSelected"
        />
      </v-col>
      <v-col class="ml-n3 h-100" lg="5" md="6" sm="12">
        <ClaimInfo />
      </v-col>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useAppStore } from '../stores/app'
  import { Claim } from '../util/entityTypes'

  const store = useAppStore()

  async function setCurrentClaim(_ : PointerEvent, { item } : { item : Claim } ) {
    if ( store.currentClaim.id && store.currentClaim.id === item.id ) {
      store.resetCurrentClaim()
      return
    }
    store.$patch({
      currentClaim: item,
    })
  }

  function highlightSelected({ item } : { item : Claim }) {
    if ( item.id === store.currentClaim.id ) {
      return {
        class : "bg-grey-lighten-1",
      }
    }
  }

  onMounted(store.fetchAllClaims)
</script>
