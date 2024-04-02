<template>
  <v-row class="mt-0 h-50 w-100">
    <v-card class="h-100 w-100 d-flex flex-column" >
      <template #title>
        <div v-if="store.currentUser.fname">
          <span class="text-h3"> {{ store.currentUser.fname }} {{ store.currentUser.lname }} </span>
          <EditActivator
            tooltipText="Edit User"
            :titleIcon="icons.USER"
            :onSave="store.saveScratchUser"
            :onCancel="store.resetScratchUser">
            <EditUserDialog />
          </EditActivator>
        </div>
        <div v-else >
          <v-icon class="mr-2" :icon="icons.USER"></v-icon>
          <span class="text-h6"> Select a user. </span>
        </div>
      </template>

      <template #subtitle>
        <div v-if="store.currentUser.fname">
          <p>
            <v-icon :icon="icons.USER"></v-icon>
            <span class="ml-2 text-h6 font-weight-light"> {{ store.currentUser.username }} </span>
          </p>
          <p>
            <v-icon :icon="icons.MAIL"></v-icon>
            <span class="ml-2 text-h6 font-weight-light"> {{ store.currentUser.email }} </span>
          </p>
        </div>
      </template>

      <div v-if="store.currentUser.fname" class="d-flex flex-grow-1 overflow-auto h-100">
        <GroupList
          :groups="store.currentGroups"
          :rowClick="setCurrentGroup"
          :rowProps="highlightSelected"
        />
      </div>

    </v-card>
  </v-row>
  <v-row class="pt-3 h-50 w-100">
    <v-card class="h-100 w-100">
      <GroupInfo />
    </v-card>
  </v-row>
</template>

<script setup lang="ts">
  import EditUserDialog from './dialogs/EditUserDialog.vue'
  import { useAppStore } from "../stores/app"
  import { Group } from "../util/entityTypes"
  import icons from '../util/icons'

  const store = useAppStore()

  async function setCurrentGroup(_ : PointerEvent, { item } : { item : Group }) {
    if ( store.currentGroup.id && store.currentGroup.id === item.id ) {
      store.resetCurrentGroup()
      return
    }
    await store.fetchClaimsForGroup(item.id)
    store.$patch({
      currentGroup: item
    })
  }

  function highlightSelected({ item } : { item : Group }) {
    if ( item.id === store.currentGroup.id ) {
      return {
        class : "bg-grey-lighten-1",
      }
    }
  }
</script>
