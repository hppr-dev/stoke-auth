<template>
  <v-row class="mt-0 h-50 w-100">
    <v-card class="h-100 w-100 d-flex flex-column" >
      <template #title>
        <div v-if="store.currentUser.fname">
          <span class="text-h3"> {{ store.currentUser.fname }} {{ store.currentUser.lname }} </span>
          <EditActivator tooltipText="Edit User" dialogTitle="Edit User">
            <EditUserDialog />
          </EditActivator>
        </div>
        <div v-else >
          <span class="text-h6"> Select a user. </span>
        </div>
      </template>

      <template #subtitle>
        <div v-if="store.currentUser.fname">
          <p>
            <v-icon icon="mdi-account"></v-icon>
            <span class="ml-2 text-h6 font-weight-light"> {{ store.currentUser.username }} </span>
          </p>
          <p>
            <v-icon icon="mdi-mail"></v-icon>
            <span class="ml-2 text-h6 font-weight-light"> {{ store.currentUser.email }} </span>
          </p>
        </div>
      </template>

      <div v-if="store.currentUser.fname" class="d-flex flex-grow-1 overflow-auto h-100">
        <GroupList :groups="store.currentGroups"/>
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

  const store = useAppStore()
</script>
