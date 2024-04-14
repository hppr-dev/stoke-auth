<template>
  <v-sheet class="px-4 pt-4 mb-n5" width="40vw">
    <v-row>
      <v-select
        variant="solo-filled"
        label="Type"
        readonly
        v-model="linkType"
        :items="['LDAP']"
      ></v-select>
    </v-row>
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Group Name"
        no-resize
        v-model="resourceSpec"
        :rules="[require('Group name')]"
        @update:modelValue="updateScratchLink"
      ></v-text-field>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref, onMounted } from "vue"
  import { useAppStore } from "../../stores/app"
  import { require } from "../../util/rules"

  const store = useAppStore()

  const linkType = ref("LDAP")
  const resourceSpec = ref("")

  function updateScratchLink() {
    store.$patch({
      scratchLink: {
        ...store.scratchLink,
        resource_spec: resourceSpec.value,
      },
    })
  }

  onMounted(() => {
    store.$patch({
      scratchLink: {
        type: "LDAP",
        claim_group: store.currentGroup.id,
      },
    })
  })
</script>
