import React, { useState } from 'react'

// this component is a bit tricky because we have to deal with the non-UTC timezone
const DateTimeInput = ({ value, label, maxDate = null, minDate = null, onChange = () => {} }) => {
  const tzOffset = (new Date()).getTimezoneOffset() * 60000
  const val = new Date(new Date(new Date(value) - tzOffset)).toISOString().substring(0, 16)
  const [localValue, setLocalValue] = useState(val)

  const props = {}
  if (maxDate) {
    props.max = maxDate
  }

  if (minDate) {
    props.min = minDate
  }

  const handleChange = e => {
    setLocalValue(e.target.value)

    if (e.target.validity.valid) {
      let evt = `${e.target.value}:00Z`
      if (e.target.value == '') {
        evt = null
      } else {
        const tzOffset = (new Date()).getTimezoneOffset() * 60000
        const newDate = new Date(new Date(evt) - Math.abs(tzOffset))
        evt = new Date(newDate.toISOString().substring(0, 16)+':00Z')
      }
      onChange(evt)
    }
  }

  return (
    <div className="field">
      <label className="label" htmlFor={`dti-${label}`}>{label}</label>
      <div className="control">
      <input
        id={`dti-${label}`}
        value={localValue}
        className="input is-fullwidth datetime-input"
        type="datetime-local"
        onChange={e => handleChange(e)}
        data-testid={`datetimeinput-${label}`}
        {...props}
      ></input>
      </div>
    </div>
  )
}

export default DateTimeInput
