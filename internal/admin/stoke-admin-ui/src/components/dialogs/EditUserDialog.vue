<template>
  <v-sheet class="d-flex flex-column px-4 pt-4 mb-n5" width="80vw" height="80vh">
    <v-row>
      <v-text-field variant="solo-filled" clearable label="Username" v-model="username" no-resize disabled></v-text-field>
      <v-btn class="ml-5 h-75" variant="tonal" color="info" stacked prepend-icon="mdi-lock-open-variant" disabled> Unlock </v-btn>
      <UpdateUserPassword />
    </v-row>
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="First Name"
        v-model="fname"
        :rules="[require('First Name'), hasAChange]"
        @update:modelValue="updateScratchUser"
      ></v-text-field>
      <v-text-field
        class="ml-4"
        variant="solo-filled"
        label="Last Name"
        v-model="lname"
        :rules="[require('Last Name'), hasAChange]"
        @update:modelValue="updateScratchUser"
      ></v-text-field>
    </v-row>
    <v-row>
      <v-text-field
        variant="solo-filled"
        label="Email"
        v-model="email"
        :rules="[require('Email'), hasAChange]"
        @update:modelValue="updateScratchUser"
      ></v-text-field>
    </v-row>
    <v-row class="mb-5 d-flex flex-grow-1 overflow-auto h-100">
      <GroupList
        showFooter
        showSearch
        addButton
        :groups="store.allGroups"
        :rowProps="setRowProps"
        :rowClick="addOrRemoveGroup"
      >
      </GroupList>
    </v-row>
  </v-sheet>
</template>

<script setup lang="ts">
  import { ref, onMounted } from "vue"
  import { useAppStore } from "../../stores/app"
  import { Group } from "../../util/entityTypes"
  import { require } from "../../util/rules"
import UpdateUserPassword from "./UpdateUserPassword.vue";

  const store = useAppStore()

  const fname = ref(store.currentUser.fname)
  const lname = ref(store.currentUser.lname)
  const username = ref(store.currentUser.username)
  const email = ref(store.currentUser.email)

  let changed = false
  function hasAChange() {
    return changed || "At least one value must be updated"
  }

  function updateScratchUser() {
    changed = true
    store.$patch({
      scratchUser: {
        ...store.scratchUser,
        fname: fname.value,
        lname: lname.value,
        email: email.value,
      },
    })
  }

  function setRowProps({ item } : { item : Group } ) {
    let isGroup = (g : Group) => g.id === item.id
    let inCurrent = store.currentGroups.find(isGroup);
    let inScratch = store.scratchGroups.find(isGroup)

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

  function addOrRemoveGroup(_ : PointerEvent,  { item } : { item : Group } ) {
    changed = true
    let isGroup = (g : Group) => g.id === item.id
    let inScratch = store.scratchGroups.find(isGroup)
    if ( inScratch ) {
      store.$patch({
        scratchGroups : store.scratchGroups.filter((v) => !isGroup(v)),
      })
    } else {
      store.$patch({
        scratchGroups : [
          ...store.scratchGroups,
          item
        ],
      })
    }
  }

  onMounted(async () => {
    await store.fetchAllGroups()
    store.$patch({
      scratchUser: {
        ...store.currentUser,
      },
      scratchGroups: [
        ...store.currentGroups,
      ],
    })
  })
</script>
