import React from 'react'
import renderer from 'react-test-renderer'
import { expect, test } from 'vitest'
import TimelineChart from './TimelineChart'

test('TimelineChart renders correctly', () => {
  const data = []
  const tc = renderer.create(
    <TimelineChart data={data}></TimelineChart>
  )
  expect(tc).toBeDefined()
})
