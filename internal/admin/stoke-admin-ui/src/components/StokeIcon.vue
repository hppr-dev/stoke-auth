<template>
  <v-icon icon="mdi-fire-circle" :size="props.size" color="red-lighten-1"></v-icon>
</template>

<script setup lang="ts">
const props = defineProps<{
  size: string,
}>()
</script>
