<template>
  <Viewport noScroll>
    <v-row class="h-screen pb-10 ml-0">
      <v-col class="ml-n3 h-100" lg="7" md="6" sm="12">
        <UserList showSearch showFooter addButton/>
      </v-col>
      <v-col class="h-100" lg="5" md="6" sm="12">
        <UserInfo />
      </v-col>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  //
</script>
