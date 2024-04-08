<template>
  <v-col cols="4">
    <v-card class="h-100" elevation="5">
      <template #title>
        <span class="text-secondary"> {{ props.name }} </span>
      </template>
      <template #subtitle>
        <p> {{ props.data?.description.help }} </p>
      </template>
      <template #text>
        <v-data-table :cell-props="{ style :'width: 2em' }" :items="tableData" :headers="tableHeaders">
          <template #item.tags="{ value }">
            <v-chip-group style="max-width : 21em">
              <v-tooltip v-for="tagKey in Object.keys(value).filter((k) => !filteredTags.includes(k))" :text="tagKey" >
                <template #activator= "{ props }">
                  <v-chip v-bind="props"> {{ value[tagKey] }} </v-chip>
                </template>
              </v-tooltip>
            </v-chip-group>
          </template>
        </v-data-table>
      </template>
    </v-card>
  </v-col>
</template>

<script setup lang="ts">
  import { ref, defineProps, onMounted } from 'vue'
  import { MetricData } from '../util/prometheus'

  const props = defineProps<{
    name: string
    data: MetricData
  }>()

  const filteredTags = [
    "part", "http_scheme", "otel_scope_name", "otel_scope_version",
    "net_host_name", "net_host_port", "net_protocol_name","net_protocol_version",
  ]

  const tableHeaders = ref([
    { title: "Tags", value: "tags" },
    { title: "Value", value: "value" },
  ])
  const tableData = ref([])

  onMounted(() => {
    if( !props.data || props.data.values.length === 0 ){
      console.error("given empty monitoring data:", props.name, props.data)
      return
    }
    let distinctTags = {}
    // We want tags: tags, le1: value1, le2: value2
    // Find which tags/what levels to show
    props.data.values.forEach((v) => {
      const objKey = JSON.stringify(v.tags)
      distinctTags[objKey] = {
        tags : v.tags,
        value : v.value,
      }
    })

    tableData.value = Object.values(distinctTags)

  })
</script>
