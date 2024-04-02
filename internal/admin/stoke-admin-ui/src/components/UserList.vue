<template>
  <EntityList
    searchIcon="mdi-account-search"
    :items="store.allUsers"
    :headers="headers"
    :showSearch="props.showSearch"
    :showFooter="props.showFooter"
    :rowClick="setCurrentUser"
    :rowProps="highlightSelected"
  >
    <template #footer-prepend>
      <AddActivator
        v-if="props.addButton"
        buttonText="Add User"
        titleIcon="mdi-account"
        :onSave="store.addScratchUser"
        :onCancel="store.resetScratchUser"
      >
        <AddUserDialog />
      </AddActivator>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps, onMounted } from "vue"
  import { useAppStore } from "../stores/app"
  import { User } from "../util/entityTypes"
  import AddUserDialog from "./dialogs/AddUserDialog.vue";

  const props= defineProps<{
    addButton?: boolean,
    showSearch?: boolean,
    showFooter?: boolean,
  }>()

  const headers = [
    { key: "id",  title: "ID" },
    { key: "fname", title: "First Name" },
    { key: "lname", title: "Last Name" },
    { key: "username", title: "Username" },
    { key: "email", title: "Email" },
  ]

  const store = useAppStore()

  async function setCurrentUser(_ : PointerEvent, { item } : { item : User }) {
    if ( store.currentUser.id && store.currentUser.id === item.id ) {
      store.resetCurrentUser()
      return
    }
    await store.fetchGroupsForUser(item.id)
    store.$patch({
      currentUser: item,
    })
  }

  function highlightSelected({ item } : { item : User }) {
    if ( item.id === store.currentUser.id ) {
      return {
        class : "bg-grey-lighten-1",
      }
    }
  }

  onMounted(store.fetchAllUsers)
</script>
