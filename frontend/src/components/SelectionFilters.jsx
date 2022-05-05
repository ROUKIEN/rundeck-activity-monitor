import React from 'react'
import DropdownFilter from './DropdownFilter'

const isChecked = (selection, property, value) =>
  property in selection === true && selection[property].includes(value)

const SelectionFilters = ({ filters, selection, onChange }) => {

  const handleChange = (field, value) => {
    if (field in selection === true) {
      if (selection[field].includes(value)) {
        selection[field] = selection[field].filter(item => item !== value)
        if (selection[field].length === 0) {
          delete selection[field]
        }
      } else {
        selection[field].push(value)
      }
    } else {
      selection[field] = [
        value
      ]
    }

    onChange(selection)
  }

  return (
    <>
      <p className="menu-label">
        Search
      </p>
      <ul className="menu-list">
        { Object.keys(filters).map(filter => {
          return (
            <DropdownFilter
              filter={filter}
              count={filters[filter].length}
              key={filter}>
              { filters[filter].map(value =>
                <li key={value}>
                  <label className="checkbox">
                    <input
                      key={Math.random()}
                      type="checkbox"
                      defaultChecked={isChecked(selection, filter, value)}
                      onChange={e => handleChange(filter, value, e)}
                    ></input>&nbsp;{value}
                  </label>
                </li>
              )}
            </DropdownFilter>
          )
        })}
        <li>
          <button
            className="button is-small is-primary is-fullwidth"
            data-testid="reset-search-btn"
            onClick={ () => {
              onChange({})
            }}
          >Reset Filters&nbsp;<span className="delete"></span></button>
        </li>
      </ul>
    </>
  )
}

export default SelectionFilters
