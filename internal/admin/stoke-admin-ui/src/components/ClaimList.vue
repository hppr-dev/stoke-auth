<template>
  <EntityList :items="claims" :headers="headers" :showSearch="props.showSearch" :showFooter="props.showFooter" :rowClick="props.rowClick">
    <template #footer-prepend>
      <AddActivator
        v-if="props.addButton"
        buttonText="Add Claim"
        :onSave="store.addScratchClaim"
        :onCancel="store.resetScratchClaim"
      >
        <EditClaimDialog add/>
      </AddActivator>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps } from "vue"
  import { useAppStore } from "../stores/app"
  import { Claim } from "../stores/entityTypes"
  import EditClaimDialog from "./dialogs/EditClaimDialog"

  const props= defineProps<{
    claims: Claim[],
    rowClick: Function,
    addButton?: boolean,
    showSearch?: boolean,
    showFooter?: boolean,
  }>()

  const headers = [
    { key : "name", title: "Claim Name" },
    { key : "description", title: "Description"},
  ]

  const store = useAppStore()

</script>
