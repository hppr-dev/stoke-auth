<template>
  <Viewport noScroll>
    <v-row class="h-screen pb-10 ml-0">
      <v-col class="ml-n3 h-100" lg="7" md="6" sm="12">
        <GroupList
          showSearch
          showFooter
          addButton
          :groups="store.allGroups"
          :rowClick="setCurrentGroup"
          :rowProps="highlightSelected"
        />
      </v-col>
      <v-col class="ml-n3 h-100" lg="5" md="6" sm="12">
        <GroupInfo />
      </v-col>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useAppStore } from '../stores/app'
  import { Group } from '../util/entityTypes'

  const store = useAppStore()

  async function setCurrentGroup(_ : PointerEvent, { item } : { item : Group } ) {
    if ( store.currentGroup.id && store.currentGroup.id === item.id ) {
      store.resetCurrentGroup()
      return
    }
    await store.fetchClaimsForGroup(item.id)
    store.$patch({
      currentGroup: item,
    })
  }

  function highlightSelected({ item } : { item : Group }) {
    if ( item.id === store.currentGroup.id ) {
      return {
        class : "bg-grey-lighten-1",
      }
    }
  }

  onMounted(store.fetchAllGroups)
</script>
