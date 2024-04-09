<template>
  <Viewport>
    <v-row class="mb-10 ml-0">
      <v-card class="h-100 w-100 mt-2 mb-5">
        <v-tabs align-tabs="center" v-model="tab" color="warning">
          <v-tab value="1">GC</v-tab>
          <v-tab value="2">Memory</v-tab>
          <v-tab value="3">Process</v-tab>
          <v-tab value="4">HTTP</v-tab>
          <v-tab v-if="store.trackedMetrics.length > 0" value="5">Chart</v-tab>
          <v-tab value="6">Logs</v-tab>
        </v-tabs>
        <v-container class="h-100">
          <div class="d-flex justify-end">
            <div class="mr-3 mt-2">
              <span class="text-body-1 font-weight-light"> Max Points: </span>
            </div>
            <div width="10em" >
              <v-select
                density="compact"
                v-model="selectedMaxPoints"
                :items="maxPointItems"
                :disabled="!store.metricsPaused"
                @update:modelValue="store.setMaxPoints"
              > </v-select>
            </div>
            <div class="mx-3 mt-2">
              <span class="text-body-1 font-weight-light"> Refresh Every: </span>
            </div>
            <div class="mx-3" width="10em" >
              <v-select
                density="compact"
                v-model="selectedTimer"
                :items="timerSelectItems"
                @update:modelValue="store.setMetricRefresh"
                :disabled="!store.metricsPaused"
              > </v-select>
            </div>
            <div class="mr-3 mt-2">
              <v-icon
                @click="onStart"
                :icon="store.metricsPaused? 'mdi-play-circle' : 'mdi-pause' "
                :color="store.metricsPaused? 'error': 'primary'"
              ></v-icon>
            </div>
          </div>
          <v-window class="h-100" v-model="tab">
            <v-window-item class="h-100" key="gc" value="1">
              <MetricFilter :metricNames="gcMetrics"/>
            </v-window-item>
            <v-window-item class="h-100" key="memstats" value="2">
              <MetricFilter :metricNames="memoryMetrics"/>
            </v-window-item>
            <v-window-item class="h-100" key="process" value="3">
              <MetricFilter :metricNames="processMetrics"/>
            </v-window-item>
            <v-window-item class="h-100" key="http" value="4">
              <MetricFilter :metricNames="httpMetrics"/>
            </v-window-item>
            <v-window-item class="h-100" key="chart" value="5">
              <MetricChart />
            </v-window-item>
            <v-window-item class="h-100" key="logs" value="6">
              <p> Logs here </p>
            </v-window-item>
          </v-window>
        </v-container>
      </v-card>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue'
  import { useAppStore } from '../stores/app'

  const store = useAppStore()

  const tab = ref(0)
  const selectedTimer = ref(30000)
  const timerSelectItems = [
    { title : "1s",  value: 1000 },
    { title : "5s",  value: 5000 },
    { title : "10s", value: 10000 },
    { title : "30s", value: 30000 },
    { title : "1m",  value: 60000 },
  ]
  const selectedMaxPoints = ref(100)
  const maxPointItems = [
    { title : "25",  value: 25 },
    { title : "50",  value: 50 },
    { title : "100", value: 100 },
    { title : "200", value: 200 },
    { title : "1000",  value: 1000 },
  ]

  function onStart() {
    store.toggleMetricPaused()
    if(! store.metricsPaused && store.trackedMetrics.length > 0) {
      tab.value = 4
    }
  }

  onMounted(async () => {
    await store.fetchMetricData()
    store.metricRefresh()
  })

  const gcMetrics = [
    { metricName: "go_gc_duration_seconds",           displayName: "GC Duration Seconds" },
    { metricName: "go_memstats_gc_sys_bytes",         displayName: "GC System Bytes" },
    { metricName: "go_memstats_last_gc_time_seconds", displayName: "Last GC Time Seconds" },
    { metricName: "go_memstats_next_gc_bytes",        displayName: "Next GC Bytes" },
    { metricName: "go_memstats_heap_alloc_bytes",     displayName: "Allocated Heap Bytes" },
    { metricName: "go_memstats_heap_idle_bytes",      displayName: "Idle Heap Bytes" },
    { metricName: "go_memstats_heap_inuse_bytes",     displayName: "Active Heap Bytes" },
    { metricName: "go_memstats_heap_released_bytes",  displayName: "Released Heap Bytes" },
    { metricName: "go_memstats_heap_sys_bytes",       displayName: "System Heap Bytes" },
    { metricName: "go_memstats_heap_objects",         displayName: "Heap Objects" },
  ]

  const memoryMetrics = [
    { metricName: "go_memstats_frees_total",         displayName: "Total Frees" },
    { metricName: "go_memstats_alloc_bytes",         displayName: "Allocated Bytes" },
    { metricName: "go_memstats_mallocs_total",       displayName: "Malloc Total" },
    { metricName: "go_memstats_alloc_bytes_total",   displayName: "Total Allocated Bytes" },
    { metricName: "go_memstats_buck_hash_sys_bytes", displayName: "System Bucket Hash Bytes" },
    { metricName: "go_memstats_lookups_total",       displayName: "Total Lookups" },
    { metricName: "go_memstats_mcache_inuse_bytes",  displayName: "Active Mcache Bytes" },
    { metricName: "go_memstats_mcache_sys_bytes",    displayName: "System Mcache Bytes" },
    { metricName: "go_memstats_mspan_inuse_bytes",   displayName: "Active Mspan Bytes" },
    { metricName: "go_memstats_mspan_sys_bytes",     displayName: "System Mspan Bytes" },
    { metricName: "go_memstats_stack_inuse_bytes",   displayName: "Active Stack Bytes" },
    { metricName: "go_memstats_stack_sys_bytes",     displayName: "System Stack Bytes" },
    { metricName: "go_memstats_other_sys_bytes",     displayName: "Other System Bytes" },
    { metricName: "go_memstats_sys_bytes",           displayName: "System Bytes" },
  ]

  const processMetrics = [
    //"go_info",
    { metricName: "go_goroutines",                    displayName: "Total Goroutines" },
    { metricName: "go_threads",                       displayName: "Total Threads" },
    { metricName: "process_cpu_seconds_total",        displayName: "Total CPU Time Seconds" },
    { metricName: "process_max_fds",                  displayName: "Max File Descriptors" },
    { metricName: "process_open_fds",                 displayName: "Open File Descriptors" },
    { metricName: "process_resident_memory_bytes",    displayName: "Resident Memory Bytes" },
    { metricName: "process_start_time_seconds",       displayName: "Start Time Seconds" },
    { metricName: "process_virtual_memory_bytes",     displayName: "Virtual Memory Bytes" },
    { metricName: "process_virtual_memory_max_bytes", displayName: "Max Virtual Memory Bytes" },
  ]

  const httpMetrics = [
    { metricName: "http_server_duration_milliseconds",        displayName: "Endpoint Millisecond Buckets" },
    { metricName: "http_server_request_size_bytes_total",     displayName: "Endpoint Request Size Bytes" },
    { metricName: "http_server_response_size_bytes_total",    displayName: "Endpoint Response Size Bytes" },
    { metricName: "stoke_database_mutation_total",            displayName: "Total Database Mutations" },
    { metricName: "stoke_database_mutation_millis_total",     displayName: "Total Database Time" },
    { metricName: "stoke_database_mutation_millis_histogram", displayName: "Total Database Time" },
  ]

</script>
