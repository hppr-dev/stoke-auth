<template>
  <EntityList :items="store.allUsers" :headers="headers" :showSearch="props.showSearch" :showFooter="props.showFooter" :rowClick="setCurrentUser">
    <template #footer-prepend>
      <v-btn v-if="props.addButton" @click="props.addButton" class="mx-auto" prepend-icon="mdi-plus" color="success"> Add User </v-btn>
    </template>
  </EntityList>
</template>

<script setup lang="ts">
  import { defineProps, onMounted } from "vue"
  import { useAppStore } from "../stores/app"
  import { User } from "../stores/entityTypes"

  const props= defineProps<{
    addButton?: Function,
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

  async function setCurrentUser(_PointerEvent, { item } : { item : User }) {
    await store.fetchGroupsForUser(item.id)
    store.$patch({
      currentUser: item,
    })
  }

  onMounted(store.fetchAllUsers)
</script>
