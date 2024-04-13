<template>
  <v-card class="d-flex flex-column w-100" height="65%" :subtitle="store.currentGroup.description">
    <template #title>
      <div v-if="store.currentGroup.name">
        <v-icon class="mr-2" size="small" :icon="icons.GROUP"></v-icon>
        <v-icon v-if="store.currentLinks.length > 0" color="success" class="mr-2" size="small" :icon="icons.LINK"></v-icon>
        <span> {{ store.currentGroup.name }} </span>
        <EditActivator
          tooltipText="Edit Group"
          :titleIcon="icons.GROUP"
          :onSave="store.saveScratchGroup"
          :onCancel="store.resetScratchGroup">
          <EditGroupDialog />
        </EditActivator>
      </div>
      <div v-else>
        <v-icon class="mr-2" size="small" :icon="icons.GROUP"></v-icon>
        <span> Select a group. </span>
      </div>
    </template>

    <div v-if="store.currentGroup.name" class="d-flex flex-grow-1 overflow-auto h-100">
      <ClaimList
        :claims="store.currentClaims"
        :rowClick="setCurrentClaim"
        :rowProps="highlightSelected"
      />
    </div>

  </v-card>
  <v-card class="my-1 w-100" height="35%">
    <ClaimInfo />
  </v-card>
</template>

<script setup lang="ts">
  import EditGroupDialog from './dialogs/EditGroupDialog.vue'
  import { useAppStore } from '../stores/app'
  import { Claim } from '../util/entityTypes'
  import icons from '../util/icons'

  const store = useAppStore()

  async function setCurrentClaim(_ : PointerEvent, { item } : { item : Claim }) {
    if ( store.currentClaim.id && store.currentClaim.id === item.id ) {
      store.resetCurrentClaim()
      return
    }
    store.$patch({
      currentClaim: item
    })
  }
  function highlightSelected({ item } : { item : Claim }) {
    if ( item.id === store.currentClaim.id ) {
      return {
        class : "bg-grey-lighten-1",
      }
    }
  }
</script>
