import React from 'react'
import DateTimeInput from './DateTimeInput'

const TimeRangeInput = ({ beginDate, endDate, onBeginSelect = () => {}, onEndSelect = () => {} }) => {
  return (
  <>
    <span className="navbar-item">
      <DateTimeInput
        value={beginDate}
        label="From"
        onChange={date => onBeginSelect(date)}
      ></DateTimeInput>
    </span>
    <span className="navbar-item">
      <DateTimeInput
        value={endDate}
        label="To"
        onChange={date => onEndSelect(date)}
      ></DateTimeInput>
    </span>
  </>
  )
}

export default TimeRangeInput
