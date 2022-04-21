import React from 'react'
import renderer from 'react-test-renderer'
import { expect, test } from 'vitest'
import TimelineChart from './TimelineChart'

test('TimelineChart renders correctly', () => {
  const executions = [
    {
      id: 142,
      rundeck_instance: 'rundeck_dev',
      rundeck_project: 'my_project',
      rundeck_job: 'my job',
      rundeck_job_id: '339164fc-b1e8-4198-babd-fff0e9055e4c',
      execution_status: 'succeeded',
      execution_id: 1244,
      execution_start: '2022-04-05T09:10:00Z',
      execution_end: '2022-04-05T09:10:01Z',
      execution_permalink: 'https://rundeck.dev/project/my_project/execution/show/1244'
    },
    {
      id: 143,
      rundeck_instance: 'rundeck_dev',
      rundeck_project: 'my_project',
      rundeck_job: 'my job',
      rundeck_job_id: '339164fc-b1e8-4198-babd-fff0e9055e4c',
      execution_status: 'succeeded',
      execution_id: 1245,
      execution_start: '2022-04-05T09:10:00Z',
      execution_end: '2022-04-05T09:10:01Z',
      execution_permalink: 'https://rundeck.dev/project/my_project/execution/show/1245'
    },
  ]
  const tc = renderer.create(
    <TimelineChart executions={executions}></TimelineChart>
  ).toJSON()
  expect(tc).toMatchSnapshot('Rendered list without any facet selected')
})
