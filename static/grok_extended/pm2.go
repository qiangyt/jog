package grok_extended 

const (
  // Pm2 ...
  Pm2 string = `PM2_DATESTAMP 20%{YEAR}-%{MONTHNUM}-%{MONTHDAY}T%{HOUR}:%{MINUTE}:%{SECOND}
PM2_MESSAGE (.*)
PM2_LOGLEVEL (log|error|debug|warn)
PM2_BOOTSTRAP %{PM2_DATESTAMP:timestamp}: PM2 %{PM2_LOGLEVEL:level}: %{PM2_MESSAGE:message}`
)
