<template>
  <EntityList
    searchIcon="icons.USER_SEARCH"
    deleteItemKey="username"
    :items="store.allUsers"
    :headers="headers"
    :showSearch="props.showSearch"
    :showFooter="props.showFooter"
    :rowClick="setCurrentUser"
    :rowProps="highlightSelected"
    :deleteClick="store.deleteUser"
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
  <template #row-icon="{ item }">
      <v-icon
        :icon="item.source == 'LDAP'? icons.LINK: icons.LOCAL"
        :color="item.source == 'LDAP'? 'success': 'warning'"
      > </v-icon>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps, onMounted } from "vue"
  import { useAppStore } from "../stores/app"
  import { User } from "../util/entityTypes"
  import icons from "../util/icons"
  import AddUserDialog from "./dialogs/AddUserDialog.vue";

  const props= defineProps<{
    addButton?: boolean,
    showSearch?: boolean,
    showFooter?: boolean,
  }>()

  const headers = [
    { key: "username", title: "Username" },
    { key: "row-icon", title:"Type", value: "source" },
    { key: "fname", title: "First Name" },
    { key: "lname", title: "Last Name" },
    { key: "email", title: "Email" },
  ]

  const store = useAppStore()

  async function setCurrentUser(_ : PointerEvent, { item } : { item : User }) {
    store.resetCurrentUser()
    if ( store.currentUser.id && store.currentUser.id === item.id ) {
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
