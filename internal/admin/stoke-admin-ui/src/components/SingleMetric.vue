<template>
  <v-col cols="4">
    <v-card elevation="5">
      <template #title>
          <span class="text-secondary"> {{ props.name }} </span>
          <span class="d-flex justify-end"> {{ props.data.values[0].value }} </span>
      </template>
      <template #subtitle>
        <p> {{ props.data.description.help }} </p>
      </template>
    </v-card>
  </v-col>
</template>

<script setup lang="ts">
  import { defineProps } from 'vue'
  import { MetricData } from '../util/prometheus'

  const props = defineProps<{
    name: string
    data: MetricData
  }>()
</script>
