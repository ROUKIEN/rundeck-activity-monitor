import React from 'react'

const TimeSinceInput = ({ onChange = () => {} }) => {
  const ONE_HOUR_AS_MS = 60*60*1000
  return (
    <span className="navbar-item">
      <div className="select is-fullwidth">
        <select onChange={ e => onChange(e.target.value) }>
          <option value={ONE_HOUR_AS_MS}>Past 1 hour</option>
          <option value={4*ONE_HOUR_AS_MS}>Past 4 hours</option>
          <option value={12*ONE_HOUR_AS_MS}>Past 12 hours</option>
          <option value={24*ONE_HOUR_AS_MS}>Past 1 day</option>
        </select>
      </div>
    </span>
  )
}

export default TimeSinceInput
