package grok_patterns 

const (
  // Mcollective ...
  Mcollective string = `
MCOLLECTIVEAUDIT %{TIMESTAMP_ISO8601:timestamp}:
`
)
