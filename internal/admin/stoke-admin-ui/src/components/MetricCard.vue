<template>
  <SingleMetric v-if="getMetricData()?.values.length === 1" :data="getMetricData()" :name="compProps.name" />
  <HistogramMetric v-else-if="getMetricData()?.description.type === 'histogram'" :data="getMetricData()" :name="compProps.name" />
  <SummaryMetric v-else-if="getMetricData()?.description.type === 'summary'" :data="getMetricData()" :name="compProps.name" />
  <MultiMetric v-else-if="getMetricData()" :data="getMetricData()" :name="compProps.name" />
</template>

<script setup lang="ts">
  import { ref, defineProps } from 'vue'
  import { useAppStore } from '../stores/app'
  import icons from '../util/icons'

  const store = useAppStore()

  const compProps = defineProps<{
    name: { metricName: string, displayName: string }
  }>()

  function getMetricData() {
    return store.metricData[compProps.name.metricName]
  }
  function isAdded() {
    return store.trackedMetrics.includes(compProps.name.metricName)
  }
</script>
