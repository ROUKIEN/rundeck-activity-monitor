import React from 'react'
import renderer from 'react-test-renderer'
import { expect, test } from 'vitest'
import TimeSearchDropdown from './TimeSearchDropdown'

test('emits a custom time range', () => {
  const fn = vi.fn(payload => {
    expect(payload).toHaveProperty('start', new Date('2022-01-01T00:00:00.000Z'))
    expect(payload).toHaveProperty('end', new Date('2022-01-02T00:00:00.000Z'))
  })

  const tc = renderer.create(
    <TimeSearchDropdown
      beginDate={new Date('2022-04-10T00:00:00')}
      endDate={new Date('2022-04-11T00:00:00')}
      onConfirm={fn}
    ></TimeSearchDropdown>
  )

  expect(tc.toJSON()).toMatchSnapshot('Defaults to since mode')

  const selectOption = tc.root.findAllByType('option')
  expect(selectOption).toHaveLength(4)

  const modeRadioInputs = tc.root.findAllByProps({
    type:'radio'
  })
  expect(modeRadioInputs).toHaveLength(2)
  modeRadioInputs[0].props.onClick({target: {value: 'range'}})

  const rangeInputs = tc.root.findAllByProps({
    type: 'datetime-local'
  })
  rangeInputs[0].props.onChange({
    target: {
      value: '2022-01-01T00:00',
      validity: {
        valid: true
      }
    }
  })

  rangeInputs[1].props.onChange({
    target: {
      value: '2022-01-02T00:00',
      validity: {
        valid: true
      }
    }
  })

  expect(tc.toJSON()).toMatchSnapshot('Switched to range mode with custom range')

  const searchButton = tc.root.findByProps({
    className: 'button is-primary'
  })

  searchButton.props.onClick()

  expect(fn).toHaveBeenCalledOnce()
})

test('emits a time range from the past 4 hours in since mode', () => {
  const fn = vi.fn(payload => {
    expect(payload).toHaveProperty('start')
    expect(payload).toHaveProperty('end')

    const { start, end } = payload
    expect(start.getHours()).toEqual(end.getHours() - 4)
  })

  const tc = renderer.create(
    <TimeSearchDropdown
      beginDate={new Date('2022-04-10T00:00:00')}
      endDate={new Date('2022-04-11T00:00:00')}
      onConfirm={fn}
    ></TimeSearchDropdown>
  )

  const select = tc.root.findByType('select')

  select.props.onChange({ target: { value: 4* 60 * 60 * 1000 }})

  const searchButton = tc.root.findByProps({
    className: 'button is-primary'
  })

  searchButton.props.onClick()

  expect(fn).toHaveBeenCalledOnce()
})
