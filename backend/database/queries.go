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
	rundeck_project VARCHAR(80) NOT NULL,
	rundeck_job VARCHAR(80) NOT NULL,
	rundeck_job_id VARCHAR(80) NOT NULL,
	execution_status VARCHAR(20) NOT NULL,
	execution_id VARCHAR(20) NOT NULL,
	execution_start TIMESTAMP NOT NULL,
	execution_end TIMESTAMP,
	execution_permalink VARCHAR(255) NOT NULL
);`
