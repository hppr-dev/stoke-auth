<template>
  <EntityList
    :items="groups"
    :headers="headers"
    :showSearch="props.showSearch"
    :searchIcon="icons.GROUP_SEARCH"
    :showFooter="props.showFooter"
    :rowClick="props.rowClick">
    <template #footer-prepend>
      <AddActivator
        v-if="props.addButton"
        buttonText="Add Group"
        titleIcon="mdi-account-multiple"
        :onSave="store.addScratchGroup"
        :onCancel="store.resetScratchGroup"
      >
        <EditGroupDialog add/>
      </AddActivator>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps } from "vue"
  import { useAppStore } from "../stores/app"
  import { Group } from "../util/entityTypes"
  import icons from '../util/icons'
  import EditGroupDialog from "./dialogs/EditGroupDialog"

  const props= defineProps<{
    groups: Group[],
    rowClick: Function,
    addButton?: boolean,
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


</script>
