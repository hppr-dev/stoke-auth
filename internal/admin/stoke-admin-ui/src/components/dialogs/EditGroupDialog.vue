<template>
  <v-sheet class="d-flex flex-column px-4 pt-4 mb-n5" width="40vw" height="70vh">
    <v-row>
      <v-text-field variant="solo-filled" clearable label="Name" v-model="name" @blur="updateCurrentGroup"></v-text-field>
    </v-row>
    <v-row>
      <v-textarea variant="solo-filled" clearable label="Description" v-model="description" no-resize @blur="updateCurrentGroup"></v-textarea>
    </v-row>
    <v-row class="mb-5 d-flex flex-grow-1 overflow-auto h-100">
      <ClaimList :claims="claims" showFooter :addButton="addClaim"/>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import { useAppStore } from "@/stores/app"

  const store = useAppStore()

  const name = ref(store.currentGroup.name)
  const description = ref(store.currentGroup.description)
  const claims = ref(store.currentClaims)

  function updateCurrentGroup() {
    store.$patch({
      currentGroup: {
        ...store.currentGroup,
        name: name.value,
        description: description.value,
      },
    })
  }

  function addClaim() {
    console.log("components/dialogs/EditGroupDialog.vue")
  }
</script>
