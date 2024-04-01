<template>
  <v-card class="h-100" :text="store.currentClaim.description">
    <template #title>
      <div v-if="store.currentClaim.name" >
        <span> {{ store.currentClaim.name }} </span>
        <EditActivator tooltipText="Edit Claim" :onSave="store.saveScratchClaim" :onCancel="store.resetScratchClaim">
          <EditClaimDialog />
        </EditActivator>
      </div>
      <div v-else >
        <span> Select a claim. </span>
      </div>
    </template>

    <template #subtitle>
      <div v-if="store.currentClaim.name">
        stk.{{ store.currentClaim.short_name }} = {{ store.currentClaim.value }}
      </div>
    </template>

  </v-card>
</template>

<script setup lang="ts">
  import EditClaimDialog from './dialogs/EditClaimDialog.vue'
  import { useAppStore } from '../stores/app'

  const store = useAppStore()
</script>
