PHONY: run-local

build:
	go build src/winter/app.go

run-local:
	make build
	LOCAL_CG=true ./app

# to run profiler on this app
profile:
	# run the app and start pprof on port 6060
	make build
	make run-local & sleep 1
	# run the profiler on port 6060 (server) and 8080 (cpu) and 8081 (mem)
	go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=5
	#go tool pprof --alloc_space -http=:8081 http://localhost:6060/debug/pprof/heap?seconds=5
