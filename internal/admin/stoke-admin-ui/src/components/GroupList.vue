<template>
  <EntityList
    deleteItemKey="name"
    :items="groups"
    :headers="headers"
    :showSearch="props.showSearch"
    :searchIcon="icons.GROUP_SEARCH"
    :showFooter="props.showFooter"
    :rowClick="props.rowClick"
    :deleteClick="props.deleteGroup"
    :totalItems="store.entityTotals.claim_groups"
    :perPage="12"
    :onNext="loadPageIfNeeded"
    >
    <template #footer-prepend>
      <AddActivator
        v-if="props.addButton"
        buttonText="Add Group"
        :titleIcon="icons.GROUP"
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
    deleteGroup?: Function
  }>()

  const headers = [
    { key : "name", title: "Group Name" },
    { key : "description", title: "Description"},
  ]

  const store = useAppStore()

  async function loadPageIfNeeded(page: number) {
    if ( store.allGroups.length < store.entityTotals.claim_groups && store.allGroups.length < (page - 1) * 12) {
      await store.fetchAllGroups(false, (store.allGroups.length/store.pageLoadSize) + 1)
    }
  }

</script>
