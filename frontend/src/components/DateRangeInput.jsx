import React from 'react'
import 'bulma-calendar/dist/css/bulma-calendar.min.css'
import bulmaCalendar from 'bulma-calendar/dist/js/bulma-calendar.min.js'

class DateRangeInput extends React.Component {
  constructor (props) {
    super(props)

    this.ref = React.createRef()

    this.state = {
      calendar: null
    }
  }

  static getDerivedStateFromProps(props, state) {

    if (!state.calendar) {
      return state
    }

    state.calendar.startDate = props.startDateTime
    state.calendar.startTime = props.startDateTime

    state.calendar.endDate = props.endDateTime
    state.calendar.endTime = props.endDateTime
    state.calendar.save()

    return state
  }

  componentDidMount() {
    const calendar = bulmaCalendar.attach(this.ref.current, {
      startDate: this.props.startDateTime,
      startTime: this.props.startDateTime,
      endDate: this.props.endDateTime,
      endTime: this.props.endDateTime,
      maxDate: new Date(),
      type: 'datetime',
      isRange: true,
      allowSameDayRange: true,
      weekStart: 1,
      dateFormat: 'yyyy/MM/dd',
      showHeader: false,
      showTodayButton: false,
      closeOnSelect: false
    })[0]

    this.setState({
      calendar
    })

    calendar.on('select', () => {
      const start = this.state.calendar.startDate
      const end = this.state.calendar.endDate

      start.setHours(this.state.calendar.startTime.getHours())
      start.setMinutes(this.state.calendar.startTime.getMinutes())
      start.setSeconds(this.state.calendar.startTime.getSeconds())

      end.setHours(this.state.calendar.endTime.getHours())
      end.setMinutes(this.state.calendar.endTime.getMinutes())
      end.setSeconds(this.state.calendar.endTime.getSeconds())
      this.props.onChange({
        start,
        end
      })
    })
  }

  componentWillUnmount() {
    this.state.calendar.destroy()
    this.setState({
      calendar: null
    })
  }

  render () {
    return <input type="date" ref={this.ref}></input>
  }
}

export default DateRangeInput
