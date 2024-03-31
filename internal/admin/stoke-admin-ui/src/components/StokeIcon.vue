<template>
  <v-icon icon="mdi-fire-circle" :size="props.size" color="error"></v-icon>
</template>

<script setup lang="ts">
const props = defineProps<{
  size: string,
}>()
</script>
