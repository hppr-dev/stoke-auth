/*
  metric_name [
    "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } [ "," ] "}"
  ] value [ timestamp ]
*/

const valueRegex = /^(?<name>\w+)({(?<tags>[\w+.,_"=\/]*)})? (?<metric>[\w.-]+)\s?(?<ts>\d+)?$/
const commentRegex = /^# (?<ltype>HELP|TYPE) (?<name>\w+) (?<value>.*)$/

interface MetricDescription {
  help : string
  type : string
}

interface MetricValue {
  tags : object,
  value: number
}

interface MetricData {
  description: MetricDescription
  values: MetricValue[]
}

export interface MetricDataMap {
  [name : string] : MetricData
}

export function parseMetricData(fullText: string) : MetricDataMap {
  let metricData : MetricDataMap = {}
  const lines = fullText.split('\n')
  lines.forEach((l) => {
    const comment = commentRegex.exec(l)
    const value = valueRegex.exec(l)
    if ( comment?.groups ) {
      if( ! metricData[comment.groups.name] ) {
        metricData[comment.groups.name] = {
          description: {} as MetricDescription,
          values : [] as MetricValue[]
        } as MetricData
      }
      if( comment.groups.ltype === "HELP") {
        metricData[comment.groups.name].description.help = comment.groups.value
      }
      if( comment.groups.ltype === "TYPE") {
        metricData[comment.groups.name].description.type = comment.groups.value
      }
      return
    }
    if ( value?.groups ) {
      if( ! metricData[value.groups.name] ) {
        metricData[value.groups.name] = {
          description: {} as MetricDescription,
          values : [] as MetricValue[]
        } as MetricData
      }
      metricData[value.groups.name].values.push({
        tags: parseTags(value.groups["tags"]),
        value: Number.parseFloat(value.groups["metric"]),
      })
    }
  })
  return metricData
}

function parseTags(tags: string) : object {
  let tagObj = {}
  if ( tags && tags.includes(",") ) {
    tags.split(",").forEach((p) => {
      const kv = p.split("=")
      tagObj[kv[0]] = kv[1].substring(1, kv[1].length - 1)
    })
  }
  return tagObj
}
