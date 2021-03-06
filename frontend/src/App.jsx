import { useEffect, useState } from 'react'

import logo from './logo.svg'

import { filterExecutions } from './utils/index'

import SelectionFilters from './components/SelectionFilters'
import RenderResults from './components/RenderResults'
import TimeSearchDropdown from './components/TimeSearchDropdown'

async function fetchFilters (options = {}) {
  const params = new URLSearchParams()
  for (const option in options) {
    params.append(option, options[option])
  }

  const response = await fetch(`/api/filters${options ? '?'+params.toString() : ''}`)
  const json = await response.json()

  return json
}

async function fetchExecutions (options = {}) {
  const params = new URLSearchParams()
  for (const option in options) {
    params.append(option, options[option])
  }

  const response = await fetch(`/api/executions${options ? '?'+params.toString() : ''}`)

  return await response.json()
}

function parseQueryString(setBeginDatefn, setEndDateFn) {
  const searchParams = new URLSearchParams(window.location.search)

  const params = {}

  if (searchParams.has('begin')) {
    const start = searchParams.get('begin')
    params.begin = start
    const newStartDate = new Date(parseInt(start))
    setBeginDatefn(newStartDate)
  }

  if (searchParams.has('end')) {
    const end = searchParams.get('end')
    params.end = end
    const newEndDate = new Date(parseInt(end))
    setEndDateFn(newEndDate)
  }

  return params
}

function App() {
  const [data, setData] = useState(undefined)
  const [filteredData, setFilteredData] = useState(undefined)
  const [filters, setFilters] = useState({})

  const beginDateTime = new Date()
  const endDateTime = new Date()
  beginDateTime.setHours(new Date().getHours() - 4)
  beginDateTime.setSeconds(0)
  endDateTime.setSeconds(0)
  const [beginDate, setBeginDate] = useState(beginDateTime)
  const [endDate, setEndDate] = useState(endDateTime)
  const [selectionFilters, setSelectionFilters] = useState({})

  const handleSelectionFiltersUpdate = (unfilteredData, newSelection) => {
    setSelectionFilters(newSelection)

    const filteredExecutions = filterExecutions(unfilteredData, newSelection)
    setFilteredData(filteredExecutions)
  }

  const handleRangeChange = async ({ start, end }) => {
    setBeginDate(start)
    setEndDate(end)

    const params = {
      begin: start.getTime(),
      end: end.getTime()
    }

    const searchParams = new URLSearchParams(window.location.search)
    searchParams.set("begin", params.begin)
    searchParams.set("end", params.end)
    const newRelativePathQuery = window.location.pathname + '?' + searchParams.toString()
    history.pushState(null, '', newRelativePathQuery)

    const filters = await fetchFilters(params)
    setFilters(filters)

    const executions = await fetchExecutions(params)
    setData(executions)
    const filteredExecutions = filterExecutions(executions, selectionFilters)
    setFilteredData(filteredExecutions)
  }

  useEffect(() => {
    const params = parseQueryString(setBeginDate, setEndDate)

    async function fetchFiltersAndExecutions(params) {
      const filters = await fetchFilters(params)
      setFilters(filters)

      const executions = await fetchExecutions(params)
      setData(executions)
      const filteredExecutions = filterExecutions(executions, selectionFilters)
      setFilteredData(filteredExecutions)
    }

    fetchFiltersAndExecutions(params)
    return () => {}
  }, [])

  return (
    <>
      <nav className="navbar" role="navigation" aria-label="main navigation">
        <div className="navbar-brand">
          <a className="navbar-item" href="/">
            <img src={logo} /> <span className="navbar-item has-text-weight-medium">Rundeck Activity Monitor</span>
          </a>
        </div>
        <div className="navbar-end">
          <TimeSearchDropdown
            beginDate={beginDate}
            endDate={endDate}
            onConfirm={range => handleRangeChange(range)}
          ></TimeSearchDropdown>
        </div>
      </nav>
      <section className="section is-small has-background-primary">
        <div className="container is-widescreen">
          <h1 className="title has-text-white">Rundeck Activity Monitor</h1>
          <h2 className="subtitle has-text-white">Having a bunch of rundeck instances ? RAM is the tool you need.</h2>
        </div>
      </section>
      <section className="section">
        <div className="container-fluid is-widescreen">
          <div className="columns is-widescreen">
            <div className="column is-one-fifth">
              <aside className="menu">
                <SelectionFilters
                  filters={filters}
                  selection={selectionFilters}
                  onChange={s => handleSelectionFiltersUpdate(data, s)}
                ></SelectionFilters>
              </aside>
            </div>
            <div className="column">
              <RenderResults
                executions={filteredData}
                start={beginDate}
                end={endDate}
              ></RenderResults>
            </div>
          </div>
        </div>
      </section>
    </>
  )
}

export default App
