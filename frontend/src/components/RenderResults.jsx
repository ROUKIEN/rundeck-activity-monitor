import React from 'react'
import TimelineChart from './TimelineChart'

function RenderResults (props) {
  if (props.executions == undefined) {
    return (
      <div className="message">
        <div className="message-body">Loading Executions...</div>
      </div>
    )
  } else {
    if (props.executions.length > 0) {
      return (
          <TimelineChart {...props}></TimelineChart>
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
