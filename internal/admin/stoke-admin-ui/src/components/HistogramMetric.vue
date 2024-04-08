<template>
  <v-col cols="8">
    <v-card class="h-100" elevation="5">
      <template #title>
        <span class="text-secondary"> {{ props.name }} </span>
      </template>
      <template #text>
        <p> {{ props.data.description.help }} </p>
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

  const tableHeaders = ref([])
  const tableData = ref([])

  onMounted(() => {
    let distinctLevels = {}
    let distinctTags = {}
    // We want tags: tags, le1: value1, le2: value2
    // Find which tags/what levels to show
    props.data.values.forEach((v) => {
      const { le, part, ...otherTags } = v.tags
      const objKey = JSON.stringify(otherTags)
      if( ! distinctTags[objKey] ) {
        distinctTags[objKey] = {
          tags : otherTags,
        }
      }
      if( part !== "bucket" ) {
        distinctLevels[part] = { title: v.tags.part, value: v.tags.part }
        distinctTags[objKey][part] = v.value
      } else {
        distinctLevels[le] = { title: le, value: le }
        distinctTags[objKey][le] = v.value
      }
    })
    tableHeaders.value = [ { title: "Tags", value: "tags" }, ...Object.values(distinctLevels)]

    tableData.value = Object.values(distinctTags)
  })
</script>
