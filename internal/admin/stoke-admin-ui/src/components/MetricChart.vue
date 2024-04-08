<template>
  <div class="ma-auto">
    <div class="chart-container" style="position: relative; height:70vh">
      <Line :data="chartData" :options="chartOptions" ref="chartRef"/>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { Line } from 'vue-chartjs'
  import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js'
  import { ChartDatasets, useAppStore } from "../stores/app"

  ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement)

  const store = useAppStore()

  const chartRef = ref({
    chart : {} as ChartJS
  })

  const chartData = {
    labels : [] as Date[],
    datasets: [],
  }

  const chartOptions = {
    responsive: true
  }

  function updateChartData(chartDatam : ChartDatasets) {
    let dates = [] as string[]
    const datasets = Object.keys(chartDatam).map((k, i) => {
      dates = chartDatam[k].times.map((t) => `${t.getHours()}:${t.getMinutes()}:${t.getSeconds()}` )
      return {
        label: chartDatam[k].displayName,
        // Have to copy to not trigger chartjs/pinia watchers again
        data: [ ...chartDatam[k].points],
        borderColor: colorWheel[i % 20],
      }
    })

    chartRef.value.chart.data.labels = dates
    chartRef.value.chart.data.datasets = datasets

    chartRef.value.chart.update('none')
  }

  watch(store.chartDatam, updateChartData, { deep : true })

  const colorWheel = [
    "#EF5350", "#337CA0", "#05A8AA", "#618985", "#DE369D",
    "#547AA5", "#EC9A29", "#9F956C", "#2D848A", "#6F5E76",
    "#FF8CC6", "#7BC950", "#837A75", "#F28F3B", "#A06CD5",
    "#6247AA", "#CA054D", "#5C946E", "#A53F2B", "#B3001B",
  ]

</script>
