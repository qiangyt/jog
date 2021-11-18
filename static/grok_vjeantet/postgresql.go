package grok_vjeantet 

const (
  // Postgresql ...
  Postgresql string = `# Default postgresql pg_log format pattern
POSTGRESQL %{DATESTAMP:timestamp} %{TZ} %{DATA:user_id} %{GREEDYDATA:connection_id} %{POSINT:pid}

`
)
