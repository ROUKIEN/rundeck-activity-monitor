import React, { useState } from 'react'
import Icon from '@mdi/react'
import { mdiArrowUp, mdiArrowDown } from '@mdi/js'

const expandIcon = ({ expanded }) =>
  <span className="icon is-small">
    { expanded === true ? <Icon path={mdiArrowUp}/> : <Icon path={mdiArrowDown}/>}
  </span>

const DropdownFilter = ({ filter, count, children }) => {
  const [expanded, setExpanded] = useState(true)
  return (
    <li>
      <button
        className="button is-small is-fullwidth"
        onClick={ () => setExpanded(!expanded) }
        style={{
          cursor: 'pointer'
        }}
      >
        <span className="is-capitalized">{filter} ({count})</span>
        <expandIcon expanded={expanded}/>
      </button>
      <ul>
        { expanded && children }
      </ul>
    </li>
  )
}

export default DropdownFilter
