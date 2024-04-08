<template>
  <v-container :class="props.noScroll? noScroll : scroll">
    <slot></slot>
  </v-container>
</template>

<script setup lang="ts">
  import { defineProps } from 'vue'

  const scroll = "d-flex h-100 pb-8"
  const noScroll = scroll + " no-scroll"

  const props = defineProps<{
    noScroll?: boolean,
  }>()
</script>
