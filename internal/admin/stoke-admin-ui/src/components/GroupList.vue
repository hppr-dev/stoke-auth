<template>
  <EntityList :items="groups" :headers="headers" :showSearch="props.showSearch" :showFooter="props.showFooter" :rowClick="setCurrentGroup">
    <template #footer-prepend>
      <v-btn v-if="props.addButton" @click="addButton" class="mx-auto" prepend-icon="mdi-plus" color="success"> Add Group </v-btn>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps } from "vue"
  import { useAppStore } from "../stores/app"
  import { Group } from "../stores/entityTypes"

  const props= defineProps<{
    groups: Group[],
    rowClick?: Function,
    addButton?: Function,
    showSearch?: boolean,
    showFooter?: boolean,
  }>()

  const headers = [
    { key : "id", title: "ID" },
    { key : "name", title: "Group Name" },
    { key : "description", title: "Description"},
    { key : "is_user_group", title: "For User?"},
  ]

  const store = useAppStore()

  async function setCurrentGroup(_ : PointerEvent, { item } : { item : Group }) {
    await store.fetchClaimsForGroup(item.id)
    store.$patch({
      currentGroup: item
    })
  }

</script>
