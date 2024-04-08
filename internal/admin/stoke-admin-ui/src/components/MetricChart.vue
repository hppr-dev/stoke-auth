<template>
  <div class="chart-container" style="position: relative; height:70vh">
    <Line :data="chartData" :options="chartOptions" ref="chartRef"/>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted } from 'vue'
  import { Line } from 'vue-chartjs'
  import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement } from 'chart.js'

  ChartJS.register(Title, Tooltip, Legend, LineElement, CategoryScale, LinearScale, PointElement)

  const chartRef = ref({
    chart : {} as ChartJS
  })

  let chartData = {
    labels: [1, 2, 3, 4, 5, 6, 7],
    datasets: [
      {
        label: "Metrics",
        data: [10, 12, 14, 32, 52, 34],
        borderColor: "rgb(75, 192, 192)",
      },
    ],
  }

  const chartOptions = {
    responsive: true
  }

  onMounted(() => {
    setTimeout(() => {
      chartRef.value.chart.data.datasets[0].data.push(5)
      chartRef.value.chart.update()
    }, 1000)
  })

</script>
