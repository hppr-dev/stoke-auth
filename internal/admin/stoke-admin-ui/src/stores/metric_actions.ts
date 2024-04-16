import { parseMetricData } from '../util/prometheus'

/*
 * dataset format
    labels: [1, 2, 3, 4, 5, 6, 7],
    datasets: [
      {
        label: "Metrics",
        data: [10, 12, 14, 32, 52, 34],
        borderColor: "rgb(75, 192, 192)",
      },
    ],
*/

export interface ChartData {
  displayName: string
  times: Date[]
  points : number[]
}

export const metricActions = {
  toggleMetric: function({ metricName, displayName } : { metricName: string, displayName: string}) {
    if( ! this.trackedMetrics.includes(metricName) ) {
      this.chartDatam[metricName] = {
        displayName: displayName,
        times: [ ],
        points: [ ],
      }
    } else {
      delete this.chartDatam[metricName]
    }
  },
  fetchMetricData: async function() {
    const response = await fetch(`${this.api_url}/metrics`, {
        method: "GET",
        headers: {
          "Content-Type" : "text/plain; version=0.0.4",
          "Authorization" : `Bearer ${this.token}`,
        },
      }
    )
    const result = await response.text();
    this.metricData = parseMetricData(result);
  },
  metricRefresh: async function() {
    await this.fetchMetricData()
    await this.fetchLogText()
    if( ! this.metricsPaused ) {
      this.metricTimeoutID = window.setTimeout(this.metricRefresh, this.metricRefreshTime)
      const timeNow = new Date();
      Object.keys(this.chartDatam).forEach((m : string) => {
        if(this.chartDatam[m].points.length >= this.maxPoints ) {
          this.chartDatam[m].times.splice(0, 1)
          this.chartDatam[m].points.splice(0, 1)
        }
        this.chartDatam[m].times.push(timeNow)
        this.chartDatam[m].points.push(this.metricData[m].values[0].value)
      })
    }
  },
  setMetricRefresh: function(millis : number) {
    this.metricRefreshTime = millis
  },
  setMaxPoints: function(max: number) {
    this.maxPoints = max
  },
  toggleMetricPaused: function() {
    this.metricsPaused = !this.metricsPaused
    if( this.metricsPaused ) {
      window.clearTimeout(this.metricTimeoutID)
    } else {
      if( this.trackedMetrics.length > 0) {
        Object.keys(this.chartDatam).forEach((m : string) => {
          this.chartDatam[m].times.splice(0, this.chartDatam[m].times.length)
          this.chartDatam[m].points.splice(0, this.chartDatam[m].points.length)
        })

      }
      this.metricTimeoutID = setTimeout(this.metricRefresh, this.metricRefreshTime)
    }
  },
  fetchLogText: async function() {
    const response = await fetch(`${this.api_url}/metrics/logs`, {
        method: "GET",
        headers: {
          "Content-Type" : "text/plain; version=0.0.4",
          "Authorization" : `Bearer ${this.token}`,
        },
      }
    )
    this.logText = await response.text();
  }
}
