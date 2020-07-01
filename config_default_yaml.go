package main

// ConfigDefaultYAML ...
const ConfigDefaultYAML = `
output:
  pattern: "${timestamp} ${level} <${thread}> ${logger}: ${message} ${others} ${stacktrace}"
  compress-logger-name: true

  colors:
    index: FgDefault, OpBold

    timestamp: FgDefault
    version: FgDefault
    message: FgDefault
    logger: FgDefault
    thread: FgDefault
    stack-trace: FgDefault
    started-line: FgGreen, OpBold
    pid: FgDefault
    host: FgDefault
    file: FgDefault
    method: FgDefault
    line: FgDefault

    levels:
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
  started-line-append: "\n\n"

input:
  ignore-conversion-error: true
  field-names:
    timestamp: "@timestamp"
    version: "@version"
    message: message
    logger: logger_name
    thread: thread_name
    level: level
    stack-trace: stack_trace
    pid: pid
    host: host
    file: file
    method: method
    line: line
`
