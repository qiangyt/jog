package main

// ConfigDefaultYAML ...
const ConfigDefaultYAML = `
output:
  pattern: "${timestamp} ${level} <${thread}> ${logger}: ${message} ${others} ${stacktrace}"
  compress-logger-name: true

  colors:
    line-no: FgDefault
    timestamp: FgDefault
    version: FgDefault
    message: FgDefault
    logger: FgDefault
    thread: FgDefault
    stack-trace: FgDefault
    started-line: FgGreen, OpBold

    debug: FgBlue,OpBold
    info: FgBlue,OpBold
    error: FgRed,OpBold
    warn: FgYellow,OpBold
    trace: FgBlue,OpBold
    fine: FgCyan,OpBold
    fatal: FgRed,OpBold

    raw: FgDefault
    others-name: FgDefault,OpBold
    others-separator: FgDefault
    others-value: FgDefault

  started-line: "Started Application in"

logstash:
  field-names:
    line-no:
    timestamp: "@timestamp"
    version: "@version"
    message: message
    logger: logger_name
    thread: thread_name
    level: level
    stack-trace: stack_trace
`
