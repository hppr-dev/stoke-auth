<template>
  <v-row>
    <MetricCard v-for="metric in props.metricNames" :name="metric"/>
  </v-row>
</template>

<script setup lang="ts">
  import { defineProps } from 'vue'
  import { useAppStore } from '../stores/app'

  const props = defineProps<{
    metricNames: string[]
  }>()

  const store = useAppStore()

</script>
