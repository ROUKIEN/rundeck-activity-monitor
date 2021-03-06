package database

const SQL_RAD_TABLE_EXISTS = `
SELECT EXISTS (
	SELECT FROM
		pg_tables
	WHERE
		tablename = 'rad_executions_history'
);`

const SQL_RAD_TABLE_STRUCTURE = `
CREATE TABLE rad_executions_history (
	id BIGSERIAL PRIMARY KEY,
	rundeck_instance VARCHAR(80) NOT NULL,
	rundeck_project VARCHAR(120) NOT NULL,
	rundeck_job VARCHAR(255) NOT NULL,
	rundeck_job_id VARCHAR(80) NOT NULL,
	execution_status VARCHAR(20) NOT NULL,
	execution_id INT NOT NULL,
	execution_start TIMESTAMP NOT NULL,
	execution_end TIMESTAMP,
	execution_permalink VARCHAR(255) NOT NULL
);`

const SQL_FIND_EXECUTION = `
SELECT
	*
FROM rad_executions_history
WHERE rundeck_instance = $1
AND	rundeck_job_id = $2
AND execution_id = $3
;`

const SQL_INSERT_EXECUTION = `
INSERT INTO rad_executions_history (
	rundeck_instance,
	rundeck_project,
	rundeck_job,
	rundeck_job_id,
	execution_status,
	execution_id,
	execution_start,
	execution_end,
	execution_permalink
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9
);`

const SQL_FIND_EXECUTIONS_BY_DATE = `
SELECT
	id,
	rundeck_instance,
	rundeck_project,
	rundeck_job,
	rundeck_job_id,
	execution_status,
	execution_id,
	execution_start,
	execution_end,
	execution_permalink
FROM rad_executions_history
WHERE
	execution_start >= $1
AND
	execution_end < $2
ORDER BY execution_start;`

const SQL_FIND_PROJECTS_ON_TIMERANGE = `
SELECT
	DISTINCT(rundeck_project) AS projects
FROM rad_executions_history
WHERE
	execution_start >= $1
AND
	execution_end < $2;`

const SQL_FIND_INSTANCES_ON_TIMERANGE = `
SELECT
	DISTINCT(rundeck_instance) AS instances
FROM rad_executions_history
WHERE
	execution_start >= $1
AND
	execution_end < $2;`

const SQL_FIND_STATUSES_ON_TIMERANGE = `
SELECT
	DISTINCT(execution_status) AS statuses
FROM rad_executions_history
WHERE
	execution_start >= $1
AND
	execution_end < $2;`
