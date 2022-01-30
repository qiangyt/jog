package static

const (
	// DefaultServer_yml ...
	DefaultServer_yml string = `listen-on:
- protocol: http
  ip: 127.0.0.1
  port: 8585
  timeout: 6s
- protocol: grpc
  ip: 127.0.0.1
  port: 9595
  timeout: 6s
users:
- name: admin
  password: pwd
`
)
