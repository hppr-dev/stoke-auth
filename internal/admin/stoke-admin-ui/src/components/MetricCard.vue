<template>
  <SingleMetric v-if="getMetricData()?.values.length === 1" :data="getMetricData()" :name="getPrettyName()" />
  <HistogramMetric v-else-if="getMetricData()?.description.type === 'histogram'" :data="getMetricData()" :name="getPrettyName()" />
  <SummaryMetric v-else-if="getMetricData()?.description.type === 'summary'" :data="getMetricData()" :name="getPrettyName()" />
  <MultiMetric v-else-if="getMetricData()" :data="getMetricData()" :name="getPrettyName()" />
  <v-tooltip :text="isAdded()? 'Remove From Chart' : 'Add To Chart' " location="top">
    <template #activator="{ props }">
      <v-icon v-bind="props" class="ml-n10 mt-5" :color="isAdded()? 'success' : 'error' " @click="addMetric" :icon="icons.ADD"></v-icon>
    </template>
  </v-tooltip>
</template>

<script setup lang="ts">
  import { ref, defineProps } from 'vue'
  import { useAppStore } from '../stores/app'
  import icons from '../util/icons'

  const store = useAppStore()

  const props = defineProps<{
    name: { metricName: string, displayName: string },
  }>()

  function getPrettyName() {
    return props.name.displayName
  }
  function getMetricData() {
    return store.metricData[props.name.metricName]
  }
  function addMetric() {
    if( ! isAdded() ) {
      store.$patch({
        trackedMetrics : [
          ...store.trackedMetrics,
          props.name.metricName
        ]
      })
    } else {
      store.$patch({
        trackedMetrics : [
          ...store.trackedMetrics.filter((s) => s !== props.name.metricName),
        ]
      })
    }
  }
  function isAdded() {
    return store.trackedMetrics.includes(props.name.metricName)
  }
</script>
