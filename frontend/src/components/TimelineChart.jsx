import React from 'react'
import { Chart } from 'react-google-charts'
import { parseExecutions } from '../utils'

const Loader = () =>
  <div className="message">
    <div className="message-body">Loading, please wait...</div>
  </div>

const FailedToLoad = () =>
  <div className="message">
    <div className="message-body">Failed to render executions.</div>
  </div>

function TimelineChart({ executions }) {
  const { data } = parseExecutions(executions)

  return (
    <Chart
      chartType="Timeline"
      data={data}
      options={
        {
          timeline: {
            showBarLabels: false
          }
        }
      }
      width="100%"
      height="600px"
      loader={<Loader></Loader>}
      errorElement={<FailedToLoad></FailedToLoad>}
      chartEvents={
        [
          {
            eventName: 'select',
            callback: (e) => {
              const selection = e.chartWrapper.getChart().getSelection()
              const execution = executions[selection[0].row]
              window.open(execution.execution_permalink, '_blank').focus()
            }
          }
        ]
      }
    />
  )
}

export default TimelineChart
