<template>
  <v-sheet class="d-flex flex-column px-4 pt-4 mb-n5" width="60vw" height="65vh">
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Name"
        v-model="name"
        :rules="[require('Name')]"
        @blur="updateScratchGroup"
      ></v-text-field>
    </v-row>
    <v-row>
      <v-textarea
        variant="solo-filled"
        label="Description"
        no-resize
        rows="2"
        v-model="description"
        :rules="[require('Description')]"
        @blur="updateScratchGroup"
      ></v-textarea>
    </v-row>
    <v-row class="mb-5 d-flex flex-grow-1 overflow-auto h-100">
      <v-col cols="9">
        <ClaimList
          showFooter
          showSearch
          addButton
          :claims="store.allClaims"
          :rowClick="addOrRemoveClaim"
          :rowProps="setRowProps"
        />
      </v-col>
      <v-col cols="3">
        <GroupLinkList v-if="!props.add" :links="store.currentLinks"/>
      </v-col>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref, onMounted, defineProps } from "vue"
  import { useAppStore } from "../../stores/app"
  import { Claim } from "../../util/entityTypes"
  import { require } from "../../util/rules"

  const props = defineProps<{
    add?: boolean,
  }>()

  const store = useAppStore()

  const name = ref(store.currentGroup.name)
  const description = ref(store.currentGroup.description)
  if ( props.add ) {
    name.value = ""
    description.value = ""
  }

  function updateScratchGroup() {
    store.$patch({
      scratchGroup: {
        ...store.scratchGroup,
        name: name.value,
        description: description.value,
      },
    })
  }

  function setRowProps({ item } : { item : Claim } ) {
    let isClaim = (g : Claim) => g.id === item.id
    let inCurrent = store.currentClaims.find(isClaim);
    let inScratch = store.scratchClaims.find(isClaim)

    if ( inCurrent && inScratch ) {
      // Existed
      return {
        class: "bg-teal-darken-2",
      }
    } else if ( !inCurrent && inScratch ) {
      // Added
      return {
        class: "bg-green-darken-3",
      }
    } else if ( inCurrent && !inScratch ) {
      // Removed
      return {
        class : "bg-red-darken-3",
      }
    }
    return {}
  }

  function addOrRemoveClaim(_ : PointerEvent,  { item } : { item : Claim } ) {
    let isClaim = (g : Claim) => g.id === item.id
    let inScratch = store.scratchClaims.find(isClaim)
    if ( inScratch ) {
      store.$patch({
        scratchClaims : store.scratchClaims.filter((v) => !isClaim(v)),
      })
    } else {
      store.$patch({
        scratchClaims : [
          ...store.scratchClaims,
          item
        ],
      })
    }
  }

  onMounted(async () => {
    await store.fetchAllClaims()
    if ( props.add ) {
      store.$patch({
        scratchGroup: {},
        currentGroup: {},
        scratchClaims: [],
        currentClaims: [],
      })
    } else {
      store.$patch({
        scratchGroup: {
          ...store.currentGroup,
        },
        scratchClaims: [
          ...store.currentClaims,
        ],
      })
    }
  })
</script>
