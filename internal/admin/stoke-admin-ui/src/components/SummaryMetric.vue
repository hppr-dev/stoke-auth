<template>
  <v-col clas="h-100" cols="3">
    <v-card class="h-100" elevation="5">
      <template #title>
        <span class="text-secondary"> {{ props.name }} </span>
        <div class="d-flex justify-end">
          <span class="text-secondary"> Sum: </span>
          <span class=""> {{ getSum() }} </span>
        </div>
        <div class="d-flex justify-end">
          <span class="text-secondary"> Count: </span>
          <span class=""> {{ getCount() }} </span>
        </div>
      </template>
      <template #text>
        <span> {{ props.data.description.help }} </span>
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

  function getCount() : string {
    return props.data.values.find((v) => v.tags["part"] === "count").value
  }

  function getSum() : string {
    return props.data.values.find((v) => v.tags["part"] === "sum").value
  }
</script>
