package grok_vjeantet 

const (
  // Redis ...
  Redis string = `
REDISTIMESTAMP %{MONTHDAY} %{MONTH} %{TIME}
REDISLOG \[%{POSINT:pid}\] %{REDISTIMESTAMP:timestamp} \* 

`
)
