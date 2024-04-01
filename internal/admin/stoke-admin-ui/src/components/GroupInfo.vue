<template>
  <v-card class="d-flex flex-column w-100" height="65%" :subtitle="store.currentGroup.description">
    <template #title>
      <div v-if="store.currentGroup.name">
        <span> {{ store.currentGroup.name }} </span>
        <EditActivator tooltipText="Edit Group" :onSave="store.saveScratchGroup" :onCancel="store.resetScratchGroup">
          <EditGroupDialog />
        </EditActivator>
      </div>
      <div v-else>
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
  import { Claim } from '../stores/entityTypes'

  const store = useAppStore()

  async function setCurrentClaim(_ : PointerEvent, { item } : { item : Claim }) {
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
