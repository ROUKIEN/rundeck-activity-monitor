import React from 'react'
import TimelineChart from './TimelineChart'

function RenderResults ({ executions }) {
  if (executions == undefined) {
    return (
      <div className="message">
        <div className="message-body">Loading Executions...</div>
      </div>
    )
  } else {
    if (executions.length > 0) {
      return (
          <TimelineChart executions={executions}></TimelineChart>
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
