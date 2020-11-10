package grok_patterns 

const (
  // Redis ...
  Redis = `
REDISTIMESTAMP %{MONTHDAY} %{MONTH} %{TIME}
REDISLOG \[%{POSINT:pid}\] %{REDISTIMESTAMP:timestamp} \* 

`
)
