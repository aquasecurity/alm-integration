module github.com/aquasecurity/postee/ui/backend

go 1.16

require (
	github.com/aquasecurity/postee v1.1.1
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/securecookie v1.1.1
	go.etcd.io/bbolt v1.3.5
)

replace github.com/aquasecurity/postee => ../../
