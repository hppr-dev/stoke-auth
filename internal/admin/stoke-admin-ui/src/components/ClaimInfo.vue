<template>
  <v-card class="h-100" :text="store.currentClaim.description">
    <template #title>
      <div v-if="store.currentClaim.name" >
        <v-icon class="mr-2" size="small" :icon="icons.CLAIM"></v-icon>
        <span> {{ store.currentClaim.name }} </span>
        <EditActivator
          tooltipText="Edit Claim"
          :titleIcon="icons.CLAIM"
          :onSave="store.saveScratchClaim"
          :onCancel="store.resetScratchClaim">
          <EditClaimDialog />
        </EditActivator>
      </div>
      <div v-else >
        <v-icon class="mr-2" size="small" :icon="icons.CLAIM"></v-icon>
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
  import icons from '../util/icons'

  const store = useAppStore()
</script>
