package main

// ConfigDefaultYAML ...
const ConfigDefaultYAML = `
output:
  pattern: "${prefix} | ${timestamp} ${level} <${thread}> ${logger}: ${message} ${others} ${stacktrace}"
  compress-logger-name: true

  colors:
    index: FgDefault, OpBold
    prefix: FgCyan
    app: FgDefault
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
    app: "name, Name, app, App, @name, @Name, @app, @App"
    timestamp: "time, Time, timestamp, Timestamp, @time, @Time, @timestamp, @Timestamp"
    version: "version, Version, @version, @Version"
    message: "msg, message, Message, @msg, @message, @Message"
    logger: "id, Id, ID, logger_name, logger-name, loggerName, LoggerName, logger, Logger, @id, @Id, @ID, @logger_name, @logger-name, @loggerName, @LoggerName, @logger, @Logger"
    thread: "thread_name, thread-name, threadName, ThreadName, thread, Thread, @thread, @Thread"
    level: "level, Level, @level, @Level"
    stack-trace: "stack_trace, stack-trace, stackTrace, StackTrace, stack, Stack, @stack_trace, @stack-trace, @stackTrace, @StackTrace, @stack, @Stack"
    pid: "pid, PID, @pid, @PID"
    host: "host, Host, @host, @Host, hostname, Hostname, hostName, HostName, @Hostname, @Hostname, @hostName, @HostName"
    file: "file, File, @file, @File"
    method: "method, Method, @method, @Method"
    line: "line, Line, @line, @Line"
`
