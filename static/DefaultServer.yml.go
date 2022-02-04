package static

const (
	// DefaultServer_yml ...
	DefaultServer_yml string = `server:
  #id: server-1 # default is host name
  #http:
  #  addr: 0.0.0.0:8585 # default is "0.0.0.0:8585"
  #  network: "" # default is ""
  #  timeout: 6s # default is "6s"
  #grpc:
  #  addr: 0.0.0.0:9595 # default is "0.0.0.0:9595"
  #  network: "" # default is ""
  #  timeout: 6s # default is "6s"
data:
  database:
    driver: mysql
    source: root:root@tcp(127.0.0.1:3306)/test
  redis:
    addr: 127.0.0.1:6379
    read_timeout: 0.2s
    write_timeout: 0.2s
`
)
