<template>
  <v-sheet class="px-4 pt-4 mb-n5" width="40vw">
    <v-row>
      <v-select
        variant="solo-filled"
        label="Provider Name"
        v-model="linkType"
        :items="store.availableProviders"
        item-title="type_spec"
        item-value="type_spec"
        :rules="[require('Provider Name')]"
      ></v-select>
    </v-row>
    <v-row v-if="selectedProviderType() == 'LDAP'" >
      <v-text-field
        variant="solo-filled"
        label="Group Name"
        no-resize
        v-model="ldapGroup"
        :rules="[require('Group name')]"
        @update:modelValue="updateScratchLink"
      ></v-text-field>
    </v-row>
    <v-row v-if="selectedProviderType() == 'OIDC'" >
      <v-text-field
        variant="solo-filled"
        label="Claim"
        no-resize
        v-model="oidcClaim"
        :rules="[require('Claim')]"
        @update:modelValue="updateScratchLink"
      ></v-text-field>
      <h1>=</h1>
      <v-text-field
        variant="solo-filled"
        label="Value"
        no-resize
        v-model="oidcValue"
        :rules="[require('Value')]"
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

  const linkType = ref("")

  const ldapGroup = ref("")

  const oidcClaim = ref("")
  const oidcValue = ref("")

  function selectedProviderType() {
    let sel = store.availableProviders.find((p) => p.type_spec == linkType.value)
    return sel? sel.provider_type : ""
  }

  function updateScratchLink() {
    let provType = selectedProviderType()
    let resourceSpec = ""
    if ( provType == "LDAP" ) {
      resourceSpec = ldapGroup.value
    } else if ( provType == "OIDC" ) {
      resourceSpec = `${ oidcClaim.value }=${ oidcValue.value }`
    }
    store.$patch({
      scratchLink: {
        ...store.scratchLink,
        type: linkType.value,
        resource_spec: resourceSpec,
      },
    })
  }

  onMounted(() => {
    store.fetchAvailableProviders()
    store.$patch({
      scratchLink: {
        type: "",
        claim_group: store.currentGroup.id,
      },
    })
  })
</script>
