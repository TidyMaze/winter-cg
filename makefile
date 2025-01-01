PHONY: run-local

build:
	go build src/winter/app.go

run-local:
	make build
	LOCAL_CG=true ./app

profile:
	make build
	make run-local & sleep 1
	go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=5

profile-mem:
	make build
	make run-local & sleep 1
	go tool pprof --alloc_space -http=:8081 http://localhost:6060/debug/pprof/heap?seconds=5
