# Rundeck Activity Monitor (RAM)

`RAM` is a tool to monitor rundeck jobs executions across multiple rundeck instances. It is composed of two main components: a scraper and a UI.

The scraper is in charge of requesting the different rundeck instances to get the execution history, and the UI is rendering the executions across all instances in a swimlane chart.

# Building

1. build the frontend app (which is basically a react SPA)
2. build the backend app and embedd the frontend inside it

## Running locally

1. Start the postgresql container (it is recommanded to run the one described in the `docker-compose.yml` file)
2. Ensure the schema is up to date by running `ram database update`
3. Start developping
