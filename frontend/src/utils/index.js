export function parseExecutions (executions) {
  const colors = []

  const columns = [
    { type: "string", id: "Job" },
    { type: "string", id: "Exec" },
    { type: "string", role: "tooltip" },
    { type: "date", id: "Start" },
    { type: "date", id: "End" },
  ]

  const colorMapping = {
    succeeded: '#77dd76',
    failed: '#ff6962',
  }

  const jobExecutions = []
  for (const execution of executions) {
    let tagLabel = 'info'
    switch (execution.execution_status) {
      case 'succeeded':
        tagLabel = 'success'
        break
      case 'failed-with-retry':
      case 'failed':
        tagLabel = 'danger'
        break
      case 'aborted':
      case 'timedout':
        tagLabel = 'warning'
        break
    }

    jobExecutions.push([
      `[${execution.rundeck_instance}] ${execution.rundeck_project}: ${execution.rundeck_job}`,
      `${execution.execution_id}`,
      `<div class="card">
        <header class="card-header">
        <p class="card-header-title">
        [${execution.rundeck_instance}]&nbsp;${execution.rundeck_job.replace(' ', '&nbsp;')}
        </p>
        </header>
        <div class="card-content">
          <div class="tags has-addons">
            <span class="tag">Status</span>
            <span class="tag is-${tagLabel}">${execution.execution_status}</span>
          </div>
          <ul>
            <li>Started:&nbsp;${execution.execution_start}</li>
            <li>Ended:&nbsp;${execution.execution_end}</li>
          </ul>
        </div>
      </div>`,
      new Date(execution.execution_start),
      new Date(execution.execution_end),
    ])

    colors.push(colorMapping[execution.execution_status] || '#FFA500')
  }

  return {
    colors,
    data: [columns, ...jobExecutions]
  }
}

// only returns executions that matches the selection
export function filterExecutions (executions, selection = {}) {
  const jobExecutions = []
  for (const execution of executions) {
    if ('instances' in selection) {
      if (!selection.instances.includes(execution.rundeck_instance)) {
        continue
      }
    }

    if ('projects' in selection) {
      if (!selection.projects.includes(execution.rundeck_project)) {
        continue
      }
    }

    if ('statuses' in selection) {
      if (!selection.statuses.includes(execution.execution_status)) {
        continue
      }
    }

    jobExecutions.push(execution)
  }

  return jobExecutions
}

export default {
  parseExecutions,
  filterExecutions
}
