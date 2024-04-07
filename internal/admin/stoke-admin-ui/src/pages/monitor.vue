<template>
  <Viewport>
    <v-row class="h-screen pb-10 ml-0">
      <v-card class="w-100 mt-2 mb-4">
        <v-tabs align-tabs="center" v-model="tab" color="warning">
          <v-tab value="0">Logs</v-tab>
          <v-tab value="1">GC</v-tab>
          <v-tab value="2">Memory</v-tab>
          <v-tab value="3">Process</v-tab>
          <v-tab value="4">HTTP</v-tab>
          <v-tab value="5">Tracing</v-tab>
        </v-tabs>
        <v-window class="h-100" v-model="tab">
          <v-container>
            <v-window-item class="h-100" key="logs" value="0">
              <p> Logs here </p>
            </v-window-item>
            <v-window-item class="h-100" key="gc" value="1">
              <p> {{ Object.keys(store.metricData).filter((n) => n.startsWith("go_gc") ) }} </p>
            </v-window-item>
            <v-window-item class="h-100" key="memstats" value="2">
              <p> {{ Object.keys(store.metricData).filter((n) => n.startsWith("go_memstats") ) }} </p>
            </v-window-item>
            <v-window-item class="h-100" key="process" value="3">
              <p> {{ Object.keys(store.metricData).filter((n) => n.startsWith("process")) }} </p>
            </v-window-item>
            <v-window-item class="h-100" key="http" value="4">
              <p> {{ Object.keys(store.metricData).filter((n) => n.startsWith("http") || n.startsWith("ogen")) }} </p>
            </v-window-item>
            <v-window-item class="h-100" key="tracing" value="5">
              <p> {{ Object.keys(store.metricData).filter((n) => n.startsWith("promhttp") || n.startsWith("otel")) }} </p>
            </v-window-item>
          </v-container>
        </v-window>
      </v-card>
    </v-row>
  </Viewport>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue'
  import { useAppStore } from '../stores/app'

  const tab = ref(0)

  const store = useAppStore()

  onMounted(async () => {
    await store.fetchMetricData()
  })
</script>
