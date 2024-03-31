<template>
  <v-sheet class="px-4 pt-4 mb-n5" width="30vw">
    <v-row>
      <v-text-field variant="solo-filled" clearable label="Name" v-model="name" @blur="updateCurrentClaim"></v-text-field>
    </v-row>
    <v-row class="d-flex justify-center">
      <v-col cols="4">
        <v-text-field variant="outlined" clearable label="Short Name" v-model="shortName" @blur="updateCurrentClaim"></v-text-field>
      </v-col>
      <v-col class="mx-n2 mt-2 text-center" cols="1">
        <span class="text-h4"> = </span>
      </v-col>
      <v-col cols="4">
        <v-text-field variant="outlined" clearable label="Value" v-model="value"></v-text-field @blur="updateCurrentClaim">
      </v-col>
    </v-row>
    <v-row>
      <v-textarea variant="solo-filled" clearable label="Description" v-model="description" no-resize @blur="updateCurrentClaim"></v-textarea>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref } from "vue"
  import { useAppStore } from "@/stores/app"

  const store = useAppStore()

  const name = ref(store.currentClaim.name)
  const shortName = ref(store.currentClaim.short_name)
  const value = ref(store.currentClaim.value)
  const description = ref(store.currentClaim.description)

  function updateCurrentClaim() {
    store.$patch({
      currentClaim: {
        ...store.currentClaim,
        name: name.value,
        short_name: shortName.value,
        value: value.value,
        description: description.value,
      },
    })
  }
</script>
