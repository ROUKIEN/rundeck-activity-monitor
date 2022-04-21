import React from 'react'
import renderer from 'react-test-renderer'
import { expect, test, vi } from 'vitest'
import SelectionFilters from './SelectionFilters'

test('SelectionFilters renders correctly', () => {
  const filters = {
    instances: ["foo", "bar", "zee"],
    statuses: ["succeeded", "failed"]
  }
  const tc = renderer.create(
    <SelectionFilters filters={ filters } selection={ {} }></SelectionFilters>
  ).toJSON()
  expect(tc).toBeDefined()
  expect(tc).toMatchSnapshot('Rendered list without any facet selected')
})

test('SelectionFilters renders with selection', () => {
  const filters = {
    instances: ["foo", "bar", "zee"],
    statuses: ["succeeded", "failed"]
  }

  const selection = {
    instances: ["foo"]
  }

  const tc = renderer.create(
    <SelectionFilters filters={ filters } selection={ selection }></SelectionFilters>
  ).toJSON()
  expect(tc).toBeDefined()
  expect(tc).toMatchSnapshot('Rendered list with foo instance selected')
})

test('SelectionFilters triggers onChange when something gets selected', () => {
  const filters = {
    instances: ['foo', 'bar', 'zee'],
    statuses: ['succeeded', 'failed']
  }

  const selection = {
    instances: ['foo']
  }

  const fn = vi.fn()

  const tc = renderer.create(
    <SelectionFilters
      filters={ filters }
      selection={ selection }
      onChange={fn}
    ></SelectionFilters>
  )

  tc.root.findAllByType('input')[3].props.onChange()
  expect(fn).toHaveBeenCalledTimes(1)
  expect(fn).toHaveBeenCalledWith({
    instances: ['foo'],
    statuses: ['succeeded'],
  })
})

test('SelectionFilters triggers onChange with empty selection when reset filters button is clicked', () => {
  const filters = {
    instances: ['foo', 'bar', 'zee'],
    statuses: ['succeeded', 'failed']
  }

  const selection = {
    instances: ['foo', 'bar'],
    statuses: ['failed']
  }

  const fn = vi.fn()

  const tc = renderer.create(
    <SelectionFilters
      filters={ filters }
      selection={ selection }
      onChange={fn}
    ></SelectionFilters>
  )

  tc.root.findByType('button').props.onClick()
  expect(fn).toHaveBeenCalledTimes(1)
  expect(fn).toHaveBeenCalledWith({})
})

test('SelectionFilters onChange triggers unselected ', () => {
  const filters = {
    instances: ['foo', 'bar', 'zee'],
    statuses: ['succeeded', 'failed']
  }

  const selection = {
    instances: ['foo', 'bar'],
    statuses: ['failed']
  }

  const fn = vi.fn()

  const tc = renderer.create(
    <SelectionFilters
      filters={ filters }
      selection={ selection }
      onChange={fn}
    ></SelectionFilters>
  )

  tc.root.findAllByType('input')[0].props.onChange()
  expect(fn).toHaveBeenCalledTimes(1)
  expect(fn).toHaveBeenCalledWith({
    instances: ['bar'],
    statuses: ['failed']
  })
})
