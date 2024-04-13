<template>
  <v-data-table
    sticky
    :headerProps="headerProps"
    :headers="headers"
    :items="store.currentLinks"
  >
    <template #no-data>
      <span> Group is not linked </span>
    </template>
    <template #bottom>
      <div class="text-center pt-2">
        <v-row>
          <v-pagination
            v-if="store.currentLinks.length >= 4"
            v-model="page"
            :length="store.currentLinks.length/4 + 1"
          ></v-pagination>
          <AddActivator
            buttonText="Link Group"
            :titleIcon="icons.LINK"
            :onSave="store.addScratchLink"
            :onCancel="store.resetScratchLink"
          >
            <AddLinkDialog/>
          </AddActivator>
        </v-row>
      </div>
    </template>
  </v-data-table>
</template>

<script setup lang="ts">
  import { ref, defineProps } from "vue"
  import { useAppStore } from "../stores/app"
  import icons from '../util/icons'
  import AddLinkDialog from './dialogs/AddLinkDialog'

  const page = ref(1)
  const headers = [
    { key : "type", title: "Link Type" },
    { key : "resource_spec", title: "Resource"},
  ]

  const headerProps = {
    class : "bg-blue-grey",
    height: "3em",
  }

  const store = useAppStore()

</script>
