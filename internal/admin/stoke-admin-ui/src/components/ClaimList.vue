<template>
  <EntityList :items="claims" :headers="headers" :showSearch="props.showSearch" :showFooter="props.showFooter" :rowClick="setCurrentClaim">
    <template #footer-prepend>
      <v-btn v-if="props.addButton" @click="addButton" class="mx-auto" prepend-icon="mdi-plus" color="success"> Add Claim </v-btn>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps } from "vue"
  import { useAppStore } from "../stores/app"
  import { Claim } from "../stores/entityTypes"

  const props= defineProps<{
    claims: Claim[],
    addButton?: Function,
    showSearch?: boolean,
    showFooter?: boolean,
  }>()

  const headers = [
    { key : "name", title: "Claim Name" },
    { key : "description", title: "Description"},
  ]

  const store = useAppStore()

  async function setCurrentClaim(_ : PointerEvent, { item } : { item : Claim }) {
    store.$patch({
      currentClaim: item
    })
  }
</script>
