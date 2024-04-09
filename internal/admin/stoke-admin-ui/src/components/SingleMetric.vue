<template>
  <v-col cols="3">
    <v-card
      class="h-100"
      elevation="5"
      :color="store.trackedMetrics.includes(props.name.metricName)? 'blue-grey-darken-2' : ''"
      @click="store.toggleMetric(props.name)"
    >
      <template #title>
          <span class="text-secondary"> {{ props.name.displayName }} </span>
          <span class="d-flex justify-end"> {{ props.data.values[0].value }} </span>
      </template>
      <template #text>
        <p> {{ props.data.description.help }} </p>
      </template>
    </v-card>
  </v-col>
</template>

<script setup lang="ts">
  import { defineProps } from 'vue'
  import { useAppStore } from '../stores/app'
  import { MetricData } from '../util/prometheus'

  const store = useAppStore()

  const props = defineProps<{
    name: { metricName: string, displayName: string }
    data: MetricData
  }>()
</script>
