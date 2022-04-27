import React, { useState } from 'react'
import TimeSinceInput from './search/TimeSinceInput'
import TimeRangeInput from './search/TimeRangeInput'

const TIME_MODE_RANGE = 'range'
const TIME_MODE_SINCE = 'since'

const TimeSearchDropdown = ({ beginDate, endDate, onConfirm = () => {} }) => {
  const [currentBeginDate, setCurrentBeginDate] = useState(beginDate)
  const [currentEndDate, setCurrentEndDate] = useState(endDate)
  const [since, setSince] = useState(0)
  const [timeMode, setTimeMode] = useState(TIME_MODE_SINCE)

  const onValueChange = e => {
    setTimeMode(e)
  }

  const handleSearch = () => {
    if (timeMode == TIME_MODE_SINCE) {
      const start = new Date(Date.now() - since)
      onConfirm({
        start,
        end: new Date(),
      })
    } else {
      onConfirm({
        start: currentBeginDate,
        end: currentEndDate,
      })
    }
  }

  return (
  <div className="navbar-item has-dropdown is-hoverable">
    <a className="navbar-link">
      Search
    </a>
    <div className="navbar-dropdown is-right">
      <span className="navbar-item">
        <div className="control">
          <label className="radio">
            <input
              type="radio"
              name="mode"
              value={TIME_MODE_SINCE}
              defaultChecked={timeMode === TIME_MODE_SINCE}
              onClick={e => onValueChange(e.target.value)}
              data-testid={`radio-${TIME_MODE_SINCE}`}
            ></input>&nbsp;Since
          </label>
          <label className="radio">
            <input
              type="radio"
              name="mode"
              value={TIME_MODE_RANGE}
              defaultChecked={timeMode === TIME_MODE_RANGE}
              onClick={e => onValueChange(e.target.value)}
              data-testid={`radio-${TIME_MODE_RANGE}`}
            ></input>&nbsp;Range
          </label>
        </div>
      </span>
      { timeMode == TIME_MODE_RANGE && (
        <TimeRangeInput
          beginDate={currentBeginDate}
          endDate={currentEndDate}
          onBeginSelect={v => setCurrentBeginDate(v)}
          onEndSelect={v => setCurrentEndDate(v)}
        ></TimeRangeInput>)
      }
      { timeMode == TIME_MODE_SINCE && <TimeSinceInput onChange={ e => setSince(e)}></TimeSinceInput> }
      <span className="navbar-item">
        <div className="field is-grouped is-grouped-centered">
          <p className="control">
            <button className="button is-primary" onClick={() => handleSearch()}>Search</button>
          </p>
        </div>
      </span>
    </div>
  </div>
  )
}

export default TimeSearchDropdown
