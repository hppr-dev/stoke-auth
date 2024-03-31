<template>
  <Viewport>
    <v-row class="h-screen pb-10 ml-0">
      <v-col class="ml-n3 h-100" md="8" sm="5">
        <GroupList showSearch showFooter :addButton="addGroup" :groups="store.allGroups"/>
      </v-col>
      <v-col class="ml-n3 h-100" md="4" offset-sm="1" offset-md="0" sm="5">
        <GroupInfo />
      </v-col>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue'
  import { useAppStore } from '../stores/app'

  const store = useAppStore()

  onMounted(store.fetchAllGroups)

  function addGroup() {
    console.log("pages/group.vue")
  }
</script>
