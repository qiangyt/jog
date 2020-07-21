package config

// DefaultYAML ...
const DefaultYAML = `
colorization: true
pattern: "${timestamp} ${level} <${thread}> ${logger}: ${message} {${others}} ${stacktrace}"
startup-line:
  color: FgGreen, OpBold
  contains: "Started Application in"

line-no:
  print: true
  color: FgDefault

unknown-line:
  print: true
  color: FgBlue

prefix:
  print: true
  color: FgBlue

fields:
  others:
    name:
      color: FgDefault,OpBold
    separator:
      label: =
      color: FgDefault
    value:
      color: FgDefault
  level:
    alias: "level, @level, severity, @severity"
    case-sensitive: false
    enums:
      case-sensitive: false
      default: WARN
      DEBUG:
        alias: debug,20
        color: FgBlue
      INFO:
        alias: info,30
        color: FgBlue,OpBold
      ERROR:
        alias: error,err,critical,50
        color: FgRed
      WARN:
        alias: warn,warning,40
        color: FgYellow,OpBold
      TRACE:
        alias: trace,10
        color: FgCyan
      FINE:
        alias: fine
        color: FgCyan
      FATAL:
        alias: fatal,60
        color: FgRed,OpBold
  app:
    alias: "name, @name, @app"
    case-sensitive: false
    color: FgDefault
  class:
    alias: "classname, class-name, @class_name, @classname, @class-name, @class_name"
    case-sensitive: false
    color: FgDefault
    compress-prefix:
      enabled: true
      separators: ., /, \
      action: remove
  env:
    alias: "environment, @env, @environment"
    case-sensitive: false
    color: FgDefault
  file:
    alias: "src, source, filename, file-name, file_name, filepath,file-path, file_path, @src, @source, @file, @filename, @file-name, @file_name, @filepath, @file-path, @file_path"
    case-sensitive: false
    color: FgDefault
    compress-prefix:
      enabled: true
      separators: /, \
      action: remove
  host:
    alias: "hostname, host-name, host_name, @host, @hostname, @host-name, @host_name"
    case-sensitive: false
    color: FgDefault
  line:
    alias: "lineno, line-no, line_no, linenum, line-num, line_num, linenumber, line-number, line_number, @lineno, @line-no, @line_no, @linenum, @line-num, @line_num, @linenumber, @line-number, @line_number"
    case-sensitive: false
    color: FgDefault
  logger:
    alias: "id, logger_name, logger-name, loggername, @id, @logger_name, @logger-name, @loggername, @logger"
    case-sensitive: false
    color: FgDefault
    compress-prefix:
      enabled: true
      separators: . , /
      white-list: com.wxcount
      action: remove-non-first-letter
  message:
    alias: "msg, @msg, @message"
    case-sensitive: false
    color: FgDefault
  method:
    alias: "methodname, method-name, method_name, func, funcname, func-name, func_name, function, functionname, function-name, function_name,  @method, @methodname, @method-name, @method_name, @func, @funcname, @func-name, @func_name, @function, @functionname, @function-name, @function_name"
    case-sensitive: false
    color: FgDefault
  pid:
    alias: "process, process-id, processid, process_id, @pid, @process, @process-id, @processid, @process_id"
    case-sensitive: false
    color: FgDefault
  request:
    alias: "req, @req, @request"
    case-sensitive: false
    color: FgDefault
  response:
    alias: "res, resp, @res, @resp, @response"
    case-sensitive: false
    color: FgDefault
  stack-trace:
    alias: "err, error, stack, stack_trace, stacktrace, @stack, @stack_trace, @stack-trace, @stacktrace"
    case-sensitive: false
    color: FgDefault
    before: "\nStack trace: \n"
  thread:
    alias: "thread_name, thread-name, threadname, @thread, @thread_name, @thread-name, @threadname"
    case-sensitive: false
    color: FgDefault
  timestamp:
    alias: "time, date, datetime, date-time, date_time, @time, @timestamp, @date, @datetime, @date-time, @date_time"
    case-sensitive: false
    color: FgDefault
  user:
    alias: "usr, username, user-name, user_name, @usr, @username, @user-name, @user_name, @user"
    case-sensitive: false
    color: FgDefault
  version:
    alias: "ver, @ver, @version"
    case-sensitive: false
    color: FgDefault
`
