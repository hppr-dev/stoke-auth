<template>
    <v-textarea class="h-100" :model-value="store.logText.split('\n').reverse().join('\n').trimStart()" auto-grow></v-textarea>
</template>

<script setup lang="ts">
  import { onMounted } from 'vue'
  import { useAppStore } from "../stores/app"

  const store = useAppStore()
  // ASNI_UP? https://www.npmjs.com/package/ansi-up

  onMounted( async () => {
    await store.fetchLogText()
  })

</script>
