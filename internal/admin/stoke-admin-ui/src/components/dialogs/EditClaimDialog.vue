<template>
  <v-sheet class="px-4 pt-4 mb-n5" width="40vw">
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Name"
        v-model="name"
        :rules="[require('Name')]"
        @update:modelValue="updateScratchClaim"
      ></v-text-field>
    </v-row>
    <v-row class="d-flex justify-center">
      <v-col cols="4">
        <v-text-field
          variant="outlined"
          label="Short Name"
          v-model="shortName"
          :rules="[require('Short name')]"
          @update:modelValue="updateScratchClaim"
        ></v-text-field>
      </v-col>
      <v-col class="mx-n2 mt-2 text-center" cols="1">
        <span class="text-h4"> = </span>
      </v-col>
      <v-col cols="4">
        <v-text-field
          variant="outlined"
          label="Value"
          v-model="value"
          :rules="[require('Value')]"
          @update:modelValue="updateScratchClaim"
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-textarea
        variant="solo-filled"
        label="Description"
        no-resize
        v-model="description"
        :rules="[require('Description')]"
        @update:modelValue="updateScratchClaim"
      ></v-textarea>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref, onMounted, defineProps } from "vue"
  import { useAppStore } from "../../stores/app"
  import { require } from "../../util/rules"

  const props = defineProps<{
    add?: boolean,
  }>()

  const store = useAppStore()

  const name = ref(store.currentClaim.name)
  const shortName = ref(store.currentClaim.short_name)
  const value = ref(store.currentClaim.value)
  const description = ref(store.currentClaim.description)

  if ( props.add ) {
    name.value = ""
    shortName.value = ""
    value.value = ""
    description.value = ""
  }

  function updateScratchClaim() {
    store.$patch({
      scratchClaim: {
        ...store.scratchClaim,
        name: name.value,
        short_name: shortName.value,
        value: value.value,
        description: description.value,
      },
    })
  }

  onMounted(() => {
    if ( props.add ) {
      store.$patch({
        scratchClaim: {},
        currentClaim: {},
      })
    } else {
      store.$patch({
        scratchClaim: {
          ...store.currentClaim,
        },
      })
    }
  })
</script>
