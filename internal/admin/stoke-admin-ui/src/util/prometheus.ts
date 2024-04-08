/*
  metric_name [
    "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } [ "," ] "}"
  ] value [ timestamp ]
*/

const valueRegex = /^(?<name>[\w.]+)({(?<tags>[\w\/\-\{\}+.,"=]+)})? (?<metric>[\w\-.+]+)\s?(?<ts>\d+)?/
const commentRegex = /^# (?<ltype>HELP|TYPE) (?<name>\w+) (?<value>.*)$/

export interface MetricDescription {
  help : string
  type : string
}

export interface MetricValue {
  tags : object,
  value: number
}

export interface MetricData {
  description: MetricDescription
  values: MetricValue[]
}

export interface MetricDataMap {
  [name : string] : MetricData
}

export function parseMetricData(fullText: string) : MetricDataMap {
  let metricData : MetricDataMap = {}
  const lines = fullText.split('\n')
  const subMetrics = ["_bucket", "_sum", "_count"]
  lines.forEach((l) => {
    const comment = commentRegex.exec(l)
    const value = valueRegex.exec(l)
    let name = ""
    if( comment?.groups ) {
      name = comment.groups.name
    } else if( value?.groups ) {
      name = value.groups.name
    } else if (l === "" ) {
      return
    }else {
      console.error(`Unmatched line! ${l}`)
      return
    }

    const isSubMetric = subMetrics.map((s: string) => name.endsWith(s))
    let addPart = ""

    if( isSubMetric.reduce( (r: boolean, b: boolean) => r || b, false ) ) {
      isSubMetric.forEach( (matches: boolean , i: number) => {
        if (matches) {
          addPart = subMetrics[i].substring(1)
        }
      })
      name = subMetrics.reduce((r : string, s: string) => r.replace(s, ""), name)
    }

    if( ! metricData[name] ) {
      metricData[name] = {
        description: {} as MetricDescription,
        values : [] as MetricValue[]
      } as MetricData
    }

    if ( comment?.groups ) {
      if( comment.groups.ltype === "HELP") {
        metricData[comment.groups.name].description.help = comment.groups.value
      }
      if( comment.groups.ltype === "TYPE") {
        metricData[comment.groups.name].description.type = comment.groups.value
      }
      return
    }

    if ( value?.groups ) {
      metricData[name].values.push({
        tags: parseTags(value.groups["tags"], addPart),
        value: Number.parseFloat(value.groups["metric"]),
      })
      return
    }
  })
  return metricData
}

function parseTags(tags: string, addPart : string) : object {
  let tagObj = {}
  if( addPart !== "" ){
    tagObj["part"] = addPart
  }
  if( tags && tags.includes(",") ) {
    tags.split(",").forEach((p) => {
      const kv = p.split("=")
      tagObj[kv[0]] = kv[1].substring(1, kv[1].length - 1)
    })
  }
  return tagObj
}
