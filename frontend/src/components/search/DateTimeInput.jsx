import React, { useState } from 'react'

const DateTimeInput = ({ value, label, maxDate = null, minDate = null, onChange = () => {} }) => {
  const tzOffset = (new Date()).getTimezoneOffset() * 60000
  const val = (new Date(value)).toISOString().substring(0, 16)
  const [localValue, setLocalValue] = useState(val)

  const valWithTz = (new Date(value - tzOffset)).toISOString().substring(0, 16)

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
        evt = new Date(evt)
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
        title={valWithTz}
        value={valWithTz}
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
