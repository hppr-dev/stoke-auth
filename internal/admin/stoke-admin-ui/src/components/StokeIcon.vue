<template>
  <v-icon :icon="mdiFireCircle" :size="props.size" color="red-lighten-1"></v-icon>
</template>

<script setup lang="ts">
import { mdiFireCircle } from "@mdi/js"

const props = defineProps<{
  size: string,
}>()
</script>
