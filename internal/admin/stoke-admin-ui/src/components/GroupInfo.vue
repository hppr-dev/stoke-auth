<template>
  <v-card class="d-flex flex-column w-100" height="65%" :subtitle="store.currentGroup.description">
    <template #title>
      <div v-if="store.currentGroup.name">
        <span> {{ store.currentGroup.name }} </span>
        <EditActivator tooltipText="Edit Group">
          <EditGroupDialog />
        </EditActivator>
      </div>
      <div v-else>
        <span> Select a group. </span>
      </div>
    </template>

    <div v-if="store.currentGroup.name" class="d-flex flex-grow-1 overflow-auto h-100">
      <ClaimList :claims="store.currentClaims"/>
    </div>

  </v-card>
  <v-card class="my-1 w-100" height="35%">
    <ClaimInfo />
  </v-card>
</template>

<script setup lang="ts">
  import EditGroupDialog from './dialogs/EditGroupDialog.vue'
  import { useAppStore } from '../stores/app'

  const store = useAppStore()
</script>
