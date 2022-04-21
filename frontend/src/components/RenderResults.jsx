import React from 'react'
import TimelineChart from './TimelineChart'

function RenderResults ({ data, colors }) {
  if (data == undefined) {
    return (
      <div className="message">
        <div className="message-body">Loading Executions...</div>
      </div>
    )
  } else {
    if (data.length > 1) {
      return (
          <TimelineChart colors={colors} data={data}></TimelineChart>
      )
    } else {
      return (
        <div className="message is-warning">
          <div className="message-body">No results !</div>
        </div>
      )
    }
  }
}

export default RenderResults
