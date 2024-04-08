<template>
  <v-row>
    <MetricCard v-for="metric in props.metricNames" :name="metric"/>
  </v-row>
</template>

<script setup lang="ts">
  import { defineProps } from 'vue'

  const props = defineProps<{
    metricNames: string[]
  }>()

</script>
