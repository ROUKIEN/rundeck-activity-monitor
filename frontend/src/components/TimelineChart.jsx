import React from 'react'
import { Chart } from 'react-google-charts'

function TimelineChart({ data, colors }) {
  return (
    <Chart
      chartType="Timeline"
      data={data}
      options={
        colors={colors}
      }
      width="100%"
      height="600px"
    />
  )
}

export default TimelineChart
