<template>
  <v-icon :icon="icons.STOKE" :size="props.size" color="red-lighten-1"></v-icon>
</template>

<script setup lang="ts">
import icons from '../util/icons'

const props = defineProps<{
  size: string,
}>()
</script>
